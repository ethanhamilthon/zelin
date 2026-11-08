import { afterEach, beforeEach, vi } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import { EditorView } from '@codemirror/view';
import App from './App.svelte';
import { agentInfo, runAgent } from './agent';
import { currentSession, listSessions } from './sessions';
import { call } from './bridge';
import { placeholderInfo } from './agent-types';

export const $ = (selector: string) => document.querySelector<HTMLElement>(selector);
export const all = (selector: string) => [...document.querySelectorAll<HTMLElement>(selector)];
export const wait = async () => { await new Promise((r) => setTimeout(r, 160)); await tick(); };
export const key = (k: string, code: string, init: KeyboardEventInit = {}) => window.dispatchEvent(new KeyboardEvent('keydown', { key: k, code, cancelable: true, ...init }));
export const button = (text: string, scope: ParentNode = document) => [...scope.querySelectorAll<HTMLButtonElement>('button')].find((b) => (b.getAttribute('aria-label') ?? b.textContent?.trim()) === text)!;
export const view = () => EditorView.findFromDOM($('.cm-editor')!)!;
export const write = (text: string) => { view().dispatch({ changes: { from: 0, to: view().state.doc.length, insert: text }, selection: { anchor: text.length } }); };
export const press = (name: string) => $('.cm-content')!.dispatchEvent(new KeyboardEvent('keydown', { key: name, bubbles: true, cancelable: true }));
export const entries = [
  { name: 'src', path: '/w/src', kind: 'dir', size: 0 },
  { name: 'notes.md', path: '/w/notes.md', kind: 'file', size: 5 },
  { name: 'data.csv', path: '/w/data.csv', kind: 'file', size: 5 },
  { name: 'link', path: '/w/link', kind: 'file', size: 0, link: '/w/real.txt' },
];
export const calls: [string, unknown[]][] = [];

export function useApp() {
  let app: ReturnType<typeof mount>;
  beforeEach(async () => {
    vi.resetAllMocks();
    calls.length = 0;
    Element.prototype.scrollIntoView = vi.fn();
    Object.defineProperties(Range.prototype, { getClientRects: { configurable: true, value: () => [] }, getBoundingClientRect: { configurable: true, value: () => new DOMRect() } });
    vi.mocked(currentSession).mockResolvedValue({ id: 'one', title: 'Task', draft: '', messages: [], updated: '', workspaceId: 'default' });
    vi.mocked(listSessions).mockResolvedValue([]);
    vi.mocked(agentInfo).mockResolvedValue(placeholderInfo('Test'));
    vi.mocked(runAgent).mockImplementation(() => new Promise(() => {}));
    vi.mocked(call).mockImplementation(async (method, ...args) => {
      calls.push([method, args]);
      if (method === 'WorkspaceService.List') return [{ id: 'default', name: 'Default', icon: '', dir: '/w' }];
      if (method === 'WorkspaceService.Browse') return args[1] === '/w/src' ? { dir: '/w/src', parent: '/w', entries: [], truncated: false } : { dir: '/w', parent: '', entries, truncated: false };
      if (method === 'WorkspaceService.Suggest') return entries.filter((e) => e.name.startsWith(String(args[2]))).slice(0, 5);
      if (method === 'FileService.Read') return '# notes';
      return '';
    });
    app = mount(App, { target: document.body }); flushSync(); await wait();
  });
  afterEach(async () => { await unmount(app); document.body.replaceChildren(); });
}
