import { defaultKeymap, history, historyKeymap, indentWithTab, insertNewlineAndIndent } from '@codemirror/commands';
import { insertNewlineContinueMarkup, markdown } from '@codemirror/lang-markdown';
import { syntaxHighlighting } from '@codemirror/language';
import { Compartment, EditorState } from '@codemirror/state';
import type { Extension } from '@codemirror/state';
import { drawSelection, dropCursor, EditorView, keymap, placeholder } from '@codemirror/view';
import { theme, highlighting } from './editor-theme';
import { codeMirrorKey } from './hotkeys';
import { imageInput, imageTokens } from './image-tokens';

export type EditorKeys = { send: string; newline: string };
export type EditorOptions = {
  label?: string; placeholder?: string; keys?: EditorKeys; onSubmit?: () => void;
  onImage?: (number: number) => void; onAttach?: (url: string) => number;
};

export function createEditor(parent: HTMLElement, doc: string, onChange: (text: string) => void, options: EditorOptions = {}) {
  const keys = new Compartment();
  const hint = new Compartment();
  let current = options.keys;
  let text = options.placeholder ?? 'Write a message…';

  function bindings(value?: EditorKeys) {
    if (!value || !options.onSubmit) return keymap.of([]);
    const submit = options.onSubmit;
    return keymap.of([
      { key: codeMirrorKey(value.send), run: () => { submit(); return true; } },
      { key: codeMirrorKey(value.newline), run: (view) => insertNewlineContinueMarkup(view) || insertNewlineAndIndent(view) },
    ]);
  }

  const extensions = (): Extension[] => [
    history(), drawSelection(), dropCursor(), markdown(),
    syntaxHighlighting(highlighting), theme, EditorView.lineWrapping,
    EditorView.contentAttributes.of({ 'aria-label': options.label ?? 'Message', spellcheck: 'false' }),
    hint.of(placeholder(text)),
    keys.of(bindings(current)),
    keymap.of([...defaultKeymap, ...historyKeymap, indentWithTab]),
    options.onImage ? imageTokens(options.onImage) : [],
    options.onAttach ? imageInput(options.onAttach) : [],
    EditorView.updateListener.of((update) => {
      if (update.docChanged) onChange(update.state.doc.toString());
    }),
  ];
  const view = new EditorView({ parent, state: EditorState.create({ doc, extensions: extensions() }) });
  return {
    view,
    focus: () => view.focus(),
    destroy: () => view.destroy(),
    text: () => view.state.doc.toString(),
    selected: () => !view.state.selection.main.empty,
    insert(text: string) {
      const { from, to } = view.state.selection.main;
      view.dispatch({ changes: { from, to, insert: text }, selection: { anchor: from + text.length } });
    },
    setPlaceholder(value: string) {
      text = value;
      view.dispatch({ effects: hint.reconfigure(placeholder(value)) });
    },
    setKeys(value: EditorKeys) {
      current = value;
      view.dispatch({ effects: keys.reconfigure(bindings(value)) });
    },
    setText(text: string) {
      view.setState(EditorState.create({ doc: text, extensions: extensions() }));
      onChange(text);
    },
  };
}

export type Editor = ReturnType<typeof createEditor>;
