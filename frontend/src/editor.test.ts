import { afterEach, describe, expect, it, vi } from 'vitest';
import { undo } from '@codemirror/commands';
import { createEditor } from './editor';

let editor: ReturnType<typeof createEditor> | undefined;
afterEach(() => { editor?.destroy(); document.body.replaceChildren(); });

function setup(submit = vi.fn()) {
  const container = document.createElement('div');
  document.body.append(container);
  editor = createEditor(container, 'Hello', vi.fn(), submit);
  editor.focus();
  return editor;
}

function key(shiftKey: boolean) {
  editor!.view.contentDOM.dispatchEvent(new KeyboardEvent('keydown', {
    key: 'Enter', code: 'Enter', shiftKey, bubbles: true, cancelable: true,
  }));
}

describe('message editor', () => {
  it('submits with Shift+Enter, while Enter inserts a newline', () => {
    const submit = vi.fn();
    const editor = setup(submit);
    key(true);
    expect(submit).toHaveBeenCalledOnce();
    expect(editor.view.state.doc.toString()).toBe('Hello');
    key(false);
    expect(editor.view.state.doc.toString()).toContain('\n');
  });

  it('keeps input editable when submission is ignored', () => {
    const editor = setup();
    key(true);
    expect(editor.view.state.readOnly).toBe(false);
    expect(editor.view.contentDOM.getAttribute('contenteditable')).toBe('true');
    editor.view.dispatch({ changes: { from: 5, insert: ' next' } });
    expect(editor.view.state.doc.toString()).toBe('Hello next');
  });

  it('clears undo history when a message has been submitted', () => {
    const editor = setup();
    editor.view.dispatch({ changes: { from: 5, insert: ' world' } });
    editor.setText('');
    expect(undo(editor.view)).toBe(false);
    expect(editor.view.state.doc.toString()).toBe('');
  });
});
