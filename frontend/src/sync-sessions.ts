import { remoteSession, updateRemote, localRuns } from './remote-session';
import type { Conversation } from './conversation';
export type SavedDraft = { draft: string; attachments: string[]; refs: string };
export async function syncSession(id: string, conversations: Record<string, Conversation>, saved: Map<string, SavedDraft>) {
 if (!id || localRuns.has(id)) return;
 const before = conversations[id]?.draft, known = saved.get(id);
 const { session, live } = await remoteSession(id);
 const target = conversations[id];
 if (!target || localRuns.has(id)) return;
 const previous = saved.get(id);
 const dirty = target.draft !== before || previous !== known || target.draft !== previous?.draft || JSON.stringify(target.attachments) !== JSON.stringify(previous?.attachments ?? []) || JSON.stringify(target.refs) !== previous?.refs;
 updateRemote(target, session, live, dirty);
 target.unread = false;
 if (!dirty) saved.set(id, { draft: session.draft, attachments: [...(session.attachments ?? [])], refs: JSON.stringify(session.refs ?? []) });
}
