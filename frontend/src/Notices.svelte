<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import type { Conversation } from './conversation';
  let { conversation, error = '', onretry }: { conversation: Conversation; error?: string; onretry: () => void } = $props();
</script>

{#if conversation.busy && !conversation.started}<p class="error" role="status">{conversation.activity}</p>{/if}
{#if conversation.error}<p class="error" role="alert">{conversation.error}{#if conversation.compactFailed && !conversation.busy} <Button data-ui="ui-notices" class="retry" onclick={onretry}>Retry compaction</Button>{/if}</p>{/if}
{#if error}<p class="error" role="alert">{error}</p>{/if}

<style>
  :global(.error) { margin: 0; padding: 12px clamp(20px, 6vw, 80px); font: var(--font-size)/1.5 var(--mono); color: var(--muted); }
  :global(.retry) { margin-left: 10px; padding: 2px 8px; border: 1px solid var(--border); border-radius: 4px; background: none; color: var(--accent); font: inherit; cursor: pointer; }
</style>
