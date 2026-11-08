<script lang="ts">
  import type { Workspace } from './workspaces';
  let { workspace, selected = false, onselect, onsettings }: {
    workspace: Workspace; selected?: boolean; onselect: () => void; onsettings: () => void;
  } = $props();
</script>

<div class="head" class:selected>
  <button class="name" onclick={onselect} aria-pressed={selected} title={workspace.dir}>
    <span class="icon" aria-hidden="true">{workspace.icon || '▸'}</span>{workspace.name}
  </button>
  <button class="more" onclick={onsettings} aria-label={`Settings for ${workspace.name}`} title="Workspace settings">⋯</button>
</div>

<style>
  .head { display: flex; align-items: center; gap: 4px; margin: 8px 0 2px; color: var(--dim); }
  .name { flex: 1; min-width: 0; display: flex; align-items: center; gap: 8px; padding: 6px 12px; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; text-align: left; font: inherit; font-weight: 700; color: inherit; background: none; border: 0; border-radius: 4px; cursor: pointer; }
  .selected .name { color: var(--foreground); }
  .icon { flex: 0 0 auto; min-width: 1em; text-align: center; }
  .more { flex: 0 0 auto; padding: 4px 10px; font: inherit; color: var(--dim); background: none; border: 0; border-radius: 4px; cursor: pointer; }
  .name:hover, .more:hover, .name:focus-visible, .more:focus-visible { color: var(--foreground); outline: none; background: color-mix(in srgb, var(--foreground) 8%, transparent); }
</style>
