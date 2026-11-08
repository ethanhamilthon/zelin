<script lang="ts">
  import Modal from './Modal.svelte';
  let { efforts, active, error = '', back = false, onclose, ondismiss, onselect }: {
    efforts: string[]; active: string; error?: string; back?: boolean;
    onclose: () => void; ondismiss?: () => void; onselect: (effort: string) => void;
  } = $props();
</script>

<Modal title="Reasoning effort" {back} {onclose} {ondismiss}>
  {#each efforts as effort}
    <button onclick={() => onselect(effort)} aria-current={effort === active ? 'true' : undefined}>{effort} <span>{effort === active ? 'Current' : ''}</span></button>
  {/each}
  {#if error}<p role="alert">{error}</p>{/if}
</Modal>

<style>
  p { margin: 0; padding: 12px; color: var(--muted); }
</style>
