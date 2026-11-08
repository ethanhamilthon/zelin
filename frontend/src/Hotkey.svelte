<script lang="ts">
 import Badge from "./lib/ui/badge.svelte";
  let { value }: { value: string } = $props();
</script>

<Badge><kbd>{value}</kbd></Badge>

<style>
  kbd { font: inherit; color: var(--faint); font-weight: 400; white-space: nowrap; }
</style>
