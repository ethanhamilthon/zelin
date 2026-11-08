import { call } from './bridge';
import type { Session } from './sessions';
import type { Conversation } from './conversation';
import type { AgentEvent } from './agent-types';
import { observeEvent } from './observe';

export const localRuns = new Set<string>();
type Live = { started: number; events: AgentEvent[] } | null;
export async function remoteSession(id: string) {
  const session = await call<Session>('SessionService.Get', id);
  const live = await call<Live>('AgentService.Active', id);
  return { session, live };
}
export function updateRemote(target: Conversation, session: Session, live: Live, dirty: boolean) {
  if (localRuns.has(session.id)) return;
  const wasBusy = target.busy;
  target.messages = session.messages ?? [];
  target.usage = session.usage ?? null;
  target.workspaceId = session.workspaceId;
  if (!dirty) {
    target.draft = session.draft;
    target.attachments = session.attachments ?? [];
    target.refs = session.refs ?? [];
  }
  target.busy = Boolean(live);
  target.started = live?.started;
  if (live) {
    target.executions = [];
    target.streaming = '';
    target.activity = 'Thinking';
    for (const event of live.events) observeEvent(target, event);
  } else if (wasBusy) {
    target.activity = 'Done';
    target.executions = [];
    target.streaming = '';
    target.unread = true;
  }
}
