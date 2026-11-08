import { expect, it } from 'vitest';
import { applyTheme } from './themes';

it('applies the dark Zen theme', () => {
  const root = document.createElement('div');
  applyTheme(root);
  expect(root.style.getPropertyValue('--background')).toBe('#0a0a0a');
  expect(root.style.getPropertyValue('--error')).toBe('#ef8a80');
  expect(root.style.colorScheme).toBe('dark');
});
