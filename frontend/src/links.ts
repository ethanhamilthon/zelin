import { isMarkdown, openFile, openUrl, resolveFile } from './files';

export type Target = { url?: string; path?: string; bases?: string[] };

function candidates(path: string, bases: string[] = []) {
  if (path.startsWith('/') || path.startsWith('~')) return [path];
  const list = [path];
  for (const base of bases) {
    const trimmed = base.replace(/\/+$/, '');
    const dir = /\.[a-z0-9]+$/i.test(trimmed) && !base.endsWith('/') ? trimmed.replace(/\/[^/]+$/, '') : trimmed;
    if (dir && (dir.startsWith('/') || dir.startsWith('~') || dir.includes('/'))) list.push(`${dir}/${path}`);
  }
  return [...new Set(list)];
}

export async function follow(target: Target, edit: (path: string) => void) {
  if (target.url) return openUrl(target.url);
  if (!target.path) return;
  let resolved = '';
  let lastError: unknown;
  for (const candidate of candidates(target.path, target.bases)) {
    try { resolved = await resolveFile(candidate); break; }
    catch (error) { lastError = error; }
  }
  if (!resolved) throw lastError;
  if (isMarkdown(resolved)) edit(resolved);
  else await openFile(resolved);
}

export function links(node: HTMLElement, open: (target: Target) => void) {
  const click = (event: MouseEvent) => {
    const anchor = (event.target as Element).closest<HTMLAnchorElement>('a[data-url], a[data-path]');
    if (!anchor || !node.contains(anchor)) return;
    event.preventDefault();
    const scope = anchor.closest('.text') ?? anchor.closest('article');
    const bases = scope ? [...scope.querySelectorAll<HTMLAnchorElement>('a[data-path]')].map((a) => a.dataset.path ?? '').filter((p) => p && p !== anchor.dataset.path) : [];
    open({ url: anchor.dataset.url, path: anchor.dataset.path, ...(bases.length ? { bases } : {}) });
  };
  node.addEventListener('click', click);
  return { destroy: () => node.removeEventListener('click', click) };
}
