<script lang="ts">
  import type { Snippet } from 'svelte';
  let { hint, label = 'Status', children }: { hint: string; label?: string; children: Snippet } = $props();
</script>

<footer aria-label={label}>
  <div class="left">{@render children()}</div>
  <span class="hint">{hint}</span>
</footer>

<style>
  footer { display: flex; align-items: center; justify-content: space-between; gap: 18px; height: 44px; padding: 0 20px; border-top: 1px solid var(--border); background: var(--background); color: var(--dim); font: var(--font-size)/1.5 var(--mono); white-space: nowrap; }
  .left { display: flex; align-items: center; gap: 10px; min-width: 0; overflow: hidden; color: var(--muted); }
  .hint { color: var(--faint); flex-shrink: 0; }
  .left :global(button) { color: inherit; font: inherit; border: 0; background: none; padding: 0; cursor: pointer; transition: color 140ms; overflow: hidden; text-overflow: ellipsis; }
  .left :global(button:hover), .left :global(button:focus-visible) { color: var(--foreground); outline: none; }
  @media (max-width: 480px) { footer { gap: 10px; padding: 0 12px; } }
</style>
