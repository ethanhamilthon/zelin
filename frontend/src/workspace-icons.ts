export const workspaceIcons = [
  { id: 'folder', name: 'Folder', path: 'M3 7V5h6l2 2h10v13H3Z' },
  { id: 'code', name: 'Code', path: 'm8 7-5 5 5 5m8-10 5 5-5 5m-3-13-2 16' },
  { id: 'rocket', name: 'Rocket', path: 'M14 5c3-3 7-2 7-2s1 4-2 7l-7 7-5-5Zm-7 7H3l4-5h5m0 10v4l5-4v-5M4 17l-1 4 4-1' },
  { id: 'leaf', name: 'Leaf', path: 'M20 3C8 3 3 7 3 13a8 8 0 0 0 8 8c6 0 9-8 9-18ZM4 20 15 9' },
  { id: 'palette', name: 'Palette', path: 'M12 3a9 9 0 1 0 0 18h2a2 2 0 0 0 0-4h-1a2 2 0 0 1 0-4h5a3 3 0 0 0 3-3c0-4-4-7-9-7ZM7 9h.01M11 6h.01M16 7h.01M6 14h.01' },
  { id: 'tool', name: 'Tool', path: 'm14 6 4 4 3-3a6 6 0 0 1-8 8l-7 7-4-4 7-7a6 6 0 0 1 8-8Z' },
  { id: 'book', name: 'Book', path: 'M12 6C9 3 5 3 2 4v15c4-1 7-1 10 2 3-3 6-3 10-2V4c-3-1-7-1-10 2Zm0 0v15' },
  { id: 'bolt', name: 'Bolt', path: 'm13 2-9 12h7l-1 8 10-13h-7Z' },
  { id: 'globe', name: 'Globe', path: 'M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM3 12h18M12 3c-5 5-5 13 0 18 5-5 5-13 0-18Z' },
  { id: 'flask', name: 'Flask', path: 'M9 3h6m-5 0v7L4 20h16l-6-10V3M7 15h10' },
];
const legacy = ['📁', '💻', '🚀', '🌱', '🎨', '🔧', '📚', '⚡', '🌍', '🧪'];
export function workspaceIcon(id: string) {
  return workspaceIcons.find((icon) => icon.id === id) ?? workspaceIcons[legacy.indexOf(id)] ?? workspaceIcons[0];
}
