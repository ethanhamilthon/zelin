import { sendMessage } from './conversation';
import { runAgent, stopAgent, steerAgent } from './agent';
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
    async steer() {
      const id = sessions.state.id, conversation = sessions.conversations[id];
      if (!conversation?.busy || !conversation.started || !conversation.draft.trim() || conversation.steering) return;
      conversation.steering = true;
      try {
        const accepted = await steerAgent(id, conversation.draft, conversation.attachments);
        if (!accepted) conversation.steering = false;
      } catch (error) { conversation.steering = false; conversation.error = String(error); }
    },
    edit() {
      const id = sessions.state.id, conversation = sessions.conversations[id];
      if (!conversation?.steering) return;
      conversation.steering = false;
      void steerAgent(id, '').catch((error) => { conversation.error = String(error); });
    },
    stop() {
      const id = sessions.state.id;
      if (!sessions.conversations[id]?.busy) return false;
      void stopAgent(id).catch(() => {});
      return true;
    },
  };
}
