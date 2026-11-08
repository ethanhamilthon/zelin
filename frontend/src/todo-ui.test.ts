import { afterEach, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import Commands from './Commands.svelte';
import SessionList from './SessionList.svelte';
import StatusBar from './StatusBar.svelte';
import ExecutionBlock from './ExecutionBlock.svelte';
import Settings from './Settings.svelte';
import { defaultSettings } from './settings.svelte';
import { fontFamily } from './appearance';
import { defaultHotkeys } from './hotkeys';
import { placeholderInfo } from './agent-types';
import { modelName, search } from './models';

vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
let app: ReturnType<typeof mount> | undefined;
afterEach(async () => { if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });

it('focuses command search and removes redundant actions', () => {
  const onopen = vi.fn();
  app = mount(Commands, { target: document.body, props: { hotkeys: defaultHotkeys, onopen, onclose: vi.fn() } });
  flushSync();
  const input = document.querySelector('input')!;
  expect(document.activeElement).toBe(input);
  const labels = [...document.querySelectorAll('button')].map((button) => button.textContent);
  expect(labels).toHaveLength(8);
  expect(labels[0]).toContain('Session actions');
  expect(labels.join(' ')).not.toMatch(/Compact|Handoff/);
  expect(labels.at(-1)).toContain('Switch session');
  expect(labels.join(' ')).not.toMatch(/New session|Stop|Attach image|Focus|Scroll|Select model|Select reasoning/);
  input.value = 'system';
  input.dispatchEvent(new Event('input', { bubbles: true }));
  flushSync();
  expect(document.querySelectorAll('button')).toHaveLength(1);
  expect(document.querySelector('button')?.textContent).toContain('Instructions');
  input.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true }));
  expect(onopen).toHaveBeenCalledWith('instructions');
});

it('closes outside the modal but not inside and keeps New session in the header', () => {
  const onclose = vi.fn();
  const onnew = vi.fn();
  app = mount(SessionList, { target: document.body, props: {
    sessions: [], active: '', status: () => '', onclose, onnew, onselect: vi.fn(),
  } });
  flushSync();
  const button = document.querySelector<HTMLButtonElement>('.title button')!;
  expect(button.textContent).toContain('New session');
  button.click();
  expect(onnew).toHaveBeenCalledOnce();
  expect(onclose).not.toHaveBeenCalled();
  document.querySelector<HTMLElement>('.overlay')!.click();
  expect(onclose).toHaveBeenCalledOnce();
});

it('shows context usage before the readable model and opens usage', () => {
  const onopen = vi.fn();
  const info = { ...placeholderInfo('claude-sonnet-5-5'), contextWindow: 1000 };
  app = mount(StatusBar, { target: document.body, props: { info, context: 250, commands: 'Ctrl+K', onopen } });
  flushSync();
  const buttons = [...document.querySelectorAll('button')];
  expect(buttons[0].getAttribute('aria-label')).toBe('Context used: 25.0%');
  expect(buttons[0].querySelector('span')?.style.background).toContain('25%');
  expect(buttons[1].textContent).toBe('Claude Sonnet 5.5');
  buttons[0].click();
  expect(onopen).toHaveBeenCalledWith('usage');
});

it('folds code to one line and hides expanded output', () => {
  app = mount(ExecutionBlock, { target: document.body, props: {
    execution: { id: 'test', code: 'const a = 1;\nconst b = 2;', output: 'Result', state: 'complete' },
  } });
  flushSync();
  document.querySelector<HTMLButtonElement>('.source')!.click();
  flushSync();
  expect(document.querySelector('.output')?.textContent).toBe('Result');
  document.querySelector<HTMLButtonElement>('.fold')!.click();
  flushSync();
  expect(document.querySelector('.source code')?.textContent).toBe('const a = 1; const b = 2;');
  expect(document.querySelector('.source')?.classList.contains('folded')).toBe(true);
  expect(document.querySelector('.output')).toBeNull();
});

it.each([
  ['claude-sonnet-5-5', 'Claude Sonnet 5.5'],
  ['gpt-4-1-mini', 'GPT 4.1 Mini'],
  ['model-1-2-3', 'Model 1.2.3'],
  ['GPT-5', 'GPT 5'],
])('formats model name %s', (raw, expected) => {
  expect(modelName(raw)).toBe(expected);
  const model = { id: raw, name: raw, contextWindow: 0, pricing: placeholderInfo('').pricing };
  expect(search([model], expected)).toEqual([model]);
});it('uses SF Pro Display without a font selector', () => {
 app = mount(Settings, { target: document.body, props: { settings: defaultSettings(), onclose: vi.fn(), onchange: vi.fn() } });
 flushSync();
 expect(document.querySelector('[aria-label="Application font"]')).toBeNull();
 expect(fontFamily).toContain('SF Pro Display');
});
