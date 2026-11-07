import { Decoration, EditorView, MatchDecorator, ViewPlugin } from '@codemirror/view';
import type { DecorationSet, ViewUpdate } from '@codemirror/view';
import { imageFiles, readImage, token } from './images';

const matcher = new MatchDecorator({ regexp: /\[Image #\d+\]/g, decoration: Decoration.mark({ class: 'cm-image-token' }) });

export function imageTokens(onOpen: (number: number) => void) {
  return ViewPlugin.fromClass(class {
    decorations: DecorationSet;
    constructor(view: EditorView) { this.decorations = matcher.createDeco(view); }
    update(update: ViewUpdate) { this.decorations = matcher.updateDeco(update, this.decorations); }
  }, {
    decorations: (plugin) => plugin.decorations,
    provide: (plugin) => EditorView.atomicRanges.of((view) => view.plugin(plugin)?.decorations ?? Decoration.none),
    eventHandlers: {
      mousedown(event) {
        const element = (event.target as HTMLElement).closest('.cm-image-token');
        if (!element) return false;
        event.preventDefault();
        onOpen(Number(element.textContent?.match(/\d+/)?.[0]));
        return true;
      },
    },
  });
}

export function imageInput(onAttach: (url: string) => number) {
  async function attach(view: EditorView, files: File[], at?: number) {
    for (const file of files) {
      const text = token(onAttach(await readImage(file)));
      const from = at ?? view.state.selection.main.from;
      const to = at ?? view.state.selection.main.to;
      view.dispatch({ changes: { from, to, insert: text }, selection: { anchor: from + text.length } });
      at = at === undefined ? undefined : from + text.length;
    }
  }
  return EditorView.domEventHandlers({
    paste(event, view) {
      const files = imageFiles(event.clipboardData?.files);
      if (!files.length) return false;
      event.preventDefault();
      void attach(view, files);
      return true;
    },
    drop(event, view) {
      const files = imageFiles(event.dataTransfer?.files);
      if (!files.length) return false;
      event.preventDefault();
      void attach(view, files, view.posAtCoords({ x: event.clientX, y: event.clientY }) ?? undefined);
      return true;
    },
  });
}
