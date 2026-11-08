<script lang="ts">
  import { onMount } from 'svelte';
  import Hotkey from './Hotkey.svelte';
  import type { Snippet } from 'svelte';
  import { fade, fly } from 'svelte/transition';
  let { title, hint = 'Close Esc', onclose, children, actions }: { title: string; hint?: string; onclose: () => void; children: Snippet; actions?: Snippet } = $props();
  let panel: HTMLDivElement;
  onMount(() => { if (!panel.contains(document.activeElement)) (panel.querySelector<HTMLButtonElement>('input:not(:disabled), textarea:not(:disabled), select:not(:disabled), button:not(:disabled)') ?? panel).focus(); });

  function navigate(event: KeyboardEvent) {
    if (event.key === 'Escape') { event.preventDefault(); event.stopPropagation(); onclose(); return; }
    if (!['ArrowDown', 'ArrowUp', 'Tab'].includes(event.key)) return;
    if (event.key !== 'Tab' && event.target instanceof HTMLSelectElement) return;
    const buttons = [...panel.querySelectorAll<HTMLElement>('input:not(:disabled), textarea:not(:disabled), select:not(:disabled), button:not(:disabled)')];
    if (!buttons.length) return;
    event.preventDefault();
    const index = buttons.indexOf(document.activeElement as HTMLElement);
    const step = event.key === 'ArrowUp' || (event.key === 'Tab' && event.shiftKey) ? -1 : 1;
    buttons[(index + step + buttons.length) % buttons.length].focus();
  }
</script>

<div class="overlay" role="presentation" onclick={(event) => { if (event.target === event.currentTarget) onclose(); }} transition:fade={{ duration: 120 }}>
  <div class="panel" role="dialog" aria-modal="true" aria-label={title} tabindex="-1" bind:this={panel} onkeydown={navigate} transition:fly={{ y: -6, duration: 160 }}>
    <div class="title"><span class="heading">{title}</span>{#if actions}{@render actions()}{:else}<span>{#if hint === 'Close Esc'}Close <Hotkey value="Esc" />{:else}{hint}{/if}</span>{/if}</div>
    {@render children()}
  </div>
</div>

<style>
  .overlay { position: fixed; inset: 0; background: color-mix(in srgb, #000 50%, transparent); z-index: 10; padding: min(20vh, 120px) 20px 20px; }
  .panel { max-width: 640px; max-height: 65vh; overflow-y: auto; margin: auto; border: 1px solid var(--border); border-radius: 8px; background: color-mix(in srgb, var(--foreground) 5%, var(--background)); padding: 8px; box-shadow: 0 20px 80px color-mix(in srgb, #000 40%, transparent); outline: none; }
  .title { padding: 10px 12px 16px; color: var(--muted); font: var(--font-size)/1.5 var(--mono); display: flex; justify-content: space-between; align-items: center; gap: 16px; }
  .title .heading { color: var(--muted); }
  .title span { color: var(--faint); font-size: var(--font-size); }
  .panel :global(button) { width: 100%; display: flex; justify-content: space-between; gap: 16px; padding: 12px; border: 0; border-radius: 4px; color: var(--foreground); background: none; text-align: left; font: var(--font-size)/1.5 var(--mono); cursor: pointer; transition: background 140ms; }
  .panel :global(button:focus), .panel :global(button:hover) { outline: none; background: var(--hover); }
  .panel .title :global(button) { width: auto; flex-shrink: 0; gap: 12px; padding: 2px 8px; color: var(--muted); }
  .panel :global(button:disabled) { opacity: 0.35; cursor: default; }
  .panel :global(button span) { color: var(--dim); font-size: var(--font-size); }
</style>
