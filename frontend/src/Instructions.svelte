<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import Modal from './Modal.svelte';
  import { documents } from './documents';
  let { back = false, onclose, ondismiss, onopen }: { back?: boolean; onclose: () => void; ondismiss?: () => void; onopen: (name: string) => void } = $props();
</script>

<Modal title="Instructions" {back} {onclose} {ondismiss}>
  {#each documents as document (document.name)}
    <Button data-ui="ui-instructions" onclick={() => onopen(document.name)}>{document.title} <span>{document.file}</span></Button>
  {/each}
</Modal>
