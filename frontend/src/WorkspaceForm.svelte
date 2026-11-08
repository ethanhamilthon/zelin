<script lang="ts">
  import { untrack } from 'svelte';
  import Input from "./lib/ui/input/input.svelte";
  import Button from "./lib/ui/button/button.svelte";
  import WorkspaceIcon from './WorkspaceIcon.svelte';
  import { workspaceIcons, workspaceIcon } from './workspace-icons';
  import { call } from './bridge';
  import Modal from './Modal.svelte';
  import { DEFAULT_WORKSPACE, type AgentsFile, type Workspace } from './workspaces';
  let { workspace, agents, back = false, onclose, ondismiss, onsave, onagents, ondelete }: {
    workspace: Workspace | null; agents: AgentsFile | null; back?: boolean; onclose: () => void; ondismiss?: () => void;
    onsave: (name: string, icon: string, dir: string) => Promise<string | null>; onagents: () => void; ondelete: () => Promise<string | null>;
  } = $props();
  let name = $state(untrack(() => workspace?.name ?? '')), icon = $state(untrack(() => workspaceIcon(workspace?.icon ?? '').id)), dir = $state(untrack(() => workspace?.dir ?? '')), error = $state(''), saving = $state(false);
  const icons = workspaceIcons;
  let choosing = $state(false);
  async function choose() {
    choosing = true;
    try { const path = await call<string>('WorkspaceService.ChooseDirectory'); if (path) dir = path; }
    catch (failure) { error = String(failure); }
    finally { choosing = false; }
  }
  let first = $state<HTMLInputElement | null>(null);
  $effect(() => { first?.focus(); });
  async function submit() {
    if (saving) return;
    saving = true;
    error = (await onsave(name, icon, dir)) ?? '';
    saving = false;
  }
  let confirming = $state(false);
  async function remove() {
    if (!confirming) { confirming = true; return; }
    error = (await ondelete()) ?? '';
    confirming = false;
  }
  const enter = (event: KeyboardEvent) => { if (event.key === 'Enter') { event.preventDefault(); void submit(); } };
</script>

<Modal title={workspace ? 'Workspace settings' : 'New workspace'} {back} {onclose} {ondismiss}>
  <label>Name<Input bind:ref={first} bind:value={name} aria-label="Workspace name" onkeydown={enter} /></label>
  <fieldset><legend>Icon</legend><div class="icons">{#each icons as item}<Button data-ui="ui-workspaceform" variant="icon" class={{ chosen: icon === item.id }} title={item.name} aria-label={`Workspace icon ${item.name}`} aria-pressed={icon === item.id} onclick={() => icon = item.id}><WorkspaceIcon icon={item.id} /></Button>{/each}</div></fieldset>
  <div class="directory"><span>Working directory</span><Button data-ui="ui-workspaceform" class="path" onclick={() => void choose()} disabled={choosing} aria-label="Choose working directory" title={dir}>{dir || 'Choose folder…'}</Button></div>
  {#if workspace}<p class="note">New sessions start in this directory. Existing sessions keep their own directory until they call cd().</p>{/if}
  {#if error}<p role="alert">{error}</p>{/if}
  <div class="form-actions">
    {#if workspace && workspace.id !== DEFAULT_WORKSPACE}<Button data-ui="ui-workspaceform" class={["delete", { confirming }]} onclick={() => { void remove(); }} onblur={() => { confirming = false; }}>{confirming ? 'Confirm delete' : 'Delete'}</Button>{/if}
    <span class="spacer"></span>
    {#if workspace && agents}<Button data-ui="ui-workspaceform" onclick={onagents}>{agents.exists ? 'Open AGENTS.md' : 'Create AGENTS.md'}</Button>{/if}
    <Button data-ui="ui-workspaceform" class="primary" onclick={() => { void submit(); }} disabled={saving}>{workspace ? 'Save' : 'Create'}</Button>
  </div>
</Modal>

<style>
  fieldset { border: 0; padding: 0; margin: 0 0 12px; color: var(--muted); }
  :global(.icons) { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 4px; }
  :global(.icons) :global(button[data-ui="ui-workspaceform"]) { width: 40px; height: 40px; padding: 6px; justify-content: center; border-radius: 8px; }
  :global(.icons) :global(button[data-ui="ui-workspaceform"]):global(.chosen) { background: var(--selection); outline: 1px solid var(--accent); }
  .directory { display: block; margin: 0 0 10px; color: var(--muted); }
  :global(button[data-ui="ui-workspaceform"].path) { display: block; width: 100%; margin-top: 4px; padding: 10px 12px; min-height: 0; border: 1px solid var(--border); border-radius: 4px; color: var(--accent); text-align: left; white-space: normal; overflow-wrap: anywhere; text-decoration: none; }
  :global(button[data-ui="ui-workspaceform"].path):hover { border-color: var(--dim); }
  label { display: block; margin: 0 0 10px; color: var(--muted); }
  :global(input[data-slot="input"]) { display: block; width: 100%; margin-top: 4px; padding: 10px 12px; border: 1px solid var(--border); border-radius: 4px; background: var(--background); color: var(--foreground); font: inherit; outline: none; }
  :global(input[data-slot="input"]:focus) { border-color: var(--dim); }
  p { margin: 0 0 10px; color: var(--muted); }
  .note { color: var(--dim); }
  .form-actions { display: flex; align-items: center; gap: 8px; margin-top: 16px; }
  .spacer { flex: 1; }
  .form-actions :global(button[data-ui="ui-workspaceform"]) { width: auto; min-height: 0; padding: 8px 12px; border: 0; color: var(--muted); }
  .form-actions :global(button[data-ui="ui-workspaceform"].primary) { color: var(--accent); }
  .form-actions :global(button[data-ui="ui-workspaceform"].delete) { color: var(--error); }
  .form-actions :global(button[data-ui="ui-workspaceform"].confirming) { text-decoration: underline; }
</style>
