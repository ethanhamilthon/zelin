<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import type { ModelOption } from './agent-types';
  import { compact } from './format';
  import { modelName } from './models';
  let { models, active, onpick, onhover }: { models: ModelOption[]; active: number; onpick: (id: string) => void; onhover: (index: number) => void } = $props();
  let list: HTMLDivElement;
  $effect(() => { list?.children[active]?.scrollIntoView({ block: 'nearest' }); });
</script>

<div class="picker" role="listbox" aria-label="Models" bind:this={list}>
  {#each models as model, index (model.id)}
    <Button data-ui="ui-modelpicker" role="option" aria-selected={index === active} class={{ active: index === active }} onclick={() => onpick(model.id)} onmousemove={() => onhover(index)}>
      <span class="name">{modelName(model.name)}</span>
      <span class="meta">{model.id}{model.contextWindow ? ` · ${compact(model.contextWindow)}` : ''}</span>
    </Button>
  {:else}
    <p>No models found</p>
  {/each}
</div>

<style>
  :global(.picker) { max-height: 40vh; overflow-y: auto; margin: 0 var(--gutter) 8px; padding: 4px; border: 1px solid var(--border); border-radius: 6px; }
  :global(button[data-ui="ui-modelpicker"]) { display: flex; width: 100%; gap: 16px; justify-content: space-between; padding: 6px 10px; border: 0; border-radius: 4px; background: none; color: var(--muted); font: inherit; line-height: 1.6; text-align: left; cursor: pointer; }
  :global(.active) { background: var(--hover); color: var(--foreground); }
  :global(.name), :global(.meta) { min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  :global(.meta) { color: var(--faint); flex-shrink: 1; }
  p { margin: 0; padding: 6px 10px; color: var(--dim); }
</style>
