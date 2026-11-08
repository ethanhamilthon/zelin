<script lang="ts">
  import Modal from './Modal.svelte';
  import Hotkey from './Hotkey.svelte';
  import type { SessionAction } from './session-actions';
  let { busy = false, hasMessages = false, back = false, onclose, ondismiss, onaction }: {
    busy?: boolean; hasMessages?: boolean; back?: boolean;
    onclose: () => void; ondismiss?: () => void; onaction: (action: SessionAction) => void;
  } = $props();
  const items = $derived<{ id: SessionAction; label: string; hint: string; off: boolean }[]>([
    { id: 'compact', label: 'Compact context', hint: 'Summarize in this session', off: busy || !hasMessages },
    { id: 'handoff', label: 'Handoff', hint: 'Continue in a new session', off: busy || !hasMessages },
    { id: 'export', label: 'Export session', hint: 'Markdown file', off: !hasMessages },
    { id: 'stop', label: 'Stop Session', hint: 'Cancel the running turn', off: !busy },
    { id: 'rename', label: 'Rename Session', hint: 'Change the title', off: false },
  ]);
</script>

<Modal title="Session actions" {back} {onclose} {ondismiss}>
  {#each items as item (item.id)}
    <button disabled={item.off} onclick={() => onaction(item.id)}>{item.label} <span>{item.hint}</span></button>
  {/each}
</Modal>
