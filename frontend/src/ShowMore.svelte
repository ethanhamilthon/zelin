<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import Icon from './Icon.svelte';
  let { count, onclick }: { count: number; onclick: () => void } = $props();
</script>

<Button data-ui="ui-showmore" variant="icon-text" class="show-more" {onclick}><Icon name="down" />Show More ({count})</Button>

<style>
  :global(.dialog-content .dialog-body > button.show-more) { width: auto; justify-content: flex-start; color: var(--dim); }
</style>
