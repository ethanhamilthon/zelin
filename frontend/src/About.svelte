<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from './Modal.svelte';
  import { loadAbout, notices, type AboutInfo } from './about';
  let { back = false, onclose, ondismiss }: { back?: boolean; onclose: () => void; ondismiss?: () => void } = $props();
  let info = $state<AboutInfo | null>(null), error = $state('');
  onMount(() => { void loadAbout().then((value) => { info = value; }).catch((failure) => { error = String(failure).replace(/^Error: /, ''); }); });
  const rows = $derived(info ? [
    ['Version', info.version], ['Build', [info.commit, info.date].filter(Boolean).join(' · ') || 'Not available'],
    ['Runtime', info.runtime], ['Platform', info.platform], ['Profile', info.profile], ['Daemon', info.daemon],
  ] : []);
</script>

<Modal title="About" {back} {onclose} {ondismiss}>
  {#if error}<p role="alert">{error}</p>
  {:else if info}
    <h2>{info.name} <small>{info.version}</small></h2>
    <p>{info.description}</p>
    <dl>{#each rows as [label, value] (label)}<dt>{label}</dt><dd>{value}</dd>{/each}</dl>
    <h3>Open source notices</h3>
    <ul>{#each notices as item (item.name)}<li><a href={item.url} data-url={item.url}>{item.name}</a> <small>{item.license}</small></li>{/each}</ul>
  {:else}<p role="status">Loading…</p>{/if}
</Modal>

<style>
  h2 { margin: 0 0 6px; font-size: 1.15em; color: var(--foreground); }
  small { color: var(--dim); font-weight: 400; }
  p { margin: 0 0 12px; color: var(--muted); }
  dl { display: grid; grid-template-columns: max-content 1fr; gap: 4px 16px; margin: 0 0 12px; }
  dt { color: var(--dim); }
  dd { margin: 0; color: var(--muted); overflow-wrap: anywhere; }
  h3 { margin: 0 0 6px; font-size: 1em; color: var(--muted); }
  ul { margin: 0; padding: 0; list-style: none; }
  li { color: var(--muted); }
  a { color: inherit; text-decoration: underline; text-decoration-color: var(--dim); }
</style>
