<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
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
  <Button data-ui="ui-stepstack" class={["toggle", { open }]} onclick={() => { open = !open; }} aria-expanded={open}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6" /></svg>
    <span class="count">{steps.length} steps</span>{#if preview}<span class="sep">|</span><span class="preview">{preview}</span>{/if}
  </Button>
  {#if open}
    <div class="items">
      <BoundedSteps {steps} {foldCode} />
    </div>
  {/if}
</div>

<style>
  :global(.stack) { margin: 0 0 8px; }
  :global(.items) { display: grid; grid-template-columns: minmax(0, 1fr); gap: 8px; margin-top: 8px; padding-left: 30px; }
  :global(.items) { position: relative; text-align: left; }
  :global(.items) > :global(*) { margin: 0; text-align: left; }
  :global(.items) > :global(.execution)::before, :global(.items) > :global(.thinking)::before { content: ""; position: absolute; left: -52px; top: -8px; height: calc(100% + 8px); border-left: 1px solid var(--border); }
  :global(.items) > :global(.execution)::after, :global(.items) > :global(.thinking)::after { content: ""; position: absolute; left: -52px; top: 14px; width: 18px; border-top: 1px solid var(--border); }
  :global(.items) > :global(.execution:last-child)::before, :global(.items) > :global(.thinking:last-child)::before { height: 22px; border-bottom-left-radius: 6px; }
  :global(.toggle) { position: relative; display: flex; align-items: center; gap: 8px; width: 100%; min-width: 0; min-height: 0; padding: 0; justify-content: flex-start; border: 0; background: none; color: var(--dim); font: inherit; line-height: 1.8; text-align: left; cursor: pointer; transition: color 140ms; }
  :global(.toggle):hover, :global(.toggle):focus-visible { color: var(--muted); outline: none; }
  :global(.count), :global(.sep) { flex-shrink: 0; }
  :global(.preview) { min-width: 0; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; color: var(--faint); }
  svg { position: absolute; left: -30px; top: 6px; transform: rotate(-90deg); transition: transform 160ms; }
  :global(.open) svg { transform: none; }
</style>
