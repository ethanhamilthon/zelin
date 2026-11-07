import { HighlightStyle } from '@codemirror/language';
import { EditorView } from '@codemirror/view';
import { tags } from '@lezer/highlight';

export const theme = EditorView.theme({
  '&': { height: '100%', backgroundColor: 'var(--background)', color: '#e0e0e0' },
  '&.cm-focused': { outline: 'none' },
  '.cm-scroller': { overflow: 'auto', fontFamily: 'inherit', fontSize: 'var(--font-size)', lineHeight: '1.8' },
  '.cm-content': { padding: '32px clamp(48px, 6vw, 80px)', minHeight: '100%', caretColor: '#eee', overflowWrap: 'anywhere' },
  '.cm-line': { padding: '0' },
  '.cm-cursor, .cm-dropCursor': { borderLeftColor: '#eee' },
  '&.cm-focused .cm-selectionBackground, .cm-selectionBackground': { backgroundColor: '#3b3b3b' },
  '.cm-content ::selection': { backgroundColor: '#3b3b3b' },
  '.cm-placeholder': { color: '#666' },
  '.cm-scroller::-webkit-scrollbar': { width: '8px', height: '8px' },
  '.cm-scroller::-webkit-scrollbar-thumb': { backgroundColor: '#404040', border: '2px solid var(--background)', borderRadius: '8px' },
}, { dark: true });

export const highlighting = HighlightStyle.define([
  { tag: tags.heading, color: '#eee', fontWeight: '600' },
  { tag: tags.strong, color: '#eee', fontWeight: '700' },
  { tag: tags.emphasis, color: '#ddd', fontStyle: 'italic' },
  { tag: [tags.link, tags.url], color: '#aaa' },
  { tag: tags.monospace, color: '#aaa' },
  { tag: tags.quote, color: '#888' },
  { tag: [tags.processingInstruction, tags.meta], color: '#666' },
]);
