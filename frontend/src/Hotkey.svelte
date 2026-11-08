<script lang="ts">
  let { value }: { value: string } = $props();
</script>

<kbd>{value}</kbd>

<style>
  kbd { font: inherit; color: var(--faint); font-weight: 400; white-space: nowrap; }
</style>
