<script lang="ts">
  import type { ModelOption } from './agent-types';
  import { compact } from './format';
  import { modelName } from './models';
  let { models, active, onpick, onhover }: { models: ModelOption[]; active: number; onpick: (id: string) => void; onhover: (index: number) => void } = $props();
  let list: HTMLDivElement;
  $effect(() => { list?.children[active]?.scrollIntoView({ block: 'nearest' }); });
</script>

<div class="picker" role="listbox" aria-label="Models" bind:this={list}>
  {#each models as model, index (model.id)}
    <button role="option" aria-selected={index === active} class:active={index === active} onclick={() => onpick(model.id)} onmousemove={() => onhover(index)}>
      <span class="name">{modelName(model.name)}</span>
      <span class="meta">{model.id}{model.contextWindow ? ` · ${compact(model.contextWindow)}` : ''}</span>
    </button>
  {:else}
    <p>No models found</p>
  {/each}
</div>

<style>
  .picker { max-height: 40vh; overflow-y: auto; margin: 0 var(--gutter) 8px; padding: 4px; border: 1px solid var(--border); border-radius: 6px; }
  button { display: flex; width: 100%; gap: 16px; justify-content: space-between; padding: 6px 10px; border: 0; border-radius: 4px; background: none; color: var(--muted); font: inherit; line-height: 1.6; text-align: left; cursor: pointer; }
  .active { background: var(--hover); color: var(--foreground); }
  .name, .meta { min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .meta { color: var(--faint); flex-shrink: 1; }
  p { margin: 0; padding: 6px 10px; color: var(--dim); }
</style>
