import { describe, expect, it, vi } from 'vitest';
import { sendMessage } from './conversation';
import type { Conversation } from './conversation';
import type { AgentResult, Execution } from './agent-types';

function initial(draft = 'Hello'): Conversation {
  return { messages: [], draft, busy: false, error: '', activity: 'Ready', executions: [] };
}
const answer = (content: string): AgentResult => ({ content, executions: [] });
const execution: Execution = { id: 'call', code: '1 + 1', output: '2', state: 'complete' };

describe('conversation', () => {
  it('keeps the next draft while rejecting duplicate sends', async () => {
    const state = initial();
    let finish!: (result: AgentResult) => void;
    const run = vi.fn(() => new Promise<AgentResult>((resolve) => { finish = resolve; }));
    const pending = sendMessage(state, run);
    expect(state.busy).toBe(true);
    expect(state.draft).toBe('');
    expect(state.messages).toEqual([{ role: 'user', content: 'Hello' }]);
    state.draft = 'Next';
    expect(await sendMessage(state, run)).toBe(false);
    finish(answer('**Hi**'));
    expect(await pending).toBe(true);
    expect(state.busy).toBe(false);
    expect(state.draft).toBe('Next');
    expect(state.messages[1].content).toBe('**Hi**');
    await sendMessage(state, async (messages) => {
      expect(messages).toHaveLength(3);
      return answer('Again');
    });
    expect(state.messages).toHaveLength(4);
    expect(run).toHaveBeenCalledTimes(1);
  });

  it('restores failed input without overwriting a new draft', async () => {
    for (const next of ['', 'New draft']) {
      const state = initial();
      state.messages = [{ role: 'assistant', content: 'Earlier' }];
      await sendMessage(state, async () => {
        state.draft = next;
        throw new Error('Offline');
      });
      expect(state.draft).toBe(next || 'Hello');
      expect(state.error).toBe('Offline');
      expect(state.busy).toBe(false);
      expect(state.messages).toEqual([{ role: 'assistant', content: 'Earlier' }]);
    }
  });

  it('updates live executions and retains them in the answer', async () => {
    const state = initial();
    await sendMessage(state, async (_, observe) => {
      observe({ kind: 'activity', activity: 'Running JavaScript' });
      observe({ kind: 'execution', execution: { ...execution, output: '', state: 'running' } });
      expect(state.activity).toBe('Running JavaScript');
      observe({ kind: 'execution', execution });
      expect(state.executions).toEqual([execution]);
      return { content: 'Done', executions: [execution] };
    });
    expect(state.messages[1].executions).toEqual([execution]);
    expect(state.executions).toEqual([]);
  });

  it('retains completed commands after an interrupted turn', async () => {
    const state = initial();
    await sendMessage(state, async () => ({ content: '', executions: [execution], error: 'Offline' }));
    expect(state.messages).toHaveLength(2);
    expect(state.messages[1].executions).toEqual([execution]);
    expect(state.error).toBe('Offline');
  });

  it('rejects empty input and empty responses', async () => {
    const run = vi.fn();
    expect(await sendMessage(initial(' \n'), run)).toBe(false);
    expect(run).not.toHaveBeenCalled();
    const state = initial();
    await sendMessage(state, async () => answer(''));
    expect(state.draft).toBe('Hello');
    expect(state.error).toContain('empty response');
  });
});
