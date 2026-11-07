<script lang="ts">
  import Modal from './Modal.svelte';
  let { efforts, active, fallback, onclose, onselect }: {
    efforts: string[]; active: string; fallback: string;
    onclose: () => void; onselect: (effort: string) => void;
  } = $props();
</script>

<Modal title="Reasoning effort" {onclose}>
  <button onclick={() => onselect('')} aria-current={active === '' ? 'true' : undefined}>Default <span>{fallback || 'Provider'}</span></button>
  {#each efforts as effort}
    <button onclick={() => onselect(effort)} aria-current={effort === active ? 'true' : undefined}>{effort} <span>{effort === active ? 'Current' : ''}</span></button>
  {/each}
</Modal>
