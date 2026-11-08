<script lang="ts">
  import type { Usage } from './agent-types';
  import TurnSummary from './TurnSummary.svelte';
  import { live } from './turn';
  let { started, usage, baseline, window }: { started: number; usage: Usage | null; baseline?: Usage | null; window: number } = $props();
  let now = $state(Date.now());
  $effect(() => {
    const timer = setInterval(() => { now = Date.now(); }, 1000);
    return () => clearInterval(timer);
  });
  const turn = $derived(live(started, Math.max(now, started), usage, baseline, window));
</script>

<TurnSummary {turn} running />
