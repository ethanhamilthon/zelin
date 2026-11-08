import type { FileAttachment } from './agent-types';
import { observeEvent } from './observe';
import type { AgentEvent, AgentResult, Execution, Usage } from './agent-types';
import type { Turn } from './turn';
import { usedRefs, type Ref } from './refs';

export const failure = (reason: string) => reason === 'Stopped' ? 'Stopped by user.' : `Error: ${reason}`;

export type Message = { role: 'user' | 'assistant'; content: string; files?: FileAttachment[]; images?: string[]; executions?: Execution[]; turn?: Turn; refs?: Ref[]; error?: boolean };
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
  streaming?: string;
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
  state.streaming = '';
  state.started = Date.now();
  state.baseline = state.usage;
  const observe = (event: AgentEvent) => observeEvent(state, event);
  let turn: AgentResult['turn'];
  let files: FileAttachment[] | undefined;
  try {
    const result = await run(state.messages.map((message) => ({ ...message })), observe, refs);
    state.executions = result.executions;
    turn = result.turn;
    files = result.files;
    if (result.error) throw new Error(result.error);
    if (!result.content.trim() && !files?.length) throw new Error('Agent returned an empty response');
    state.messages = [...state.messages, {
      role: 'assistant', content: result.content, files, executions: state.executions, ...(turn && { turn }),
    }];
    return true;
  } catch (error) {
    const reason = error instanceof Error ? error.message : String(error);
    state.messages = [...state.messages, {
      role: 'assistant', content: failure(reason), error: true, executions: state.executions, ...(turn && { turn }),
    }];
    return false;
  } finally {
    state.executions = [];
    state.streaming = '';
    state.started = undefined;
    state.activity = state.messages.at(-1)?.error ? 'Error' : 'Ready';
    state.busy = false;
    state.steering = false;
  }
}
