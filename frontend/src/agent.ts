import type { Message } from './conversation';
import type { AgentEvent, AgentInfo, AgentResult, ModelOption } from './agent-types';
import { placeholderInfo } from './agent-types';
import { call, desktop, listen } from './bridge';

export async function agentInfo(sessionId = ''): Promise<AgentInfo> {
  if (!desktop()) return placeholderInfo('Desktop only');
  return call<AgentInfo>('AgentService.Info', sessionId);
}

export const listModels = () => call<ModelOption[]>('AgentService.Models');
export const configureModel = (sessionId: string, id: string, effort: string) => call<AgentInfo>('AgentService.Configure', sessionId, id, effort);
export const setTitleModel = (id: string) => call<AgentInfo>('AgentService.SetTitleModel', id);
export const steerAgent = (sessionId: string, content: string, images: string[] = []) => call<boolean>('AgentService.Steer', sessionId, { role: 'user', content, images });
export const stopAgent = (sessionId: string) => call<boolean>('AgentService.Stop', sessionId);

export async function runAgent(
  messages: Message[], onEvent: (event: AgentEvent) => void, sessionId: string,
): Promise<AgentResult> {
  if (!desktop()) throw new Error('The agent is available in the desktop app. Start Zash with go run .');
  const off = await listen('agent:update', (data) => {
    const event = data as AgentEvent & { session?: string };
    if (event.session === sessionId) onEvent(event);
  });
  try {
    return await call<AgentResult>('AgentService.Run', sessionId, messages.map(({ role, content, images }) => ({ role, content, images })));
  } finally { off(); }
}
