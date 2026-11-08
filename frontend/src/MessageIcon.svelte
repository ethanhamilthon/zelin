<script lang="ts">
  import PendingDots from './PendingDots.svelte';
  let { kind }: { kind: 'assistant' | 'loading' | 'done' | 'thinking' | 'error' } = $props();
  const filled = $derived(kind === 'assistant');
</script>

<span class={["message-icon", "icon-slot", { loading: kind === 'loading', error: kind === 'error' }]} aria-label={kind === 'loading' ? 'Request running' : kind === 'done' ? 'Turn complete' : undefined} aria-hidden={kind !== 'loading' && kind !== 'done'}>
  {#if kind === 'loading'}<PendingDots label="Agent working" />{:else}
  <svg width="16" height="16" viewBox="0 0 24 24" fill={filled ? 'currentColor' : 'none'} stroke={filled ? 'none' : 'currentColor'} stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round">
    {#if kind === 'error'}
      <circle cx="12" cy="12" r="9" /><path d="M12 7.5v5" /><path d="M12 16.5h.01" />
    {:else if kind === 'done'}
      <path d="m4.5 12.5 5 5 10-11" />
    {:else if kind === 'thinking'}
      <path d="M12 5a3 3 0 1 0-6 .1 4 4 0 0 0-2.5 5.8 4 4 0 0 0 .5 6.6A4 4 0 1 0 12 18Z" />
      <path d="M12 5a3 3 0 1 1 6 .1 4 4 0 0 1 2.5 5.8 4 4 0 0 1-.5 6.6A4 4 0 1 1 12 18Z" />
      <path d="M15 13a4.5 4.5 0 0 1-3-4 4.5 4.5 0 0 1-3 4" />
    {:else}
      <path d="m12 2 3 7 7 3-7 3-3 7-3-7-7-3 7-3Z" />
    {/if}
  </svg>
  {/if}
</span>

<style>
  .message-icon { color: var(--dim); }
  .error { color: var(--error); }

</style>
