import { expect, it } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import ExecutionBlock from './ExecutionBlock.svelte';
import ThinkingBlock from './ThinkingBlock.svelte';
import { oneLine } from './preview';

const long = Array.from({ length: 5000 }, (_, i) => `line_${i}();`).join('\n');

it('bounds previews and joins lines', () => {
  expect(oneLine('a\n  b\n')).toBe('a b');
  expect(oneLine(long).length).toBeLessThanOrEqual(240);
  expect(oneLine('a   b', /\s+/g)).toBe('a b');
});

it('renders a bounded folded block and the full source when unfolded', async () => {
  const app = mount(ExecutionBlock, { target: document.body, props: { execution: { id: '1', code: long, output: '', state: 'complete' }, folded: true } });
  flushSync();
  expect(document.querySelector('code')!.textContent!.length).toBeLessThanOrEqual(240);
  document.querySelector<HTMLElement>('[aria-label="Unfold code"]')!.click(); await tick();
  expect(document.querySelector('code')!.textContent).toBe(long);
  await unmount(app); document.body.replaceChildren();
});

it('keeps full thinking text after unfolding', async () => {
  const app = mount(ThinkingBlock, { target: document.body, props: { text: long, folded: true } });
  flushSync();
  expect(document.querySelector('p')!.textContent!.length).toBeLessThanOrEqual(240);
  document.querySelector<HTMLElement>('[aria-label="Unfold thinking"]')!.click(); await tick();
  expect(document.querySelector('p')!.textContent).toBe(long);
  await unmount(app); document.body.replaceChildren();
});
