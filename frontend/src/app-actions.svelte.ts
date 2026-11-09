import { runSessionAction, saveTitle, type SessionAction } from './session-actions';
import type { createSessionController } from './session-controller.svelte';
import type { createToasts } from './toasts.svelte';
import type { createModalStack } from './modal-stack.svelte';

type Deps = {
  sessions: ReturnType<typeof createSessionController>; toasts: ReturnType<typeof createToasts>; modals: ReturnType<typeof createModalStack>;
  close: () => void; select: (id?: string) => Promise<void>;
};

export function createAppActions({ sessions, toasts, modals, close, select }: Deps) {
  const state = $state({ target: '' });
  return {
    state,
    aim() { state.target = sessions.state.id; },
    async run(action: SessionAction) {
      if (action === 'rename' || action === 'tasks') { modals.open(action); return; }
      close();
      const next = await runSessionAction(sessions, toasts, action, state.target);
      if (next) await select(next);
    },
    async rename(title: string) {
      const error = await saveTitle(sessions, state.target, title);
      if (!error) close();
      return error;
    },
  };
}
