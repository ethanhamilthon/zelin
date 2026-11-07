<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from './Modal.svelte';
  import type { ModelOption } from './agent-types';
  import { compact, perMillion } from './format';
  let { models, active, loading, error, onclose, onselect }: {
    models: ModelOption[]; active: string; loading: boolean; error: string;
    onclose: () => void; onselect: (id: string) => void;
  } = $props();
  let query = $state('');
  let input: HTMLInputElement;
  const matches = $derived(models.filter((model) => `${model.id} ${model.name}`.toLowerCase().includes(query.trim().toLowerCase())).slice(0, 100));
  onMount(() => input.focus());

  function submit(event: KeyboardEvent) {
    if (event.key === 'Enter' && matches.length) { event.preventDefault(); onselect(matches[0].id); }
  }
</script>

<Modal title="Models" {onclose}>
  <input bind:this={input} bind:value={query} onkeydown={submit} placeholder="Search models…" aria-label="Search models" />
  {#if loading}<p>Loading models…</p>{:else if error}<p role="alert">{error}</p>{:else if !matches.length}<p>No models found</p>{/if}
  {#each matches as model (model.id)}
    <button onclick={() => onselect(model.id)} aria-current={model.id === active ? 'true' : undefined} title={model.id}>
      <strong>{model.name}</strong>
      <span>{model.contextWindow ? compact(model.contextWindow) : '?'} · {perMillion(model.pricing.prompt)} in · {perMillion(model.pricing.completion)} out</span>
    </button>
  {/each}
</Modal>

<style>
  input { width: 100%; margin: 0 0 8px; padding: 10px 12px; border: 1px solid #353535; border-radius: 4px; background: var(--background); color: #ddd; font: inherit; outline: none; }
  input:focus { border-color: #555; }
  p { margin: 0; padding: 12px; color: #777; }
  strong { font-weight: 400; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
  button[aria-current='true'] strong { color: #fff; }
</style>
