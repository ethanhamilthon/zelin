import type { Message } from './conversation';
import type { AgentEvent, AgentInfo, AgentResult, ModelOption } from './agent-types';
import { placeholderInfo } from './agent-types';

type BridgeWindow = Window & {
  chrome?: { webview?: { postMessage?: unknown } };
  webkit?: { messageHandlers?: { external?: { postMessage?: unknown } } };
  wails?: { invoke?: unknown };
};

function desktop() {
  const bridge = window as BridgeWindow;
  return Boolean(bridge.webkit?.messageHandlers?.external?.postMessage ||
    bridge.chrome?.webview?.postMessage || bridge.wails?.invoke);
}

async function call<T>(method: string, ...args: unknown[]): Promise<T> {
  if (!desktop()) throw new Error('Available in the desktop app only');
  const { Call } = await import('@wailsio/runtime');
  return Call.ByName(`main.AgentService.${method}`, ...args);
}

export async function agentInfo(): Promise<AgentInfo> {
  if (!desktop()) return placeholderInfo('Desktop only');
  return call<AgentInfo>('Info');
}

export const listModels = () => call<ModelOption[]>('Models');
export const configureModel = (id: string, effort: string) => call<AgentInfo>('Configure', id, effort);

export async function runAgent(
  messages: Message[], onEvent: (event: AgentEvent) => void, sessionId: string,
): Promise<AgentResult> {
  if (!desktop()) throw new Error('The agent is available in the desktop app. Start Zash with go run .');
  const { Call, Events } = await import('@wailsio/runtime');
  const off = Events.On('agent:update', (event) => onEvent(event.data as AgentEvent));
  try {
    return await Call.ByName('main.AgentService.Run', sessionId, messages.map(({ role, content }) => ({ role, content })));
  } finally { off(); }
}
