import { sendMessage } from './conversation';
import { runAgent, stopAgent } from './agent';
import type { createSessionController } from './session-controller.svelte';

type Sessions = ReturnType<typeof createSessionController>;

export function createChat(sessions: Sessions, onFinish: (id: string, succeeded: boolean) => void) {
  return {
    async submit() {
      const id = sessions.state.id;
      const conversation = sessions.conversations[id];
      if (sessions.state.loading || !conversation || conversation.busy || !conversation.draft.trim()) return null;
      const succeeded = await sendMessage(conversation, async (messages, observe) => {
        await sessions.save(id);
        return runAgent(messages, observe, id);
      });
      if (id !== sessions.state.id) {
        conversation.unread = true;
        onFinish(id, succeeded);
      }
      await sessions.save(id).catch(() => {});
      return id;
    },
    stop() {
      const id = sessions.state.id;
      if (!sessions.conversations[id]?.busy) return false;
      void stopAgent(id).catch(() => {});
      return true;
    },
  };
}
