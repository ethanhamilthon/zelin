import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import { EditorView } from '@codemirror/view';
import App from './App.svelte';
import { agentInfo, runAgent, stopAgent } from './agent';
import { currentSession, createSession, openSession, listSessions, saveSession } from './sessions';
import { call } from './bridge';
import type { AgentEvent, AgentResult, Execution } from './agent-types';
import { placeholderInfo } from './agent-types';

vi.mock('./agent', () => ({ runAgent: vi.fn(), agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn(), setTitleModel: vi.fn(), stopAgent: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn() }));
vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => false, listen: vi.fn(async () => () => {}) }));
vi.mock('./sound', () => ({ chime: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
let app: ReturnType<typeof mount> | undefined;
const execution: Execution = { id: 'call', code: 'bash("pwd", 10)', output: '"/tmp"', state: 'complete' };
const thinking: Execution = { id: 'think', kind: 'thinking', code: 'Plan the command', output: '', state: 'complete' };
const wait = async () => { await new Promise((resolve) => setTimeout(resolve, 0)); await tick(); };
const session = (id = 'first', draft = 'Hello **world**') => ({ id, title: `${id} title`, draft, messages: [], updated: new Date().toISOString() });
const $ = (selector: string) => document.querySelector<HTMLElement>(selector);

beforeEach(() => {
  vi.resetAllMocks();
  Element.prototype.scrollIntoView = vi.fn();
  Object.defineProperties(Range.prototype, {
    getClientRects: { configurable: true, value: () => [] },
    getBoundingClientRect: { configurable: true, value: () => new DOMRect() },
  });
  vi.mocked(currentSession).mockResolvedValue(session());
  vi.mocked(createSession).mockResolvedValue(session('new', 'Second task'));
  vi.mocked(openSession).mockImplementation(async (id) => session(id, ''));
  vi.mocked(listSessions).mockResolvedValue([session(), session('new')]);
  vi.mocked(saveSession).mockResolvedValue();
  vi.mocked(stopAgent).mockResolvedValue(true);
  vi.mocked(call).mockResolvedValue('');
  vi.mocked(agentInfo).mockResolvedValue({ ...placeholderInfo('Test Model'), model: 'test/model', effort: 'medium', efforts: ['none', 'medium'], contextWindow: 128000 });
});
afterEach(async () => { if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });
async function start() { app = mount(App, { target: document.body }); flushSync(); await wait(); }
function key(key: string, code: string, init: KeyboardEventInit = {}) { window.dispatchEvent(new KeyboardEvent('keydown', { key, code, cancelable: true, ...init })); }
function submit() { $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', code: 'Enter', bubbles: true, cancelable: true })); }
function draft(text: string) {
  const view = EditorView.findFromDOM($('.cm-editor')!)!;
  view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: text } });
}
function pending() {
  const runs: { finish: (result: AgentResult) => void; notify: (event: AgentEvent) => void }[] = [];
  vi.mocked(runAgent).mockImplementation((_, notify) => new Promise((finish) => { runs.push({ finish, notify }); }));
  return runs;
}
const button = (text: string) => [...document.querySelectorAll<HTMLButtonElement>('button')].find((item) => item.textContent?.includes(text))!;

describe('agent UI', () => {
  it('sends with Enter, shows thinking and code, and renders minimal Markdown', async () => {
    const runs = pending();
    await start(); submit(); await wait();
    expect($('article .text')?.textContent).toBe('Hello **world**');
    expect($('[aria-label="Request running"]')).not.toBeNull();
    draft('Next draft'); submit();
    expect(runAgent).toHaveBeenCalledOnce();
    runs[0].notify({ kind: 'execution', execution: thinking });
    runs[0].notify({ kind: 'execution', execution: { ...execution, state: 'running', output: '' } });
    await tick();
    expect($('[aria-label="Thinking"]')?.textContent).toContain('Plan the command');
    runs[0].finish({ content: '**Done** with `code`', executions: [thinking, execution] }); await wait();
    expect($('article.assistant strong')?.textContent).toBe('Done');
    expect($('article.assistant .text')?.textContent).toBe('Done with code');
    expect($('article [aria-label="Turn complete"]')).not.toBeNull();
    button('bash').click(); await tick();
    expect($('.output')?.textContent).toBe('"/tmp"');
    $('[aria-label="Fold code"]')!.click(); await tick();
    expect($('[aria-label="Unfold code"]')).not.toBeNull();
    expect($('.cm-content')?.textContent).toBe('Next draft');
    expect($('footer')?.textContent).toContain('Test Model');
    expect($('footer')?.textContent).toContain('medium');
    expect($('footer .right button')?.textContent).toBe('Commands Ctrl+K');
    expect($('footer [role="status"]')).toBeNull();
    expect(saveSession).toHaveBeenLastCalledWith('first', expect.arrayContaining([expect.objectContaining({ content: '**Done** with `code`' })]), 'Next draft', []);
  });

  it('preserves the next draft on failure', async () => {
    let fail!: (reason: Error) => void;
    vi.mocked(runAgent).mockImplementation(() => new Promise((_, reject) => { fail = reject; }));
    await start(); submit(); await wait(); draft('Keep this'); fail(new Error('Offline')); await wait();
    expect($('[role="alert"]')?.textContent).toBe('Offline');
    expect($('.cm-content')?.textContent).toBe('Keep this');
  });

  it('opens commands with a Russian layout and closes with Escape', async () => {
    await start(); key('л', 'KeyK', { ctrlKey: true }); await tick();
    expect($('[role="dialog"]')?.getAttribute('aria-label')).toBe('Commands');
    key('Escape', 'Escape');
    await vi.waitFor(() => expect($('[role="dialog"]')).toBeNull());
    expect(document.activeElement).toBe($('.cm-content'));
  });

  it('stops the current run with Ctrl+C', async () => {
    pending();
    await start(); submit(); await wait();
    key('с', 'KeyC', { ctrlKey: true });
    expect(stopAgent).toHaveBeenCalledWith('first');
  });

  it('runs sessions in parallel and toasts when a background session finishes', async () => {
    const runs = pending();
    await start(); submit(); await wait();
    key('n', 'KeyN', { ctrlKey: true }); await wait();
    expect($('.cm-content')?.textContent).toBe('Second task');
    submit(); await wait();
    expect(runAgent).toHaveBeenCalledTimes(2);
    expect(vi.mocked(runAgent).mock.calls.map((args) => args[2])).toEqual(['first', 'new']);
    key('s', 'KeyS', { ctrlKey: true }); await wait();
    expect(button('first title').textContent).toContain('Running');
    key('Escape', 'Escape'); await wait();
    runs[0].finish({ content: 'First done', executions: [] }); await wait(); await wait();
    expect($('.toasts button')?.textContent).toBe('first title · finished');
    $('.toasts button')!.click(); await wait();
    expect(openSession).toHaveBeenLastCalledWith('first');
    expect(agentInfo).toHaveBeenLastCalledWith('first');
    expect($('article.assistant .text')?.textContent).toBe('First done');
    expect($('[aria-label="Request running"]')).toBeNull();
  });

  it('edits the system prompt and saves with Ctrl+S', async () => {
    vi.mocked(call).mockImplementation(async (method) => method === 'SettingsService.Read' ? 'You are Zash.' : '');
    await start(); key('k', 'KeyK', { ctrlKey: true }); await tick();
    button('Instructions').click(); await tick();
    button('System Prompt').click(); await wait();
    expect($('footer[aria-label="System Prompt"]')?.textContent).toContain('Save Ctrl+S');
    const editor = [...document.querySelectorAll<HTMLElement>('.cm-content')].find((item) => item.getAttribute('aria-label') === 'System Prompt')!;
    expect(editor.textContent).toBe('You are Zash.');
    const view = EditorView.findFromDOM(editor.parentElement!.parentElement!)!;
    view.dispatch({ changes: { from: view.state.doc.length, insert: ' Be brief.' } });
    key('ы', 'KeyS', { ctrlKey: true }); await wait();
    expect(call).toHaveBeenCalledWith('SettingsService.Write', 'system', 'You are Zash. Be brief.');
    expect($('[aria-label="System Prompt"] [role="status"]')?.textContent).toBe('Saved');
  });

  it('hides code blocks from settings', async () => {
    vi.mocked(currentSession).mockResolvedValue({ ...session(), messages: [{ role: 'assistant', content: 'Done', executions: [execution, thinking] }] });
    await start();
    expect($('[aria-label="JavaScript execution"]')).not.toBeNull();
    key('k', 'KeyK', { ctrlKey: true }); await tick();
    button('Settings').click(); await tick();
    button('Show code').click(); await tick();
    expect($('[aria-label="JavaScript execution"]')).toBeNull();
    expect($('[aria-label="Thinking"]')).not.toBeNull();
    expect(call).toHaveBeenCalledWith('SettingsService.Save', expect.stringContaining('"showCode":false'));
  });
});
