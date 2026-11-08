<script lang="ts">
  import type { Execution } from './agent-types';
  import StepStack from './StepStack.svelte';
  import BoundedSteps from './BoundedSteps.svelte';
  let { executions, running = false, showCode = true, foldCode = false, stack = false, started, activity = '' }: {
    executions: Execution[]; running?: boolean; showCode?: boolean; foldCode?: boolean; stack?: boolean; started?: number; activity?: string;
  } = $props();
  const visible = $derived(executions.filter((execution) => execution.kind !== 'message' && (showCode || execution.kind === 'thinking')));
</script>

{#if running}
  <StepStack steps={visible} running {foldCode} {started} {activity} expanded={!stack} />
{:else if stack && visible.length > 1}
  <StepStack steps={visible} {foldCode} />
{:else}
  <BoundedSteps steps={visible} {foldCode} />
{/if}
