<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import Modal from './Modal.svelte';
  import type { SessionAction } from './session-actions';
  import { listTasks } from './tasks';
  let { busy = false, hasMessages = false, session = '', back = false, onclose, ondismiss, onaction }: {
    busy?: boolean; hasMessages?: boolean; session?: string; back?: boolean;
    onclose: () => void; ondismiss?: () => void; onaction: (action: SessionAction) => void;
  } = $props();
  let tasks = $state(0);
  $effect(() => { if (session) void listTasks(session).then((list) => { tasks = list.length; }).catch(() => {}); });
  const items = $derived<{ id: SessionAction; label: string; hint: string; off: boolean }[]>([
    ...(tasks ? [{ id: 'tasks' as const, label: `Tasks (${tasks})`, hint: 'Running in the background', off: false }] : []),
    { id: 'compact', label: 'Compact context', hint: 'Summarize in this session', off: busy || !hasMessages },
    { id: 'handoff', label: 'Handoff', hint: 'Continue in a new session', off: busy || !hasMessages },
    { id: 'export', label: 'Export session', hint: 'Markdown file', off: !hasMessages },
    { id: 'stop', label: 'Stop Session', hint: 'Cancel the running turn', off: !busy },
    { id: 'rename', label: 'Rename Session', hint: 'Change the title', off: false },
  ]);
</script>

<Modal title="Session actions" {back} {onclose} {ondismiss}>
  {#each items as item (item.id)}
    <Button data-ui="ui-sessionactions" disabled={item.off} onclick={() => onaction(item.id)}>{item.label} <span>{item.hint}</span></Button>
  {/each}
</Modal>
