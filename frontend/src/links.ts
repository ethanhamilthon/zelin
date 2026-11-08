import { isMarkdown, openFile, openUrl, resolveFile } from './files';

export type Target = { url?: string; path?: string };

export async function follow(target: Target, edit: (path: string) => void) {
  if (target.url) return openUrl(target.url);
  if (!target.path) return;
  const path = await resolveFile(target.path);
  if (isMarkdown(path)) edit(path);
  else await openFile(path);
}

export function links(node: HTMLElement, open: (target: Target) => void) {
  const click = (event: MouseEvent) => {
    const anchor = (event.target as Element).closest<HTMLAnchorElement>('a[data-url], a[data-path]');
    if (!anchor || !node.contains(anchor)) return;
    event.preventDefault();
    open({ url: anchor.dataset.url, path: anchor.dataset.path });
  };
  node.addEventListener('click', click);
  return { destroy: () => node.removeEventListener('click', click) };
}
