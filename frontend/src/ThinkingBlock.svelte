<script lang="ts">
  import MessageIcon from './MessageIcon.svelte';
  let { text }: { text: string } = $props();
</script>

<div class="thinking" aria-label="Thinking">
  <MessageIcon kind="thinking" />
  <p>{text}</p>
</div>

<style>
  .thinking { position: relative; margin: 0 0 8px; color: var(--dim); }
  p { margin: 0; white-space: pre-wrap; overflow-wrap: anywhere; line-height: 1.8; max-height: 320px; overflow-y: auto; }
</style>
