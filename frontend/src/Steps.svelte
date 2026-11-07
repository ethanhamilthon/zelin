<script lang="ts">
  import type { Execution } from './agent-types';
  import ExecutionBlock from './ExecutionBlock.svelte';
  import ThinkingBlock from './ThinkingBlock.svelte';
  import StepStack from './StepStack.svelte';
  let { executions, showCode = true, foldCode = false, stack = false }: { executions: Execution[]; showCode?: boolean; foldCode?: boolean; stack?: boolean } = $props();
  const visible = $derived(executions.filter((execution) => showCode || execution.kind === 'thinking'));
</script>

{#if stack && visible.length > 1}
  <StepStack steps={visible} {foldCode} />
{:else}
  {#each visible as execution (execution.id)}
    {#if execution.kind === 'thinking'}
      <ThinkingBlock text={execution.code} />
    {:else}
      <ExecutionBlock {execution} folded={foldCode} />
    {/if}
  {/each}
{/if}
