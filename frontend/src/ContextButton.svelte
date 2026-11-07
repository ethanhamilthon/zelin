<script lang="ts">
  import { fly } from 'svelte/transition';
  import type { Usage } from './agent-types';
  import { percent } from './format';
  import UsageView from './UsageView.svelte';
  let { usage, contextWindow }: { usage: Usage | null; contextWindow: number } = $props();
  let open = $state(false);
  let root: HTMLDivElement;
  const ratio = $derived(usage?.context != null && contextWindow > 0 ? Math.min(usage.context / contextWindow, 1) : 0);
  const label = $derived(`Context and usage: ${percent(usage?.context ?? (usage?.requests ? null : 0), contextWindow)} of context used`);

  function outside(event: MouseEvent) { if (open && !root.contains(event.target as Node)) open = false; }
  function escape(event: KeyboardEvent) {
    if (open && event.key === 'Escape') { event.stopPropagation(); open = false; root.querySelector('button')?.focus(); }
  }
</script>

<svelte:window onclick={outside} />

<div class="context" bind:this={root} onkeydown={escape} role="presentation">
  <button onclick={() => { open = !open; }} aria-expanded={open} aria-label={label} title={label}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke-width="3">
      <circle cx="12" cy="12" r="8" stroke="#404040" />
      <circle cx="12" cy="12" r="8" stroke="currentColor" pathLength="100" stroke-dasharray={`${ratio * 100} 100`} transform="rotate(-90 12 12)" />
    </svg>
  </button>
  {#if open}
    <div class="popover" role="dialog" aria-label="Context and usage" transition:fly={{ y: 6, duration: 140 }}>
      <UsageView {usage} {contextWindow} />
    </div>
  {/if}
</div>

<style>
  .context { position: relative; display: flex; }
  button { display: flex; color: inherit; border: 0; background: none; padding: 0; cursor: pointer; transition: color 140ms; }
  button:hover, button:focus-visible, button[aria-expanded='true'] { color: #ddd; }
  .popover { position: absolute; bottom: 32px; left: -12px; z-index: 5; min-width: 300px; padding: 12px 4px 4px; border: 1px solid #353535; border-radius: 8px; background: #1c1c1c; box-shadow: 0 20px 80px #0006; cursor: default; }
</style>
