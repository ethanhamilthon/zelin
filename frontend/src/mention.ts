import { EditorView } from '@codemirror/view';
import type { Extension } from '@codemirror/state';

export type Mention = { from: number; to: number; query: string; x: number; y: number } | null;

const trigger = /(?:^|[\s(])@([^\s@\[\]]*)$/;

export function findMention(before: string) {
  const match = trigger.exec(before);
  return match ? { query: match[1], length: match[1].length + 1 } : null;
}

export function mentionWatcher(onChange: (mention: Mention) => void): Extension {
  return EditorView.updateListener.of((update) => {
    if (!update.docChanged && !update.selectionSet) return;
    const { head, empty } = update.state.selection.main;
    if (!empty) return onChange(null);
    const line = update.state.doc.lineAt(head);
    const found = findMention(line.text.slice(0, head - line.from));
    if (!found) return onChange(null);
    const box = update.view.coordsAtPos(head);
    onChange({ from: head - found.length, to: head, query: found.query, x: box?.left ?? 0, y: box?.bottom ?? 0 });
  });
}
