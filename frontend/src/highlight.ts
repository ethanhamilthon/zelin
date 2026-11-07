import { parser } from '@lezer/javascript';
import { highlightCode, classHighlighter } from '@lezer/highlight';

function escape(text: string): string {
  return text.replace(/[&<>"']/g, (char) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[char]!);
}

export function highlightJavaScript(code: string): string {
  let html = '';
  highlightCode(code, parser.parse(code), classHighlighter,
    (text, classes) => { html += classes ? `<span class="${classes}">${escape(text)}</span>` : escape(text); },
    () => { html += '\n'; });
  return html;
}
