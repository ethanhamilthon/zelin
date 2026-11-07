<script lang="ts">
  import Modal from './Modal.svelte';
  import SearchInput from './SearchInput.svelte';
  import type { Action } from './hotkeys';
  export type Panel = 'sessions' | 'models' | 'efforts' | 'usage' | 'manager' | 'settings' | 'hotkeys' | 'instructions';
  let { hotkeys, onclose, onopen }: {
    hotkeys: Record<Action, string>; onclose: () => void; onopen: (name: Panel) => void;
  } = $props();
  let query = $state('');
  const commands: { panel: Panel; label: string; hint: string }[] = [
    { panel: 'manager', label: 'Models', hint: 'Visibility and titles' },
    { panel: 'instructions', label: 'Instructions', hint: 'System prompt' },
    { panel: 'settings', label: 'Settings', hint: 'Theme, toasts, code' },
    { panel: 'hotkeys', label: 'Hotkeys', hint: 'Keyboard' },
    { panel: 'usage', label: 'Context and usage', hint: 'Usage' },
    { panel: 'sessions', label: 'Switch session', hint: '' },
  ];
  const matches = $derived(commands.filter(({ label, hint }) => `${label} ${hint}`.toLowerCase().includes(query.trim().toLowerCase())));
</script>

<Modal title="Commands" {onclose}>
  <SearchInput bind:value={query} label="Search commands" onkeydown={(event) => {
    if (event.key === 'Enter' && matches.length) { event.preventDefault(); onopen(matches[0].panel); }
  }} />
  {#each matches as command (command.panel)}
    <button onclick={() => onopen(command.panel)}>{command.label} <span>{command.panel === 'sessions' ? hotkeys.sessions : command.hint}</span></button>
  {:else}
    <p>No commands found</p>
  {/each}
</Modal>
