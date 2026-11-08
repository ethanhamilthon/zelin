import type { AgentEvent, AgentResult, Execution, Usage } from './agent-types';
import type { Turn } from './turn';

export type Message = { role: 'user' | 'assistant'; content: string; images?: string[]; executions?: Execution[]; turn?: Turn };
export type Conversation = {
  messages: Message[];
  draft: string;
  attachments: string[];
  unread: boolean;
  busy: boolean;
  error: string;
  activity: string;
  executions: Execution[];
  usage: Usage | null;
  started?: number;
  baseline?: Usage | null;
};

export async function sendMessage(
  state: Conversation,
  run: (messages: Message[], onEvent: (event: AgentEvent) => void) => Promise<AgentResult>,
): Promise<boolean> {
  if (state.busy || !state.draft.trim()) return false;
  const draft = state.draft;
  const attachments = state.attachments;
  const previous = state.messages;
  const images = attachments.map((url, index) => draft.includes(`[Image #${index + 1}]`) ? url : '');
  state.messages = [...previous, images.some(Boolean) ? { role: 'user', content: draft, images } : { role: 'user', content: draft }];
  state.draft = '';
  state.attachments = [];
  state.busy = true;
  state.error = '';
  state.activity = 'Thinking';
  state.executions = [];
  state.started = Date.now();
  state.baseline = state.usage;
  const observe = (event: AgentEvent) => {
    if (event.kind === 'activity') state.activity = event.activity;
    else if (event.kind === 'usage') state.usage = event.usage;
    else if (event.kind === 'execution') {
      const execution = event.execution;
      const index = state.executions.findIndex((item) => item.id === execution.id);
      if (index < 0) state.executions = [...state.executions, execution];
      else state.executions = state.executions.map((item, i) => i === index ? execution : item);
    }
  };
  let turn: AgentResult['turn'];
  try {
    const result = await run(state.messages.map((message) => ({ ...message })), observe);
    state.executions = result.executions;
    turn = result.turn;
    if (result.error) throw new Error(result.error);
    if (!result.content.trim()) throw new Error('Agent returned an empty response');
    state.messages = [...state.messages, {
      role: 'assistant', content: result.content, executions: state.executions, ...(turn && { turn }),
    }];
    return true;
  } catch (error) {
    if (state.executions.length) {
      state.messages = [...state.messages, {
        role: 'assistant', content: 'The agent stopped before completing this turn.', executions: state.executions, ...(turn && { turn }),
      }];
    } else {
      state.messages = previous;
      if (!state.draft) { state.draft = draft; state.attachments = attachments; }
    }
    state.error = error instanceof Error ? error.message : String(error);
    return false;
  } finally {
    state.executions = [];
    state.started = undefined;
    state.activity = state.error ? 'Error' : 'Ready';
    state.busy = false;
  }
}
