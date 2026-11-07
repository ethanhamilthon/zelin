<script lang="ts">
  import { fly } from 'svelte/transition';
  import type { Toast } from './toasts.svelte';
  let { items, onopen }: { items: Toast[]; onopen: (toast: Toast) => void } = $props();
</script>

<div class="toasts" aria-live="polite">
  {#each items as toast (toast.id)}
    <button onclick={() => onopen(toast)} transition:fly={{ y: -8, duration: 160 }}>{toast.text}</button>
  {/each}
</div>

<style>
  .toasts { position: fixed; top: 16px; left: 50%; transform: translateX(-50%); z-index: 20; display: flex; flex-direction: column; gap: 8px; align-items: center; pointer-events: none; }
  button { pointer-events: auto; max-width: min(560px, 90vw); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; padding: 10px 16px; border: 1px solid var(--border); border-radius: 8px; background: color-mix(in srgb, var(--foreground) 7%, var(--background)); color: var(--foreground); font: var(--font-size)/1.5 var(--mono); cursor: pointer; box-shadow: 0 12px 40px color-mix(in srgb, #000 35%, transparent); }
  button:hover { border-color: var(--accent); }
</style>
