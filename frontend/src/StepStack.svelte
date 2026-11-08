<script lang="ts">
  import type { Execution } from './agent-types';
  import ExecutionBlock from './ExecutionBlock.svelte';
  import ThinkingBlock from './ThinkingBlock.svelte';
  let { steps, foldCode = false }: { steps: Execution[]; foldCode?: boolean } = $props();
  let open = $state(false);
  const thinking = $derived(steps.filter((step) => step.kind === 'thinking').length);
  const running = $derived(steps.some((step) => step.state === 'running'));
  const label = $derived([
    steps.length - thinking && `${steps.length - thinking} code`,
    thinking && `${thinking} thinking`,
  ].filter(Boolean).join(' · '));
</script>

<div class="stack" aria-label="Steps">
  <button class="toggle" class:open onclick={() => { open = !open; }} aria-expanded={open}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6" /></svg>
    <span>{running ? 'Working' : 'Steps'} · {label}</span>
  </button>
  {#if open}
    <div class="items">
      {#each steps as step (step.id)}
        {#if step.kind === 'thinking'}<ThinkingBlock text={step.code} />{:else}<ExecutionBlock execution={step} folded={foldCode} />{/if}
      {/each}
    </div>
  {/if}
</div>

<style>
  .stack { margin: 0 0 8px; }
  .items { display: grid; grid-template-columns: minmax(0, 1fr); gap: 8px; margin-top: 8px; padding-left: 30px; }
  .items > :global(*) { margin: 0; }
  .toggle { position: relative; display: block; width: 100%; padding: 0; border: 0; background: none; color: var(--dim); font: inherit; line-height: 1.8; text-align: left; cursor: pointer; transition: color 140ms; }
  .toggle:hover, .toggle:focus-visible { color: var(--muted); outline: none; }
  svg { position: absolute; left: -30px; top: 6px; transform: rotate(-90deg); transition: transform 160ms; }
  .open svg { transform: none; }
</style>
