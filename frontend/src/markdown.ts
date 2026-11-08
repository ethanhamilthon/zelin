import { Marked, type Tokens } from 'marked';
import DOMPurify from 'dompurify';
import { escape, highlightJavaScript } from './highlight';
import { anchor, isUrl, linkPaths, localTarget } from './paths';

export { isPath, isUrl } from './paths';

const HIGHLIGHT_LIMIT = 20000;
const CACHE_LIMIT = 200;
const languages = new Set(['js', 'javascript', 'ts', 'typescript', 'jsx', 'tsx', 'json']);
const cache = new Map<string, string>();
let depth = 0;

function code({ text, lang }: Tokens.Code) {
  const name = (lang ?? '').trim().split(/\s+/)[0].replace(/[^\w+#.-]/g, '');
  const body = languages.has(name.toLowerCase()) && text.length <= HIGHLIGHT_LIMIT ? highlightJavaScript(text) : escape(text);
  const label = name ? `<span class="lang">${escape(name)}</span>` : '';
  return `<div class="code-block">${label}<pre><code${name ? ` class="language-${escape(name)}"` : ''}>${body}</code></pre></div>`;
}

const marked = new Marked({
  gfm: true,
  breaks: false,
  renderer: {
    html: ({ text }) => escape(text),
    image: ({ text }) => escape(text),
    code,
    codespan(token) {
      const html = `<code>${escape(token.text)}</code>`;
      return depth === 0 && localTarget(token.text) ? anchor(token.text, html) : html;
    },
    link(token) {
      depth++;
      const inner = this.parser.parseInline(token.tokens);
      depth--;
      if (isUrl(token.href)) return anchor(token.href, inner);
      return localTarget(token.href) ? anchor(token.href, inner) : inner;
    },
    text(token) {
      if ('tokens' in token && token.tokens) return this.parser.parseInline(token.tokens);
      return depth === 0 ? linkPaths(token.text) : token.text;
    },
  },
});

const config = { ADD_ATTR: ['data-url', 'data-path'], FORBID_TAGS: ['style', 'img', 'form', 'iframe', 'object', 'embed'], FORBID_ATTR: ['style'] };

function build(text: string) {
  const html = DOMPurify.sanitize(marked.parse(text, { async: false }), config);
  return html.replace(/<table>/g, '<div class="table-wrap"><table>').replace(/<\/table>/g, '</table></div>').trim();
}

export function renderMarkdown(text: string): string {
  const hit = cache.get(text);
  if (hit !== undefined) return hit;
  const html = build(text);
  if (cache.size >= CACHE_LIMIT) cache.delete(cache.keys().next().value!);
  cache.set(text, html);
  return html;
}
