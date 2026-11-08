import { HighlightStyle } from '@codemirror/language';
import { EditorView } from '@codemirror/view';
import { tags } from '@lezer/highlight';

export const theme = EditorView.theme({
  '&': { width: '100%', maxWidth: '100%', height: '100%', backgroundColor: 'var(--background)', color: 'var(--foreground)' },
  '&.cm-focused': { outline: 'none' },
  '.cm-scroller': { overflowX: 'hidden', overflowY: 'auto', fontFamily: 'inherit', fontSize: 'var(--font-size)', lineHeight: '1.8' },
  '.cm-content': { width: '100%', maxWidth: '100%', boxSizing: 'border-box', padding: '32px clamp(20px, 6vw, 80px)', minHeight: '100%', caretColor: 'var(--foreground)', whiteSpace: 'pre-wrap', overflowWrap: 'anywhere' },
  '.cm-line': { padding: '0' },
  '.cm-cursor, .cm-dropCursor': { borderLeftColor: 'var(--foreground)' },
  '.cm-selectionBackground, &.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground': { background: 'var(--selection)' },
  '.cm-content ::selection': { backgroundColor: 'var(--selection)' },
  '.cm-placeholder': { color: 'var(--faint)' },
  '.cm-image-token': { color: 'var(--accent)', cursor: 'pointer' },
  '.cm-ref-token': { color: 'var(--accent)', cursor: 'pointer', background: 'var(--code-bg)', borderRadius: '4px', padding: '0 4px' },
  '.cm-scroller::-webkit-scrollbar': { width: '8px', height: '8px' },
  '.cm-scroller::-webkit-scrollbar-thumb': { backgroundColor: 'var(--border)', border: '2px solid var(--background)', borderRadius: '8px' },
});

export const highlighting = HighlightStyle.define([
  { tag: tags.heading, color: 'var(--foreground)', fontWeight: '600' },
  { tag: tags.strong, color: 'var(--foreground)', fontWeight: '700' },
  { tag: tags.emphasis, color: 'var(--foreground)', fontStyle: 'italic' },
  { tag: [tags.link, tags.url], color: 'var(--muted)' },
  { tag: tags.monospace, color: 'var(--muted)' },
  { tag: tags.quote, color: 'var(--dim)' },
  { tag: [tags.processingInstruction, tags.meta], color: 'var(--faint)' },
]);
