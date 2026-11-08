import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { createSessionController } from './session-controller.svelte';
import { currentSession, openSession, saveAttachments, saveDraft, saveSession } from './sessions';

vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));

const record = (id: string, draft = '') => ({ id, title: id, draft, attachments: [], updated: '', messages: [{ role: 'user' as const, content: 'hi' }] });

async function started() {
  vi.mocked(currentSession).mockResolvedValue(record('one'));
  const sessions = createSessionController();
  await sessions.start();
  return sessions;
}

beforeEach(() => {
  vi.useFakeTimers();
  vi.resetAllMocks();
  vi.mocked(saveDraft).mockResolvedValue();
  vi.mocked(saveAttachments).mockResolvedValue();
  vi.mocked(saveSession).mockResolvedValue();
});
afterEach(() => { vi.useRealTimers(); });

it('saves only the draft after a typing pause', async () => {
  const sessions = await started();
  sessions.current.draft = 'abc'; sessions.schedule();
  await vi.advanceTimersByTimeAsync(300);
  expect(saveDraft).toHaveBeenCalledExactlyOnceWith('one', 'abc');
  expect(saveSession).not.toHaveBeenCalled();
  expect(saveAttachments).not.toHaveBeenCalled();
});

it('coalesces rapid edits into one write with the latest text', async () => {
  const sessions = await started();
  for (const text of ['a', 'ab', 'abc']) { sessions.current.draft = text; sessions.schedule(); await vi.advanceTimersByTimeAsync(100); }
  await vi.advanceTimersByTimeAsync(300);
  expect(saveDraft).toHaveBeenCalledExactlyOnceWith('one', 'abc');
});

it('skips writes when nothing changed', async () => {
  const sessions = await started();
  await sessions.flush();
  expect(saveDraft).not.toHaveBeenCalled();
  expect(saveAttachments).not.toHaveBeenCalled();
});

it('saves attachments separately from the draft', async () => {
  const sessions = await started();
  sessions.attach('data:image/png;base64,AA');
  await vi.advanceTimersByTimeAsync(300);
  expect(saveAttachments).toHaveBeenCalledExactlyOnceWith('one', ['data:image/png;base64,AA']);
  expect(saveDraft).not.toHaveBeenCalled();
});

it('keeps write order and retries after a failure', async () => {
  const sessions = await started();
  vi.mocked(saveDraft).mockRejectedValueOnce(new Error('disk'));
  sessions.current.draft = 'one'; await expect(sessions.flush()).rejects.toThrow('disk');
  expect(sessions.state.error).toContain('disk');
  sessions.current.draft = 'two'; await sessions.flush();
  expect(vi.mocked(saveDraft).mock.calls.map((call) => call[1])).toEqual(['one', 'two']);
  expect(sessions.state.error).toBe('');
});

it('flushes the pending draft before switching sessions', async () => {
  const sessions = await started();
  vi.mocked(openSession).mockResolvedValue(record('two'));
  sessions.current.draft = 'unsaved'; sessions.schedule();
  await sessions.select('two');
  expect(saveDraft).toHaveBeenCalledExactlyOnceWith('one', 'unsaved');
  expect(vi.mocked(saveDraft).mock.invocationCallOrder[0]).toBeLessThan(vi.mocked(openSession).mock.invocationCallOrder[0]);
  await vi.advanceTimersByTimeAsync(300);
  expect(saveDraft).toHaveBeenCalledOnce();
});

it('full save records the draft so it is not written twice', async () => {
  const sessions = await started();
  sessions.current.draft = 'sent'; await sessions.save();
  await sessions.flush();
  expect(saveSession).toHaveBeenCalledOnce();
  expect(saveDraft).not.toHaveBeenCalled();
});
