<script lang="ts">
  import Modal from './Modal.svelte';
  let { efforts, active, fallback, error = '', onclose, onselect }: {
    efforts: string[]; active: string; fallback: string; error?: string;
    onclose: () => void; onselect: (effort: string) => void;
  } = $props();
</script>

<Modal title="Reasoning effort" {onclose}>
  <button onclick={() => onselect('')} aria-current={active === '' ? 'true' : undefined}>Default <span>{fallback || 'Provider default'}</span></button>
  {#each efforts as effort}
    <button onclick={() => onselect(effort)} aria-current={effort === active ? 'true' : undefined}>{effort} <span>{effort === active ? 'Current' : ''}</span></button>
  {/each}
  {#if error}<p role="alert">{error}</p>{/if}
</Modal>

<style>
  p { margin: 0; padding: 12px; color: #aaa; }
</style>
