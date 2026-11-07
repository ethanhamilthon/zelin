<script lang="ts">
  import { onMount } from 'svelte';
  import type { Snippet } from 'svelte';
  import { fade, fly } from 'svelte/transition';
  let { title, onclose, children }: { title: string; onclose: () => void; children: Snippet } = $props();
  let panel: HTMLDivElement;
  onMount(() => (panel.querySelector<HTMLButtonElement>('button:not(:disabled)') ?? panel).focus());

  function navigate(event: KeyboardEvent) {
    if (event.key === 'Escape') { event.preventDefault(); event.stopPropagation(); onclose(); return; }
    if (!['ArrowDown', 'ArrowUp', 'Tab'].includes(event.key)) return;
    const buttons = [...panel.querySelectorAll<HTMLButtonElement>('button:not(:disabled)')];
    if (!buttons.length) return;
    event.preventDefault();
    const index = buttons.indexOf(document.activeElement as HTMLButtonElement);
    const step = event.key === 'ArrowUp' || (event.key === 'Tab' && event.shiftKey) ? -1 : 1;
    buttons[(index + step + buttons.length) % buttons.length].focus();
  }
</script>

<div class="overlay" transition:fade={{ duration: 120 }}>
  <div class="panel" role="dialog" aria-modal="true" aria-label={title} tabindex="-1" bind:this={panel} onkeydown={navigate} transition:fly={{ y: -6, duration: 160 }}>
    <div class="title">{title}<span>Esc to close</span></div>
    {@render children()}
  </div>
</div>

<style>
  .overlay { position: fixed; inset: 0; background: #0008; z-index: 10; padding: min(20vh, 120px) 20px 20px; }
  .panel { max-width: 500px; max-height: 65vh; overflow-y: auto; margin: auto; border: 1px solid #353535; border-radius: 8px; background: #1c1c1c; padding: 8px; box-shadow: 0 20px 80px #0006; outline: none; }
  .title { padding: 10px 12px 16px; color: #999; font: 13px/1.5 var(--mono); display: flex; justify-content: space-between; }
  .title span { color: #666; font-size: 12px; }
  .panel :global(button) { width: 100%; display: flex; justify-content: space-between; gap: 16px; padding: 12px; border: 0; border-radius: 4px; color: #ddd; background: none; text-align: left; font: 14px/1.5 var(--mono); cursor: pointer; transition: background 140ms; }
  .panel :global(button:focus), .panel :global(button:hover) { outline: none; background: #2b2b2b; }
  .panel :global(button:disabled) { opacity: 0.35; cursor: default; }
  .panel :global(button span) { color: #777; font-size: 12px; }
</style>
