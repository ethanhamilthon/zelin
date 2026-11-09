import type { Conversation } from './conversation';
import type { AgentEvent } from './agent-types';
import { usedRefs } from './refs';

export function observeEvent(state: Conversation, event: AgentEvent) {
    if (event.kind === 'steer') {
      const steerRefs = usedRefs(event.message.content, state.refs);
      const steps = state.executions.filter((item) => item.state === 'complete');
      state.messages = [...state.messages, { role: 'user', content: event.message.content, images: event.message.images, ...(steerRefs.length && { refs: steerRefs }), ...(steps.length && { executions: steps }) }];
      state.executions = state.executions.filter((item) => item.state !== 'complete');
      if (state.steering) { state.draft = ''; state.attachments = []; state.refs = []; }
      state.steering = false;
    }
    else if (event.kind === 'notice') {
      const steps = state.executions.filter((item) => item.state === 'complete');
      state.messages = [...state.messages, { role: 'user', kind: 'task', content: event.message.content, ...(steps.length && { executions: steps }) }];
      state.executions = state.executions.filter((item) => item.state !== 'complete');
    }
    else if (event.kind === 'activity') state.activity = event.activity;
    else if (event.kind === 'compact_error') { state.error = event.activity; state.compactFailed = true; }
    else if (event.kind === 'usage') state.usage = event.usage;
    else if (event.kind === 'stream') state.streaming = event.message.content;
    else if (event.kind === 'execution') {
      const execution = event.execution;
      const index = state.executions.findIndex((item) => item.id === execution.id);
      if (index < 0) state.executions = [...state.executions, execution];
      else state.executions = state.executions.map((item, i) => i === index ? execution : item);
    }
}
