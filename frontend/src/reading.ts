export function readToEnd(scroller: HTMLElement): boolean {
  const last = scroller.querySelector('article:last-of-type');
  if (!last) return true;
  if (document.visibilityState === 'hidden') return false;
  return last.getBoundingClientRect().bottom <= scroller.getBoundingClientRect().bottom + 1;
}
