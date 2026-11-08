<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  let { count, label, onclick }: { count: number; label: string; onclick: (event: MouseEvent) => void } = $props();
</script>

<Button data-ui="ui-morebutton" class="more" {onclick}>Show {count} earlier {label}</Button>

<style>
  :global(.more) { display: block; margin: 0 0 12px; padding: 0; border: 0; background: none; color: var(--dim); font: inherit; line-height: 1.8; cursor: pointer; transition: color 140ms; }
  :global(.more):hover, :global(.more):focus-visible { color: var(--muted); outline: none; }
</style>
