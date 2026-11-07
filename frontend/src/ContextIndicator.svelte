<script lang="ts">
  let { used, limit, onclick }: { used: number | null; limit: number; onclick: () => void } = $props();
  const percent = $derived(used === null || limit <= 0 ? null : Math.min(100, Math.max(0, used / limit * 100)));
  const label = $derived(percent === null ? 'Context usage unavailable' : `Context used: ${percent.toFixed(1)}%`);
</script>

<button class="context" {onclick} title={label} aria-label={label}>
  <span style:background={`conic-gradient(#fff ${(percent ?? 0)}%, #383838 0)`}></span>
</button>

<style>
  .context { display: inline-flex; flex-shrink: 0; align-items: center; padding: 4px; border: 0; background: none; cursor: pointer; }
  span { width: 14px; height: 14px; border-radius: 50%; }
</style>
