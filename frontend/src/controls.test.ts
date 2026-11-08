import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import { EditorView } from '@codemirror/view';
import App from './App.svelte';
import { agentInfo, runAgent, steerAgent, stopAgent } from './agent';
import { currentSession, openSession, saveSession } from './sessions';
import { call } from './bridge';
import { placeholderInfo, type AgentEvent, type AgentResult } from './agent-types';

vi.mock('./agent', () => ({ runAgent: vi.fn(), agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn(), setTitleModel: vi.fn(), steerAgent: vi.fn(), stopAgent: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));
vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => false, listen: vi.fn(async () => () => {}) }));
vi.mock('./sound', () => ({ chime: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
let app: ReturnType<typeof mount>;
let notify: (event: AgentEvent) => void, finish: (result: AgentResult) => void;
const $ = (selector: string) => document.querySelector<HTMLElement>(selector);
const wait = async () => { await new Promise((r) => setTimeout(r, 0)); await tick(); };
const button = (text: string) => [...document.querySelectorAll<HTMLButtonElement>('button')].find((b) => b.textContent?.includes(text))!;
function send() { $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', code: 'Enter', bubbles: true, cancelable: true })); }
function type(text: string) { const view = EditorView.findFromDOM($('.cm-editor')!)!; view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: text } }); }
function key() { window.dispatchEvent(new KeyboardEvent('keydown', { key: 'c', code: 'KeyC', ctrlKey: true, cancelable: true })); }
beforeEach(async () => {
  vi.resetAllMocks();
  Element.prototype.scrollIntoView = vi.fn();
  Object.defineProperties(Range.prototype, { getClientRects: { configurable: true, value: () => [] }, getBoundingClientRect: { configurable: true, value: () => new DOMRect() } });
  vi.mocked(currentSession).mockResolvedValue({ id: 'one', title: 'Task', draft: 'start', messages: [], updated: '' });
  vi.mocked(saveSession).mockResolvedValue();
  vi.mocked(agentInfo).mockResolvedValue(placeholderInfo('Test'));
  vi.mocked(call).mockResolvedValue('');
  vi.mocked(steerAgent).mockResolvedValue(true);
  vi.mocked(stopAgent).mockResolvedValue(true);
  vi.mocked(runAgent).mockImplementation((_, observe) => { notify = observe; return new Promise((resolve) => { finish = resolve; }); });
  app = mount(App, { target: document.body }); flushSync(); await wait();
});
afterEach(async () => { window.getSelection()?.removeAllRanges(); await unmount(app); document.body.replaceChildren(); });

it('queues steering, cancels on edit, retains text, and clears it on consumption', async () => {
  send(); await wait(); type('change direction'); send(); await wait();
  expect(steerAgent).toHaveBeenLastCalledWith('one', 'change direction', [], []);
  expect($('[aria-label="Steering queued"]')).not.toBeNull();
  expect($('.cm-content')?.textContent).toBe('change direction');
  type('change direction again'); await tick();
  expect(steerAgent).toHaveBeenLastCalledWith('one', '');
  expect($('[aria-label="Steering queued"]')).toBeNull();
  send(); await wait(); notify({ kind: 'steer', message: { content: 'change direction again' } }); await wait();
  expect(EditorView.findFromDOM($('.cm-editor')!)!.state.doc.toString()).toBe('');
  expect(document.querySelectorAll('article')).toHaveLength(2);
  finish({ content: 'done', executions: [] }); await wait();
  expect(document.querySelectorAll('article')).toHaveLength(3);
});

it('shows Stop before Commands only when busy; selection keeps Ctrl+C for copying', async () => {
  expect($('footer')?.textContent).not.toContain('Stop');
  send(); await wait();
  expect([...document.querySelectorAll('footer .right button')].map((b) => b.textContent)).toEqual(['Stop Ctrl+C', 'Commands Ctrl+K']);
  const range = document.createRange(); range.selectNodeContents($('article .text')!);
  window.getSelection()!.removeAllRanges(); window.getSelection()!.addRange(range); key(); expect(stopAgent).not.toHaveBeenCalled();
  window.getSelection()!.removeAllRanges(); key(); expect(stopAgent).toHaveBeenCalledWith('one');
  vi.mocked(stopAgent).mockClear(); button('Stop').click(); expect(stopAgent).toHaveBeenCalledWith('one');
  finish({ content: 'done', executions: [] }); await wait();
  expect($('footer')?.textContent).not.toContain('Stop');
});

it('shows Working in the status bar while a run is active', async () => {
  send(); await wait();
  expect($('footer [role="status"]')?.textContent?.trim()).toBe('1 Working');
  finish({ content: 'done', executions: [] }); await wait();
  expect($('footer [role="status"]')).toBeNull();
});

it('runs manual compaction from Commands and updates usage', async () => {
  send(); await wait(); finish({ content: 'done', executions: [] }); await wait();
  vi.mocked(call).mockImplementation(async (method) => method === 'AgentService.Compact' ? { context: 25 } : '');
  button('Commands').click(); await tick(); button('Session actions').click(); await tick(); button('Compact context').click(); await wait();
  expect(call).toHaveBeenCalledWith('AgentService.Compact', 'one');
  expect($('[role="dialog"]')).toBeNull();
});

it('hands off and opens the new session', async () => {
  send(); await wait(); finish({ content: 'done', executions: [] }); await wait();
  vi.mocked(call).mockImplementation(async (method) => method === 'AgentService.Handoff' ? { id: 'next' } : '');
  vi.mocked(openSession).mockResolvedValue({ id: 'next', title: 'Handoff', draft: '', messages: [{ role: 'user', content: 'Handoff summary' }], updated: '' });
  button('Commands').click(); await tick(); button('Session actions').click(); await tick(); button('Handoff').click(); await wait(); await wait();
  expect(call).toHaveBeenCalledWith('AgentService.Handoff', 'one');
  expect(openSession).toHaveBeenCalledWith('next');
  expect($('article .text')?.textContent).toBe('Handoff summary');
});

it('persists the auto-compaction toggle and threshold', async () => {
  button('Commands').click(); await tick(); button('Settings').click(); await tick();
  expect(button('Auto compaction').textContent).toContain('On');
  button('Auto compaction').click(); await tick();
  const input = $('[aria-label="Compaction threshold"]') as HTMLInputElement;
  expect(input.value).toBe('60'); input.value = '75'; input.dispatchEvent(new Event('change', { bubbles: true })); await tick();
  expect(call).toHaveBeenLastCalledWith('SettingsService.Save', expect.stringContaining('"compactThreshold":75'));
  const saved = vi.mocked(call).mock.calls.filter(([m]) => m === 'SettingsService.Save').at(-1)!;
  expect(JSON.parse(saved[1] as string).autoCompact).toBe(false);
});
