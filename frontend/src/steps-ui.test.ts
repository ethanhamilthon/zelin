import { afterEach, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import Steps from './Steps.svelte';
import History from './History.svelte';
import SessionList from './SessionList.svelte';
import { elapsed, summary } from './turn';
import type { Execution } from './agent-types';

vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
let app: ReturnType<typeof mount> | undefined;
afterEach(async () => { if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });
const code = (id: string): Execution => ({ id, code: `run(${id})\nnext()`, output: id, state: 'complete' });
const thought: Execution = { id: 'think', kind: 'thinking', code: 'Plan', output: '', state: 'complete' };
const $$ = (selector: string) => document.querySelectorAll(selector);

it('stacks consecutive steps and keeps each block foldable', () => {
  app = mount(Steps, { target: document.body, props: { executions: [thought, code('a'), code('b')], stack: true } });
  flushSync();
  const toggle = document.querySelector<HTMLButtonElement>('.toggle')!;
  expect(toggle.textContent).toContain('2 code · 1 thinking');
  expect($$('.execution, .thinking')).toHaveLength(0);
  toggle.click();
  flushSync();
  expect($$('.execution')).toHaveLength(2);
  expect($$('.thinking')).toHaveLength(1);
  $$('.fold')[0].dispatchEvent(new MouseEvent('click', { bubbles: true }));
  flushSync();
  expect($$('.source')[0].classList.contains('folded')).toBe(true);
  expect($$('.source')[1].classList.contains('folded')).toBe(false);
  toggle.click();
  flushSync();
  expect($$('.execution')).toHaveLength(0);
});

it('renders steps flat when stacking is off or there is one step', () => {
  app = mount(Steps, { target: document.body, props: { executions: [thought, code('a')] } });
  flushSync();
  expect($$('.toggle')).toHaveLength(0);
  expect($$('.execution, .thinking')).toHaveLength(2);
});

it('formats the turn summary', () => {
  expect(elapsed(4400)).toBe('4s');
  expect(elapsed(125000)).toBe('2m 05s');
  expect(elapsed(3720000)).toBe('1h 02m');
  const finished = new Date(2026, 10, 8, 4, 42).toISOString();
  expect(summary({ finished, duration: 65000, cost: 0.0123, context: 42000, window: 200000 })).toBe('04:42 · 1m 05s · $0.01 · 42K / 200K');
  expect(summary({ finished, duration: 1000, cost: null, context: null, window: 0 })).toBe('04:42 · 1s');
});

it('shows a summary divider after the final answer', () => {
  const turn = { finished: new Date().toISOString(), duration: 3000, cost: 0.5, context: 1000, window: 2000 };
  app = mount(History, { target: document.body, props: { messages: [{ role: 'user', content: 'Hi' }, { role: 'assistant', content: 'Done', turn }] } });
  flushSync();
  const divider = document.querySelector('[aria-label="Turn summary"]');
  expect(divider?.textContent).toContain('3s · $0.50 · 1K / 2K');
  expect(divider?.closest('article')?.getAttribute('aria-label')).toBe('Assistant message');
});

it('shows the new session hotkey in the header button', () => {
  app = mount(SessionList, { target: document.body, props: { sessions: [], active: '', hotkey: 'Ctrl+N', status: () => '', onclose: vi.fn(), onselect: vi.fn(), onnew: vi.fn() } });
  flushSync();
  expect(document.querySelector('.title button')?.textContent).toBe('New session Ctrl+N');
});
