import type { createSettings } from './settings.svelte';
import type { createModelController } from './model-controller.svelte';
import { listen } from './bridge';
import type { createSessionController } from './session-controller.svelte';
import type { createWorkspaces } from './workspace-controller.svelte';

export function listenToBackend(sessions: ReturnType<typeof createSessionController>, workspaces: ReturnType<typeof createWorkspaces>, settings: ReturnType<typeof createSettings>, models: ReturnType<typeof createModelController>) {
  const preferences = listen('settings:changed', () => { void settings.start(); });
  const model = listen('model:changed', () => { void models.start(); });
  const title = listen('session:title', (data) => {
    const { id, title } = data as { id: string; title: string };
    sessions.rename(id, title);
  });
  const session = listen('session:changed', () => { void sync(); });
  const changed = listen('workspace:changed', () => { void workspaces.load(); });
  let syncing = false;
  async function sync() {
    if (syncing) return;
    syncing = true;
    try { await sessions.sync(); await sessions.refresh(); } catch {}
    finally { syncing = false; }
  }
  const timer = setInterval(() => { void sync(); }, 700);
  return () => { clearInterval(timer); void preferences.then((off) => off()); void model.then((off) => off()); void session.then((off) => off()); void title.then((off) => off()); void changed.then((off) => off()); };
}
