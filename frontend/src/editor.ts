import { defaultKeymap, history, historyKeymap, indentWithTab } from '@codemirror/commands';
import { markdown } from '@codemirror/lang-markdown';
import { syntaxHighlighting } from '@codemirror/language';
import { EditorState } from '@codemirror/state';
import { drawSelection, dropCursor, EditorView, keymap, placeholder } from '@codemirror/view';
import { theme, highlighting } from './editor-theme';

export function createEditor(
  parent: HTMLElement, doc: string, onChange: (text: string) => void,
  onSubmit: () => void = () => {},
) {
  const extensions = [
    history(), drawSelection(), dropCursor(), markdown(),
    syntaxHighlighting(highlighting), theme, EditorView.lineWrapping,
    EditorView.contentAttributes.of({ 'aria-label': 'Message', spellcheck: 'false' }),
    placeholder('Write a message…'),
    keymap.of([
      { key: 'Shift-Enter', run: () => { onSubmit(); return true; } },
      ...defaultKeymap, ...historyKeymap, indentWithTab,
    ]),
    EditorView.updateListener.of((update) => {
      if (update.docChanged) onChange(update.state.doc.toString());
    }),
  ];
  const view = new EditorView({ parent, state: EditorState.create({ doc, extensions }) });
  return {
    view,
    focus: () => view.focus(),
    destroy: () => view.destroy(),
    setText(text: string) {
      view.setState(EditorState.create({ doc: text, extensions }));
      onChange(text);
    },
  };
}
