import { call } from './bridge';
import type { Usage } from './agent-types';
import type { Session } from './sessions';
import type { createSessionController } from './session-controller.svelte';

export async function contextAction(sessions: ReturnType<typeof createSessionController>, action: 'compact' | 'handoff', id = sessions.state.id) {
  const conversation = sessions.conversations[id];
  if (!conversation || conversation.busy || !conversation.messages.length) return null;
  conversation.busy = true;
  conversation.error = '';
  conversation.compactFailed = false;
  conversation.activity = action === 'compact' ? 'Compacting context' : 'Preparing handoff';
  try {
    await sessions.save(id);
    if (action === 'compact') {
      conversation.usage = await call<Usage>('AgentService.Compact', id);
      return null;
    }
    const session = await call<Session>('AgentService.Handoff', id);
    return session.id;
  } catch (error) {
    conversation.error = String(error);
    if (action === 'compact') conversation.compactFailed = true;
    return null;
  }
  finally { conversation.busy = false; conversation.activity = 'Ready'; }
}
