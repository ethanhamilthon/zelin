import { tick } from 'svelte';

export async function keepScroll(from: Element | null, change: () => void) {
  const scroller = from?.closest('main') ?? null;
  const height = scroller?.scrollHeight ?? 0, top = scroller?.scrollTop ?? 0;
  change();
  await tick();
  if (scroller) scroller.scrollTop = top + scroller.scrollHeight - height;
}
