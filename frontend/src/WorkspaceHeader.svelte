<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import Icon from './Icon.svelte';
  import WorkspaceIcon from './WorkspaceIcon.svelte';
  import type { Workspace } from './workspaces';
  let { workspace, selected = false, folded = false, onselect, onsettings, onnew }: {
    workspace: Workspace; selected?: boolean; folded?: boolean; onselect: () => void; onsettings: () => void; onnew: () => void;
  } = $props();
</script>

<div class="head" class:selected>
  <Button data-ui="ui-workspaceheader" variant="icon-text" class="name" onclick={onselect} aria-pressed={selected} aria-expanded={!folded} title={workspace.dir}>
    <span class="chevron" class:folded><Icon name="down" /></span>
    <span class="icon"><WorkspaceIcon icon={workspace.icon} /></span><span class="label">{workspace.name}</span>
  </Button>
  <Button data-ui="ui-workspaceheader" variant="icon" class="add" onclick={onnew} aria-label={`New session in ${workspace.name}`} title="New session"><Icon name="plus" /></Button>
  <Button data-ui="ui-workspaceheader" variant="icon" class="settings" onclick={onsettings} aria-label={`Settings for ${workspace.name}`} title="Workspace settings"><Icon name="dots" /></Button>
</div>

<style>
  .head { display: flex; align-items: center; gap: 2px; margin: 8px 0 2px; color: var(--dim); }
  .head.selected :global(.name) { color: var(--foreground); }
  .head :global(.name) { flex: 1; min-width: 0; justify-content: flex-start; padding: 6px 0; font-weight: 700; }
  .head :global(.add), .head :global(.settings) { flex: 0 0 auto; color: var(--dim); }
  .chevron { display: inline-flex; flex-shrink: 0; transition: transform 140ms; }
  .chevron.folded { transform: rotate(-90deg); }
  .label { overflow: hidden; text-overflow: ellipsis; }
  .icon { color: var(--accent); flex: 0 0 auto; display: inline-flex; align-items: center; }
</style>
