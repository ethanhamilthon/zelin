import { suggest } from './browser';
import type { Entry } from './refs';
import type { Mention } from './mention';

export function createMentions(workspace: () => string, delay = 120) {
  const state = $state({ open: false, items: [] as Entry[], active: 0, loading: false, mention: null as Mention });
  let timer: ReturnType<typeof setTimeout> | undefined;
  let serial = 0;

  function close() {
    clearTimeout(timer);
    serial++;
    state.open = false; state.mention = null; state.items = []; state.active = 0; state.loading = false;
  }

  function update(mention: Mention) {
    if (!mention || !workspace()) { if (state.open) close(); return; }
    state.mention = mention;
    state.open = true;
    state.loading = true;
    clearTimeout(timer);
    const mine = ++serial;
    timer = setTimeout(async () => {
      try {
        const items = await suggest(workspace(), '', mention.query);
        if (mine !== serial) return;
        state.items = items.slice(0, 5);
        state.active = 0;
      } catch { if (mine === serial) state.items = []; }
      finally { if (mine === serial) state.loading = false; }
    }, delay);
  }

  return {
    state,
    close,
    update,
    move(step: number) {
      if (!state.items.length) return;
      state.active = (state.active + step + state.items.length) % state.items.length;
    },
    get chosen() { return state.items[state.active]; },
    dispose() { clearTimeout(timer); serial++; },
  };
}
