import { expect, it } from 'vitest';
import { blank } from './session-controller.svelte';
import { observeEvent } from './observe';
import { syncSession } from './sync-sessions';
import { localRuns } from './remote-session';
import { vi } from 'vitest';

vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => false, listen: vi.fn() }));

it('moves finished steps before the steering message', () => {
  const state = { ...blank(), busy: true, steering: true, draft: 'turn left' };
  const step = (id: string, done: string) => ({ id, code: id, output: '', state: done as 'complete' | 'running' });
  state.executions = [step('a', 'complete'), step('b', 'complete'), step('c', 'running')];
  observeEvent(state, { kind: 'steer', message: { role: 'user', content: 'turn left' } } as never);
  expect(state.messages).toHaveLength(1);
  expect(state.messages[0].executions?.map((item) => item.id)).toEqual(['a', 'b']);
  expect(state.executions.map((item) => item.id)).toEqual(['c']);
});

it('keeps text typed while a remote refresh is in flight', async () => {
  const { call } = await import('./bridge');
  const target = { ...blank(), draft: '' };
  const conversations = { one: target };
  vi.mocked(call).mockImplementation(async (method: string) => {
    if (method === 'SessionService.Get') { target.draft = 'typed'; return { id: 'one', title: '', messages: [], draft: '', updated: '' }; }
    return null;
  });
  localRuns.clear();
  await syncSession('one', conversations, new Map([['one', { draft: 'typed', attachments: [], refs: '[]' }]]));
  expect(target.draft).toBe('typed');
});

it('keeps the answer given right before a steer', () => {
  const state = { ...blank(), busy: true, steering: true, draft: 'more' };
  observeEvent(state, { kind: 'execution', execution: { id: 't-note', kind: 'message', code: 'First answer', output: '', state: 'complete' } });
  observeEvent(state, { kind: 'steer', message: { role: 'user', content: 'more' } } as never);
  expect(state.messages[0].executions?.[0].code).toBe('First answer');
  expect(state.messages[0].content).toBe('more');
});
