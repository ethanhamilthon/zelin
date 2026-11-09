<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import type { Execution } from './agent-types';
  import { oneLine } from './preview';
  import BoundedSteps from './BoundedSteps.svelte';
  import Elapsed from './Elapsed.svelte';
  import Spinner from './Spinner.svelte';
  let { steps, running = false, started, activity = '' }: {
    steps: Execution[]; running?: boolean; started?: number; activity?: string;
  } = $props();
  let open = $state(false);
  const active = $derived(running || steps.some((step) => step.state === 'running'));
  const current = $derived(steps.find((step) => step.state === 'running') ?? steps.at(-1));
  const preview = $derived(active ? (current ? oneLine(current.code, /\s+/g) : activity) : '');
  const label = $derived(`${steps.length} ${steps.length === 1 ? 'step' : 'steps'}`);
</script>

<div class="stack" aria-label="Steps">
  <Button data-ui="ui-stepstack" class={["toggle", { open }]} onclick={() => { open = !open; }} aria-expanded={open}>
    <span class="icon-slot">{#if active}<Spinner />{:else}<svg class="chevron" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6" /></svg>{/if}</span>
    <span class="count">{label}</span>
    {#if started}<span class="sep">|</span><span class="count"><Elapsed {started} /></span>{/if}
    {#if preview}<span class="sep">|</span><span class="preview">{preview}</span>{/if}
  </Button>
  {#if open && steps.length}
    <div class="items">
      <BoundedSteps {steps} />
    </div>
  {/if}
</div>

<style>
  :global(.stack) { margin: 0 0 8px; }
  :global(.items) { display: grid; grid-template-columns: minmax(0, 1fr); gap: 8px; margin-top: 8px; padding-left: 30px; }
  :global(.items) { position: relative; text-align: left; }
  :global(.items) > :global(*) { margin: 0; text-align: left; }
  :global(.items) > :global(.execution)::before, :global(.items) > :global(.thinking)::before { content: ""; position: absolute; left: calc(var(--icon-x) - 22px); top: -8px; height: calc(100% + 8px); border-left: 1px solid var(--border); }
  :global(.items) > :global(.execution)::after, :global(.items) > :global(.thinking)::after { content: ""; position: absolute; left: calc(var(--icon-x) - 22px); top: 14px; width: 18px; border-top: 1px solid var(--border); }
  :global(.items) > :global(.execution:last-child)::before, :global(.items) > :global(.thinking:last-child)::before { height: 22px; border-bottom-left-radius: 6px; }
  :global(.toggle) { position: relative; display: flex; align-items: center; gap: 8px; width: 100%; min-width: 0; min-height: 0; padding: 0; justify-content: flex-start; border: 0; background: none; color: var(--dim); font: inherit; line-height: 1.8; text-align: left; cursor: pointer; transition: color 140ms; }
  :global(.toggle):hover, :global(.toggle):focus-visible { color: var(--muted); outline: none; }
  :global(.count), :global(.sep) { flex-shrink: 0; }
  :global(.preview) { min-width: 0; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; color: var(--faint); }
  .chevron { transform: rotate(-90deg); transition: transform 160ms; }
  :global(.open) .chevron { transform: none; }
</style>
