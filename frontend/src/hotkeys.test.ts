import { describe, expect, it } from 'vitest';
import { codeMirrorKey, combo, matches } from './hotkeys';

const press = (key: string, code: string, init: KeyboardEventInit = {}) => new KeyboardEvent('keydown', { key, code, ...init });

describe('hotkeys', () => {
  it('matches physical keys on any layout', () => {
    expect(matches(press('л', 'KeyK', { ctrlKey: true }), 'Ctrl+K')).toBe(true);
    expect(matches(press('с', 'KeyC', { ctrlKey: true }), 'Ctrl+C')).toBe(true);
    expect(matches(press('k', 'KeyK'), 'Ctrl+K')).toBe(false);
  });

  it('formats combos and ignores lone modifiers', () => {
    expect(combo(press('Enter', 'Enter', { shiftKey: true }))).toBe('Shift+Enter');
    expect(combo(press('1', 'Digit1', { ctrlKey: true, altKey: true }))).toBe('Ctrl+Alt+1');
    expect(combo(press('Control', 'ControlLeft', { ctrlKey: true }))).toBe('');
  });

  it('converts to CodeMirror key names', () => {
    expect(codeMirrorKey('Shift+Enter')).toBe('Shift-Enter');
    expect(codeMirrorKey('Ctrl+K')).toBe('Ctrl-k');
  });
});
