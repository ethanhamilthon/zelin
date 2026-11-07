import DOMPurify from 'dompurify';
import { marked } from 'marked';

export function renderMarkdown(text: string): string {
  return DOMPurify.sanitize(marked.parse(text, { async: false }), {
    FORBID_TAGS: ['img', 'video', 'audio', 'iframe', 'form', 'input', 'style'],
    FORBID_ATTR: ['style'],
  });
}
