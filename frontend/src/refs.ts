export type Ref = { id: string; kind: 'file' | 'dir' | 'image'; label: string; path?: string; data?: string; workspaceId?: string; sessionId?: string };
export type Entry = { name: string; path: string; kind: 'file' | 'dir'; link?: string; broken?: boolean; size: number };

export const refPattern = /\[@([A-Za-z0-9_-]{4,40})(?:\|([^\]\n]*))?\]/g;

export const refToken = (ref: Pick<Ref, 'id' | 'label'>) => `[@${ref.id}|${ref.label.replace(/[\]\n|]/g, ' ')}]`;

export function newId() {
  const bytes = crypto.getRandomValues(new Uint8Array(6));
  return [...bytes].map((byte) => byte.toString(16).padStart(2, '0')).join('');
}

export function entryRef(entry: Entry, workspaceId: string, sessionId: string): Ref {
  return { id: newId(), kind: entry.kind === 'dir' ? 'dir' : isImageName(entry.name) ? 'image' : 'file', label: entry.name, path: entry.path, workspaceId, sessionId };
}

export const isImageName = (name: string) => /\.(png|jpe?g|gif|webp)$/i.test(name);

export function usedRefs(text: string, refs: Ref[]): Ref[] {
  const ids = new Set([...text.matchAll(refPattern)].map((match) => match[1]));
  return refs.filter((ref) => ids.has(ref.id));
}

export function unresolved(text: string, refs: Ref[]) {
  const known = new Set(refs.map((ref) => ref.id));
  return [...text.matchAll(refPattern)].filter((match) => !known.has(match[1])).map((match) => match[2] || match[1]);
}

export const tagged = (text: string) => text.replace(refPattern, (_, id: string, label?: string) => '@' + (label || id));

const extensions: Record<string, string> = { 'image/jpeg': 'jpg', 'image/gif': 'gif', 'image/webp': 'webp' };

export function imageName(refs: Ref[], url: string, name?: string) {
  const taken = new Set(refs.map((ref) => ref.label));
  const clean = (name ?? '').replace(/[\s\]|]+/g, '_');
  if (isImageName(clean) && !taken.has(clean)) return clean;
  const extension = extensions[/^data:([^;,]+)/.exec(url)?.[1] ?? ''] ?? 'png';
  let number = refs.filter((ref) => ref.data).length + 1;
  while (taken.has(`image-${number}.${extension}`)) number++;
  return `image-${number}.${extension}`;
}
