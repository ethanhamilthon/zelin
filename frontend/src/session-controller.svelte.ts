import { createSharedActivity } from './shared-activity.svelte';
import { syncSession } from './sync-sessions';
import type { Conversation } from './conversation';
import { currentSession, createSession, openSession, listSessions, saveSession, saveDraft, saveAttachments, saveRefs } from './sessions';
import type { Session } from './sessions';
import { evictable, stale } from './retention';
import { workspaceSession } from './workspaces';
export type SessionStatus = 'running' | 'unread' | 'error' | '';
export const blank = (): Conversation => ({
  messages: [], draft: '', attachments: [], refs: [], unread: false, busy: false, error: '', activity: 'Ready', executions: [], usage: null,
});
const same = (a: string[], b: string[]) => a.length === b.length && a.every((item, index) => item === b[index]);
export function createSessionController() {
  const shared = createSharedActivity();
  const state = $state({ id: '', items: [] as Session[], loading: true, error: '' });
  const conversations = $state<Record<string, Conversation>>({});
  const timers = new Map<string, ReturnType<typeof setTimeout>>();
  const saved = new Map<string, { draft: string; attachments: string[]; refs: string }>();
  const order: string[] = [];
  let queue = Promise.resolve();
  function retain() {
    for (const id of stale(order, state.id, (item) => evictable(conversations[item], saved.get(item)))) {
      delete conversations[id]; saved.delete(id); order.splice(order.indexOf(id), 1);
    }
  }
  function enqueue(write: () => Promise<void>) {
    const run = queue.catch(() => {}).then(write);
    queue = run;
    return run.then(() => { state.error = ''; }, (error) => {
      state.error = `Could not save session: ${String(error)}`;
      throw error;
    });
  }
  function apply(session: Session) {
    conversations[session.id] ??= {
      ...blank(), messages: session.messages ?? [], draft: session.draft,
      attachments: session.attachments ?? [], refs: session.refs ?? [], usage: session.usage ?? null,
    };
    conversations[session.id].workspaceId = session.workspaceId;
    saved.set(session.id, { draft: session.draft, attachments: [...(session.attachments ?? [])], refs: JSON.stringify(session.refs ?? []) });
    state.id = session.id;
    conversations[session.id].unread = false;
    shared.read(session.id);
    const known = order.indexOf(session.id);
    if (known >= 0) order.splice(known, 1);
    order.push(session.id);
    retain();
  }
  async function save(id = state.id) {
    clearTimeout(timers.get(id));
    timers.delete(id);
    const conversation = conversations[id];
    if (!conversation) return;
    const snapshot = JSON.parse(JSON.stringify([conversation.messages, conversation.attachments, conversation.refs]));
    const draft = conversation.draft;
    await enqueue(() => saveSession(id, snapshot[0], draft, snapshot[1]));
    await enqueue(() => saveRefs(id, snapshot[2]));
    saved.set(id, { draft, attachments: [...snapshot[1]], refs: JSON.stringify(snapshot[2]) });
  }
  async function flush(id = state.id) {
    clearTimeout(timers.get(id));
    timers.delete(id);
    const conversation = conversations[id];
    if (!conversation) return;
    const previous = saved.get(id) ?? { draft: '', attachments: [], refs: '[]' };
    const draft = conversation.draft, attachments = [...conversation.attachments];
    if (draft !== previous.draft) {
      await enqueue(() => saveDraft(id, draft));
      saved.set(id, { ...(saved.get(id) ?? previous), draft });
    }
    if (!same(attachments, previous.attachments)) {
      await enqueue(() => saveAttachments(id, attachments));
      saved.set(id, { ...(saved.get(id) ?? previous), attachments });
    }
    const refs = JSON.stringify(conversation.refs);
    if (refs !== previous.refs) {
      await enqueue(() => saveRefs(id, JSON.parse(refs)));
      saved.set(id, { ...(saved.get(id) ?? previous), refs });
    }
  }
  function schedule(id = state.id) {
    clearTimeout(timers.get(id));
    timers.set(id, setTimeout(() => { void flush(id).catch(() => {}); }, 250));
  }
  return {
    state,
    conversations,
    get current() { return conversations[state.id]; },
    save,
    flush,
    status(id: string): SessionStatus {
      const remote = shared.state.items[id];
      if (remote) return remote.running ? 'running' : remote.unread ? 'unread' : '';
      const conversation = conversations[id];
      if (!conversation) return '';
      return conversation.busy ? 'running' : conversation.unread ? 'unread' : conversation.error ? 'error' : '';
    },
    get activity() {
      let working = 0, done = 0;
      for (const id of new Set([...Object.keys(conversations), ...Object.keys(shared.state.items)])) {
        const remote = shared.state.items[id], local = conversations[id];
        const item = remote ? { busy: remote.running, unread: remote.unread, error: '' } : local;
        if (item.busy) working++;
        else if (item.unread && !item.error) done++;
      }
      return { working, done };
    },
    steps(id: string): number {
      return conversations[id]?.executions.length ?? 0;
    },
    schedule,
    attach: (url: string) => {
      const conversation = conversations[state.id];
      if (!conversation) return 0;
      conversation.attachments = [...conversation.attachments, url];
      schedule();
      return conversation.attachments.length;
    },
    rename(id: string, title: string) {
      state.items = state.items.map((item) => item.id === id ? { ...item, title } : item);
    },
    async start() {
      try { apply(await currentSession()); }
      catch (error) { state.error = `Could not load session: ${String(error)}`; }
      finally { state.loading = false; }
    },
    async sync() {
      await shared.sync();
      shared.read(state.id);
      await queue.catch(() => {});
      await syncSession(state.id, conversations, saved);
    },
    async refresh() {
      try { state.items = await listSessions(); }
      catch (error) { state.error = `Could not list sessions: ${String(error)}`; }
    },
    async select(id?: string, workspace = '') {
      if (state.loading) return false;
      state.loading = true;
      try {
        if (state.id && !conversations[state.id]?.busy) await flush();
        apply(id ? await openSession(id) : workspace ? await workspaceSession(workspace) : await createSession());
        return true;
      } catch (error) { state.error = `Could not open session: ${String(error)}`; return false; }
      finally { state.loading = false; }
    },
    dispose() { timers.forEach((timer) => clearTimeout(timer)); },
  };
}
