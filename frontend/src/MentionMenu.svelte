<script lang="ts">
  import type { Entry } from './refs';
  let { items, active, x, y, loading = false, onpick, onhover }: {
    items: Entry[]; active: number; x: number; y: number; loading?: boolean; onpick: (entry: Entry) => void; onhover: (index: number) => void;
  } = $props();
</script>

<div class="menu" role="listbox" aria-label="File suggestions" style:left={`${x}px`} style:top={`${y + 4}px`}>
  {#each items as entry, index (entry.path)}
    <button role="option" aria-selected={index === active} class:active={index === active} onmousedown={(event) => { event.preventDefault(); onpick(entry); }} onmouseenter={() => onhover(index)}>
      <span>{entry.kind === 'dir' ? '▸ ' : ''}{entry.name}</span>
    </button>
  {:else}
    <p>{loading ? 'Searching…' : 'No matches'}</p>
  {/each}
</div>

<style>
  .menu { position: fixed; z-index: 20; min-width: 220px; max-width: min(420px, 90vw); padding: 4px; border: 1px solid var(--border); border-radius: 6px; background: color-mix(in srgb, var(--foreground) 5%, var(--background)); box-shadow: 0 10px 40px color-mix(in srgb, #000 40%, transparent); font: var(--font-size)/1.5 var(--mono); }
  button { display: block; width: 100%; padding: 4px 8px; border: 0; border-radius: 4px; background: none; color: var(--muted); font: inherit; text-align: left; cursor: pointer; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
  button.active { background: color-mix(in srgb, var(--foreground) 10%, transparent); color: var(--foreground); }
  p { margin: 0; padding: 4px 8px; color: var(--dim); }
</style>
