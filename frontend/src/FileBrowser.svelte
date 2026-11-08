<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from './Modal.svelte';
  import { browse, crumbs, isMarkdownName, type Listing } from './browser';
  import type { Entry } from './refs';
  let { workspace, root, start = '', back = false, onclose, ondismiss, onopenfile, onopenapp, ontag }: {
    workspace: string; root: string; start?: string; back?: boolean; onclose: () => void; ondismiss?: () => void;
    onopenfile: (entry: Entry) => void; onopenapp: (entry: Entry) => void; ontag: (entry: Entry) => void;
  } = $props();
  let listing = $state<Listing | null>(null), error = $state(''), loading = $state(true), token = 0;

  async function enter(dir: string) {
    const mine = ++token;
    loading = true; error = '';
    try {
      const next = await browse(workspace, dir);
      if (mine === token) listing = next;
    } catch (failure) { if (mine === token) error = String(failure).replace(/^Error: /, ''); }
    finally { if (mine === token) loading = false; }
  }
  onMount(() => { void enter(start || root); });
  const trail = $derived(listing ? crumbs(root, listing.dir) : []);
</script>

<Modal title="Files" {back} {onclose} {ondismiss}>
  {#snippet top()}
    <nav class="crumbs" aria-label="Path">
      {#if listing?.parent}<button class="up" onclick={() => { void enter(listing!.parent); }} aria-label="Go up">←</button>{/if}
      {#each trail as crumb, index (crumb.path)}
        {#if index}<span aria-hidden="true">/</span>{/if}
        <button onclick={() => { void enter(crumb.path); }}>{crumb.name}</button>
      {/each}
    </nav>
  {/snippet}
  {#if loading}<p role="status">Loading…</p>
  {:else if error}<p role="alert">{error}</p>
  {:else if listing && !listing.entries.length}<p>This folder is empty</p>
  {:else if listing}
    {#each listing.entries as entry (entry.path)}
      <div class="row" class:broken={entry.broken}>
        <span class="name" title={entry.link ?? entry.path}>{entry.kind === 'dir' ? '▸ ' : ''}{entry.name}{#if entry.link}<small> → {entry.link}</small>{/if}{#if entry.broken}<small> (broken link)</small>{/if}</span>
        <span class="acts">
          {#if entry.kind === 'dir'}
            <button onclick={() => { void enter(entry.link ?? entry.path); }}>Open folder</button>
          {:else if isMarkdownName(entry.name)}
            <button onclick={() => onopenfile(entry)}>Open in app</button>
          {:else}
            <button onclick={() => onopenapp(entry)} disabled={entry.broken}>Open in default app</button>
          {/if}
          <button onclick={() => ontag(entry)} disabled={entry.broken}>Tag</button>
        </span>
      </div>
    {/each}
    {#if listing.truncated}<p>Showing the first entries only</p>{/if}
  {/if}
</Modal>

<style>
  p { margin: 0; padding: 12px; color: var(--dim); }
  .crumbs { display: flex; flex-wrap: wrap; align-items: center; gap: 4px; margin: 0 0 8px; padding: 0 12px; color: var(--dim); }
  .crumbs button { width: auto; padding: 2px 6px; }
  .row { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 4px 12px; }
  .name { min-width: 0; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
  small { color: var(--dim); }
  .acts { display: flex; gap: 8px; flex-shrink: 0; }
  .acts button { width: auto; padding: 2px 8px; }
  .broken .name { color: var(--dim); }
</style>
