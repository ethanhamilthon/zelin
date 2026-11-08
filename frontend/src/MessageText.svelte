<script lang="ts">
  import { highlightMarkdown } from './highlight';
  import { splitTokens } from './images';
  import Markdown from './Markdown.svelte';
  import { tagged } from './refs';
  import type { Ref } from './refs';
  let { content, markdown = false, images = [], refs = [], onimage, onref }: {
    content: string; markdown?: boolean; images?: string[]; refs?: Ref[];
    onimage?: (url: string, number: number) => void; onref?: (ref: Ref) => void;
  } = $props();
  const parts = $derived(splitTokens(content));
  const found = (id: string) => refs.find((ref) => ref.id === id);
</script>

{#if markdown}
  <Markdown content={tagged(content)} />
{:else}
  <div class="text">{#each parts as part}{#if part.image && images[part.image - 1]}<button class="token" onclick={() => onimage?.(images[part.image! - 1], part.image!)}>{part.text}</button>{:else if part.ref}{@const ref = found(part.ref.id)}{#if ref}<button class="token" onclick={() => onref?.(ref)}>{part.text}</button>{:else}{part.text}{/if}{:else}{@html highlightMarkdown(part.text)}{/if}{/each}</div>
{/if}

<style>
  .text { font-size: var(--font-size); line-height: 1.8; white-space: pre-wrap; }
  .token { font: inherit; color: var(--accent); background: none; border: 0; padding: 0; cursor: pointer; }
  .token:hover { text-decoration: underline; }
</style>
