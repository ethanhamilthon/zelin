import type { Message } from './conversation';
import type { AgentEvent, AgentInfo, AgentResult } from './agent-types';

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

export async function agentInfo(): Promise<AgentInfo> {
  if (!desktop()) return { model: 'Desktop only', contextWindow: 0 };
  const { Call } = await import('@wailsio/runtime');
  return Call.ByName('main.AgentService.Info');
}

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
