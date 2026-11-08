import { escape } from './highlight';

const token = /(?<!`)`([^`\n]+)`(?!`)|\[([^\]\n]+)\]\(([^)\s]+)\)|\bhttps?:\/\/[^\s<>"'`]+|(?<![\w/.~])~?\/[\w.-]+(?:\/[\w.-]+)*/g;
const bold = /\*\*(?=\S)([^\n]*?\S)\*\*/g;
const italic = /(?<![*\w])\*(?=[^\s*])([^*\n]*?[^\s*])\*(?![*\w])|(?<![_\w])_(?=[^\s_])([^_\n]*?[^\s_])_(?![_\w])/g;
const filename = /^[\w-][\w.-]*\.(md|go|ts|js|svelte|json|ya?ml|txt|css|html|toml|sh|py|rs|mod|sum|lock)$/i;
const codePath = /^(~|\.{1,2})?\/?[\w.-]+(\/[\w.-]+)*\/?$/;

export const isUrl = (target: string) => /^(https?:\/\/|mailto:)/i.test(target);

export function isPath(text: string) {
  return codePath.test(text) && (text.includes('/') || filename.test(text)) && !/^\/+$/.test(text);
}

function emphasis(text: string) {
  return escape(text)
    .replace(bold, '<strong>$1</strong>')
    .replace(italic, (_, star, underscore) => `<em>${star ?? underscore}</em>`);
}

function link(target: string, html: string) {
  const kind = isUrl(target) ? 'url' : 'path';
  return `<a href="${escape(target)}" data-${kind}="${escape(target)}">${html}</a>`;
}

function trim(text: string) {
  const end = text.match(/[.,;:!?)\]]+$/)?.[0] ?? '';
  return [text.slice(0, text.length - end.length), end];
}

function render(match: RegExpMatchArray) {
  const [whole, code, label, target] = match;
  if (code !== undefined) return isPath(code) ? link(code, `<code>${escape(code)}</code>`) : `<code>${escape(code)}</code>`;
  if (label !== undefined) return link(target, emphasis(label));
  const [value, rest] = trim(whole);
  return (value.length > 1 ? link(value, escape(value)) : escape(value)) + escape(rest);
}

export function renderMarkdown(text: string): string {
  let html = '';
  let last = 0;
  for (const match of text.matchAll(token)) {
    html += emphasis(text.slice(last, match.index)) + render(match);
    last = match.index + match[0].length;
  }
  return html + emphasis(text.slice(last));
}
