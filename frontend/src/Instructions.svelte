<script lang="ts">
  import Modal from './Modal.svelte';
  import { documents } from './documents';
  let { onclose, onopen }: { onclose: () => void; onopen: (name: string) => void } = $props();
</script>

<Modal title="Instructions" {onclose}>
  {#each documents as document (document.name)}
    <button onclick={() => onopen(document.name)}>{document.title} <span>{document.file}</span></button>
  {/each}
</Modal>
