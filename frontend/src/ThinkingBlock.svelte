<script lang="ts">
  import MessageIcon from './MessageIcon.svelte';
  let { text }: { text: string } = $props();
</script>

<div class="thinking" aria-label="Thinking">
  <MessageIcon kind="thinking" />
  <p>{text}</p>
</div>

<style>
  .thinking { position: relative; margin: 0 0 8px; padding: 4px 14px; border-radius: 6px; color: var(--dim); }
  .thinking :global(.icon) { top: 10px; }
  p { margin: 0; white-space: pre-wrap; overflow-wrap: anywhere; line-height: 1.8; max-height: 320px; overflow-y: auto; }
</style>
