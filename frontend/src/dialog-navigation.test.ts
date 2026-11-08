import { expect, it, vi } from 'vitest';
import { mount, unmount, tick } from 'svelte';
import Commands from './Commands.svelte';
import { defaultHotkeys } from './hotkeys';
it('moves focus through dialog options with arrows and activates the chosen option', async () => {
 const onopen = vi.fn();
 const app = mount(Commands, { target: document.body, props: { hotkeys: defaultHotkeys, onclose: vi.fn(), onopen } });
 try {
  await tick();
  const input = document.querySelector<HTMLInputElement>('input')!;
  input.focus();
  input.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowDown', bubbles: true, cancelable: true }));
  const first = document.activeElement as HTMLButtonElement;
  expect(first.tagName).toBe('BUTTON');
  expect(first.textContent).toContain('Session actions');
  first.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowDown', bubbles: true, cancelable: true }));
  const second = document.activeElement as HTMLButtonElement;
  expect(second.textContent).toContain('Models');
  second.click();
  expect(onopen).toHaveBeenCalledWith('manager');
  second.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowUp', bubbles: true, cancelable: true }));
  expect(document.activeElement).toBe(first);
 } finally { await unmount(app); document.body.replaceChildren(); }
});
