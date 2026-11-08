<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { createEditor, type Editor } from './editor';
  import { createModelController } from './model-controller.svelte';
  import { blank, createSessionController } from './session-controller.svelte';
  import { createSettings } from './settings.svelte';
  import { createToasts } from './toasts.svelte';
  import { createChat } from './chat.svelte';
  import { matches } from './hotkeys';
  import { token } from './images';
  import ImagePicker from './ImagePicker.svelte';
  import { listen } from './bridge';
  import { chime } from './sound';
  import { fileDocument, namedDocument, type Document } from './documents';
  import { follow, links, type Target } from './links';
  import History from './History.svelte';
  import StatusBar from './StatusBar.svelte';
  import Panels, { type ModalName } from './Panels.svelte';
  import Toasts from './Toasts.svelte';
  import DocumentEditor from './DocumentEditor.svelte';

  let container: HTMLDivElement, scroller: HTMLElement, picker: ImagePicker;
  let editor: Editor | undefined;
  let documentEditor = $state<DocumentEditor>();
  let modal = $state<ModalName | null>(null);
  let image = $state<{ url: string; number: number } | null>(null);
  let doc = $state<Document | null>(null);
  const settings = createSettings();
  const sessions = createSessionController();
  const models = createModelController();
  const toasts = createToasts();
  const chat = createChat(sessions, finished);
  const conversation = $derived(sessions.current ?? blank());
  const keys = $derived(settings.state.hotkeys);

  async function finished(id: string, succeeded: boolean) {
    if (settings.state.sound) chime();
    if (!settings.state.toast) return;
    await sessions.refresh();
    const title = sessions.state.items.find((item) => item.id === id)?.title ?? 'Session';
    toasts.push(id, `${title} · ${succeeded ? 'finished' : 'stopped'}`);
  }

  const openLink = (target: Target) => follow(target, (path) => { modal = null; doc = fileDocument(path); }).catch((error) => toasts.push('', String(error)));
  function scrollLatest() { scroller.querySelector('article:last-of-type')?.scrollIntoView({ block: 'center', behavior: 'smooth' }); }
  function closeModal() { modal = null; image = null; editor?.focus(); }
  function open(name: ModalName) {
    modal = name;
    if (name === 'sessions') void sessions.refresh();
    if (name === 'models' || name === 'manager') void models.load();
  }
  function toggle(name: ModalName) { if (modal === name) closeModal(); else open(name); }
  function showImage(url: string | undefined, number: number) { if (url) { modal = null; image = { url, number }; } }

  async function selectSession(id?: string) {
    if (!(await sessions.select(id))) return;
    editor?.setText(conversation.draft);
    closeModal();
    await tick();
    scroller.scrollTop = 0;
  }

  function stop() {
    if (window.getSelection()?.toString() || editor?.selected()) return false;
    return chat.stop();
  }

  function keyboard(event: KeyboardEvent) {
    if (event.key === 'Escape' && (modal || image)) { event.preventDefault(); closeModal(); return; }
    if (event.key === 'Escape' && doc) { event.preventDefault(); doc = null; editor?.focus(); return; }
    if (doc && matches(event, keys.save)) { event.preventDefault(); void documentEditor?.commit(); return; }
    if (doc) return;
    if (matches(event, keys.commands)) toggle('commands');
    else if (matches(event, keys.sessions)) toggle('sessions');
    else if (matches(event, keys.newSession)) void selectSession();
    else if (!(matches(event, keys.stop) && stop())) return;
    event.preventDefault();
  }

  async function submit() {
    const current = sessions.current;
    if (!editor || !current || current.busy) return;
    const pending = chat.submit();
    editor.setText(current.draft);
    await tick();
    scrollLatest();
    const id = await pending;
    if (id && id === sessions.state.id && editor.text() !== current.draft) editor.setText(current.draft);
    if (!modal && !doc) editor.focus();
  }

  $effect(() => { editor?.setKeys({ send: keys.send, newline: keys.newline }); });

  onMount(() => {
    let mounted = true;
    const off = listen('session:title', (data) => { const { id, title } = data as { id: string; title: string }; sessions.rename(id, title); });
    void Promise.all([settings.start(), sessions.start()]).then(() => {
      if (!mounted) return;
      editor = createEditor(container, conversation.draft, (text) => {
        if (!sessions.current) return;
        sessions.current.draft = text;
        if (!sessions.state.loading) sessions.schedule();
      }, {
        keys: { send: keys.send, newline: keys.newline }, onSubmit: () => { void submit(); },
        onImage: (number) => showImage(sessions.current?.attachments[number - 1], number), onAttach: sessions.attach,
      });
      editor.focus();
    });
    void models.start();
    return () => { mounted = false; sessions.dispose(); editor?.destroy(); void off.then((stop) => stop()); };
  });
</script>

<svelte:window onkeydown={keyboard} />

<div class="app" class:hidden={!!doc}>
  <main bind:this={scroller} use:links={openLink}>
    {#if conversation.messages.length || conversation.busy}
      <History messages={conversation.messages} executions={conversation.executions} busy={conversation.busy}
        started={conversation.started} usage={conversation.usage} baseline={conversation.baseline} window={models.state.info.contextWindow}
        showCode={settings.state.showCode} foldCode={settings.state.foldCode} stack={settings.state.stackSteps} onimage={(url, number) => showImage(url, number)} />
    {/if}
    <div class="editor" class:empty={!conversation.messages.length} bind:this={container}></div>
    {#if conversation.error}<p class="error" role="alert">{conversation.error}</p>{/if}
    {#if sessions.state.error}<p class="error" role="alert">{sessions.state.error}</p>{/if}
    {#if conversation.messages.length}<div class="scroll-space" aria-hidden="true"></div>{/if}
  </main>
  <StatusBar context={conversation.usage?.context ?? null} info={models.state.info} commands={keys.commands} onopen={open} />
</div>
{#if doc}
  {#key doc}<DocumentEditor bind:this={documentEditor} title={doc.title} load={doc.load} save={doc.save} saveKey={keys.save}
    onclose={() => { doc = null; editor?.focus(); }} />{/key}
{/if}
<ImagePicker bind:this={picker} onpick={(url) => { editor?.insert(token(sessions.attach(url))); editor?.focus(); }} />
<Panels {modal} {image} {settings} {sessions} {models} {conversation} onclose={closeModal} {open}
  onnew={() => { void selectSession(); }} onselect={(id) => { void selectSession(id); }}
  onstop={() => { closeModal(); chat.stop(); }} onattach={() => { closeModal(); picker.open(); }}
  onlatest={() => { closeModal(); scrollLatest(); }} ondocument={(name) => { modal = null; doc = namedDocument(name); }} />
<Toasts items={toasts.items} onopen={(toast) => { toasts.dismiss(toast.id); if (toast.session) void selectSession(toast.session); }} />

<style>
  .app { height: 100%; display: grid; grid-template-rows: minmax(0, 1fr) auto; }
  .hidden { display: none; }
  main { min-height: 0; overflow-x: hidden; overflow-y: auto; }
  .editor { min-height: 180px; }
  .empty { min-height: calc(100vh - 44px); }
  .scroll-space { height: 50vh; }
  .error { margin: 0; padding: 12px clamp(48px, 6vw, 80px); font: var(--font-size)/1.5 var(--mono); color: var(--muted); }
</style>
