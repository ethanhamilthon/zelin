import { chime } from './sound';
import type { createSettings } from './settings.svelte';
import type { createSessionController } from './session-controller.svelte';
import type { createToasts } from './toasts.svelte';

export function finish(settings: ReturnType<typeof createSettings>, sessions: ReturnType<typeof createSessionController>, toasts: ReturnType<typeof createToasts>) {
  return async (id: string, succeeded: boolean) => {
    if (settings.state.sound) chime();
    if (!settings.state.toast) return;
    await sessions.refresh();
    const title = sessions.state.items.find((item) => item.id === id)?.title ?? 'Session';
    toasts.push(id, `${title} · ${succeeded ? 'finished' : 'stopped'}`);
  };
}
