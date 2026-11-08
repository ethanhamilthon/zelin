import type { ModalName } from './Panels.svelte';

export function createModalStack() {
  let items = $state<ModalName[]>([]);
  return {
    get current() { return items.at(-1) ?? null; },
    get canBack() { return items.length > 1; },
    open(name: ModalName) {
      if (items.at(-1) !== name) items = [...items, name];
    },
    toggle(name: ModalName) {
      if (items.at(-1) === name) items = [];
      else items = [name];
    },
    back() {
      if (!items.length) return false;
      items = items.slice(0, -1);
      return true;
    },
    clear() { items = []; },
  };
}
