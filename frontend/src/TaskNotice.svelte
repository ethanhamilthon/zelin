<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import { oneLine } from './preview';
  let { content }: { content: string } = $props();
  let folded = $state(true);
  const text = $derived(folded ? oneLine(content.split('\n').filter((line) => line.startsWith('- ')).join(', ').replaceAll('- ', '') || content) : content);
</script>

<div class="task-notice" aria-label="Task notice">
  <Button data-ui="ui-tasknotice" variant="icon" class={["fold", "icon-slot", { folded }]} onclick={() => { folded = !folded; }} aria-expanded={!folded} aria-label={folded ? 'Unfold tasks' : 'Fold tasks'}>
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round"><path d="m6 9 6 6 6-6" /></svg>
  </Button>
  <Button data-ui="ui-tasknotice" class={["body", { folded }]} onclick={() => { if (!window.getSelection()?.toString()) folded = !folded; }} aria-expanded={!folded}>
    <p>{text}</p>
  </Button>
</div>

<style>
  .task-notice { position: relative; color: var(--dim); }
  :global(button[data-ui="ui-tasknotice"].body) { display: block; width: 100%; min-height: 0; padding: 0; border: 0; background: none; color: inherit; font: inherit; text-align: left; cursor: pointer; transition: color 160ms; }
  :global(button[data-ui="ui-tasknotice"].body):hover { color: var(--muted); }
  p { margin: 0; white-space: pre-wrap; overflow-wrap: anywhere; line-height: 1.8; }
  :global(button[data-ui="ui-tasknotice"].body.folded) p { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
