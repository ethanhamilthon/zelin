import { afterEach, expect, it } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import History from './History.svelte';
import type { Message } from './conversation';
import type { Execution } from './agent-types';

let app: ReturnType<typeof mount> | undefined;
afterEach(async () => { if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });

const step = (id: number): Execution => ({ id: String(id), code: `step_${id}()`, output: '', state: 'complete' });
const history = (count: number): Message[] => Array.from({ length: count }, (_, i) => ({ role: i % 2 ? 'assistant' : 'user', content: `message ${i}` }));
const articles = () => document.querySelectorAll('article').length;
const more = () => [...document.querySelectorAll('button')].find((item) => item.textContent?.includes('earlier'));

function show(messages: Message[]) { app = mount(History, { target: document.body, props: { messages } }); flushSync(); }

it('mounts only the latest messages and reveals older ones on request', async () => {
  show(history(100));
  expect(articles()).toBe(30);
  expect(document.body.textContent).toContain('message 99');
  expect(document.body.textContent).not.toContain('message 69');
  expect(more()?.textContent).toBe('Show 30 earlier messages');
  more()!.click(); await tick(); await tick();
  expect(articles()).toBe(60);
  expect(document.body.textContent).toContain('message 40');
  more()!.click(); await tick(); await tick();
  expect(articles()).toBe(90);
  more()!.click(); await tick(); await tick();
  expect(articles()).toBe(100);
  expect(more()).toBeUndefined();
});

it('renders short histories in full without a button', () => {
  show(history(30));
  expect(articles()).toBe(30);
  expect(more()).toBeUndefined();
});

it('bounds execution blocks per message and keeps full source reachable', async () => {
  const messages: Message[] = [{ role: 'assistant', content: 'done', executions: Array.from({ length: 80 }, (_, i) => step(i)) }];
  show(messages);
  const count = () => document.querySelectorAll('[aria-label="JavaScript execution"]').length;
  expect(count()).toBe(10);
  expect(document.body.textContent).toContain('step_79()');
  expect(document.body.textContent).not.toContain('step_69()');
  expect(more()?.textContent).toBe('Show 10 earlier steps');
  for (let i = 0; i < 7; i++) { more()!.click(); await tick(); await tick(); }
  expect(count()).toBe(80);
  expect(more()).toBeUndefined();
  expect(document.body.textContent).toContain('step_0()');
});
