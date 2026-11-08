export type FileKind = 'image' | 'code' | 'text' | 'data' | 'archive' | 'pdf' | 'file';

const kinds: [RegExp, FileKind][] = [
  [/\.(png|jpe?g|gif|webp)$/i, 'image'],
  [/\.(ts|tsx|js|jsx|mjs|go|py|rs|svelte|vue|sh|c|cpp|h|java|rb|css|html)$/i, 'code'],
  [/\.(md|txt|log|rst)$/i, 'text'],
  [/\.(json|ya?ml|toml|csv|xml|sql|db)$/i, 'data'],
  [/\.(zip|tar|gz|tgz|xz|7z|rar)$/i, 'archive'],
  [/\.pdf$/i, 'pdf'],
];

export const fileKind = (path: string): FileKind => kinds.find(([pattern]) => pattern.test(path))?.[1] ?? 'file';
