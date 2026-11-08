import { expect, it, vi } from 'vitest';
import { call } from './bridge';
import { runAgent } from './agent';
import { $, all, press, useApp, view, wait, write } from './files-fixture';

vi.mock('./agent', () => ({ runAgent: vi.fn(), agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn(), setTitleModel: vi.fn(), steerAgent: vi.fn(), stopAgent: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));
vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => false, listen: vi.fn(async () => () => {}) }));
vi.mock('./sound', () => ({ chime: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));

useApp();

it('traverses folders with Tab and accepts a child without sending', async () => {
  const base = vi.mocked(call).getMockImplementation()!;
  vi.mocked(call).mockImplementation(async (method, ...args) => {
    if (method === 'WorkspaceService.Suggest' && args[2] === 'src/') return [{ name: 'main.go', path: '/w/src/main.go', kind: 'file', size: 1 }];
    return base(method, ...args);
  });
  write('@s'); await wait(); press('Tab'); await wait();
  expect(view().state.doc.toString()).toBe('@src/');
  expect($('[role="listbox"]')).not.toBeNull();
  expect(all('[role="option"]')[0].textContent).toBe('main.go');
  press('Tab'); await wait();
  expect($('.cm-ref-token')?.textContent).toBe('@main.go');
  expect($('[role="listbox"]')).toBeNull();
  expect(runAgent).not.toHaveBeenCalled();
});

it('keeps more than five results reachable with arrow keys', async () => {
  const base = vi.mocked(call).getMockImplementation()!;
  vi.mocked(call).mockImplementation(async (method, ...args) => {
    if (method === 'WorkspaceService.Suggest') return Array.from({ length: 12 }, (_, i) => ({ name: 'f' + i, path: '/w/f' + i, kind: 'file', size: 1 }));
    return base(method, ...args);
  });
  write('@'); await wait();
  expect(all('[role="option"]')).toHaveLength(12);
  for (let i = 0; i < 7; i++) press('ArrowDown');
  await wait();
  expect(all('[role="option"]')[7].getAttribute('aria-selected')).toBe('true');
  press('Tab'); await wait();
  expect($('.cm-ref-token')?.textContent).toBe('@f7');
});
