<script lang="ts">
  import { untrack } from 'svelte';
  import Input from "./lib/ui/input/input.svelte";
  import Modal from './Modal.svelte';
  let { title, back = false, onclose, ondismiss, onsave }: {
    title: string; back?: boolean; onclose: () => void; ondismiss?: () => void; onsave: (title: string) => Promise<string | null>;
  } = $props();
  let value = $state(untrack(() => title)), error = $state(''), saving = $state(false);
  let input = $state<HTMLInputElement | null>(null);
  $effect(() => { input?.focus(); input?.select(); });
  async function submit() {
    if (saving) return;
    saving = true;
    error = (await onsave(value)) ?? '';
    saving = false;
  }
</script>

<Modal title="Rename session" {back} {onclose} {ondismiss}>
  <Input bind:ref={input} bind:value aria-label="Session title" disabled={saving}
    onkeydown={(event) => { if (event.key === 'Enter') { event.preventDefault(); void submit(); } }} />
  {#if error}<p role="alert">{error}</p>{/if}
</Modal>

<style>
  :global(input[data-slot="input"]) { width: 100%; padding: 10px 12px; border: 1px solid var(--border); border-radius: 4px; background: var(--background); color: var(--foreground); font: inherit; outline: none; }
  :global(input[data-slot="input"]:focus) { border-color: var(--dim); }
  p { margin: 8px 0 0; color: var(--muted); }
</style>
