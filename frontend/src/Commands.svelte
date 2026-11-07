<script lang="ts">
  import Modal from './Modal.svelte';
  let { onclose, onfocus, onlatest, onnew, onsessions, busy = false }: {
    onclose: () => void; onfocus: () => void; onlatest: () => void;
    onnew: () => void; onsessions: () => void; busy?: boolean;
  } = $props();
</script>

<Modal title="Commands" {onclose}>
  <button onclick={onnew} disabled={busy}>New session <span>Create</span></button>
  <button onclick={onsessions} disabled={busy}>Switch session <span>Ctrl + S</span></button>
  <button onclick={onfocus}>Focus message editor <span>Write</span></button>
  <button onclick={onlatest}>Scroll to latest message <span>Navigate</span></button>
  <button onclick={onclose}>Close commands <span>Esc</span></button>
</Modal>
