import { expect, it, vi } from 'vitest';
import { createSessionController } from './session-controller.svelte';
import { currentSession, openSession, saveDraft } from './sessions';
import { KEEP } from './retention';

vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));

const record = (id: string) => ({ id, title: id, draft: '', attachments: [], updated: '', messages: [{ role: 'user' as const, content: id }] });

async function begin() {
  vi.mocked(saveDraft).mockResolvedValue();
  vi.mocked(currentSession).mockResolvedValue(record('s0'));
  vi.mocked(openSession).mockImplementation(async (id) => record(id));
  const sessions = createSessionController();
  await sessions.start();
  return sessions;
}
async function visit(sessions: Awaited<ReturnType<typeof begin>>, from: number, to: number) {
  for (let i = from; i < to; i++) await sessions.select('s' + i);
}

it('evicts the oldest idle conversations beyond the limit', async () => {
  const sessions = await begin();
  await visit(sessions, 1, KEEP + 8);
  expect(Object.keys(sessions.conversations).length).toBe(KEEP);
  expect(sessions.conversations['s0']).toBeUndefined();
  expect(sessions.current.messages[0].content).toBe('s' + (KEEP + 7));
});

it('keeps running conversations', async () => {
  const sessions = await begin();
  sessions.conversations['s0'].busy = true;
  await visit(sessions, 1, KEEP + 8);
  expect(sessions.conversations['s0']).toBeDefined();
  expect(sessions.conversations['s1']).toBeUndefined();
});

it('keeps unread and failed conversations', async () => {
  const sessions = await begin();
  await sessions.select('s1'); await sessions.select('s2');
  sessions.conversations['s0'].unread = true;
  sessions.conversations['s1'].error = 'failed';
  await visit(sessions, 3, KEEP + 8);
  expect(sessions.conversations['s0']).toBeDefined();
  expect(sessions.conversations['s1']).toBeDefined();
  expect(sessions.conversations['s2']).toBeUndefined();
});

it('keeps conversations with unsaved drafts and saves the draft when leaving', async () => {
  const sessions = await begin();
  await sessions.select('s1');
  sessions.conversations['s0'].draft = 'unsaved';
  await visit(sessions, 2, KEEP + 8);
  expect(sessions.conversations['s0'].draft).toBe('unsaved');
  expect(saveDraft).not.toHaveBeenCalledWith('s0', 'unsaved');
});

it('flushes a draft before leaving so the conversation can be evicted safely', async () => {
  const sessions = await begin();
  sessions.current.draft = 'typed';
  await visit(sessions, 1, KEEP + 8);
  expect(saveDraft).toHaveBeenCalledWith('s0', 'typed');
  expect(sessions.conversations['s0']).toBeUndefined();
});

it('reloads an evicted session from the backend when reopened', async () => {
  const sessions = await begin();
  await visit(sessions, 1, KEEP + 8);
  await sessions.select('s0');
  expect(sessions.current.messages[0].content).toBe('s0');
});
