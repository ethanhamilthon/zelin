import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import History from './History.svelte';
import Onboarding from './Onboarding.svelte';
import { links } from './links';
import type { Message } from './conversation';
import { call } from './bridge';

vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => true, listen: vi.fn(async () => () => {}) }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
let app: ReturnType<typeof mount> | undefined;
const wait = async () => { for (let i = 0; i < 5; i++) { await new Promise((resolve) => setTimeout(resolve, 0)); await tick(); } };
beforeEach(() => {
  vi.resetAllMocks();
  Element.prototype.scrollIntoView = vi.fn();
  Object.defineProperties(Range.prototype, {
    getClientRects: { configurable: true, value: () => [] },
    getBoundingClientRect: { configurable: true, value: () => new DOMRect() },
  });
});
afterEach(async () => { vi.useRealTimers(); if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });

it('shows steps, elapsed time and streamed text while busy, then the turn time', async () => {
  vi.useFakeTimers({ now: 10_000 });
  const execution = { id: 'a', code: 'bash("ls")', output: '', state: 'running' as const };
  const props = $state({ messages: [{ role: 'user' as const, content: 'Hi' }] as Message[], busy: true, started: 10_000, executions: [execution], streaming: '' });
  app = mount(History, { target: document.body, props });
  flushSync();
  const row = () => document.querySelector('.toggle')?.textContent ?? '';
  expect(row()).toMatch(/1 step\s*\|\s*0s\s*\|\s*bash\("ls"\)/);
  expect(document.querySelector('.spinner')).not.toBeNull();
  vi.advanceTimersByTime(3000); flushSync();
  expect(row()).toContain('3s');
  props.streaming = 'Partial **answer**'; flushSync();
  expect(document.querySelector('[aria-label="Streaming answer"]')?.textContent).toContain('Partial answer');
  props.busy = false;
  props.messages = [...props.messages, { role: 'assistant', content: 'Done', turn: { finished: '', duration: 46000, cost: null, context: null, window: 0 } }];
  flushSync();
  expect(document.querySelector('[aria-label="Streaming answer"]')).toBeNull();
  expect(document.querySelector('[aria-label="Turn time"]')?.textContent).toBe('46s');
  expect(document.querySelector('[role="separator"]')).toBeNull();
});

it('routes link clicks to the handler', () => {
  const node = document.createElement('div');
  node.innerHTML = '<a href="https://x.dev" data-url="https://x.dev">x</a><a href="README.md" data-path="README.md"><code>README.md</code></a>';
  document.body.append(node);
  const open = vi.fn();
  links(node, open);
  node.querySelector('a')!.click();
  node.querySelector('code')!.click();
  expect(open.mock.calls).toEqual([[{ url: 'https://x.dev', path: undefined }], [{ url: undefined, path: 'README.md' }]]);
});

it('walks through onboarding and enables only the main model', async () => {
  const models = [{ id: 'openai/gpt-5', name: 'GPT-5' }, { id: 'openai/gpt-5-mini', name: 'GPT-5 Mini' }, { id: 'google/gemini', name: 'Gemini' }];
  const checks: string[] = [];
  vi.mocked(call).mockImplementation(async (method: string, ...args: unknown[]) => {
    if (method === 'SettingsService.Load') return '';
    if (method === 'SetupService.Connect') return models.map((model) => ({ ...model, contextWindow: 0, pricing: {} }));
    if (method === 'SetupService.Check') { checks.push(args[0] as string); return { model: args[0], reply: 'OK', error: '', millis: 800 }; }
    if (method === 'SetupService.Finish') return {};
    return undefined;
  });
  const ondone = vi.fn(); vi.useFakeTimers({ shouldAdvanceTime: true });
  app = mount(Onboarding, { target: document.body, props: { ondone } });
  flushSync(); await wait();
  const footer = () => document.querySelector('footer')!.textContent;
  expect(footer()).toContain('Welcome to Zelin');
  expect(footer()).toContain('Setup 1/4');
  const { EditorView } = await import('@codemirror/view');
  const send = async (text: string) => {
    const view = EditorView.findFromDOM(document.querySelector('.cm-editor')!)!;
    view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: text } });
    document.querySelector('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', code: 'Enter', bubbles: true, cancelable: true }));
    await wait();
  };
  await send('');
  await send('sk-or-1234567890abcd');
  expect(call).toHaveBeenCalledWith('SetupService.Connect', 'https://openrouter.ai/api/v1', 'sk-or-1234567890abcd');
  expect(document.body.textContent).not.toContain('sk-or-1234567890abcd');
  expect(footer()).toContain('Setup 2/4');
  const options = () => [...document.querySelectorAll('[role="option"]')].map((item) => item.textContent);
  expect(options()).toHaveLength(3);
  const type = async (text: string) => { const view = EditorView.findFromDOM(document.querySelector('.cm-editor')!)!; view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: text } }); await wait(); };
  await type('gpt-5');
  expect(options()).toHaveLength(2);
  await send('gpt-5');
  expect(footer()).toContain('Setup 3/4');
  expect(options()[0]).toContain('openai/gpt-5');
  await type('gpt');
  document.querySelector('.editor')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowDown', bubbles: true, cancelable: true })); await wait();
  expect(document.querySelector('[aria-selected="true"]')?.textContent).toContain('gpt-5-mini');
  document.querySelector<HTMLButtonElement>('[aria-selected="true"]')!.click(); await wait();
  expect(checks).toEqual(['openai/gpt-5', 'openai/gpt-5-mini']);
  expect(call).toHaveBeenCalledWith('SetupService.Finish', 'openai/gpt-5', 'openai/gpt-5-mini');
  expect(document.body.textContent).toContain('All set');
  const saved = vi.mocked(call).mock.calls.filter(([method]) => method === 'SettingsService.Save').pop()!;
  expect(JSON.parse(saved[1] as string).hiddenModels).toEqual(['openai/gpt-5-mini', 'google/gemini']);
  vi.advanceTimersByTime(1000);
  expect(ondone).toHaveBeenCalled();
});
