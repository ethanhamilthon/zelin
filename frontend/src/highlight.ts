import { parser } from '@lezer/javascript';
import { markdownLanguage } from '@codemirror/lang-markdown';
import { highlightCode, classHighlighter } from '@lezer/highlight';
import type { Highlighter } from '@lezer/highlight';

type Tree = ReturnType<typeof parser.parse>;
import { highlighting } from './editor-theme';

export function escape(text: string): string {
  return text.replace(/[&<>"']/g, (char) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[char]!);
}

function render(code: string, tree: Tree, highlighter: Highlighter): string {
  let html = '';
  highlightCode(code, tree, highlighter,
    (text, classes) => { html += classes ? `<span class="${classes}">${escape(text)}</span>` : escape(text); },
    () => { html += '\n'; });
  return html;
}

export const highlightJavaScript = (code: string) => render(code, parser.parse(code), classHighlighter);
export const highlightMarkdown = (text: string) => render(text, markdownLanguage.parser.parse(text), highlighting);
