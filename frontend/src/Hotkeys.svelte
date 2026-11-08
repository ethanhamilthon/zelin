<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import Modal from './Modal.svelte';
  import Hotkey from './Hotkey.svelte';
  import { actions, combo, defaultHotkeys } from './hotkeys';
  import type { Action } from './hotkeys';
  let { hotkeys, back = false, onclose, ondismiss, onbind }: { hotkeys: Record<Action, string>; back?: boolean; onclose: () => void; ondismiss?: () => void; onbind: (action: Action, hotkey: string) => void } = $props();
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

<Modal title="Hotkeys" hint={recording ? 'Press keys' : undefined} {back} {onclose} {ondismiss}>
  {#each actions as action (action.id)}
    <Button data-ui="ui-hotkeys" onclick={() => { recording = action.id; }} aria-pressed={recording === action.id}>
      {action.label} <span class={{ recording: recording === action.id }}>{#if recording === action.id}Press keys…{:else}<Hotkey value={hotkeys[action.id]} />{/if}</span>
    </Button>
  {/each}
  <Button data-ui="ui-hotkeys" onclick={() => actions.forEach(({ id }) => onbind(id, defaultHotkeys[id]))}>Reset to defaults <span>Defaults</span></Button>
</Modal>

<style>
  :global(.recording) { color: var(--accent) !important; }
</style>
