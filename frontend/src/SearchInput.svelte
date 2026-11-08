<script lang="ts">
  import Input from './lib/ui/input/input.svelte';
  import { onMount } from 'svelte';
  let { value = $bindable(''), label, onkeydown }: { value?: string; label: string; onkeydown?: (event: KeyboardEvent) => void } = $props();
  let input = $state<HTMLInputElement>(null!);
  onMount(() => input.focus());
</script>

<Input class="search-input" bind:ref={input} bind:value {onkeydown} placeholder={`${label}…`} aria-label={label} />

<style>
  :global(input.search-input) { width: 100%; margin: 0 0 8px; padding: 10px 12px; border: 1px solid var(--border); border-radius: 4px; background: var(--background); color: var(--foreground); font: inherit; outline: none; }
  :global(input.search-input:focus) { border-color: var(--dim); }
</style>
