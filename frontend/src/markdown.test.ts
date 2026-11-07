import { describe, expect, it } from 'vitest';
import { renderMarkdown } from './markdown';

describe('read-only Markdown', () => {
  it('renders formatting and code blocks', () => {
    const html = renderMarkdown('**hello**\n\n```js\nconst x = 1;\n```');
    expect(html).toContain('<strong>hello</strong>');
    expect(html).toContain('<pre><code');
  });

  it('strips scripts, event handlers, remote images and unsafe links', () => {
    const html = renderMarkdown('<script>alert(1)</script><img src="https://example.com/tracker" onerror="alert(1)">\n\n[click](javascript:alert%281%29)');
    expect(html).not.toContain('<script');
    expect(html).not.toContain('<img');
    expect(html).not.toContain('onerror');
    expect(html).not.toContain('javascript:');
  });
});
