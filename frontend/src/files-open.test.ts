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

it('opens Markdown in the in-app editor and other files through the default app', async () => {
  key('f', 'KeyF', { ctrlKey: true }); await wait();
  button('Open in default app').click(); await wait();
  expect(call).toHaveBeenCalledWith('WorkspaceService.Open', 'default', '/w/data.csv');
  button('Open in app').click(); await wait(); await wait();
  expect($('.cm-content[aria-label="/w/notes.md"]')?.textContent).toBe('# notes');
});

it('tags a file at the cursor without submitting and returns focus to the input', async () => {
  write('look at ');
  key('f', 'KeyF', { ctrlKey: true }); await wait();
  all('.row')[1].querySelectorAll('button')[1].click(); await wait();
  expect($('[role="dialog"]')).toBeNull();
  expect(view().state.doc.toString()).toMatch(/^look at \[@[0-9a-f]{12}\|notes\.md\] $/);
  expect($('.cm-ref-token')?.textContent).toBe('@notes.md');
  expect(document.activeElement).toBe($('.cm-content'));
  expect(runAgent).not.toHaveBeenCalled();
});

it('offers suggestions after @, accepts with Enter without sending', async () => {
  write('@'); await wait();
  expect($('[role="listbox"]')).not.toBeNull();
  expect(all('[role="option"]').length).toBeLessThanOrEqual(5);
  const down = new KeyboardEvent('keydown', { key: 'ArrowDown', bubbles: true, cancelable: true });
  $('.cm-content')!.dispatchEvent(down); await tick();
  expect(all('[role="option"]')[1].getAttribute('aria-selected')).toBe('true');
  $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true, cancelable: true })); await wait();
  expect(runAgent).not.toHaveBeenCalled();
  expect(view().state.doc.toString()).toMatch(/^\[@[0-9a-f]{12}\|notes\.md\] $/);
  expect($('[role="listbox"]')).toBeNull();
});

it('dismisses the menu with Escape and never hijacks literal @ in text', async () => {
  write('mail me@example'); await wait();
  expect($('[role="listbox"]')).toBeNull();
  write('hi @'); await wait();
  expect($('[role="listbox"]')).not.toBeNull();
  $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true })); await tick();
  expect($('[role="listbox"]')).toBeNull();
  expect(view().state.doc.toString()).toBe('hi @');
});

it('treats a token as one unit for the cursor and Backspace', async () => {
  write(''); key('f', 'KeyF', { ctrlKey: true }); await wait();
  all('.row')[0].querySelectorAll('button')[1].click(); await wait();
  expect($('.cm-ref-token')?.textContent).toBe('@src');
  const full = view().state.doc.toString();
  expect(full).toMatch(/^\[@[0-9a-f]{12}\|src\] $/);
  view().dispatch({ selection: { anchor: full.length - 1 } });
  $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Backspace', bubbles: true, cancelable: true }));
  expect(view().state.doc.toString()).toBe(' ');
  expect($('.cm-ref-token')).toBeNull();
});

it('sends structured references and blocks unknown ones', async () => {
  write('check '); key('f', 'KeyF', { ctrlKey: true }); await wait();
  all('.row')[2].querySelectorAll('button')[1].click(); await wait();
  $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true, cancelable: true })); await wait();
  expect(runAgent).toHaveBeenCalledTimes(1);
  const [messages, , , refs] = vi.mocked(runAgent).mock.calls[0];
  expect(messages[0].content).toMatch(/^check \[@[0-9a-f]{12}\|data\.csv\]/);
  expect(refs).toEqual([expect.objectContaining({ kind: 'file', path: '/w/data.csv', label: 'data.csv', workspaceId: 'default', sessionId: 'one' })]);
  expect($('article .token')?.textContent).toBe('@data.csv');
});

it('rejects a draft with an unknown reference and keeps the text', async () => {
  write('see [@deadbeef0000|ghost.txt]'); await wait();
  $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true, cancelable: true })); await wait();
  expect(runAgent).not.toHaveBeenCalled();
  expect($('[role="alert"]')?.textContent).toContain('Unknown reference: @ghost.txt');
  expect(view().state.doc.toString()).toContain('ghost.txt');
});

it('opens Markdown tokens in-app and other tokens in the default application', async () => {
  write('@'); await wait();
  $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'ArrowDown', bubbles: true, cancelable: true })); await tick();
  $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true, cancelable: true })); await wait();
  $('.cm-ref-token')!.dispatchEvent(new MouseEvent('mousedown', { bubbles: true, cancelable: true })); await wait(); await wait();
  expect($('.cm-content[aria-label="/w/notes.md"]')).not.toBeNull();
});
