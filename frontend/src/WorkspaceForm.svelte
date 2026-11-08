<script lang="ts">
  import { untrack } from 'svelte';
  import Input from "./lib/ui/input/input.svelte";
  import Button from "./lib/ui/button/button.svelte";
  import WorkspaceIcon from './WorkspaceIcon.svelte';
  import { workspaceIcons, workspaceIcon } from './workspace-icons';
  import { call } from './bridge';
  import Modal from './Modal.svelte';
  import type { AgentsFile, Workspace } from './workspaces';
  let { workspace, agents, back = false, onclose, ondismiss, onsave, onagents }: {
    workspace: Workspace | null; agents: AgentsFile | null; back?: boolean; onclose: () => void; ondismiss?: () => void;
    onsave: (name: string, icon: string, dir: string) => Promise<string | null>; onagents: () => void;
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
  const enter = (event: KeyboardEvent) => { if (event.key === 'Enter') { event.preventDefault(); void submit(); } };
</script>

<Modal title={workspace ? 'Workspace settings' : 'New workspace'} {back} {onclose} {ondismiss}>
  <label>Name<Input bind:ref={first} bind:value={name} aria-label="Workspace name" onkeydown={enter} /></label>
  <fieldset><legend>Icon</legend><div class="icons">{#each icons as item}<Button data-ui="ui-workspaceform" variant="icon" class={{ chosen: icon === item.id }} title={item.name} aria-label={`Workspace icon ${item.name}`} aria-pressed={icon === item.id} onclick={() => icon = item.id}><WorkspaceIcon icon={item.id} /></Button>{/each}</div></fieldset>
  <div class="directory"><span>Working directory</span><Button data-ui="ui-workspaceform" onclick={() => void choose()} disabled={choosing} aria-label="Choose working directory" title={dir}>{dir || 'Choose folder…'}</Button></div>
  {#if workspace}<p class="note">New sessions start in this directory. Existing sessions keep their own directory until they call cd().</p>{/if}
  {#if error}<p role="alert">{error}</p>{/if}
  <div class="row">
    {#if workspace && agents}<Button data-ui="ui-workspaceform" onclick={onagents}>{agents.exists ? 'Open AGENTS.md' : 'Create AGENTS.md'}</Button>{/if}
    <Button data-ui="ui-workspaceform" onclick={() => { void submit(); }} disabled={saving}>{workspace ? 'Save' : 'Create'}</Button>
  </div>
</Modal>

<style>
  fieldset { border: 0; padding: 0; margin: 0 0 12px; color: var(--muted); }
  :global(.icons) { display: flex; flex-wrap: wrap; gap: 6px; }
  :global(.icons) :global(button[data-ui="ui-workspaceform"]) { width: 40px; height: 40px; padding: 6px; justify-content: center; }
  :global(.icons) :global(button[data-ui="ui-workspaceform"]):global(.chosen) { background: var(--selection); outline: 1px solid var(--accent); }
  :global(.directory) { margin-bottom: 12px; color: var(--muted); }
  :global(.directory) :global(button[data-ui="ui-workspaceform"]) { margin-top: 4px; border: 1px solid var(--border); overflow-wrap: anywhere; color: var(--accent); }
  label { display: block; margin: 0 0 10px; color: var(--muted); }
  :global(input[data-slot="input"]) { display: block; width: 100%; margin-top: 4px; padding: 10px 12px; border: 1px solid var(--border); border-radius: 4px; background: var(--background); color: var(--foreground); font: inherit; outline: none; }
  :global(input[data-slot="input"]:focus) { border-color: var(--dim); }
  p { margin: 0 0 10px; color: var(--muted); }
  :global(.note) { color: var(--dim); }
  :global(.row) { display: flex; justify-content: flex-end; gap: 8px; }
  :global(.row) :global(button[data-ui="ui-workspaceform"]) { width: auto; }
</style>
