import { call } from './bridge';
import type { Entry } from './refs';

export type Listing = { dir: string; parent: string; entries: Entry[]; truncated: boolean };

export const browse = (workspace: string, dir = '') => call<Listing>('WorkspaceService.Browse', workspace, dir);
export const suggest = (workspace: string, dir: string, query: string) => call<Entry[]>('WorkspaceService.Suggest', workspace, dir, query);
export const openPath = (workspace: string, path: string) => call<void>('WorkspaceService.Open', workspace, path);
export const isMarkdownName = (name: string) => /\.(md|markdown)$/i.test(name);

export function crumbs(root: string, dir: string) {
  const parts = dir.slice(root.length).split('/').filter(Boolean);
  const out = [{ name: root.split('/').filter(Boolean).pop() ?? '/', path: root }];
  parts.forEach((part, index) => out.push({ name: part, path: root + '/' + parts.slice(0, index + 1).join('/') }));
  return out;
}
