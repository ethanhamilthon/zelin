<script lang="ts">
  import { onMount } from 'svelte';
  import { createEditor } from './editor';
  import type { Editor } from './editor';
  import { createSettings } from './settings.svelte';
  import { createSetup, total } from './setup.svelte';
  import History from './History.svelte';
  import ModelPicker from './ModelPicker.svelte';
  import BottomBar from './BottomBar.svelte';
  let { ondone }: { ondone: () => void } = $props();
  let container: HTMLDivElement, scroller: HTMLElement;
  let editor: Editor | undefined;
  const settings = createSettings();
  const setup = createSetup((model, models) => {
    settings.update({ hiddenModels: models.map((item) => item.id).filter((id) => id !== model) });
    setTimeout(ondone, 900);
  });

  async function act(task: () => Promise<unknown>) {
    if (!editor || setup.state.busy) return;
    const pending = task();
    editor.setText('');
    await pending;
    scroller.querySelector('.editor')?.scrollIntoView({ block: 'end', behavior: 'smooth' });
    editor.focus();
  }
  const submit = () => { const text = editor?.text() ?? ''; return act(() => setup.submit(text)); };

  function arrows(event: KeyboardEvent) {
    if (!setup.picking || (event.key !== 'ArrowDown' && event.key !== 'ArrowUp')) return;
    event.preventDefault();
    event.stopPropagation();
    setup.move(event.key === 'ArrowDown' ? 1 : -1);
  }

  $effect(() => { editor?.setPlaceholder(setup.placeholder); });

  onMount(() => {
    void settings.start();
    editor = createEditor(container, '', (text) => { if (setup.picking) setup.filter(text); }, {
      label: 'Setup answer', placeholder: setup.placeholder, keys: { send: 'Enter', newline: 'Shift+Enter' }, onSubmit: () => { void submit(); },
    });
    editor.focus();
    return () => editor?.destroy();
  });
</script>

<div class="app">
  <main bind:this={scroller}>
    <History messages={setup.state.messages} busy={setup.state.busy} />
    <div class={["editor", { secret: setup.state.step === 'key' }]} bind:this={container} onkeydowncapture={arrows}></div>
    {#if setup.picking}<ModelPicker models={setup.options} active={setup.state.active} onhover={(index) => { setup.state.active = index; }} onpick={(id) => act(() => setup.pick(id))} />{/if}
  </main>
  <BottomBar hint={`Setup ${setup.number}/${total}`} label="Setup">
    <span>Welcome to Zelin</span>
  </BottomBar>
</div>

<style>
  .app { width: 100%; height: 100%; min-width: 0; display: grid; grid-template-columns: minmax(0, 1fr); grid-template-rows: minmax(0, 1fr) auto; }
  main { --gutter: clamp(48px, 6vw, 80px); min-height: 0; overflow-x: hidden; overflow-y: auto; }
  .editor { min-height: 60px; } .editor :global(.cm-content) { padding-left: var(--gutter); padding-right: var(--gutter); }
  .secret :global(.cm-line) { -webkit-text-security: disc; }
</style>
