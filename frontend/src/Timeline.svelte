<script lang="ts">
  import type { Execution } from './agent-types';
  import Steps from './Steps.svelte';
  import MessageIcon from './MessageIcon.svelte';
  import MessageText from './MessageText.svelte';
  import { segments } from './segments';
  let { executions, running = false, started, activity = '' }: {
    executions: Execution[]; running?: boolean; started?: number; activity?: string;
  } = $props();
  const parts = $derived(segments(executions));
  const last = $derived(parts.at(-1));
</script>

{#each parts as part (part.id)}
  {#if part.note}
    <div class="message note"><MessageIcon kind="assistant" /><MessageText content={part.note.code} markdown /></div>
  {:else}
    <Steps executions={part.steps} running={running && part === last} {started} {activity} />
  {/if}
{/each}
{#if running && (!last || last.note)}<Steps executions={[]} running {started} {activity} />{/if}

<style>
  .note { position: relative; margin: 0 0 16px; color: var(--muted); }
</style>
