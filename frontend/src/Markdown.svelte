<script lang="ts">
  import { renderMarkdown } from './markdown';
  let { content }: { content: string } = $props();
  const html = $derived(renderMarkdown(content));
</script>

<div class="text md">{@html html}</div>

<style>
  .text { font-size: var(--font-size); line-height: 1.8; overflow-wrap: anywhere; }
  .md :global(> :first-child) { margin-top: 0; }
  .md :global(> :last-child) { margin-bottom: 0; }
  .md :global(p), .md :global(ul), .md :global(ol), .md :global(blockquote), .md :global(.code-block), .md :global(.table-wrap), .md :global(hr) { margin: 0 0 0.8em; }
  .md :global(h1), .md :global(h2), .md :global(h3), .md :global(h4), .md :global(h5), .md :global(h6) { margin: 1.1em 0 0.4em; color: var(--foreground); font-size: 1em; font-weight: 700; line-height: 1.5; }
  .md :global(h1) { font-size: 1.25em; }
  .md :global(h2) { font-size: 1.12em; }
  .md :global(ul), .md :global(ol) { padding-left: 1.6em; }
  .md :global(li > ul), .md :global(li > ol), .md :global(li > p) { margin: 0; }
  .md :global(li) { margin: 0.15em 0; }
  .md :global(input[type='checkbox']) { margin: 0 0.5em 0 -1.2em; }
  .md :global(blockquote) { padding-left: 1em; border-left: 2px solid var(--border); color: var(--dim); }
  .md :global(hr) { border: 0; border-top: 1px solid var(--border); }
  .md :global(strong) { font-weight: 700; color: var(--foreground); }
  .md :global(a) { color: inherit; text-decoration: underline; text-decoration-color: var(--dim); text-underline-offset: 3px; cursor: pointer; }
  .md :global(a:hover) { color: var(--accent); text-decoration-color: currentColor; }
  .md :global(code) { font-family: inherit; color: var(--foreground); background: var(--code-bg); border-radius: 4px; padding: 1px 4px; }
  .md :global(.code-block) { position: relative; }
  .md :global(.code-block .lang) { position: absolute; top: 4px; right: 8px; color: var(--faint); font-size: 0.85em; }
  .md :global(pre) { margin: 0; padding: 10px 12px; overflow-x: auto; white-space: pre; border-radius: 4px; background: var(--code-bg); line-height: 1.6; }
  .md :global(pre code) { padding: 0; background: none; }
  .md :global(.table-wrap) { max-width: 100%; overflow-x: auto; }
  .md :global(table) { border-collapse: collapse; min-width: 100%; }
  .md :global(th), .md :global(td) { padding: 4px 10px; border: 1px solid var(--border); text-align: left; vertical-align: top; }
  .md :global(th) { color: var(--foreground); font-weight: 700; }
  .md :global([align='center']) { text-align: center; }
  .md :global([align='right']) { text-align: right; }
  .md :global(.tok-keyword), .md :global(.tok-operator) { color: var(--keyword); }
  .md :global(.tok-string), .md :global(.tok-number) { color: var(--string); }
  .md :global(.tok-comment) { color: var(--comment); font-style: italic; }
  .md :global(.tok-variableName), .md :global(.tok-propertyName) { color: var(--muted); }
</style>
