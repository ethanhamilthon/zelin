<script lang="ts">
  let { count, label, onclick }: { count: number; label: string; onclick: (event: MouseEvent) => void } = $props();
</script>

<button class="more" {onclick}>Show {count} earlier {label}</button>

<style>
  .more { display: block; margin: 0 0 12px; padding: 0; border: 0; background: none; color: var(--dim); font: inherit; line-height: 1.8; cursor: pointer; transition: color 140ms; }
  .more:hover, .more:focus-visible { color: var(--muted); outline: none; }
</style>
