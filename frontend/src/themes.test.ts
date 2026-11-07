import { describe, expect, it } from 'vitest';
import { applyTheme, themes } from './themes';

describe('themes', () => {
  it('offers 10 to 15 unique themes', () => {
    expect(themes.length).toBeGreaterThanOrEqual(10);
    expect(themes.length).toBeLessThanOrEqual(15);
    expect(new Set(themes.map((theme) => theme.id)).size).toBe(themes.length);
  });

  it('applies colors and falls back to the default', () => {
    const root = document.createElement('div');
    applyTheme('nord', root);
    expect(root.style.getPropertyValue('--background')).toBe('#2e3440');
    applyTheme('missing', root);
    expect(root.style.getPropertyValue('--background')).toBe('#111111');
    expect(root.style.getPropertyValue('--thinking-bg')).toContain('color-mix');
  });
});
