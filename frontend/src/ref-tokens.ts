import { Decoration, EditorView, MatchDecorator, ViewPlugin, WidgetType } from '@codemirror/view';
import type { DecorationSet, ViewUpdate } from '@codemirror/view';
import { refPattern } from './refs';

class Chip extends WidgetType {
  constructor(readonly id: string, readonly label: string) { super(); }
  eq(other: Chip) { return other.id === this.id && other.label === this.label; }
  toDOM() {
    const chip = document.createElement('span');
    chip.className = 'cm-ref-token';
    chip.dataset.ref = this.id;
    chip.textContent = '@' + (this.label || this.id);
    return chip;
  }
  ignoreEvent() { return false; }
}

const matcher = new MatchDecorator({
  regexp: new RegExp(refPattern.source, 'g'),
  decoration: (match) => Decoration.replace({ widget: new Chip(match[1], match[2] ?? '') }),
});

export function refTokens(onOpen: (id: string) => void) {
  return ViewPlugin.fromClass(class {
    decorations: DecorationSet;
    constructor(view: EditorView) { this.decorations = matcher.createDeco(view); }
    update(update: ViewUpdate) { this.decorations = matcher.updateDeco(update, this.decorations); }
  }, {
    decorations: (plugin) => plugin.decorations,
    provide: (plugin) => EditorView.atomicRanges.of((view) => view.plugin(plugin)?.decorations ?? Decoration.none),
    eventHandlers: {
      mousedown(event) {
        const chip = (event.target as HTMLElement).closest<HTMLElement>('.cm-ref-token');
        if (!chip?.dataset.ref) return false;
        event.preventDefault();
        onOpen(chip.dataset.ref);
        return true;
      },
    },
  });
}
