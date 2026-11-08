<script lang="ts">
  import type { Snippet } from 'svelte';
  let { hint = '', label = 'Status', children, actions, center }: { hint?: string; label?: string; children: Snippet; actions?: Snippet; center?: Snippet } = $props();
</script>

<footer aria-label={label} class:centered={!!center}>
  <div class="left">{@render children()}</div>
  {#if center}<div class="center">{@render center()}</div>{/if}
  <div class="right">{#if actions}{@render actions()}{:else}<span class="hint">{hint}</span>{/if}</div>
</footer>

<style>
  footer { display: flex; align-items: center; justify-content: space-between; gap: 16px; width: 100%; min-width: 0; height: 44px; padding: 0 20px; border-top: 1px solid var(--border); background: var(--background); color: var(--dim); font: var(--font-size)/1.5 var(--mono); white-space: nowrap; }
  .left { display: flex; align-items: center; gap: 10px; flex: 1; min-width: 0; overflow: hidden; color: var(--dim); }
  .left > :global(span) { min-width: 0; overflow: hidden; text-overflow: ellipsis; }
  .right { display: flex; gap: 16px; align-items: center; flex-shrink: 0; }
  .centered { display: grid; grid-template-columns: minmax(0, 1fr) auto minmax(max-content, 1fr); }
  .centered .right { justify-self: end; }
  .center { display: flex; align-items: center; justify-content: center; min-width: 0; }
  .hint { color: var(--faint); flex-shrink: 0; }
  footer :global(button) { color: var(--dim); font: inherit; border: 0; background: none; padding: 0; cursor: pointer; transition: color 140ms; overflow: hidden; text-overflow: ellipsis; }
  .right :global(button) { flex-shrink: 0; overflow: visible; }
  footer :global(button:hover), footer :global(button:focus-visible) { color: var(--foreground); outline: none; }
  @media (max-width: 560px) { footer { gap: 10px; padding: 0 12px; } .right { gap: 12px; } .centered { grid-template-columns: minmax(0, 1fr) auto auto; } }
</style>
