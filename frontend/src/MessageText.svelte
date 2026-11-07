<script lang="ts">
  import { highlightMarkdown } from './highlight';
  import { renderMarkdown } from './markdown';
  import { splitTokens } from './images';
  let { content, markdown = false, images = [], onimage }: {
    content: string; markdown?: boolean; images?: string[]; onimage?: (url: string, number: number) => void;
  } = $props();
  const parts = $derived(splitTokens(content));
</script>

<div class="text">{#if markdown}{@html renderMarkdown(content)}{:else}{#each parts as part}{#if part.image && images[part.image - 1]}<button class="token" onclick={() => onimage?.(images[part.image! - 1], part.image!)}>{part.text}</button>{:else}{@html highlightMarkdown(part.text)}{/if}{/each}{/if}</div>

<style>
  .text { font-size: var(--font-size); line-height: 1.8; white-space: pre-wrap; }
  .text :global(strong) { font-weight: 700; color: var(--foreground); }
  .text :global(code) { font-family: inherit; color: var(--foreground); background: var(--code-bg); border-radius: 4px; padding: 1px 4px; }
  .token { font: inherit; color: var(--accent); background: none; border: 0; padding: 0; cursor: pointer; }
  .token:hover { text-decoration: underline; }
</style>
