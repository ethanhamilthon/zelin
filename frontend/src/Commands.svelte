<script lang="ts">
  import Modal from './Modal.svelte';
  import Hotkey from './Hotkey.svelte';
  import SearchInput from './SearchInput.svelte';
  import type { Action } from './hotkeys';
  export type Panel = 'sessions' | 'actions' | 'about' | 'models' | 'efforts' | 'usage' | 'manager' | 'settings' | 'hotkeys' | 'instructions';
  let { hotkeys, back = false, onclose, ondismiss, onopen }: {
    back?: boolean; hotkeys: Record<Action, string>; onclose: () => void; ondismiss?: () => void; onopen: (name: Panel) => void;
  } = $props();
  let query = $state('');
  const commands: { panel: Panel; label: string; hint: string; key?: Action }[] = [
    { panel: 'actions', label: 'Session actions', hint: '', key: 'actions' },
    { panel: 'manager', label: 'Models', hint: 'Visibility and utility model' },
    { panel: 'instructions', label: 'Instructions', hint: 'System prompt' },
    { panel: 'settings', label: 'Settings', hint: 'Theme, toasts, code' },
    { panel: 'hotkeys', label: 'Hotkeys', hint: 'Keyboard' },
    { panel: 'usage', label: 'Context and usage', hint: 'Usage' },
    { panel: 'about', label: 'About', hint: 'Version and licenses' },
    { panel: 'sessions', label: 'Switch session', hint: '', key: 'sessions' },
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
    <button onclick={() => onopen(command.panel)}>{command.label} <span>{#if command.key}<Hotkey value={hotkeys[command.key]} />{:else}{command.hint}{/if}</span></button>
  {:else}
    <p>No commands found</p>
  {/each}
</Modal>
