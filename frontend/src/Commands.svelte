<script lang="ts">
  import Modal from './Modal.svelte';
  import Hotkey from './Hotkey.svelte';
  import SearchInput from './SearchInput.svelte';
  import type { Action } from './hotkeys';
  export type Panel = 'sessions' | 'models' | 'efforts' | 'usage' | 'manager' | 'settings' | 'hotkeys' | 'instructions';
  let { hotkeys, onclose, onopen, onaction = () => {}, busy = false, hasMessages = false }: {
    onaction?: (action: 'compact' | 'handoff') => void; busy?: boolean; hasMessages?: boolean;
    hotkeys: Record<Action, string>; onclose: () => void; onopen: (name: Panel) => void;
  } = $props();
  let query = $state('');
  const commands: { panel: Panel | 'compact' | 'handoff'; label: string; hint: string }[] = [
    { panel: 'compact', label: 'Compact context', hint: 'Summarize in this session' },
    { panel: 'handoff', label: 'Handoff', hint: 'Continue in a new session' },
    { panel: 'manager', label: 'Models', hint: 'Visibility and titles' },
    { panel: 'instructions', label: 'Instructions', hint: 'System prompt' },
    { panel: 'settings', label: 'Settings', hint: 'Theme, toasts, code' },
    { panel: 'hotkeys', label: 'Hotkeys', hint: 'Keyboard' },
    { panel: 'usage', label: 'Context and usage', hint: 'Usage' },
    { panel: 'sessions', label: 'Switch session', hint: '' },
  ];
  function select(panel: Panel | 'compact' | 'handoff') {
    if (panel === 'compact' || panel === 'handoff') { if (!busy && hasMessages) onaction(panel); } else onopen(panel);
  }
  const matches = $derived(commands.filter(({ label, hint }) => `${label} ${hint}`.toLowerCase().includes(query.trim().toLowerCase())));
</script>

<Modal title="Commands" {onclose}>
  <SearchInput bind:value={query} label="Search commands" onkeydown={(event) => {
    if (event.key === 'Enter' && matches.length) { event.preventDefault(); select(matches[0].panel); }
  }} />
  {#each matches as command (command.panel)}
    <button disabled={(command.panel === 'compact' || command.panel === 'handoff') && (busy || !hasMessages)} onclick={() => select(command.panel)}>{command.label} <span>{#if command.panel === 'sessions'}<Hotkey value={hotkeys.sessions} />{:else}{command.hint}{/if}</span></button>
  {:else}
    <p>No commands found</p>
  {/each}
</Modal>
