<script lang="ts">
  import { onMount } from 'svelte';
  import Modal from './Modal.svelte';
  import type { AgentsFile, Workspace } from './workspaces';
  let { workspace, agents, back = false, onclose, ondismiss, onsave, onagents }: {
    workspace: Workspace | null; agents: AgentsFile | null; back?: boolean; onclose: () => void; ondismiss?: () => void;
    onsave: (name: string, icon: string, dir: string) => Promise<string | null>; onagents: () => void;
  } = $props();
  let name = $state(''), icon = $state(''), dir = $state(''), error = $state(''), saving = $state(false);
  let first: HTMLInputElement;
  onMount(() => { name = workspace?.name ?? ''; icon = workspace?.icon ?? ''; dir = workspace?.dir ?? ''; first.focus(); });
  async function submit() {
    if (saving) return;
    saving = true;
    error = (await onsave(name, icon, dir)) ?? '';
    saving = false;
  }
  const enter = (event: KeyboardEvent) => { if (event.key === 'Enter') { event.preventDefault(); void submit(); } };
</script>

<Modal title={workspace ? 'Workspace settings' : 'New workspace'} {back} {onclose} {ondismiss}>
  <label>Name<input bind:this={first} bind:value={name} aria-label="Workspace name" onkeydown={enter} /></label>
  <label>Icon<input bind:value={icon} aria-label="Workspace icon" maxlength="16" placeholder="Emoji or text" onkeydown={enter} /></label>
  <label>Working directory<input bind:value={dir} aria-label="Working directory" placeholder="/path/to/project or ~/project" onkeydown={enter} /></label>
  {#if workspace}<p class="note">New sessions start in this directory. Existing sessions keep their own directory until they call cd().</p>{/if}
  {#if error}<p role="alert">{error}</p>{/if}
  <div class="row">
    {#if workspace && agents}<button onclick={onagents}>{agents.exists ? 'Open AGENTS.md' : 'Create AGENTS.md'}</button>{/if}
    <button onclick={() => { void submit(); }} disabled={saving}>{workspace ? 'Save' : 'Create'}</button>
  </div>
</Modal>

<style>
  label { display: block; margin: 0 0 10px; color: var(--muted); }
  input { display: block; width: 100%; margin-top: 4px; padding: 10px 12px; border: 1px solid var(--border); border-radius: 4px; background: var(--background); color: var(--foreground); font: inherit; outline: none; }
  input:focus { border-color: var(--dim); }
  p { margin: 0 0 10px; color: var(--muted); }
  .note { color: var(--dim); }
  .row { display: flex; justify-content: flex-end; gap: 8px; }
  .row button { width: auto; }
</style>
