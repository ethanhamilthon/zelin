<script lang="ts">
  import Modal from './Modal.svelte';
  import Hotkey from './Hotkey.svelte';
  import { actions, combo, defaultHotkeys } from './hotkeys';
  import type { Action } from './hotkeys';
  let { hotkeys, onclose, onbind }: { hotkeys: Record<Action, string>; onclose: () => void; onbind: (action: Action, hotkey: string) => void } = $props();
  let recording = $state<Action | null>(null);

  function capture(event: KeyboardEvent) {
    if (!recording) return;
    event.preventDefault();
    event.stopImmediatePropagation();
    if (event.key === 'Escape') { recording = null; return; }
    const hotkey = combo(event);
    if (!hotkey) return;
    onbind(recording, hotkey);
    recording = null;
  }
</script>

<svelte:window onkeydowncapture={capture} />

<Modal title="Hotkeys" hint={recording ? 'Press keys, Cancel Esc' : 'Close Esc'} {onclose}>
  {#each actions as action (action.id)}
    <button onclick={() => { recording = action.id; }} aria-pressed={recording === action.id}>
      {action.label} <span class:recording={recording === action.id}>{#if recording === action.id}Press keys…{:else}<Hotkey value={hotkeys[action.id]} />{/if}</span>
    </button>
  {/each}
  <button onclick={() => actions.forEach(({ id }) => onbind(id, defaultHotkeys[id]))}>Reset to defaults <span>Defaults</span></button>
</Modal>

<style>
  .recording { color: var(--accent) !important; }
</style>
