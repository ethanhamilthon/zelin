<script lang="ts">
  import { onMount } from 'svelte';
  let { value = $bindable(''), label, onkeydown }: { value?: string; label: string; onkeydown?: (event: KeyboardEvent) => void } = $props();
  let input: HTMLInputElement;
  onMount(() => input.focus());
</script>

<input bind:this={input} bind:value {onkeydown} placeholder={`${label}…`} aria-label={label} />

<style>
  input { width: 100%; margin: 0 0 8px; padding: 10px 12px; border: 1px solid var(--border); border-radius: 4px; background: var(--background); color: var(--foreground); font: inherit; outline: none; }
  input:focus { border-color: var(--dim); }
</style>
