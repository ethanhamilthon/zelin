<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import type { FileAttachment } from './agent-types';
  import { call } from './bridge';
  import { fileKind } from './file-kind';
  import FileIcon from './FileIcon.svelte';
  let { files, onopen }: { files: FileAttachment[]; onopen: (file: FileAttachment) => void } = $props();
  let previews = $state<Record<string, string>>({});
  $effect(() => {
    for (const file of files) {
      if (fileKind(file.path) !== 'image' || previews[file.path] !== undefined) continue;
      previews[file.path] = '';
      void call<string>('FileService.Image', file.path).then((url) => { previews[file.path] = url; }).catch(() => {});
    }
  });
</script>
<div class="attachments" aria-label="Attached files">
  {#each files as file (file.path)}
    {@const kind = fileKind(file.path)}
    <Button data-ui="ui-attachments" onclick={() => onopen(file)} title={file.path}>
      {#if kind === 'image' && previews[file.path]}<img src={previews[file.path]} alt="" />{:else}<FileIcon {kind} />{/if}
      <span>{file.name}</span>
    </Button>
  {/each}
</div>
<style>
  :global(.attachments) { display: flex; flex-wrap: wrap; gap: 8px; margin: 12px 0 0; }
  :global(button[data-ui="ui-attachments"]) { display: inline-flex; align-items: center; gap: 8px; max-width: 100%; font: inherit; font-size: .85em; padding: 5px 10px 5px 6px; border: 1px solid var(--border); border-radius: 8px; color: var(--foreground); background: transparent; cursor: pointer; }
  :global(button[data-ui="ui-attachments"]):hover { background: var(--surface); }
  :global(button[data-ui="ui-attachments"]) :global(svg) { margin-left: 4px; color: var(--dim); }
  :global(button[data-ui="ui-attachments"]) span { min-width: 0; overflow: hidden; text-overflow: ellipsis; }
  img { width: 28px; height: 28px; object-fit: cover; border-radius: 4px; }
</style>
