<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
 import type { FileAttachment } from './agent-types';
 let { files, onopen }: { files: FileAttachment[]; onopen: (file: FileAttachment) => void } = $props();
</script>
<div class="attachments" aria-label="Attached files">
 {#each files as file}<Button data-ui="ui-attachments" onclick={() => onopen(file)} title={file.path}>{file.name}</Button>{/each}
</div>
<style>
 :global(.attachments) { display: flex; flex-wrap: wrap; gap: 8px; margin: 12px 0 0 22px; }
 :global(button[data-ui="ui-attachments"]) { font: inherit; font-size: .85em; padding: 5px 10px; border: 1px solid var(--border); border-radius: 8px; color: var(--foreground); background: transparent; cursor: pointer; }
 :global(button[data-ui="ui-attachments"]):hover { background: var(--surface); }
</style>
