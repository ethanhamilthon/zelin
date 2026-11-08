import { expect, it, vi } from 'vitest';
import { tick } from 'svelte';
import { runAgent } from './agent';
import { call } from './bridge';
import { $, all, button, calls, key, press, useApp, view, wait, write } from './files-fixture';

vi.mock('./agent', () => ({ runAgent: vi.fn(), agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn(), setTitleModel: vi.fn(), steerAgent: vi.fn(), stopAgent: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));
vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => false, listen: vi.fn(async () => () => {}) }));
vi.mock('./sound', () => ({ chime: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));

useApp();

it('opens the file browser with Ctrl+F and shows the right two actions per row', async () => {
  key('f', 'KeyF', { ctrlKey: true }); await wait();
  expect($('[role="dialog"]')?.getAttribute('aria-label')).toBe('Files');
  const rows = all('.row').map((row) => [row.querySelector('.name')?.textContent?.trim(), ...[...row.querySelectorAll('.acts button')].map((b) => b.textContent)]);
  expect(rows[0]).toEqual(['▸ src', 'Open folder', 'Tag']);
  expect(rows[1]).toEqual(['notes.md', 'Open in app', 'Tag']);
  expect(rows[2]).toEqual(['data.csv', 'Open in default app', 'Tag']);
  expect(rows[3][0]).toContain('→ /w/real.txt');
  expect(calls.filter(([m]) => m === 'WorkspaceService.Browse')).toHaveLength(1);
  key('f', 'KeyF', { ctrlKey: true }); await wait();
  expect($('[role="dialog"]')).toBeNull();
});

it('navigates into a folder and back with breadcrumbs', async () => {
  key('f', 'KeyF', { ctrlKey: true }); await wait();
  button('Open folder').click(); await wait();
  expect($('.crumbs')?.textContent).toContain('src');
  expect($('[role="dialog"]')?.textContent).toContain('This folder is empty');
  $('[aria-label="Go up"]')!.click(); await wait();
  expect(all('.row')).toHaveLength(4);
  key('Escape', 'Escape'); await wait();
  expect($('[role="dialog"]')).toBeNull();
});
