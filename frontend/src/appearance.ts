export const fonts = [
  { id: 'jetbrains', name: 'JetBrains Mono', family: '"JetBrains Mono Variable", monospace' },
  { id: 'system-mono', name: 'System monospace', family: 'ui-monospace, "SFMono-Regular", Consolas, "Liberation Mono", monospace' },
  { id: 'system', name: 'System sans-serif', family: 'system-ui, sans-serif' },
  { id: 'serif', name: 'System serif', family: 'Georgia, "Liberation Serif", serif' },
];

export function applyFont(id: string) {
  document.documentElement.style.setProperty('--mono', (fonts.find((font) => font.id === id) ?? fonts[0]).family);
}
