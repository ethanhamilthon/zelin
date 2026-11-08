import { describe, expect, it } from 'vitest';
import { isPath, renderMarkdown } from './markdown';

const dom = (text: string) => { const box = document.createElement('div'); box.innerHTML = renderMarkdown(text); return box; };

describe('markdown', () => {
  it('renders inline formatting', () => {
    expect(renderMarkdown('**bold** and *it* or _it_ with `a **b**`')).toBe('<p><strong>bold</strong> and <em>it</em> or <em>it</em> with <code>a **b**</code></p>');
    expect(renderMarkdown('~~gone~~')).toBe('<p><del>gone</del></p>');
  });

  it('renders headings, lists, task lists, quotes and rules', () => {
    const box = dom('# Title\n\n- a\n  - nested\n- [x] done\n- [ ] todo\n\n1. one\n2. two\n\n> quote\n\n---');
    expect(box.querySelector('h1')?.textContent).toBe('Title');
    expect(box.querySelectorAll('ul ul li')).toHaveLength(1);
    expect(box.querySelectorAll('input[type=checkbox]')).toHaveLength(2);
    expect(box.querySelector<HTMLInputElement>('input[type=checkbox]')?.checked).toBe(true);
    expect(box.querySelectorAll('ol li')).toHaveLength(2);
    expect(box.querySelector('blockquote')?.textContent?.trim()).toBe('quote');
    expect(box.querySelector('hr')).not.toBeNull();
  });

  it('renders GFM tables with alignment and escaped pipes inside a scroll container', () => {
    const box = dom('| a | b |\n|:--|--:|\n| 1 \\| 2 | **x** |');
    expect(box.querySelector('.table-wrap > table')).not.toBeNull();
    expect(box.querySelectorAll('th')).toHaveLength(2);
    expect(box.querySelector('th')?.getAttribute('align')).toBe('left');
    expect(box.querySelectorAll('th')[1].getAttribute('align')).toBe('right');
    expect(box.querySelector('td')?.textContent).toBe('1 | 2');
    expect(box.querySelector('td strong')?.textContent).toBe('x');
  });

  it('renders fenced code with a language label, highlighting and nested fences', () => {
    const box = dom('```js\nconst a = "<b>";\n```\n\n````md\n```js\nx\n```\n````');
    const blocks = box.querySelectorAll('.code-block');
    expect(blocks[0].querySelector('.lang')?.textContent).toBe('js');
    expect(blocks[0].querySelector('.tok-keyword')?.textContent).toBe('const');
    expect(blocks[0].querySelector('b')).toBeNull();
    expect(blocks[1].querySelector('code')?.textContent).toBe('```js\nx\n```');
    expect(dom('```\nplain\n```').querySelector('.lang')).toBeNull();
  });

  it('does not highlight oversized code blocks', () => {
    const big = 'let x = 1;\n'.repeat(3000);
    expect(dom('```js\n' + big + '```').querySelector('.tok-keyword')).toBeNull();
  });

  it('keeps incomplete markdown literal', () => {
    expect(renderMarkdown('**unclosed and `tick')).toBe('<p>**unclosed and `tick</p>');
  });

  it('renders links, bare urls and mail addresses', () => {
    const box = dom('see [**docs**](https://a.dev/x) or https://b.dev/y. mail user@example.com');
    const links = [...box.querySelectorAll('a')];
    expect(links.map((a) => a.dataset.url)).toEqual(['https://a.dev/x', 'https://b.dev/y', 'mailto:user@example.com']);
    expect(links[0].querySelector('strong')?.textContent).toBe('docs');
  });

  it('links local paths in text, inline code and relative links without touching code blocks', () => {
    const box = dom('open ~/Projects/zelin/TODO.md, `README.md`, `npm test`, [rel](./src/App.svelte) and a/b\n\n```\n/etc/hosts\n```');
    expect([...box.querySelectorAll('a')].map((a) => a.dataset.path)).toEqual(['~/Projects/zelin/TODO.md', 'README.md', './src/App.svelte']);
    expect(box.querySelector('a code')?.textContent).toBe('README.md');
    expect(box.querySelector('pre a')).toBeNull();
  });

  it('never creates anchors inside link text', () => {
    expect(dom('[see /etc/hosts here](https://a.dev)').querySelectorAll('a')).toHaveLength(1);
  });

  it('disables raw html and unsafe links and never loads remote images', () => {
    const html = renderMarkdown('<script>alert(1)</script> <img src=x onerror=alert(1)>');
    expect(html).toContain('&lt;script&gt;');
    expect(html).not.toMatch(/<script|<img/);
    const box = dom('[x](javascript:alert(1)) [y](data:text/html;base64,AA) ![i](http://e.com/a.png) <iframe src=x></iframe>');
    expect(box.querySelector('a')).toBeNull();
    expect(box.querySelector('img, iframe')).toBeNull();
    expect(box.innerHTML).not.toContain('javascript:');
  });

  it('escapes attributes in language labels', () => {
    const box = dom('```"><script>x</script>\ncode\n```');
    expect(box.querySelector('script')).toBeNull();
  });

  it('reuses rendered output for unchanged text', () => {
    expect(renderMarkdown('same *text*')).toBe(renderMarkdown('same *text*'));
  });

  it('renders a large response quickly', () => {
    const text = Array.from({ length: 2000 }, (_, i) => '## Part ' + i + '\n\nSome *text* with `code` and a [link](https://a.dev/' + i + ').\n\n- one\n- two').join('\n\n');
    const began = performance.now();
    renderMarkdown(text);
    expect(performance.now() - began).toBeLessThan(2000);
  });
});

describe('path detection', () => {
  it('detects paths conservatively', () => {
    expect(['src/App.svelte', './a.ts', '/etc/hosts', 'go.mod', '~/x'].every(isPath)).toBe(true);
    expect(['npm test', 'a.b', '/', 'foo', 'x => y', 'zash'].some(isPath)).toBe(false);
  });
});
