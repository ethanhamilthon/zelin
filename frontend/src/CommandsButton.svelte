<script lang="ts">
  import { Command } from '@lucide/svelte';
  let { onclick }: { onclick: () => void } = $props();
</script>

<button class="commands" aria-label="Commands" title="Commands" {onclick}><Command size={18} /></button>

<style>
  .commands { position: fixed; top: 12px; right: 16px; z-index: 20; display: flex; align-items: center; justify-content: center; width: 32px; height: 32px; padding: 0; border: 0; border-radius: 8px; background: none; color: var(--dim); cursor: pointer; transition: color 140ms; }
  .commands:hover, .commands:focus-visible { color: var(--foreground); outline: none; }
</style>
