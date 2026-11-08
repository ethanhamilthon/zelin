<script lang="ts">
  let { working = 0, done = 0 }: { working?: number; done?: number } = $props();
  const label = $derived([working && `${working} working`, done && `${done} done`].filter(Boolean).join(', '));
</script>

{#if working > 0 || done > 0}
  <div class="activity" role="status" aria-label={label}>
    {#if working > 0}
      <span class="working"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" aria-hidden="true"><path d="M21 12a9 9 0 1 1-6.2-8.55" /></svg>{working} Working</span>
    {/if}
    {#if working > 0 && done > 0}<span class="sep" aria-hidden="true">|</span>{/if}
    {#if done > 0}
      <span class="done"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="m4.5 12.5 5 5 10-11" /></svg>{done} Done</span>
    {/if}
  </div>
{/if}

<style>
  .activity { display: inline-flex; align-items: center; gap: 10px; min-width: 0; white-space: nowrap; color: var(--dim); }
  span { display: inline-flex; align-items: center; gap: 6px; }
  .sep { color: var(--faint); }
  .working svg { animation: spin 1.2s linear infinite; }
  .done { color: var(--dim); }
  .done svg { color: var(--success); }
  .working svg { color: var(--warning); }
  @keyframes spin { to { transform: rotate(360deg); } }
  @media (prefers-reduced-motion: reduce) { .working svg { animation: none; } }
</style>
