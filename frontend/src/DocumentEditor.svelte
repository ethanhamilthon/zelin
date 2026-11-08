<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import { onMount } from 'svelte';
  import { createEditor } from './editor';
  import type { Editor } from './editor';
  import BottomBar from './BottomBar.svelte';
  import Hotkey from './Hotkey.svelte';
  let { title, save, load, saveKey, closeKey = 'Esc', onclose }: {
    title: string; load: () => Promise<string>; save: (text: string) => Promise<void>;
    saveKey: string; closeKey?: string; onclose: () => void;
  } = $props();
  let container: HTMLDivElement;
  let editor: Editor | undefined;
  let saved = $state('');
  let text = $state('');
  let status = $state('');

  export async function commit() {
    try { await save(text); saved = text; status = 'Saved'; }
    catch (error) { status = `Could not save: ${String(error)}`; }
  }

  onMount(() => {
    let mounted = true;
    void load().then((value) => {
      if (!mounted) return;
      saved = text = value;
      editor = createEditor(container, value, (next) => { text = next; status = ''; }, { label: title, placeholder: `Write the ${title.toLowerCase()}…` });
      editor.focus();
    }).catch((error) => { status = `Could not load: ${String(error)}`; });
    return () => { mounted = false; editor?.destroy(); };
  });
</script>

<div class="screen">
  <div class="editor" bind:this={container}></div>
  <BottomBar label={title}>
    <span title={title}>{title}{text !== saved ? ' •' : ''}</span>
    {#if status}<span class="status" role="status">{status}</span>{/if}
    {#snippet actions()}
      <Button data-ui="ui-documenteditor" onclick={onclose} title="Close editor">Close <Hotkey value={closeKey} /></Button>
      <Button data-ui="ui-documenteditor" onclick={() => { void commit(); }}>Save <Hotkey value={saveKey} /></Button>
    {/snippet}
  </BottomBar>
</div>

<style>
  :global(.screen) { width: 100%; height: 100%; min-width: 0; display: grid; grid-template-columns: minmax(0, 1fr); grid-template-rows: minmax(0, 1fr) auto; }
  :global(.editor) { min-width: 0; min-height: 0; overflow: hidden; }
  :global(.status) { flex-shrink: 0; color: var(--faint) !important; }
</style>
