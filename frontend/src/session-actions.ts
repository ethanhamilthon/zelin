import { call } from './bridge';
import { stopAgent } from './agent';
import { contextAction } from './context-actions';
import type { createSessionController } from './session-controller.svelte';
import type { createToasts } from './toasts.svelte';

export type SessionAction = 'compact' | 'handoff' | 'export' | 'stop' | 'rename' | 'tasks';
type Sessions = ReturnType<typeof createSessionController>;

const reason = (error: unknown) => String(error).replace(/^Error: /, '');

export async function saveTitle(sessions: Sessions, id: string, title: string): Promise<string | null> {
  if (!title.trim()) return 'Title cannot be empty';
  try { sessions.rename(id, await call<string>('SessionService.Rename', id, title)); return null; }
  catch (error) { return reason(error); }
}

export async function runSessionAction(
  sessions: Sessions, toasts: ReturnType<typeof createToasts>, action: Exclude<SessionAction, 'rename' | 'tasks'>, id: string,
): Promise<string | null> {
  const conversation = sessions.conversations[id];
  if (action === 'stop') {
    if (conversation?.busy) await stopAgent(id).catch(() => {});
    return null;
  }
  if (action !== 'export') return contextAction(sessions, action, id);
  try {
    if (conversation && !conversation.busy) await sessions.save(id);
    const path = await call<string>('SessionService.Export', id);
    if (path) toasts.push('', `Exported to ${path}`);
  } catch (error) { toasts.push('', `Export failed: ${reason(error)}`); }
  return null;
}
