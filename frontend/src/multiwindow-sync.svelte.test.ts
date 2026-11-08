import { beforeEach, expect, it, vi } from 'vitest';
import { createSessionController } from './session-controller.svelte';
import { currentSession, saveDraft } from './sessions';
import { call } from './bridge';
vi.mock('./bridge', () => ({ call: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));
const record = { id: 'one', title: 'Chat', draft: '', attachments: [], refs: [], updated: '', messages: [] };
let draft = '';
let activity: Record<string, { running: boolean; unread: boolean }>;
beforeEach(() => {
 vi.clearAllMocks(); draft = ''; activity = {};
 vi.mocked(currentSession).mockResolvedValue({ ...record });
 vi.mocked(saveDraft).mockImplementation(async (_id, value) => { draft = value; });
 vi.mocked(call).mockImplementation(async (method, ...args) => {
  if (method === 'AgentService.Activity') return structuredClone(activity);
  if (method === 'AgentService.Read') { activity[String(args[0])].unread = false; return; }
  if (method === 'SessionService.Get') return { ...record, draft };
  if (method === 'AgentService.Active') return null;
  throw new Error(method);
 });
});
it('propagates a saved draft between two independent window controllers', async () => {
 const first = createSessionController(), second = createSessionController();
 await first.start(); await second.start();
 first.current.draft = 'shared input'; await first.flush(); await second.sync();
 expect(second.current.draft).toBe('shared input');
 second.current.draft = 'changed in second'; await second.flush(); await first.sync();
 expect(first.current.draft).toBe('changed in second');
 second.current.draft = 'unsaved'; draft = 'remote'; await second.sync();
 expect(second.current.draft).toBe('unsaved');
});
it('keeps typed text when a save finishes while a stale remote draft is loading', async () => {
 const window = createSessionController();
 await window.start();
 let release!: () => void;
 vi.mocked(call).mockImplementation(async (method) => {
  if (method === 'AgentService.Activity') return {};
  if (method === 'SessionService.Get') { const stale = { ...record, draft: 'hel' }; await new Promise<void>((resolve) => { release = resolve; }); return stale; }
  if (method === 'AgentService.Active') return null;
  throw new Error(method);
 });
 window.current.draft = 'hello';
 const syncing = window.sync();
 await new Promise((resolve) => setTimeout(resolve, 0));
 await window.flush();
 release();
 await syncing;
 expect(window.current.draft).toBe('hello');
});
it('counts backend runs and unread completions without opening their sessions', async () => {
 const first = createSessionController(), second = createSessionController();
 await first.start(); await second.start();
 activity = { other: { running: true, unread: false }, third: { running: false, unread: true } };
 await first.sync(); await second.sync();
 expect(first.activity).toEqual({ working: 1, done: 1 });
 expect(second.activity).toEqual(first.activity);
 expect(first.status('other')).toBe('running');
 activity.other = { running: false, unread: true };
 await second.sync(); expect(second.activity).toEqual({ working: 0, done: 2 });
});
it('clears unread globally when the completed session is visible', async () => {
 const sessions = createSessionController(); await sessions.start();
 activity.one = { running: false, unread: true };
 await sessions.sync();
 expect(sessions.activity).toEqual({ working: 0, done: 0 });
 expect(call).toHaveBeenCalledWith('AgentService.Read', 'one');
});
