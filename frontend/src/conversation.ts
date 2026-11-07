import type { AgentEvent, AgentResult, Execution } from './agent-types';

export type Message = { role: 'user' | 'assistant'; content: string; executions?: Execution[] };
export type Conversation = {
  messages: Message[];
  draft: string;
  busy: boolean;
  error: string;
  activity: string;
  executions: Execution[];
};

export async function sendMessage(
  state: Conversation,
  run: (messages: Message[], onEvent: (event: AgentEvent) => void) => Promise<AgentResult>,
): Promise<boolean> {
  if (state.busy || !state.draft.trim()) return false;
  const draft = state.draft;
  const previous = state.messages;
  state.messages = [...previous, { role: 'user', content: draft }];
  state.draft = '';
  state.busy = true;
  state.error = '';
  state.activity = 'Thinking';
  state.executions = [];
  const observe = (event: AgentEvent) => {
    if (event.kind === 'activity') state.activity = event.activity;
    else if (event.kind === 'execution') {
      const execution = event.execution;
      const index = state.executions.findIndex((item) => item.id === execution.id);
      if (index < 0) state.executions = [...state.executions, execution];
      else state.executions = state.executions.map((item, i) => i === index ? execution : item);
    }
  };
  try {
    const result = await run(state.messages.map((message) => ({ ...message })), observe);
    state.executions = result.executions;
    if (result.error) throw new Error(result.error);
    if (!result.content.trim()) throw new Error('Agent returned an empty response');
    state.messages = [...state.messages, {
      role: 'assistant', content: result.content, executions: state.executions,
    }];
    return true;
  } catch (error) {
    if (state.executions.length) {
      state.messages = [...state.messages, {
        role: 'assistant', content: 'The agent stopped before completing this turn.', executions: state.executions,
      }];
    } else {
      state.messages = previous;
      if (!state.draft) state.draft = draft;
    }
    state.error = error instanceof Error ? error.message : String(error);
    return false;
  } finally {
    state.executions = [];
    state.activity = state.error ? 'Error' : 'Ready';
    state.busy = false;
  }
}
