<script lang="ts">
  import type { Execution } from './agent-types';
  import StepStack from './StepStack.svelte';
  import BoundedSteps from './BoundedSteps.svelte';
  let { executions, running = false, started, activity = '' }: {
    executions: Execution[]; running?: boolean; started?: number; activity?: string;
  } = $props();
  const visible = $derived(executions.filter((execution) => execution.kind !== 'message'));
</script>

{#if running || visible.length > 1}
  <StepStack steps={visible} {running} {started} {activity} />
{:else}
  <BoundedSteps steps={visible} />
{/if}
