import { afterEach, describe, expect, it, vi } from 'vitest';
import { deleteCharBackward, undo } from '@codemirror/commands';
import { createEditor } from './editor';
import type { EditorOptions } from './editor';

let editor: ReturnType<typeof createEditor> | undefined;
afterEach(() => { editor?.destroy(); document.body.replaceChildren(); });

function setup(options: EditorOptions = {}, doc = 'Hello') {
  const container = document.createElement('div');
  document.body.append(container);
  editor = createEditor(container, doc, vi.fn(), { keys: { send: 'Enter', newline: 'Shift+Enter' }, onSubmit: vi.fn(), ...options });
  editor.focus();
  return editor;
}

function key(shiftKey: boolean, ctrlKey = false) {
  editor!.view.contentDOM.dispatchEvent(new KeyboardEvent('keydown', {
    key: 'Enter', code: 'Enter', shiftKey, ctrlKey, bubbles: true, cancelable: true,
  }));
}

describe('message editor', () => {
  it('submits with Enter, while Shift+Enter inserts a newline', () => {
    const onSubmit = vi.fn();
    const editor = setup({ onSubmit });
    key(false);
    expect(onSubmit).toHaveBeenCalledOnce();
    expect(editor.text()).toBe('Hello');
    key(true);
    expect(editor.text()).toContain('\n');
  });

  it('applies changed send keys', () => {
    const onSubmit = vi.fn();
    const editor = setup({ onSubmit });
    editor.setKeys({ send: 'Ctrl+Enter', newline: 'Enter' });
    key(false);
    expect(onSubmit).not.toHaveBeenCalled();
    key(false, true);
    expect(onSubmit).toHaveBeenCalledOnce();
  });

  it('keeps input editable when submission is ignored', () => {
    const editor = setup();
    key(false);
    expect(editor.view.state.readOnly).toBe(false);
    expect(editor.view.contentDOM.getAttribute('contenteditable')).toBe('true');
    editor.view.dispatch({ changes: { from: 5, insert: ' next' } });
    expect(editor.text()).toBe('Hello next');
  });

  it('clears undo history when a message has been submitted', () => {
    const editor = setup();
    editor.view.dispatch({ changes: { from: 5, insert: ' world' } });
    editor.setText('');
    expect(undo(editor.view)).toBe(false);
    expect(editor.text()).toBe('');
  });

  it('treats image tokens as one unit and opens them on click', () => {
    const onImage = vi.fn();
    const editor = setup({ onImage, onAttach: () => 1 }, 'See [Image #12]');
    editor.view.dispatch({ selection: { anchor: editor.text().length } });
    const token = editor.view.contentDOM.querySelector('.cm-image-token');
    expect(token?.textContent).toBe('[Image #12]');
    token!.dispatchEvent(new MouseEvent('mousedown', { bubbles: true, cancelable: true }));
    expect(onImage).toHaveBeenCalledWith(12);
    deleteCharBackward(editor.view);
    expect(editor.text()).toBe('See ');
  });
});
