<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  let { used, limit, onclick }: { used: number | null; limit: number; onclick: () => void } = $props();
  const percent = $derived(used === null || limit <= 0 ? null : Math.min(100, Math.max(0, used / limit * 100)));
  const label = $derived(percent === null ? 'Context usage unavailable' : `Context used: ${percent.toFixed(1)}%`);
</script>

<Button data-ui="ui-contextindicator" variant="icon" class="context" {onclick} title={label} aria-label={label}>
  <span style:background={`conic-gradient(var(--dim) ${(percent ?? 0)}%, var(--border) 0)`}></span>
</Button>

<style>
  :global(.context) { min-width: 0; min-height: 0; padding: 4px; }
  span { width: 14px; height: 14px; border-radius: 50%; }
</style>
