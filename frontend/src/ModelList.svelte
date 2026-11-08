<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import Modal from './Modal.svelte';
  import SearchInput from './SearchInput.svelte';
  import type { ModelOption } from './agent-types';
  import { compact, perMillion } from './format';
  import { search, modelName } from './models';
  let { models, active, hidden = [], loading, error, back = false, onclose, ondismiss, onselect }: {
    models: ModelOption[]; active: string; hidden?: string[]; loading: boolean; error: string; back?: boolean;
    onclose: () => void; ondismiss?: () => void; onselect: (id: string) => void;
  } = $props();
  let query = $state('');
  const matches = $derived(search(models.filter((model) => !hidden.includes(model.id)), query));

  function submit(event: KeyboardEvent) {
    if (event.key === 'Enter' && matches.length) { event.preventDefault(); onselect(matches[0].id); }
  }
</script>

<Modal title="Select model" {back} {onclose} {ondismiss}>
  {#snippet top()}<SearchInput bind:value={query} onkeydown={submit} label="Search models" />{/snippet}
  {#if loading}<p>Loading models…</p>{:else if error}<p role="alert">{error}</p>{:else if !matches.length}<p>No models found</p>{/if}
  {#each matches as model (model.id)}
    <Button data-ui="ui-modellist" onclick={() => onselect(model.id)} aria-current={model.id === active ? 'true' : undefined} title={model.id}>
      <strong>{modelName(model.name)}</strong>
      <span>{model.contextWindow ? compact(model.contextWindow) : '?'} · {perMillion(model.pricing.prompt)} in · {perMillion(model.pricing.completion)} out</span>
    </Button>
  {/each}
</Modal>

<style>
  p { margin: 0; padding: 12px; color: var(--dim); }
  strong { font-weight: 400; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
  :global(button[aria-current='true']) strong { color: var(--accent); }
</style>
