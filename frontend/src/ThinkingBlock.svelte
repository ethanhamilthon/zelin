<script lang="ts">
  import { oneLine } from './preview';
  let { text, folded: initial = false }: { text: string; folded?: boolean } = $props();
  let folded = $state(false);
  $effect.pre(() => { folded = initial; });
  const display = $derived(folded ? oneLine(text) : text);
</script>

<div class="thinking" aria-label="Thinking">
  <button class="think-fold" class:folded onclick={() => { folded = !folded; }} aria-expanded={!folded} aria-label={folded ? 'Unfold thinking' : 'Fold thinking'}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round">
      <path d="M12 5a3 3 0 1 0-6 .1 4 4 0 0 0-2.5 5.8 4 4 0 0 0 .5 6.6A4 4 0 1 0 12 18Z" />
      <path d="M12 5a3 3 0 1 1 6 .1 4 4 0 0 1 2.5 5.8 4 4 0 0 1-.5 6.6A4 4 0 1 1 12 18Z" />
      <path d="M15 13a4.5 4.5 0 0 1-3-4 4.5 4.5 0 0 1-3 4" />
    </svg>
  </button>
  <button class="body" class:folded onclick={() => { if (!window.getSelection()?.toString()) folded = !folded; }} aria-expanded={!folded}>
    <p>{display}</p>
  </button>
</div>

<style>
  .thinking { position: relative; margin: 0 0 8px; color: var(--dim); }
  .think-fold { position: absolute; left: -30px; top: 6px; display: inline-flex; padding: 0; border: 0; background: none; color: var(--dim); cursor: pointer; transition: color 140ms, opacity 140ms; }
  .think-fold.folded { opacity: 0.65; }
  .think-fold:hover, .think-fold:focus-visible { color: var(--foreground); opacity: 1; outline: none; }
  .body { display: block; width: 100%; padding: 0; border: 0; background: none; color: inherit; font: inherit; text-align: left; cursor: pointer; transition: color 160ms; }
  .body:hover { color: var(--muted); }
  .body:focus-visible { outline: 1px solid var(--dim); outline-offset: 5px; }
  p { margin: 0; white-space: pre-wrap; overflow-wrap: anywhere; line-height: 1.8; }
  .body.folded p { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
