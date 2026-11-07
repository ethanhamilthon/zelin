import { expect, it } from 'vitest';
import { highlightJavaScript, highlightMarkdown } from './highlight';

it('adds syntax classes while escaping untrusted code', () => {
  const html = highlightJavaScript('var text = "<script>alert(1)</script>";\ntext');
  expect(html).toContain('tok-keyword');
  expect(html).toContain('tok-string');
  expect(html).not.toContain('<script>');
  expect(html).toContain('&lt;script&gt;');
  expect(html).toContain('\n');
});

it('keeps Markdown markers literal and escapes HTML', () => {
  const source = '## Title\n\n- **bold** item\n\n<img src=x onerror=alert(1)>\n\n[link](javascript:alert)';
  const html = highlightMarkdown(source);
  const element = document.createElement('div');
  element.innerHTML = html;
  expect(element.textContent).toBe(source);
  expect(element.querySelector('img, a, strong, ul, h2')).toBeNull();
  expect(element.querySelectorAll('span').length).toBeGreaterThan(3);
});
