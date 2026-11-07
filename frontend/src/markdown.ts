import { escape } from './highlight';

const code = /(?<!`)`([^`\n]+)`(?!`)/g;
const bold = /\*\*(?=\S)([^\n]*?\S)\*\*/g;
const italic = /(?<![*\w])\*(?=[^\s*])([^*\n]*?[^\s*])\*(?![*\w])|(?<![_\w])_(?=[^\s_])([^_\n]*?[^\s_])_(?![_\w])/g;

function emphasis(text: string) {
  return escape(text)
    .replace(bold, '<strong>$1</strong>')
    .replace(italic, (_, star, underscore) => `<em>${star ?? underscore}</em>`);
}

export function renderMarkdown(text: string): string {
  let html = '';
  let last = 0;
  for (const match of text.matchAll(code)) {
    html += emphasis(text.slice(last, match.index)) + `<code>${escape(match[1])}</code>`;
    last = match.index + match[0].length;
  }
  return html + emphasis(text.slice(last));
}
