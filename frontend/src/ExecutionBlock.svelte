<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import type { Execution } from './agent-types';
  import { highlightJavaScript } from './highlight';
  import { oneLine } from './preview';
  import { fade } from 'svelte/transition';
  let { execution, folded: initial = false }: { execution: Execution; folded?: boolean } = $props();
  let expanded = $state(false);
  let folded = $state(false);
  $effect.pre(() => { folded = initial; });
  const code = $derived(folded ? oneLine(execution.code) : execution.code);
</script>

<div class="execution" aria-label="JavaScript execution" aria-busy={execution.state === 'running'}>
  <Button data-ui="ui-executionblock" class={["fold", { folded }]} onclick={() => { folded = !folded; }} aria-expanded={!folded} aria-label={folded ? 'Unfold code' : 'Fold code'}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6" /></svg>
  </Button>
  <Button data-ui="ui-executionblock" class={["source", { folded }]} onclick={() => { if (!window.getSelection()?.toString()) expanded = !expanded; }} aria-expanded={expanded && !folded} aria-label="Toggle execution output">
    <pre><code>{@html highlightJavaScript(code)}</code></pre>
  </Button>
  {#if !folded && expanded && execution.output}
    <pre class="output" transition:fade={{ duration: 120 }}>{execution.output}</pre>
  {/if}
</div>

<style>
  :global(.execution) { position: relative; margin: 0 0 8px; }
  :global(.fold) { position: absolute; left: -30px; top: 6px; display: inline-flex; padding: 0; border: 0; background: none; color: var(--accent); cursor: pointer; transition: transform 160ms, color 140ms; }
  :global(.fold):hover, :global(.fold):focus-visible { color: var(--foreground); outline: none; }
  :global(.fold.folded) { transform: rotate(-90deg); }
  :global(.source) { display: block; width: 100%; text-align: left; color: var(--dim); background: none; border: 0; padding: 0; cursor: pointer; font: inherit; transition: color 160ms; }
  :global(.source.folded) pre { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  :global(.source):hover { color: var(--muted); }
  :global(.source):focus-visible { outline: 1px solid var(--dim); outline-offset: 5px; }
  pre { margin: 0; overflow-x: hidden; white-space: pre-wrap; overflow-wrap: anywhere; font: inherit; font-size: var(--font-size); line-height: 1.8; }
  code { font-family: inherit; }
  :global(.source) :global(.tok-keyword), :global(.source) :global(.tok-operator) { color: var(--keyword); }
  :global(.source) :global(.tok-string), :global(.source) :global(.tok-number) { color: var(--string); }
  :global(.source) :global(.tok-comment) { color: var(--comment); font-style: italic; }
  :global(.source) :global(.tok-variableName), :global(.source) :global(.tok-propertyName) { color: var(--muted); }
  :global(.output) { margin-top: 8px; color: var(--faint); font-size: var(--font-size); max-height: 320px; overflow-y: auto; }
</style>
