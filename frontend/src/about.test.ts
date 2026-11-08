import { expect, it, vi } from 'vitest';
import { tick } from 'svelte';
import { call } from './bridge';
import { $, all, button, key, useApp, wait } from './files-fixture';
import { notices } from './about';

vi.mock('./agent', () => ({ runAgent: vi.fn(), agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn(), setTitleModel: vi.fn(), steerAgent: vi.fn(), stopAgent: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));
vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => false, listen: vi.fn(async () => () => {}) }));
vi.mock('./sound', () => ({ chime: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
useApp();

const info = { name: 'Zelin', version: '0.1.0', commit: 'abc123def456', date: '2026-11-08', description: 'Zen AI agent', runtime: 'go1.26', platform: 'linux/arm64', profile: '/data/zelin', daemon: 'Not used' };
const open = async () => { key('k', 'KeyK', { ctrlKey: true }); await tick(); button('About Version and licenses').click(); await wait(); };

it('opens About from Commands and shows only safe build facts', async () => {
  const base = vi.mocked(call).getMockImplementation()!;
  vi.mocked(call).mockImplementation(async (method, ...args) => method === 'AboutService.Info' ? info : base(method, ...args));
  await open();
  const text = $('[role="dialog"]')!.textContent!;
  expect(text).toContain('Zelin');
  for (const part of ['0.1.0', 'abc123def456', '2026-11-08', 'go1.26', 'linux/arm64', '/data/zelin', 'Not used']) expect(text).toContain(part);
  expect(text).not.toMatch(/api.?key|secret|token/i);
  expect(all('.body li')).toHaveLength(notices.length);
  expect(all('.body a').every((a) => a.dataset.url?.startsWith('https://'))).toBe(true);
});

it('returns to Commands with Escape and shows a load error without crashing', async () => {
  const base = vi.mocked(call).getMockImplementation()!;
  vi.mocked(call).mockImplementation(async (method, ...args) => { if (method === 'AboutService.Info') throw new Error('unavailable'); return base(method, ...args); });
  await open();
  expect($('[role="alert"]')?.textContent).toBe('unavailable');
  key('Escape', 'Escape'); await wait();
  expect($('[role="dialog"]')?.getAttribute('aria-label')).toBe('Commands');
});

it('lists the open source dependencies with verified licenses', () => {
  expect(notices.map((item) => item.name)).toEqual(expect.arrayContaining(['Wails', 'Goja', 'Svelte', 'marked', 'DOMPurify']));
  expect(notices.find((item) => item.name === 'DOMPurify')?.license).toBe('MPL-2.0 OR Apache-2.0');
});
