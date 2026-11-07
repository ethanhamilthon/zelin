import type { Conversation } from './conversation';
import { currentSession, createSession, openSession, listSessions, saveSession } from './sessions';
import type { Session } from './sessions';

export type SessionStatus = 'running' | 'unread' | 'error' | '';

export const blank = (): Conversation => ({
  messages: [], draft: '', attachments: [], unread: false, busy: false, error: '', activity: 'Ready', executions: [], usage: null,
});

export function createSessionController() {
  const state = $state({ id: '', items: [] as Session[], loading: true, error: '' });
  const conversations = $state<Record<string, Conversation>>({});
  const timers = new Map<string, ReturnType<typeof setTimeout>>();
  let queue = Promise.resolve();

  function apply(session: Session) {
    conversations[session.id] ??= {
      ...blank(), messages: session.messages ?? [], draft: session.draft,
      attachments: session.attachments ?? [], usage: session.usage ?? null,
    };
    state.id = session.id;
    conversations[session.id].unread = false;
  }

  async function save(id = state.id) {
    clearTimeout(timers.get(id));
    const conversation = conversations[id];
    if (!conversation) return;
    const snapshot = JSON.parse(JSON.stringify([conversation.messages, conversation.attachments]));
    const draft = conversation.draft;
    queue = queue.catch(() => {}).then(() => saveSession(id, snapshot[0], draft, snapshot[1]));
    try { await queue; state.error = ''; }
    catch (error) { state.error = `Could not save session: ${String(error)}`; throw error; }
  }

  function schedule(id = state.id) {
    clearTimeout(timers.get(id));
    if (!conversations[id]?.busy) timers.set(id, setTimeout(() => { void save(id).catch(() => {}); }, 250));
  }

  return {
    state,
    conversations,
    get current() { return conversations[state.id]; },
    save,
    status(id: string): SessionStatus {
      const conversation = conversations[id];
      if (!conversation) return '';
      return conversation.busy ? 'running' : conversation.unread ? 'unread' : conversation.error ? 'error' : '';
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
    async refresh() {
      try { state.items = await listSessions(); }
      catch (error) { state.error = `Could not list sessions: ${String(error)}`; }
    },
    async select(id?: string) {
      if (state.loading) return false;
      state.loading = true;
      try {
        if (state.id && !conversations[state.id]?.busy) await save();
        apply(id ? await openSession(id) : await createSession());
        return true;
      } catch (error) { state.error = `Could not open session: ${String(error)}`; return false; }
      finally { state.loading = false; }
    },
    dispose() { timers.forEach((timer) => clearTimeout(timer)); },
  };
}
