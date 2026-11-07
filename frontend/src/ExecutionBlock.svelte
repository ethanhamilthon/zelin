<script lang="ts">
  import type { Execution } from './agent-types';
  import { highlightJavaScript } from './highlight';
  import { fade } from 'svelte/transition';
  let { execution, folded: initial = false }: { execution: Execution; folded?: boolean } = $props();
  let expanded = $state(false);
  let folded = $state(false);
  $effect.pre(() => { folded = initial; });
  const lines = $derived(execution.code.split('\n'));
  const code = $derived(folded ? lines[0] + (lines.length > 1 ? ' …' : '') : execution.code);
</script>

<div class="execution" aria-label="JavaScript execution" aria-busy={execution.state === 'running'}>
  <button class="fold" class:folded onclick={() => { folded = !folded; }} aria-expanded={!folded} aria-label={folded ? 'Unfold code' : 'Fold code'}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6" /></svg>
  </button>
  <button class="source" class:folded onclick={() => { expanded = !expanded; }} aria-expanded={expanded && !folded} aria-label="Toggle execution output">
    <pre><code>{@html highlightJavaScript(code)}</code></pre>
  </button>
  {#if !folded && expanded && execution.output}
    <pre class="output" transition:fade={{ duration: 120 }}>{execution.output}</pre>
  {/if}
</div>

<style>
  .execution { position: relative; margin: 0 0 8px; padding: 4px 14px; border-radius: 6px; }
  .fold { position: absolute; left: -30px; top: 10px; display: inline-flex; padding: 0; border: 0; background: none; color: var(--dim); cursor: pointer; transition: transform 160ms, color 140ms; }
  .fold:hover, .fold:focus-visible { color: var(--foreground); outline: none; }
  .fold.folded { transform: rotate(-90deg); }
  .source { display: block; width: 100%; text-align: left; color: var(--dim); background: none; border: 0; padding: 0; cursor: pointer; font: inherit; transition: color 160ms; }
  .source.folded pre { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .source:hover { color: var(--muted); }
  .source:focus-visible { outline: 1px solid var(--dim); outline-offset: 5px; }
  pre { margin: 0; overflow-x: hidden; white-space: pre-wrap; overflow-wrap: anywhere; font: inherit; font-size: var(--font-size); line-height: 1.8; }
  code { font-family: inherit; }
  .source :global(.tok-keyword), .source :global(.tok-operator) { color: var(--keyword); }
  .source :global(.tok-string), .source :global(.tok-number) { color: var(--string); }
  .source :global(.tok-comment) { color: var(--comment); font-style: italic; }
  .source :global(.tok-variableName), .source :global(.tok-propertyName) { color: var(--muted); }
  .output { margin-top: 8px; color: var(--faint); font-size: var(--font-size); max-height: 320px; overflow-y: auto; }
</style>
