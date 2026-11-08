<script lang="ts">
  import Modal from './Modal.svelte';
  import Hotkey from './Hotkey.svelte';
  import type { Session } from './sessions';
  import type { SessionStatus } from './session-controller.svelte';
  let { sessions, active, status, hotkey = 'Ctrl+N', onclose, onselect, onnew }: {
    sessions: Session[]; active: string; hotkey?: string; status: (id: string) => SessionStatus;
    onclose: () => void; onselect: (id: string) => void; onnew: () => void;
  } = $props();
  const labels: Record<SessionStatus, string> = { running: 'Running', unread: 'Finished', error: 'Failed', '': '' };
</script>

<Modal title="Sessions" {onclose}>
  {#snippet actions()}<button onclick={onnew}>New session <Hotkey value={hotkey} /></button>{/snippet}
  {#each sessions as session (session.id)}
    {@const state = status(session.id)}
    <button onclick={() => onselect(session.id)} aria-current={session.id === active ? 'true' : undefined}>
      <strong><i class={state} aria-hidden="true"></i>{session.title}</strong>
      <span>{labels[state] || (session.id === active ? 'Current' : new Date(session.updated).toLocaleDateString('en-US'))}</span>
    </button>
  {/each}
</Modal>

<style>
  strong { display: flex; align-items: center; gap: 10px; min-width: 0; font-weight: 400; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
  i { width: 7px; height: 7px; flex: 0 0 7px; border-radius: 50%; background: var(--border); }
  i.running { background: var(--accent); animation: pulse 1.2s ease-in-out infinite; }
  i.unread { background: var(--accent); }
  i.error { background: #d16969; }
  @keyframes pulse { 50% { opacity: 0.3; } }
</style>
