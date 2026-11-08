import { afterEach, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ActivitySummary from './ActivitySummary.svelte';
import { createSessionController } from './session-controller.svelte';
import { currentSession, openSession } from './sessions';

vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));

let app: ReturnType<typeof mount> | undefined;
afterEach(async () => { if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });
const text = () => document.body.textContent!.replace(/\s+/g, ' ').trim();
const show = (working: number, done: number) => { app = mount(ActivitySummary, { target: document.body, props: { working, done } }); flushSync(); };
const record = (id: string) => ({ id, title: id, draft: '', attachments: [], updated: '', messages: [{ role: 'user' as const, content: id }] });

it('renders nothing for zero counts', () => {
  show(0, 0);
  expect(document.body.textContent).toBe('');
  expect(document.querySelector('[role="status"]')).toBeNull();
});

it('shows only the segments with counts and the separator for both', async () => {
  show(2, 0);
  expect(text()).toBe('2 Working');
  await unmount(app!); document.body.replaceChildren();
  show(0, 3);
  expect(text()).toBe('3 Done');
  await unmount(app!); document.body.replaceChildren();
  show(2, 3);
  expect(text()).toBe('2 Working | 3 Done');
  expect(document.querySelector('[role="status"]')?.getAttribute('aria-label')).toBe('2 working, 3 done');
});

it('counts sessions, not steps, and never counts one session twice', async () => {
  vi.mocked(currentSession).mockResolvedValue(record('a'));
  vi.mocked(openSession).mockImplementation(async (id) => record(id));
  const sessions = createSessionController();
  await sessions.start();
  await sessions.select('b'); await sessions.select('c');
  expect(sessions.activity).toEqual({ working: 0, done: 0 });
  sessions.conversations['a'].busy = true;
  sessions.conversations['a'].executions = [{ id: '1', code: '', output: '', state: 'running' }, { id: '2', code: '', output: '', state: 'running' }];
  sessions.conversations['b'].busy = true; sessions.conversations['b'].unread = true;
  expect(sessions.activity).toEqual({ working: 2, done: 0 });
  sessions.conversations['b'].busy = false;
  expect(sessions.activity).toEqual({ working: 1, done: 1 });
  sessions.conversations['b'].error = 'failed';
  expect(sessions.activity).toEqual({ working: 1, done: 0 });
  sessions.conversations['b'].error = '';
  await sessions.select('b');
  expect(sessions.activity).toEqual({ working: 1, done: 0 });
});
