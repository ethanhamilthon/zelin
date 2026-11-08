import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import App from './App.svelte';
import { agentInfo, runAgent, stopAgent } from './agent';
import { currentSession, listSessions } from './sessions';
import { call } from './bridge';
import { placeholderInfo, type AgentResult } from './agent-types';

vi.mock('./agent', () => ({ runAgent: vi.fn(), agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn(), setTitleModel: vi.fn(), steerAgent: vi.fn(), stopAgent: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));
vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => false, listen: vi.fn(async () => () => {}) }));
vi.mock('./sound', () => ({ chime: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));

let app: ReturnType<typeof mount>;
let finish: (result: AgentResult) => void;
const $ = (selector: string) => document.querySelector<HTMLElement>(selector);
const wait = async () => { await new Promise((r) => setTimeout(r, 0)); await tick(); };
const button = (text: string) => [...document.querySelectorAll<HTMLButtonElement>('button')].find((b) => b.textContent?.includes(text))!;
const key = (k: string, code: string, init: KeyboardEventInit = {}) => window.dispatchEvent(new KeyboardEvent('keydown', { key: k, code, cancelable: true, ...init }));
const send = () => $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', code: 'Enter', bubbles: true, cancelable: true }));
const record = { id: 'one', title: 'Task', draft: 'start', messages: [], updated: '' };

beforeEach(async () => {
  vi.resetAllMocks();
  Element.prototype.scrollIntoView = vi.fn();
  Object.defineProperties(Range.prototype, { getClientRects: { configurable: true, value: () => [] }, getBoundingClientRect: { configurable: true, value: () => new DOMRect() } });
  vi.mocked(currentSession).mockResolvedValue(record);
  vi.mocked(listSessions).mockResolvedValue([{ ...record }]);
  vi.mocked(agentInfo).mockResolvedValue(placeholderInfo('Test'));
  vi.mocked(stopAgent).mockResolvedValue(true);
  vi.mocked(call).mockResolvedValue('');
  vi.mocked(runAgent).mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  app = mount(App, { target: document.body }); flushSync(); await wait();
});
afterEach(async () => { await unmount(app); document.body.replaceChildren(); });

it('opens with Ctrl+E, toggles, and returns to Commands with Escape', async () => {
  key('e', 'KeyE', { ctrlKey: true }); await tick();
  expect($('[role="dialog"]')?.getAttribute('aria-label')).toBe('Session actions');
  key('Escape', 'Escape'); await vi.waitFor(() => expect($('[role="dialog"]')).toBeNull());
  button('Commands').click(); await tick(); button('Session actions').click(); await tick();
  expect($('[role="dialog"]')?.getAttribute('aria-label')).toBe('Session actions');
  key('Escape', 'Escape'); await vi.waitFor(() => expect($('[role="dialog"]')?.getAttribute('aria-label')).toBe('Commands'));
});

it('disables actions by state', async () => {
  key('e', 'KeyE', { ctrlKey: true }); await tick();
  const off = (name: string) => button(name).disabled;
  expect([off('Compact'), off('Handoff'), off('Export'), off('Stop Session'), off('Rename')]).toEqual([true, true, true, true, false]);
  key('Escape', 'Escape'); await wait();
  send(); await wait();
  key('e', 'KeyE', { ctrlKey: true }); await tick();
  expect([off('Compact'), off('Handoff'), off('Export'), off('Stop Session'), off('Rename')]).toEqual([true, true, false, false, false]);
  button('Stop Session').click(); await wait();
  expect(stopAgent).toHaveBeenCalledWith('one');
  finish({ content: 'done', executions: [] }); await wait();
  key('e', 'KeyE', { ctrlKey: true }); await tick();
  expect([off('Compact'), off('Handoff'), off('Export'), off('Stop Session')]).toEqual([false, false, false, true]);
});

it('exports the target session and reports the saved path', async () => {
  send(); await wait(); finish({ content: 'done', executions: [] }); await wait();
  vi.mocked(call).mockImplementation(async (method) => method === 'SessionService.Export' ? '/tmp/task.md' : '');
  key('e', 'KeyE', { ctrlKey: true }); await tick(); button('Export session').click(); await wait();
  expect(call).toHaveBeenCalledWith('SessionService.Export', 'one');
  expect($('.toasts')?.textContent).toContain('/tmp/task.md');
});

it('does not toast when the export dialog is cancelled and reports failures', async () => {
  send(); await wait(); finish({ content: 'done', executions: [] }); await wait();
  key('e', 'KeyE', { ctrlKey: true }); await tick(); button('Export session').click(); await wait();
  expect($('.toasts')?.textContent ?? '').not.toContain('Exported');
  vi.mocked(call).mockImplementation(async (method) => { if (method === 'SessionService.Export') throw new Error('disk full'); return ''; });
  key('e', 'KeyE', { ctrlKey: true }); await tick(); button('Export session').click(); await wait();
  expect($('.toasts')?.textContent).toContain('Export failed: disk full');
});

it('renames with validation and updates the list', async () => {
  vi.mocked(call).mockImplementation(async (method, _id, title) => {
    if (method !== 'SessionService.Rename') return '';
    vi.mocked(listSessions).mockResolvedValue([{ ...record, title: String(title).trim() }]);
    return String(title).trim();
  });
  key('e', 'KeyE', { ctrlKey: true }); await tick(); button('Rename Session').click(); await tick();
  const input = $('[aria-label="Session title"]') as HTMLInputElement;
  expect(input.value).toBe('Task');
  input.value = '   '; input.dispatchEvent(new Event('input', { bubbles: true }));
  input.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true })); await wait();
  expect($('[role="alert"]')?.textContent).toBe('Title cannot be empty');
  expect(call).not.toHaveBeenCalledWith('SessionService.Rename', expect.anything(), expect.anything());
  input.value = '  Better  '; input.dispatchEvent(new Event('input', { bubbles: true }));
  input.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true })); await wait();
  expect(call).toHaveBeenCalledWith('SessionService.Rename', 'one', '  Better  ');
  expect($('[role="dialog"]')).toBeNull();
  key('s', 'KeyS', { ctrlKey: true }); await wait();
  expect($('[role="dialog"]')?.textContent).toContain('Better');
});

it('keeps the rename dialog open and shows the backend error', async () => {
  vi.mocked(call).mockImplementation(async (method) => { if (method === 'SessionService.Rename') throw new Error('Session not found'); return ''; });
  key('e', 'KeyE', { ctrlKey: true }); await tick(); button('Rename Session').click(); await tick();
  $('[aria-label="Session title"]')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true })); await wait();
  expect($('[role="alert"]')?.textContent).toContain('Session not found');
  expect($('[aria-label="Session title"]')).not.toBeNull();
});
