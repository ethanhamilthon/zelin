import { matches } from './hotkeys';
import type { Action } from './hotkeys';
import type { ModalName } from './Panels.svelte';

type Deps = {
  keys: () => Record<Action, string>; blocked: () => boolean; doc: () => boolean; hasFiles: () => boolean;
  stepBack: () => void; commit: () => void; toggle: (name: ModalName) => void; files: () => void; create: () => void; stop: () => unknown; zoom: (step: number) => void;
};

const panels: [Action, ModalName][] = [['commands', 'commands'], ['models', 'models'], ['sessions', 'sessions'], ['actions', 'actions']];

export function createKeyboard(deps: Deps) {
  return (event: KeyboardEvent) => {
    const keys = deps.keys();
    if (event.key === 'Escape' && deps.blocked()) { if (event.defaultPrevented) return; event.preventDefault(); deps.stepBack(); return; }
    if (deps.doc()) {
      if (matches(event, keys.save)) { event.preventDefault(); deps.commit(); }
      return;
    }
    const panel = panels.find(([action]) => matches(event, keys[action]));
    if (panel) deps.toggle(panel[1]);
    else if (matches(event, keys.files)) { if (!deps.hasFiles()) return; deps.files(); }
    else if (matches(event, keys.newSession)) deps.create();
    else if (matches(event, keys.zoomIn)) deps.zoom(1);
    else if (matches(event, keys.zoomOut)) deps.zoom(-1);
    else if (!(matches(event, keys.stop) && deps.stop())) return;
    event.preventDefault();
  };
}
