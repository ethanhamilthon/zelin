<script lang="ts">
  import type { Execution } from './agent-types';
  import ExecutionBlock from './ExecutionBlock.svelte';
  import ThinkingBlock from './ThinkingBlock.svelte';
  let { executions, showCode = true, foldCode = false }: { executions: Execution[]; showCode?: boolean; foldCode?: boolean } = $props();
</script>

{#each executions as execution (execution.id)}
  {#if execution.kind === 'thinking'}
    <ThinkingBlock text={execution.code} />
  {:else if showCode}
    <ExecutionBlock {execution} folded={foldCode} />
  {/if}
{/each}
