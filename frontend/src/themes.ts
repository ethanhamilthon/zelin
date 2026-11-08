const mix = (color: string, amount: number) => `color-mix(in srgb, ${color} ${amount}%, var(--background))`;

const background = '#0a0a0a', foreground = '#e0e0e0';

export function applyTheme(root: HTMLElement = document.documentElement) {
  const values: Record<string, string> = {
    background, foreground, accent: '#679bd4', keyword: '#b5b5b5', string: '#7c7c7c', comment: '#606060',
    success: '#8bd5a0', warning: '#edc777', error: '#ef8a80', violet: '#c4a7ef',
    muted: mix(foreground, 78), dim: mix(foreground, 48), faint: mix(foreground, 36),
    border: mix(foreground, 14), surface: mix(foreground, 7), hover: mix(foreground, 13),
    selection: 'rgb(255 255 255 / 0.18)', 'code-bg': mix(foreground, 4), 'thinking-bg': mix('#000', 32),
  };
  for (const [name, value] of Object.entries(values)) root.style.setProperty(`--${name}`, value);
  root.style.colorScheme = 'dark';
}
