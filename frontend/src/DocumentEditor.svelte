<script lang="ts">
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
    <span>{title}{text !== saved ? ' •' : ''}</span>
    {#if status}<span class="status" role="status">{status}</span>{/if}
    <button onclick={onclose} title="Close editor">Close <Hotkey value={closeKey} /></button>
    {#snippet actions()}<button onclick={() => { void commit(); }}>Save <Hotkey value={saveKey} /></button>{/snippet}
  </BottomBar>
</div>

<style>
  .screen { height: 100%; display: grid; grid-template-rows: minmax(0, 1fr) auto; }
  .editor { min-height: 0; overflow: hidden; }
  .status { color: var(--faint) !important; }
</style>
