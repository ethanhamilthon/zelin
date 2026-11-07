import { describe, expect, it } from 'vitest';
import { renderMarkdown } from './markdown';

describe('minimal markdown', () => {
  it('renders bold, italic and inline code without markers', () => {
    expect(renderMarkdown('**bold** and *it* or _it_ with `a **b**`'))
      .toBe('<strong>bold</strong> and <em>it</em> or <em>it</em> with <code>a **b**</code>');
  });

  it('keeps everything else literal and escaped', () => {
    const text = '# Title\n- item\n* bullet\n2 * 3 * 4\nsnake_case_name\n```js\nx\n```\n[link](url) <b>';
    expect(renderMarkdown(text)).toBe(text.replace('<b>', '&lt;b&gt;'));
  });
});
