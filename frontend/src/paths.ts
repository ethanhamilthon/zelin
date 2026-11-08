import { escape } from './highlight';

const filename = /^[\w-][\w.-]*\.(md|go|ts|js|svelte|json|ya?ml|txt|css|html|toml|sh|py|rs|mod|sum|lock)$/i;
const codePath = /^(~|\.{1,2})?\/?[\w.-]+(\/[\w.-]+)*\/?$/;
const barePath = /(?<![\w/.~])~?\/[\w.-]+(?:\/[\w.-]+)*/g;

export const isUrl = (target: string) => /^(https?:\/\/|mailto:)/i.test(target);

export function isPath(text: string) {
  return codePath.test(text) && (text.includes('/') || filename.test(text)) && !/^\/+$/.test(text);
}

export const localTarget = (href: string) => isPath(href) || /^(~|\.{1,2})\/[^\s]+$/.test(href);

export function anchor(target: string, html: string) {
  const kind = isUrl(target) ? 'url' : 'path';
  return `<a href="${escape(target)}" data-${kind}="${escape(target)}">${html}</a>`;
}

export function linkPaths(escaped: string) {
  let html = '', last = 0;
  for (const match of escaped.matchAll(barePath)) {
    if (match[0].length < 2) continue;
    html += escaped.slice(last, match.index) + anchor(match[0], match[0]);
    last = match.index + match[0].length;
  }
  return html + escaped.slice(last);
}
