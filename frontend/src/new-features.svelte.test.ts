import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import App from './App.svelte';
import Steps from './Steps.svelte';
import ThinkingBlock from './ThinkingBlock.svelte';
import SessionList from './SessionList.svelte';
import { agentInfo, listModels } from './agent';
import { currentSession, listSessions, openSession, saveSession } from './sessions';
import { call } from './bridge';
import { placeholderInfo, type Execution } from './agent-types';

vi.mock('./agent', () => ({ runAgent: vi.fn(), agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn(), setTitleModel: vi.fn(), steerAgent: vi.fn(), stopAgent: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));
vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => false, listen: vi.fn(async () => () => {}) }));
vi.mock('./sound', () => ({ chime: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
let app: ReturnType<typeof mount> | undefined;
const wait = async () => { await new Promise((r) => setTimeout(r, 0)); await tick(); };
const q = (s: string) => document.querySelector<HTMLElement>(s);
const press = (key: string, code: string, init: KeyboardEventInit = {}) => window.dispatchEvent(new KeyboardEvent('keydown', { key, code, bubbles: true, cancelable: true, ...init }));

beforeEach(() => {
  vi.resetAllMocks();
  Element.prototype.scrollIntoView = vi.fn();
  Object.defineProperties(Range.prototype, { getClientRects: { configurable: true, value: () => [] }, getBoundingClientRect: { configurable: true, value: () => new DOMRect() } });
  vi.mocked(currentSession).mockResolvedValue({ id: 'one', title: 'Alpha', draft: '', messages: [{ role: 'user', content: 'Hi' }, { role: 'assistant', content: 'Hello' }], updated: new Date().toISOString() });
  vi.mocked(openSession).mockResolvedValue({ id: 'two', title: 'Beta', draft: 'draft 2', messages: [{ role: 'user', content: 'Old' }], updated: new Date().toISOString() });
  vi.mocked(listSessions).mockResolvedValue([{ id: 'one', title: 'Alpha', draft: '', messages: [], updated: new Date().toISOString() }, { id: 'two', title: 'Beta', draft: '', messages: [], updated: '2025-05-01T10:00:00Z' }]);
  vi.mocked(saveSession).mockResolvedValue();
  vi.mocked(agentInfo).mockResolvedValue(placeholderInfo('Test'));
  vi.mocked(listModels).mockResolvedValue([{ id: 'm1', name: 'Model One', contextWindow: 128000, pricing: placeholderInfo('').pricing }]);
  vi.mocked(call).mockResolvedValue('Prompt');
});
afterEach(async () => { if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });

it('opens model selector with Ctrl+M and navigates back step-by-step with Escape', async () => {
  app = mount(App, { target: document.body }); flushSync(); await wait();
  expect(Element.prototype.scrollIntoView).toHaveBeenCalled();
  expect(document.activeElement).toBe(q('.cm-content'));
  press('m', 'KeyM', { ctrlKey: true }); await wait();
  expect(q('[role="dialog"]')?.getAttribute('aria-label')).toBe('Select model');
  press('Escape', 'Escape'); await wait();
  expect(q('[role="dialog"]')).toBeNull();
  press('k', 'KeyK', { ctrlKey: true }); await tick();
  [...document.querySelectorAll<HTMLButtonElement>('button')].find((b) => b.textContent?.includes('Instructions'))!.click(); await tick();
  expect(q('[role="dialog"]')?.getAttribute('aria-label')).toBe('Instructions');
  expect(q('.title')?.textContent).not.toContain('Esc');
  [...document.querySelectorAll<HTMLButtonElement>('button')].find((b) => b.textContent?.includes('System Prompt'))!.click(); await wait();
  expect(q('footer[aria-label="System Prompt"]')).not.toBeNull();
  press('Escape', 'Escape'); await tick();
  expect(q('[role="dialog"]')?.getAttribute('aria-label')).toBe('Instructions');
  (q('[role="dialog"]') as HTMLElement).dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true })); await tick();
  expect(q('[role="dialog"]')?.getAttribute('aria-label')).toBe('Commands');
  (q('[role="dialog"]') as HTMLElement).dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true })); await tick();
  expect(q('[role="dialog"]')).toBeNull();
  expect(document.activeElement).toBe(q('.cm-content'));
});

it('filters sessions, pins search outside .body, and shows time/date plus running step count', () => {
  const today = new Date(); today.setHours(14, 5, 0, 0);
  const onselect = vi.fn();
  app = mount(SessionList, { target: document.body, props: {
    sessions: [{ id: 'a', title: 'Deploy service', draft: '', messages: [], updated: today.toISOString() }, { id: 'b', title: 'Fix tests', draft: '', messages: [], updated: '2024-03-10T12:00:00Z' }],
    active: 'a', status: (id) => id === 'a' ? 'running' : '', steps: (id) => id === 'a' ? 4 : 0, onclose: vi.fn(), onnew: vi.fn(), onselect,
  } });
  flushSync();
  expect(q('.top input')).not.toBeNull();
  expect(q('.body input')).toBeNull();
  const rows = [...document.querySelectorAll<HTMLButtonElement>('.body button')];
  expect(rows[0].textContent).toContain('14:05 (4 Steps)');
  expect(rows[1].textContent).toContain('Mar 10, 2024');
  const input = q('.top input') as HTMLInputElement;
  input.value = 'fix'; input.dispatchEvent(new Event('input', { bubbles: true })); flushSync();
  expect(document.querySelectorAll('.body button')).toHaveLength(1);
  input.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true }));
  expect(onselect).toHaveBeenCalledWith('b');
});

it('shows live preview in StepStack while running and folds ThinkingBlock to one line', () => {
  const executions: Execution[] = [
    { id: '1', code: 'fs.read("a.go")', output: 'ok', state: 'complete' },
    { id: '2', kind: 'thinking', code: 'Checking modal\nnavigation state', output: '', state: 'complete' },
  ];
  const props = $state({ executions, running: true, stack: true, foldCode: true });
  app = mount(Steps, { target: document.body, props }); flushSync();
  expect(q('.toggle')?.textContent).toContain('2 steps');
  expect(q('.preview')?.textContent).toBe('Checking modal navigation state');
  props.running = false; flushSync();
  expect(q('.preview')).toBeNull();
  expect(q('.toggle')?.textContent?.trim()).toBe('2 steps');
});

it('folds and unfolds ThinkingBlock without inner scroll', () => {
  app = mount(ThinkingBlock, { target: document.body, props: { text: 'First line\nSecond line', folded: true } }); flushSync();
  expect(q('.body')?.classList.contains('folded')).toBe(true);
  expect(q('.thinking p')?.textContent).toBe('First line Second line');
  q('.body')!.click(); flushSync();
  expect(q('.body')?.classList.contains('folded')).toBe(false);
  expect(q('.thinking p')?.textContent).toBe('First line\nSecond line');
});
