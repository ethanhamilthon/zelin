import { expect, it } from 'vitest';
import { highlightJavaScript } from './highlight';

it('adds syntax classes while escaping untrusted code', () => {
  const html = highlightJavaScript('var text = "<script>alert(1)</script>";\ntext');
  expect(html).toContain('tok-keyword');
  expect(html).toContain('tok-string');
  expect(html).not.toContain('<script>');
  expect(html).toContain('&lt;script&gt;');
  expect(html).toContain('\n');
});
