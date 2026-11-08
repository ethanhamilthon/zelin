import { afterEach, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import Steps from './Steps.svelte';
import History from './History.svelte';
import SessionList from './SessionList.svelte';
import { elapsed } from './turn';
import Timeline from './Timeline.svelte';
import { segments } from './segments';
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
  expect(toggle.textContent?.trim()).toBe('3 steps');
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

it('renders interim assistant text as a regular answer between step groups', () => {
  const step = (id: string) => ({ id, code: id, output: '', state: 'complete' as const });
  app = mount(Timeline, { target: document.body, props: { stack: true, executions: [step('a'), step('b'), { id: 'n', kind: 'message', code: 'Checking **files**', output: '', state: 'complete' }, step('c'), step('d')] } });
  flushSync();
  const note = document.querySelector('.note')!;
  expect(note.textContent).toContain('Checking files');
  expect(note.querySelector('.icon-slot')).not.toBeNull();
  expect(note.closest('.items')).toBeNull();
  expect([...document.querySelectorAll('.toggle .count')].map((item) => item.textContent)).toEqual(['2 steps', '2 steps']);
});

it('splits executions into segments around messages', () => {
  const step = (id: string, kind?: 'message') => ({ id, kind, code: id, output: '', state: 'complete' as const });
  expect(segments([step('a'), step('n', 'message'), step('b')]).map((part) => part.note ? 'note' : part.steps.length)).toEqual([1, 'note', 1]);
});

it('formats elapsed time', () => {
  expect(elapsed(4400)).toBe('4s');
  expect(elapsed(125000)).toBe('2m 05s');
  expect(elapsed(3720000)).toBe('1h 02m');
});

it('titles the list Workspaces with a new workspace button', () => {
  app = mount(SessionList, { target: document.body, props: { sessions: [], active: '', status: () => '', onclose: vi.fn(), onselect: vi.fn() } });
  flushSync();
  expect(document.querySelector('.heading')?.textContent).toBe('Workspaces');
  expect(document.querySelector('.title button')?.textContent).toBe('New workspace');
});
