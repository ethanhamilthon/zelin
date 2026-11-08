<script lang="ts">
  import type { Execution } from './agent-types';
  import StepStack from './StepStack.svelte';
  import BoundedSteps from './BoundedSteps.svelte';
  let { executions, running = false, showCode = true, foldCode = false, stack = false }: {
    executions: Execution[]; running?: boolean; showCode?: boolean; foldCode?: boolean; stack?: boolean;
  } = $props();
  const visible = $derived(executions.filter((execution) => showCode || execution.kind === 'thinking'));
</script>

{#if stack && visible.length > 1}
  <StepStack steps={visible} {running} {foldCode} />
{:else}
  <BoundedSteps steps={visible} {foldCode} />
{/if}
