import { afterEach, expect, it, vi } from 'vitest';
import { createEditor } from './editor';

let editor: ReturnType<typeof createEditor>;
afterEach(() => { editor?.destroy(); document.body.replaceChildren(); });

it.each([
  ['- First', '- First\n- '],
  ['1. First', '1. First\n2. '],
  ['  - Nested', '  - Nested\n  - '],
  ['- [ ] Task', '- [ ] Task\n- [ ] '],
  ['- ', ''],
  ['Plain', 'Plain\n'],
])('continues markup with Shift+Enter: %s', (doc, expected) => {
  const parent = document.createElement('div');
  document.body.append(parent);
  const onSubmit = vi.fn();
  editor = createEditor(parent, doc, vi.fn(), { keys: { send: 'Enter', newline: 'Shift+Enter' }, onSubmit });
  editor.view.dispatch({ selection: { anchor: doc.length } });
  editor.focus();
  editor.view.contentDOM.dispatchEvent(new KeyboardEvent('keydown', {
    key: 'Enter', code: 'Enter', shiftKey: true, bubbles: true, cancelable: true,
  }));
  expect(editor.text()).toBe(expected);
  expect(onSubmit).not.toHaveBeenCalled();
});
