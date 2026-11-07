<script lang="ts">
  let { kind }: { kind: 'user' | 'assistant' | 'code' | 'loading' | 'done' } = $props();
  const filled = $derived(kind === 'user' || kind === 'assistant');
</script>

<span class="icon" class:loading={kind === 'loading'} aria-label={kind === 'loading' ? 'Request running' : kind === 'done' ? 'Turn complete' : undefined} aria-hidden={kind !== 'loading' && kind !== 'done'}>
  <svg width="16" height="16" viewBox="0 0 24 24" fill={filled ? 'currentColor' : 'none'} stroke={filled ? 'none' : 'currentColor'} stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round">
    {#if kind === 'user'}
      <circle cx="12" cy="7" r="4" /><path d="M4 21a8 8 0 0 1 16 0Z" />
    {:else if kind === 'loading'}
      <path d="M21 12a9 9 0 1 1-9-9" />
    {:else if kind === 'done'}
      <path d="m4.5 12.5 5 5 10-11" />
    {:else if kind === 'code'}
      <path d="m8 7-5 5 5 5m8-10 5 5-5 5m-3-14-2 18" />
    {:else}
      <path d="m12 2 3 7 7 3-7 3-3 7-3-7-7-3 7-3Z" />
    {/if}
  </svg>
</span>

<style>
  .icon { position: absolute; left: -30px; top: 6px; display: inline-flex; color: #737373; }
  .loading { animation: spin 1s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
  @media (prefers-reduced-motion: reduce) { .loading { animation: none; } }
</style>
