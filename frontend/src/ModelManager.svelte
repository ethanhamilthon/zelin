<script lang="ts">
  import Modal from './Modal.svelte';
  import SearchInput from './SearchInput.svelte';
  import type { ModelOption } from './agent-types';
  import { search, modelName } from './models';
  let { models, hidden, titleModel, loading, error, onclose, ontoggle, ontitle }: {
    models: ModelOption[]; hidden: string[]; titleModel: string; loading: boolean; error: string;
    onclose: () => void; ontoggle: (id: string) => void; ontitle: (id: string) => void;
  } = $props();
  let query = $state('');
  const matches = $derived(search(models, query, 500));
</script>

<Modal title="Models" {onclose}>
  <SearchInput bind:value={query} label="Search models" />
  {#if loading}<p>Loading models…</p>{:else if error}<p role="alert">{error}</p>{/if}
  {#each matches as model (model.id)}
    {@const visible = !hidden.includes(model.id)}
    <div class="row" class:hidden={!visible}>
      <button class="eye" onclick={() => ontoggle(model.id)} aria-pressed={visible} aria-label={`${visible ? 'Hide' : 'Show'} ${modelName(model.name)}`} title={visible ? 'Shown in the bottom bar' : 'Hidden'}>
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round">
          <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12Z" /><circle cx="12" cy="12" r="3" />
          {#if !visible}<path d="m3 3 18 18" />{/if}
        </svg>
        <strong>{modelName(model.name)}</strong>
      </button>
      <button class="title" class:current={model.id === titleModel} onclick={() => ontitle(model.id)} aria-pressed={model.id === titleModel}>
        <span>{model.id === titleModel ? 'Title model' : 'Use for titles'}</span>
      </button>
    </div>
  {/each}
</Modal>

<style>
  p { margin: 0; padding: 12px; color: var(--dim); }
  .row { display: flex; }
  .row .eye { flex: 1; min-width: 0; justify-content: flex-start !important; align-items: center; }
  .row .title { width: auto !important; flex-shrink: 0; }
  .eye svg { flex-shrink: 0; color: var(--dim); }
  .hidden strong, .hidden svg { opacity: 0.45; }
  strong { font-weight: 400; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
  .title.current span { color: var(--accent); }
</style>
