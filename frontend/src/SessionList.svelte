<script lang="ts">
  import Modal from './Modal.svelte';
  import type { Session } from './sessions';
  let { sessions, active, busy, onclose, onselect, onnew }: {
    sessions: Session[]; active: string; busy: boolean;
    onclose: () => void; onselect: (id: string) => void; onnew: () => void;
  } = $props();
</script>

<Modal title="Sessions" {onclose}>
  <button onclick={onnew} disabled={busy}>New session <span>+</span></button>
  {#each sessions as session}
    <button onclick={() => onselect(session.id)} disabled={busy} aria-current={session.id === active ? 'true' : undefined}>
      <strong>{session.title}</strong><span>{session.id === active ? 'Current' : new Date(session.updated).toLocaleDateString('en-US')}</span>
    </button>
  {/each}
  <button onclick={onclose}>Close sessions <span>Esc</span></button>
</Modal>

<style>
  strong { font-weight: 400; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
</style>
