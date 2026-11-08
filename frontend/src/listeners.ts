import { listen } from './bridge';
import type { createSessionController } from './session-controller.svelte';
import type { createWorkspaces } from './workspace-controller.svelte';

export function listenToBackend(sessions: ReturnType<typeof createSessionController>, workspaces: ReturnType<typeof createWorkspaces>) {
  const title = listen('session:title', (data) => {
    const { id, title } = data as { id: string; title: string };
    sessions.rename(id, title);
  });
  const changed = listen('workspace:changed', () => { void workspaces.load(); });
  return () => { void title.then((off) => off()); void changed.then((off) => off()); };
}
