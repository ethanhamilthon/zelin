<script lang="ts">
  import Modal from './Modal.svelte';
  import SearchInput from './SearchInput.svelte';
  import type { ModelOption } from './agent-types';
  import { search, modelName } from './models';
  let { models, hidden, titleModel, loading, error, back = false, onclose, ondismiss, ontoggle, ontitle }: {
    models: ModelOption[]; hidden: string[]; titleModel: string; loading: boolean; error: string; back?: boolean;
    onclose: () => void; ondismiss?: () => void; ontoggle: (id: string) => void; ontitle: (id: string) => void;
  } = $props();
  let query = $state('');
  const matches = $derived(search(models, query, 500));
</script>

<Modal title="Models" {back} {onclose} {ondismiss}>
  {#snippet top()}<SearchInput bind:value={query} label="Search models" />{/snippet}
  {#if loading}<p>Loading models…</p>{:else if error}<p role="alert">{error}</p>{:else if !matches.length}<p>No models found</p>{/if}
  {#each matches as model (model.id)}
    {@const visible = !hidden.includes(model.id)}
    {@const name = modelName(model.name)}
    {@const utility = model.id === titleModel}
    <div class="row" class:hidden={!visible}>
      <strong title={model.id}>{name}</strong>
      <div class="actions">
        <button class="icon eye" onclick={() => ontoggle(model.id)} aria-pressed={visible} aria-label={`${visible ? 'Hide' : 'Show'} ${name}`} title={visible ? 'Visible in model picker' : 'Hidden from model picker'}>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round">
            <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12Z" /><circle cx="12" cy="12" r="3" />
            {#if !visible}<path d="m3 3 18 18" />{/if}
          </svg>
        </button>
        <button class="icon utility" class:current={utility} onclick={() => ontitle(model.id)} aria-pressed={utility} aria-label={`Use ${name} for background tasks`} title={utility ? 'Utility model (titles, compaction, handoff)' : 'Use for background tasks (titles, compaction, handoff)'}>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round">
            <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76Z" />
          </svg>
        </button>
      </div>
    </div>
  {/each}
</Modal>

<style>
  p { margin: 0; padding: 12px; color: var(--dim); }
  .row { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 4px 4px 4px 12px; border-radius: 4px; }
  .row:hover { background: var(--hover); }
  strong { min-width: 0; flex: 1; font-weight: 400; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
  .hidden strong, .hidden .eye { opacity: 0.45; }
  .actions { display: flex; align-items: center; gap: 4px; flex-shrink: 0; }
  .icon { width: auto !important; padding: 8px !important; color: var(--dim) !important; justify-content: center !important; }
  .icon:hover, .icon:focus-visible { color: var(--foreground) !important; }
  .utility.current { color: var(--accent) !important; }
</style>
