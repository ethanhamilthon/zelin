import { describe, expect, it } from 'vitest';
import { isPath, renderMarkdown } from './markdown';

describe('minimal markdown', () => {
  it('renders bold, italic and inline code without markers', () => {
    expect(renderMarkdown('**bold** and *it* or _it_ with `a **b**`'))
      .toBe('<strong>bold</strong> and <em>it</em> or <em>it</em> with <code>a **b**</code>');
  });

  it('keeps everything else literal and escaped', () => {
    const text = '# Title\n- item\n* bullet\n2 * 3 * 4\nsnake_case_name\n```js\nx\n```\nand/or <b>';
    expect(renderMarkdown(text)).toBe(text.replace('<b>', '&lt;b&gt;'));
  });

  it('renders markdown links and bare urls', () => {
    expect(renderMarkdown('see [**docs**](https://a.dev/x) or https://b.dev/y.'))
      .toBe('see <a href="https://a.dev/x" data-url="https://a.dev/x"><strong>docs</strong></a> or <a href="https://b.dev/y" data-url="https://b.dev/y">https://b.dev/y</a>.');
  });

  it('links file paths in text and inline code', () => {
    expect(renderMarkdown('open ~/Projects/zash/TODO.md, `README.md` and `npm test`'))
      .toBe('open <a href="~/Projects/zash/TODO.md" data-path="~/Projects/zash/TODO.md">~/Projects/zash/TODO.md</a>, <a href="README.md" data-path="README.md"><code>README.md</code></a> and <code>npm test</code>');
  });

  it('detects paths conservatively', () => {
    expect(['src/App.svelte', './a.ts', '/etc/hosts', 'go.mod', '~/x'].every(isPath)).toBe(true);
    expect(['npm test', 'a.b', '/', 'foo', 'x => y', 'zash'].some(isPath)).toBe(false);
  });
});
