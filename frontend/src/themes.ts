export type Theme = { id: string; name: string; dark: boolean; colors: [string, string, string, string, string, string] };

const theme = (id: string, name: string, dark: boolean, ...colors: Theme['colors']): Theme => ({ id, name, dark, colors });

export const themes: Theme[] = [
  theme('zen', 'Zen', true, '#111111', '#e0e0e0', '#679bd4', '#b5b5b5', '#7c7c7c', '#606060'),
  theme('zen-light', 'Zen Light', false, '#f7f7f5', '#2a2a2a', '#3f74b5', '#454545', '#6f6f6f', '#9a9a9a'),
  theme('gruvbox', 'Gruvbox Dark', true, '#282828', '#ebdbb2', '#fabd2f', '#fb4934', '#b8bb26', '#928374'),
  theme('nord', 'Nord', true, '#2e3440', '#d8dee9', '#88c0d0', '#81a1c1', '#a3be8c', '#616e88'),
  theme('dracula', 'Dracula', true, '#282a36', '#f8f8f2', '#bd93f9', '#ff79c6', '#f1fa8c', '#6272a4'),
  theme('solarized-dark', 'Solarized Dark', true, '#002b36', '#93a1a1', '#268bd2', '#859900', '#2aa198', '#586e75'),
  theme('solarized-light', 'Solarized Light', false, '#fdf6e3', '#586e75', '#268bd2', '#859900', '#2aa198', '#93a1a1'),
  theme('tokyo-night', 'Tokyo Night', true, '#1a1b26', '#c0caf5', '#7aa2f7', '#bb9af7', '#9ece6a', '#565f89'),
  theme('catppuccin-mocha', 'Catppuccin Mocha', true, '#1e1e2e', '#cdd6f4', '#89b4fa', '#cba6f7', '#a6e3a1', '#6c7086'),
  theme('catppuccin-latte', 'Catppuccin Latte', false, '#eff1f5', '#4c4f69', '#1e66f5', '#8839ef', '#40a02b', '#9ca0b0'),
  theme('one-dark', 'One Dark', true, '#282c34', '#abb2bf', '#61afef', '#c678dd', '#98c379', '#5c6370'),
  theme('rose-pine', 'Rosé Pine', true, '#191724', '#e0def4', '#c4a7e7', '#31748f', '#f6c177', '#6e6a86'),
  theme('kanagawa', 'Kanagawa', true, '#1f1f28', '#dcd7ba', '#7e9cd8', '#957fb8', '#98bb6c', '#727169'),
  theme('everforest', 'Everforest', true, '#2d353b', '#d3c6aa', '#a7c080', '#e67e80', '#dbbc7f', '#859289'),
  theme('monokai', 'Monokai', true, '#272822', '#f8f8f2', '#66d9ef', '#f92672', '#e6db74', '#75715e'),
];

const mix = (color: string, amount: number, base = 'var(--background)') => `color-mix(in srgb, ${color} ${amount}%, ${base})`;

export function applyTheme(id: string, root: HTMLElement = document.documentElement) {
  const current = themes.find((item) => item.id === id) ?? themes[0];
  const [background, foreground, accent, keyword, string, comment] = current.colors;
  const shade = current.dark ? '#000' : 'var(--foreground)';
  const values: Record<string, string> = {
    background, foreground, accent, keyword, string, comment,
    muted: mix(foreground, 78), dim: mix(foreground, 48), faint: mix(foreground, 36),
    border: mix(foreground, 14), surface: mix(foreground, 7), hover: mix(foreground, 13),
    selection: mix(accent, 32), 'code-bg': mix(foreground, current.dark ? 4 : 3),
    'thinking-bg': mix(shade, current.dark ? 32 : 8),
  };
  for (const [name, value] of Object.entries(values)) root.style.setProperty(`--${name}`, value);
  root.style.colorScheme = current.dark ? 'dark' : 'light';
}
