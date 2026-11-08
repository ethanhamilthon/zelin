<script lang="ts">
  import WorkspaceIcon from './WorkspaceIcon.svelte';
  import type { Workspace } from './workspaces';
  let { workspace, selected = false, folded = false, onselect, onsettings }: {
    workspace: Workspace; selected?: boolean; folded?: boolean; onselect: () => void; onsettings: () => void;
  } = $props();
</script>

<div class="head" class:selected>
  <button class="name" onclick={onselect} aria-pressed={selected} aria-expanded={!folded} title={workspace.dir}>
    <svg class="chevron" class:folded width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="m6 9 6 6 6-6" /></svg>
    <span class="icon"><WorkspaceIcon icon={workspace.icon} /></span><span class="label">{workspace.name}</span>
  </button>
  <button class="more" onclick={onsettings} aria-label={`Settings for ${workspace.name}`} title="Workspace settings">⋯</button>
</div>

<style>
  .chevron { flex-shrink: 0; transition: transform 140ms; }
  .chevron.folded { transform: rotate(-90deg); }
  .label { overflow: hidden; text-overflow: ellipsis; }
  .head { display: flex; align-items: center; gap: 4px; margin: 8px 0 2px; color: var(--dim); }
  .head .name { width: auto; flex: 1; min-width: 0; display: flex; align-items: center; justify-content: flex-start; gap: 8px; padding: 6px 0; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; text-align: left; font: inherit; font-weight: 700; color: inherit; background: none; border: 0; border-radius: 4px; cursor: pointer; }
  .selected .name { color: var(--foreground); }
  .icon { color: var(--accent); flex: 0 0 auto; min-width: 1em; display: inline-flex; align-items: center; justify-content: center; }
  .head .more { width: 32px; height: 32px; justify-content: center; flex: 0 0 auto; padding: 4px 10px; font: inherit; color: var(--dim); background: none; border: 0; border-radius: 4px; cursor: pointer; }
  .name:hover, .more:hover, .name:focus-visible, .more:focus-visible { color: var(--foreground); outline: none; background: color-mix(in srgb, var(--foreground) 8%, transparent); }
</style>
