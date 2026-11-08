<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import Modal from './Modal.svelte';
  import SearchInput from './SearchInput.svelte';
  export type Panel = 'sessions' | 'actions' | 'about' | 'models' | 'efforts' | 'usage' | 'manager' | 'settings' | 'hotkeys' | 'instructions' | 'mcp';
  let { back = false, onclose, ondismiss, onopen }: {
    back?: boolean; onclose: () => void; ondismiss?: () => void; onopen: (name: Panel) => void;
  } = $props();
  let query = $state('');
  const commands: { panel: Panel; label: string; hint: string }[] = [
    { panel: 'actions', label: 'Session actions', hint: 'Compact, handoff, export' },
    { panel: 'manager', label: 'Models', hint: 'Visibility and utility model' },
    { panel: 'instructions', label: 'Instructions', hint: 'System prompt' },
    { panel: 'settings', label: 'Settings', hint: 'Toasts, code, compaction' },
    { panel: 'hotkeys', label: 'Hotkeys', hint: 'Keyboard' },
    { panel: 'usage', label: 'Context and usage', hint: 'Usage' },
    { panel: 'mcp', label: 'MCP servers', hint: 'Add and remove tools' },
    { panel: 'about', label: 'About', hint: 'Version and licenses' },
    { panel: 'sessions', label: 'Switch session', hint: 'Sessions and workspaces' },
  ];
  const matches = $derived(commands.filter(({ label, hint }) => `${label} ${hint}`.toLowerCase().includes(query.trim().toLowerCase())));
</script>

<Modal title="Commands" {back} {onclose} {ondismiss}>
  {#snippet top()}
    <SearchInput bind:value={query} label="Search commands" onkeydown={(event) => {
      if (event.key === 'Enter' && matches.length) { event.preventDefault(); onopen(matches[0].panel); }
    }} />
  {/snippet}
  {#each matches as command (command.panel)}
    <Button data-ui="ui-commands" onclick={() => onopen(command.panel)}>{command.label} <span>{command.hint}</span></Button>
  {:else}
    <p>No commands found</p>
  {/each}
</Modal>
