import { expect, it, vi } from 'vitest';
import { zoomed } from './zoom';
import { createKeyboard } from './app-keyboard';
import { defaultHotkeys } from './hotkeys';
import { observeEvent } from './observe';
import type { Conversation } from './conversation';

it('clamps chat zoom', () => {
  expect(zoomed(16, 1)).toBe(17);
  expect(zoomed(26, 1)).toBe(26);
  expect(zoomed(12, -1)).toBe(12);
  expect(zoomed(0, -1)).toBe(15);
});

it('zooms the chat with Ctrl+= and Ctrl+-', () => {
  const zoom = vi.fn();
  const keyboard = createKeyboard({ keys: () => defaultHotkeys, blocked: () => false, doc: () => false, hasFiles: () => false,
    stepBack: vi.fn(), commit: vi.fn(), toggle: vi.fn(), files: vi.fn(), create: vi.fn(), stop: () => false, zoom });
  keyboard(new KeyboardEvent('keydown', { key: '=', code: 'Equal', ctrlKey: true }));
  keyboard(new KeyboardEvent('keydown', { key: '-', code: 'Minus', ctrlKey: true }));
  expect(zoom.mock.calls).toEqual([[1], [-1]]);
});

it('tracks streamed answer text', () => {
  const state = { streaming: '' } as Conversation;
  observeEvent(state, { kind: 'stream', message: { content: 'Hel' } });
  expect(state.streaming).toBe('Hel');
});
