<script lang="ts">
  import type { Execution } from './agent-types';
  import { highlightJavaScript } from './highlight';
  import MessageIcon from './MessageIcon.svelte';
  import { fade } from 'svelte/transition';
  let { execution }: { execution: Execution } = $props();
  let expanded = $state(false);
</script>

<div class="execution" aria-label="JavaScript execution" aria-busy={execution.state === 'running'}>
  <MessageIcon kind="code" />
  <button class="source" onclick={() => { expanded = !expanded; }} aria-expanded={expanded} aria-label="Toggle execution output">
    <pre><code>{@html highlightJavaScript(execution.code)}</code></pre>
  </button>
  {#if expanded && execution.output}
    <pre class="output" transition:fade={{ duration: 120 }}>{execution.output}</pre>
  {/if}
</div>

<style>
  .execution { position: relative; margin: 0 0 24px; }
  .source { display: block; width: 100%; text-align: left; color: #909090; background: none; border: 0; padding: 0; cursor: pointer; font: inherit; transition: color 160ms; }
  .source:hover { color: #b0b0b0; }
  .source:focus-visible { outline: 1px solid #737373; outline-offset: 5px; }
  pre { margin: 0; overflow-x: auto; white-space: pre; font: inherit; font-size: 14px; line-height: 1.8; }
  code { font-family: inherit; }
  .source :global(.tok-keyword), .source :global(.tok-operator) { color: #b5b5b5; }
  .source :global(.tok-string), .source :global(.tok-number) { color: #7c7c7c; }
  .source :global(.tok-comment) { color: #606060; font-style: italic; }
  .source :global(.tok-variableName), .source :global(.tok-propertyName) { color: #999; }
  .output { margin-top: 8px; color: #606060; font-size: 13px; max-height: 320px; }
</style>
