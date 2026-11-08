<script lang="ts">
  import type { Execution } from './agent-types';
  import { oneLine } from './preview';
  import BoundedSteps from './BoundedSteps.svelte';
  let { steps, running = false, foldCode = false }: { steps: Execution[]; running?: boolean; foldCode?: boolean } = $props();
  let open = $state(false);
  const active = $derived(running || steps.some((step) => step.state === 'running'));
  const current = $derived(steps.find((step) => step.state === 'running') ?? steps.at(-1));
  const preview = $derived(active && current ? oneLine(current.code, /\s+/g) : '');
</script>

<div class="stack" aria-label="Steps">
  <button class="toggle" class:open onclick={() => { open = !open; }} aria-expanded={open}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6" /></svg>
    <span class="count">{steps.length} steps</span>{#if preview}<span class="sep">|</span><span class="preview">{preview}</span>{/if}
  </button>
  {#if open}
    <div class="items">
      <BoundedSteps {steps} {foldCode} />
    </div>
  {/if}
</div>

<style>
  .stack { margin: 0 0 8px; }
  .items { display: grid; grid-template-columns: minmax(0, 1fr); gap: 8px; margin-top: 8px; padding-left: 30px; }
  .items > :global(*) { margin: 0; }
  .toggle { position: relative; display: flex; align-items: center; gap: 8px; width: 100%; min-width: 0; padding: 0; border: 0; background: none; color: var(--dim); font: inherit; line-height: 1.8; text-align: left; cursor: pointer; transition: color 140ms; }
  .toggle:hover, .toggle:focus-visible { color: var(--muted); outline: none; }
  .count, .sep { flex-shrink: 0; }
  .preview { min-width: 0; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; color: var(--faint); }
  svg { position: absolute; left: -30px; top: 6px; transform: rotate(-90deg); transition: transform 160ms; }
  .open svg { transform: none; }
</style>
