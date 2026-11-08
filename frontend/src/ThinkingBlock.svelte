<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import { oneLine } from './preview';
  let { text, folded: initial = false }: { text: string; folded?: boolean } = $props();
  let folded = $state(false);
  $effect.pre(() => { folded = initial; });
  const display = $derived(folded ? oneLine(text) : text);
</script>

<div class="thinking" aria-label="Thinking">
  <Button data-ui="ui-thinkingblock" variant="icon" class={["think-fold", { folded }]} onclick={() => { folded = !folded; }} aria-expanded={!folded} aria-label={folded ? 'Unfold thinking' : 'Fold thinking'}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round">
      <path d="M12 5a3 3 0 1 0-6 .1 4 4 0 0 0-2.5 5.8 4 4 0 0 0 .5 6.6A4 4 0 1 0 12 18Z" />
      <path d="M12 5a3 3 0 1 1 6 .1 4 4 0 0 1 2.5 5.8 4 4 0 0 1-.5 6.6A4 4 0 1 1 12 18Z" />
      <path d="M15 13a4.5 4.5 0 0 1-3-4 4.5 4.5 0 0 1-3 4" />
    </svg>
  </Button>
  <Button data-ui="ui-thinkingblock" class={["body", { folded }]} onclick={() => { if (!window.getSelection()?.toString()) folded = !folded; }} aria-expanded={!folded}>
    <p>{display}</p>
  </Button>
</div>

<style>
  :global(.thinking) { position: relative; margin: 0 0 8px; color: var(--dim); }
  :global(.think-fold) { position: absolute; left: -30px; top: 6px; display: inline-flex; min-height: 0; padding: 0; border: 0; background: none; color: var(--violet); cursor: pointer; transition: color 140ms, opacity 140ms; }
  :global(.think-fold.folded) { opacity: 0.65; }
  :global(.think-fold):hover, :global(.think-fold):focus-visible { color: var(--foreground); opacity: 1; outline: none; }
  :global(button[data-ui="ui-thinkingblock"].body) { display: block; width: 100%; min-height: 0; padding: 0; border: 0; background: none; color: inherit; font: inherit; text-align: left; cursor: pointer; transition: color 160ms; }
  :global(button[data-ui="ui-thinkingblock"].body):hover { color: var(--muted); }
  :global(button[data-ui="ui-thinkingblock"].body):focus-visible { outline: 1px solid var(--dim); outline-offset: 5px; }
  p { margin: 0; white-space: pre-wrap; overflow-wrap: anywhere; line-height: 1.8; }
  :global(.body.folded) p { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
