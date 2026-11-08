<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import { onMount } from 'svelte';
  import SearchInput from './SearchInput.svelte';
  import Modal from './Modal.svelte';
  import { browse, crumbs, isMarkdownName, type Listing } from './browser';
  import type { Entry } from './refs';
  let { workspace, root, start = '', back = false, onclose, ondismiss, onopenfile, onopenapp, ontag }: {
    workspace: string; root: string; start?: string; back?: boolean; onclose: () => void; ondismiss?: () => void;
    onopenfile: (entry: Entry) => void; onopenapp: (entry: Entry) => void; ontag: (entry: Entry) => void;
  } = $props();
  let query = $state('');
  const needle = $derived(query.trim().toLowerCase());
  let listing = $state<Listing | null>(null), error = $state(''), loading = $state(true), token = 0;
  const entries = $derived((listing?.entries ?? []).filter((entry) => entry.name.toLowerCase().includes(needle)));

  async function enter(dir: string) {
    const mine = ++token;
    loading = true; error = ''; query = '';
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
  {#snippet actions()}
    <nav class="crumbs" aria-label="Project path" title={listing?.dir ?? root}>
      {#if listing?.parent}<Button data-ui="ui-filebrowser" variant="icon" class="up" onclick={() => { void enter(listing!.parent); }} aria-label="Go up">←</Button>{/if}
      {#each trail as crumb, index (crumb.path)}
        {#if index}<span aria-hidden="true">/</span>{/if}
        <Button data-ui="ui-filebrowser" onclick={() => { void enter(crumb.path); }}>{index === 0 ? root : crumb.name}</Button>
      {/each}
    </nav>
  {/snippet}
  {#snippet top()}<SearchInput bind:value={query} label="Search files" />{/snippet}
  {#if loading}<p role="status">Loading…</p>
  {:else if error}<p role="alert">{error}</p>
  {:else if listing && !listing.entries.length}<p>This folder is empty</p>
  {:else if listing}
    {#each entries as entry (entry.path)}
      <div class={["row", { broken: entry.broken }]}>
        <span class="name" title={entry.link ?? entry.path}>{entry.name}{#if entry.link}<small> → {entry.link}</small>{/if}{#if entry.broken}<small> (broken link)</small>{/if}</span>
        <span class="acts">
          {#if entry.kind === 'dir'}
            <Button data-ui="ui-filebrowser" variant="icon" onclick={() => { void enter(entry.link ?? entry.path); }} aria-label="Open folder" title="Open folder"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><path d="M3 7V5h6l2 2h10v13H3Z" /></svg></Button>
          {:else if isMarkdownName(entry.name)}
            <Button data-ui="ui-filebrowser" variant="icon" onclick={() => onopenfile(entry)} aria-label="Open in app" title="Open in app">↗</Button>
          {:else}
            <Button data-ui="ui-filebrowser" variant="icon" onclick={() => onopenapp(entry)} disabled={entry.broken} aria-label="Open in default app" title="Open in default app">↗</Button>
          {/if}
          <Button data-ui="ui-filebrowser" variant="icon" onclick={() => ontag(entry)} disabled={entry.broken} aria-label="Tag" title="Tag">@</Button>
        </span>
      </div>
    {:else}<p>No matching files</p>
    {/each}
    {#if listing.truncated}<p>Showing the first entries only</p>{/if}
  {/if}
</Modal>

<style>
  p { margin: 0; padding: 12px; color: var(--dim); }
  :global(.crumbs) { display: flex; align-items: center; gap: 4px; min-width: 0; overflow-x: auto; color: var(--dim); }
  :global(.crumbs) :global(button[data-ui="ui-filebrowser"]) { width: auto; padding: 2px 0; white-space: nowrap; flex-shrink: 0; }
  :global(.row) { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 4px 12px; }
  :global(.name) { min-width: 0; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
  small { color: var(--dim); }
  :global(.acts) { display: flex; gap: 8px; flex-shrink: 0; }
  :global(.acts) :global(button[data-ui="ui-filebrowser"]) { width: 32px; height: 32px; padding: 4px; align-items: center; justify-content: center; color: var(--dim); }
  :global(.broken) :global(.name) { color: var(--dim); }
</style>
