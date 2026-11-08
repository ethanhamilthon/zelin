import type { AgentEvent, AgentResult, Execution, Usage } from './agent-types';
import type { Turn } from './turn';
import { usedRefs, type Ref } from './refs';

export type Message = { role: 'user' | 'assistant'; content: string; images?: string[]; executions?: Execution[]; turn?: Turn; refs?: Ref[] };
export type Conversation = {
  messages: Message[];
  draft: string;
  attachments: string[];
  refs: Ref[];
  unread: boolean;
  busy: boolean;
  error: string;
  activity: string;
  executions: Execution[];
  usage: Usage | null;
  started?: number;
  baseline?: Usage | null;
  steering?: boolean;
  compactFailed?: boolean;
  workspaceId?: string;
};

export async function sendMessage(
  state: Conversation,
  run: (messages: Message[], onEvent: (event: AgentEvent) => void, refs: Ref[]) => Promise<AgentResult>,
): Promise<boolean> {
  if (state.busy || !state.draft.trim()) return false;
  const draft = state.draft;
  const attachments = state.attachments;
  const refs = usedRefs(draft, state.refs);
  const allRefs = state.refs;
  const previous = state.messages;
  const images = attachments.map((url, index) => draft.includes(`[Image #${index + 1}]`) ? url : '');
  const sent: Message = images.some(Boolean) ? { role: 'user', content: draft, images } : { role: 'user', content: draft };
  if (refs.length) sent.refs = refs;
  state.messages = [...previous, sent];
  state.draft = '';
  state.attachments = [];
  state.refs = [];
  state.busy = true;
  state.error = '';
  state.activity = 'Thinking';
  state.executions = [];
  state.started = Date.now();
  state.baseline = state.usage;
  const observe = (event: AgentEvent) => {
    if (event.kind === 'steer') {
      const steerRefs = usedRefs(event.message.content, state.refs);
      state.messages = [...state.messages, { role: 'user', content: event.message.content, images: event.message.images, ...(steerRefs.length && { refs: steerRefs }) }];
      if (state.steering) { state.draft = ''; state.attachments = []; state.refs = []; }
      state.steering = false;
    }
    else if (event.kind === 'activity') state.activity = event.activity;
    else if (event.kind === 'compact_error') { state.error = event.activity; state.compactFailed = true; }
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
    const result = await run(state.messages.map((message) => ({ ...message })), observe, refs);
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
      if (!state.draft) { state.draft = draft; state.attachments = attachments; state.refs = allRefs; }
    }
    state.error = error instanceof Error ? error.message : String(error);
    return false;
  } finally {
    state.executions = [];
    state.started = undefined;
    state.activity = state.error ? 'Error' : 'Ready';
    state.busy = false;
    state.steering = false;
  }
}
