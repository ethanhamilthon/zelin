<script lang="ts">
  import { effortName } from './models';
  import Button from "./lib/ui/button/button.svelte";
  import Modal from './Modal.svelte';
  let { efforts, active, error = '', back = false, onclose, ondismiss, onselect }: {
    efforts: string[]; active: string; error?: string; back?: boolean;
    onclose: () => void; ondismiss?: () => void; onselect: (effort: string) => void;
  } = $props();
</script>

<Modal title="Reasoning effort" {back} {onclose} {ondismiss}>
  {#each efforts as effort}
    <Button data-ui="ui-effortlist" onclick={() => onselect(effort)} aria-current={effort === active ? 'true' : undefined}>{effortName(effort)} <span>{effort === active ? 'Current' : ''}</span></Button>
  {/each}
  {#if error}<p role="alert">{error}</p>{/if}
</Modal>

<style>
  p { margin: 0; padding: 12px; color: var(--muted); }
</style>
