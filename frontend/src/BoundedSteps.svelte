<script lang="ts">
  import type { Execution } from './agent-types';
  import ExecutionBlock from './ExecutionBlock.svelte';
  import ThinkingBlock from './ThinkingBlock.svelte';
  import MoreButton from './MoreButton.svelte';
  import { STEP_PAGE, more } from './bounded';
  import { keepScroll } from './scroll-anchor';
  let { steps, foldCode = false }: { steps: Execution[]; foldCode?: boolean } = $props();
  let limit = $state(STEP_PAGE);
  const hidden = $derived(Math.max(0, steps.length - limit));
  const shown = $derived(hidden ? steps.slice(hidden) : steps);
  const reveal = (event: MouseEvent) => {
    const target = event.currentTarget as HTMLElement;
    target.blur();
    void keepScroll(target, () => { limit += STEP_PAGE; });
  };
</script>

{#if hidden}<MoreButton count={more(hidden, STEP_PAGE)} label="steps" onclick={reveal} />{/if}
{#each shown as step (step.id)}
  {#if step.kind === 'thinking'}
    <ThinkingBlock text={step.code} folded={foldCode} />
  {:else}
    <ExecutionBlock execution={step} folded={foldCode} />
  {/if}
{/each}
