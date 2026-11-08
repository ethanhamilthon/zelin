import { afterEach, expect, it, vi } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import SessionList from './SessionList.svelte';

vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
let app: ReturnType<typeof mount> | undefined;
afterEach(async () => { if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });

const make = (count: number) => Array.from({ length: count }, (_, i) => ({ id: String(i), title: `Session ${i}`, messages: [], draft: '', updated: '2026-01-01T10:00:00Z' }));
const rows = () => document.querySelectorAll('.body > button[aria-current], .body > button').length;
function show(count: number) {
  app = mount(SessionList, { target: document.body, props: { sessions: make(count), active: '', status: () => '' as const, onclose: vi.fn(), onselect: vi.fn(), onnew: vi.fn() } });
  flushSync();
}

it('mounts a bounded page for 10,000 sessions and shows more on request', async () => {
  const began = performance.now();
  show(10000);
  expect(performance.now() - began).toBeLessThan(1500);
  expect(rows()).toBe(201);
  const more = [...document.querySelectorAll('button')].find((item) => item.textContent?.includes('Show more'))!;
  expect(more.textContent).toContain('9800 left');
  more.click(); await tick();
  expect(rows()).toBe(401);
});

it('searches across all sessions, not only the mounted page', async () => {
  show(10000);
  const input = document.querySelector<HTMLInputElement>('input')!;
  input.value = 'Session 9999'; input.dispatchEvent(new Event('input', { bubbles: true })); await tick();
  expect(document.body.textContent).toContain('Session 9999');
  expect(document.body.textContent).not.toContain('Show more');
});

it('shows three sessions per workspace and expands on request', async () => {
  app = mount(SessionList, { target: document.body, props: {
    sessions: make(8).map((session) => ({ ...session, workspaceId: 'default' })),
    workspaces: [{ id: 'default', name: 'Default', icon: '', dir: '/w' }],
    active: '', status: () => '' as const, onclose: vi.fn(), onselect: vi.fn(), onnew: vi.fn(),
  } });
  flushSync();
  expect(document.querySelectorAll('.body > button[aria-current]')).toHaveLength(0);
  expect(document.body.textContent).toContain('Session 2');
  expect(document.body.textContent).not.toContain('Session 3');
  const more = [...document.querySelectorAll('button')].find((item) => item.textContent?.includes('Show More'))!;
  more.click(); await tick();
  expect(document.body.textContent).toContain('Session 7');
  expect(document.body.textContent).not.toContain('Show More');
});
