import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import { EditorView } from '@codemirror/view';
import App from './App.svelte';
import { agentInfo, configureModel, listModels, runAgent } from './agent';
import { currentSession, createSession, openSession, listSessions, saveSession } from './sessions';
import type { AgentEvent, AgentResult, Execution } from './agent-types';
import { placeholderInfo } from './agent-types';

vi.mock('./agent', () => ({ runAgent: vi.fn(), agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
let app: ReturnType<typeof mount> | undefined;
const execution: Execution = { id: 'call', code: 'bash("pwd", 10)', output: '"/tmp"', state: 'complete' };
const wait = async () => { await new Promise((resolve) => setTimeout(resolve, 0)); await tick(); };
const session = (id = 'first', draft = 'Hello **world**') => ({ id, title: id, draft, messages: [], updated: new Date().toISOString() });

beforeEach(() => {
  vi.resetAllMocks();
  Element.prototype.scrollIntoView = vi.fn();
  Object.defineProperties(Range.prototype, {
    getClientRects: { configurable: true, value: () => [] },
    getBoundingClientRect: { configurable: true, value: () => new DOMRect() },
  });
  vi.mocked(currentSession).mockResolvedValue(session());
  vi.mocked(createSession).mockResolvedValue(session('new', ''));
  vi.mocked(openSession).mockResolvedValue(session('second', 'Saved draft'));
  vi.mocked(listSessions).mockResolvedValue([session(), session('second')]);
  vi.mocked(saveSession).mockResolvedValue();
  vi.mocked(agentInfo).mockResolvedValue({ ...placeholderInfo('Test Model'), model: 'test/model', contextWindow: 128000 });
});

afterEach(async () => { if (app) await unmount(app); document.body.replaceChildren(); });
async function start() { app = mount(App, { target: document.body }); flushSync(); await wait(); }
function key(key: string, ctrlKey = false) { window.dispatchEvent(new KeyboardEvent('keydown', { key, ctrlKey, cancelable: true })); }
function submit() {
  document.querySelector('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', {
    key: 'Enter', code: 'Enter', shiftKey: true, bubbles: true, cancelable: true,
  }));
}
function draft(text: string) {
  const view = EditorView.findFromDOM(document.querySelector('.cm-editor')!)!;
  view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: text } });
}

describe('agent UI', () => {
  it('keeps drafts editable, renders plain code, and persists answers', async () => {
    let finish!: (result: AgentResult) => void;
    let notify!: (event: AgentEvent) => void;
    vi.mocked(runAgent).mockImplementation((_, observe) => {
      notify = observe;
      return new Promise((resolve) => { finish = resolve; });
    });
    await start(); submit(); await wait();
    expect(document.querySelector('article strong')?.textContent).toBe('world');
    expect(document.querySelector('article header')).toBeNull();
    expect(document.querySelector('[aria-label="Request running"]')).not.toBeNull();
    expect(document.querySelector('.cm-content')?.getAttribute('contenteditable')).toBe('true');
    draft('Next draft'); submit();
    expect(runAgent).toHaveBeenCalledOnce();
    notify({ kind: 'activity', activity: 'Running JavaScript' });
    notify({ kind: 'execution', execution: { ...execution, state: 'running', output: '' } });
    await tick();
    expect(document.querySelector('footer [role="status"]')?.getAttribute('aria-label')).toBe('Running JavaScript');
    finish({ content: '## Done', executions: [execution] }); await wait();
    expect(document.querySelectorAll('article')).toHaveLength(2);
    expect(document.querySelector('[aria-label="Request running"]')).toBeNull();
    expect(document.querySelector('article [aria-label="Turn complete"]')).not.toBeNull();
    expect(document.querySelectorAll('[aria-label="JavaScript execution"]')).toHaveLength(1);
    expect(document.querySelector('[aria-label="JavaScript execution"] header')).toBeNull();
    expect(document.querySelector('.output')).toBeNull();
    document.querySelector<HTMLButtonElement>('[aria-label="Toggle execution output"]')!.click(); await tick();
    expect(document.querySelector('.output')?.textContent).toBe('"/tmp"');
    expect(document.querySelector('.cm-content')?.textContent).toBe('Next draft');
    expect(document.querySelector('footer')?.textContent).toContain('128,000 tokens');
    expect(document.querySelector('footer')?.textContent).not.toContain('Shift + Enter');
    expect(saveSession).toHaveBeenLastCalledWith('first', expect.arrayContaining([expect.objectContaining({ content: '## Done' })]), 'Next draft');
  });

  it('preserves the next draft on failure', async () => {
    let fail!: (reason: Error) => void;
    vi.mocked(runAgent).mockImplementation(() => new Promise((_, reject) => { fail = reject; }));
    await start(); submit(); await wait(); draft('Keep this'); fail(new Error('Offline')); await wait();
    expect(document.querySelector('[role="alert"]')?.textContent).toBe('Offline');
    expect(document.querySelector('.cm-content')?.textContent).toBe('Keep this');
  });

  it('closes commands with Escape even if focus leaves the menu', async () => {
    await start(); key('k', true); await tick();
    expect(document.querySelector('[role="dialog"]')).not.toBeNull();
    key('Escape');
    await vi.waitFor(() => expect(document.querySelector('[role="dialog"]')).toBeNull());
    expect(document.activeElement).toBe(document.querySelector('.cm-content'));
  });

  it('shows unread completion until the user scrolls to the end', async () => {
    let finish!: (result: AgentResult) => void;
    vi.mocked(runAgent).mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
    await start(); submit(); await wait();
    let reached = false;
    const rect = vi.spyOn(Element.prototype, 'getBoundingClientRect').mockImplementation(function (this: Element) {
      return new DOMRect(0, 0, 100, this.tagName === 'ARTICLE' && !reached ? 900 : 500);
    });
    finish({ content: 'Done', executions: [] }); await wait();
    expect(document.querySelector('footer [role="status"]')?.getAttribute('aria-label')).toBe('Unread response');
    reached = true;
    document.querySelector('main')!.dispatchEvent(new Event('scroll')); await tick();
    expect(document.querySelector('footer [role="status"]')?.getAttribute('aria-label')).toBe('Ready');
    rect.mockRestore();
  });

  it('creates sessions and switches saved sessions with Ctrl+S', async () => {
    await start(); key('k', true); await tick();
    [...document.querySelectorAll('button')].find((button) => button.textContent?.includes('New session'))!.click();
    await wait();
    expect(createSession).toHaveBeenCalledOnce();
    expect(document.querySelector('.cm-content')?.textContent).not.toContain('Hello');
    key('s', true); await wait();
    expect(document.querySelector('[aria-label="Sessions"]')).not.toBeNull();
    [...document.querySelectorAll('button')].find((button) => button.textContent?.includes('second'))!.click();
    await wait();
    expect(openSession).toHaveBeenCalledWith('second');
    expect(document.querySelector('.cm-content')?.textContent).toBe('Saved draft');
  });
});
