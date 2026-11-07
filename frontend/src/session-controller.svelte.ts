import type { Conversation } from './conversation';
import { currentSession, createSession, openSession, listSessions, saveSession } from './sessions';
import type { Session } from './sessions';

export function createSessionController(conversation: Conversation) {
  const state = $state({ id: '', items: [] as Session[], loading: true, error: '' });
  let timer: ReturnType<typeof setTimeout> | undefined;
  let queue = Promise.resolve();

  function apply(session: Session) {
    state.id = session.id;
    conversation.messages = session.messages ?? [];
    conversation.draft = session.draft;
    conversation.error = '';
    conversation.activity = 'Ready';
    conversation.executions = [];
    conversation.usage = session.usage ?? null;
  }

  async function save() {
    clearTimeout(timer);
    if (!state.id) return;
    const id = state.id;
    const messages = JSON.parse(JSON.stringify(conversation.messages));
    const draft = conversation.draft;
    queue = queue.catch(() => {}).then(() => saveSession(id, messages, draft));
    try { await queue; state.error = ''; }
    catch (error) { state.error = `Could not save session: ${String(error)}`; throw error; }
  }

  return {
    state,
    save,
    schedule() {
      clearTimeout(timer);
      if (!conversation.busy) timer = setTimeout(() => { void save().catch(() => {}); }, 250);
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
      if (conversation.busy || state.loading) return false;
      state.loading = true;
      try {
        await save();
        apply(id ? await openSession(id) : await createSession());
        return true;
      } catch (error) { state.error = `Could not open session: ${String(error)}`; return false; }
      finally { state.loading = false; }
    },
    dispose() { clearTimeout(timer); },
  };
}
