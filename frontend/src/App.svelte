<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { createEditor, type Editor } from './editor';
  import { createModelController } from './model-controller.svelte';
  import { blank, createSessionController } from './session-controller.svelte';
  import { createSettings } from './settings.svelte';
  import { createToasts } from './toasts.svelte';
  import { createChat } from './chat.svelte';
  import { createModalStack } from './modal-stack.svelte';
  import { createKeyboard } from './app-keyboard';
  import { token } from './images';
  import ImagePicker from './ImagePicker.svelte';
  import { listenToBackend } from './listeners';
  import { finish } from './finish';
  import { contextAction } from './context-actions';
  import { createMentions } from './mention-controller.svelte';
  import { createFileRefs } from './file-refs.svelte';
  import FileOverlays from './FileOverlays.svelte';
  import { unresolved } from './refs';
  import { createAppActions } from './app-actions.svelte';
  import PendingDots from './PendingDots.svelte';
  import { agentsDocument, fileDocument, namedDocument, type Document } from './documents';
  import { createWorkspaces } from './workspace-controller.svelte';
  import { follow, links, type Target } from './links';
  import History from './History.svelte';
  import StatusBar from './StatusBar.svelte';
  import Panels, { type ModalName } from './Panels.svelte';
  import Toasts from './Toasts.svelte';
  import Notices from './Notices.svelte';
  import DocumentEditor from './DocumentEditor.svelte';
  let container: HTMLDivElement, scroller: HTMLElement, picker: ImagePicker; let editor = $state.raw<Editor>();
  let documentEditor = $state<DocumentEditor>();
  let image = $state<{ url: string; number: number } | null>(null);
  let doc = $state<Document | null>(null);
  const modals = createModalStack(), settings = createSettings(), sessions = createSessionController();
  const workspaces = createWorkspaces(), mentions = createMentions(() => workspaces.state.selected);
  const models = createModelController(() => sessions.state.id), toasts = createToasts();
  const chat = createChat(sessions, finish(settings, sessions, toasts));
  const files = createFileRefs({ sessions, workspaces, mentions, modals, toasts, editor: () => editor, openDoc: (next) => { doc = next; }, showImage: (url, number) => showImage(url, number) });
  const actions = createAppActions({ sessions, toasts, modals, close: () => closeModal(), select: (id) => selectSession(id) });
  const conversation = $derived(sessions.current ?? blank());
  const keys = $derived(settings.state.hotkeys), modal = $derived(doc || files.state.browsing ? null : modals.current);
  const shown = $derived(conversation.messages.length || conversation.busy);
  const openLink = (target: Target) => follow(target, (path) => { modals.clear(); doc = fileDocument(path); }).catch((e) => toasts.push('', String(e)));
  const scrollLatest = () => scroller.querySelector('article:last-of-type')?.scrollIntoView({ block: 'center', behavior: 'smooth' });
  const closeModal = () => { modals.clear(); image = null; files.state.browsing = false; editor?.focus(); };
  const showImage = (url: string | undefined, number: number) => { if (url) { modals.clear(); image = { url, number }; } };
  const stop = () => !(window.getSelection()?.toString() || editor?.selected()) && chat.stop();
  $effect(() => { const draft = conversation.draft; editor?.setText(draft, false); });
  const onChange = (text: string) => {
    if (!sessions.current) return;
    chat.edit(); sessions.current.draft = text;
    if (!sessions.state.loading) sessions.schedule();
  };
  function loadModal(name: ModalName) {
    if (name === 'actions') actions.aim();
    if (name === 'sessions' || name === 'actions') void sessions.refresh();
    if (name === 'sessions') void workspaces.load();
    if (name === 'models' || name === 'manager') void models.load();
  }
  function open(name: ModalName) { modals.open(name); loadModal(name); }
  function toggle(name: ModalName) { modals.toggle(name); if (modals.current) loadModal(name); else editor?.focus(); }
  function stepBack() {
    if (image) image = null; else if (doc) doc = null; else if (files.state.browsing) files.state.browsing = false; else modals.back();
    if (!doc && !modals.current && !files.state.browsing) editor?.focus();
  }
  async function focusInput() {
    await tick();
    if (conversation.messages.length) container?.scrollIntoView({ block: 'center' });
    else if (scroller) scroller.scrollTop = 0;
    editor?.focus();
  }
  async function selectSession(id?: string) {
    if (!(await sessions.select(id, id ? '' : workspaces.state.selected))) return;
    workspaces.follow(conversation.workspaceId);
    editor?.setText(conversation.draft); closeModal(); await focusInput();
  }
  const keyboard = createKeyboard({
    keys: () => keys, blocked: () => !!(modal || image || doc || files.state.browsing), doc: () => !!doc, stepBack: () => stepBack(),
    commit: () => { void documentEditor?.commit(); }, toggle: (name) => toggle(name), files: () => files.toggle(),
    create: () => { void selectSession(); }, stop: () => stop(), hasFiles: () => !!files.root,
  });
  async function submit() {
    const current = sessions.current;
    if (!editor || !current) return;
    const missing = unresolved(current.draft, current.refs);
    if (missing.length) { current.error = `Unknown reference: @${missing[0]}. Remove it and tag the file again.`; return; }
    if (current.busy) { await chat.steer(); return; }
    const pending = chat.submit(); editor.setText(current.draft); await tick(); scrollLatest();
    const id = await pending;
    if (id && id === sessions.state.id && editor.text() !== current.draft) editor.setText(current.draft);
    if (!modal && !doc) editor.focus();
  }
  $effect(() => { if (!conversation.steering && !conversation.draft && editor?.text()) editor.setText(''); });
  $effect(() => { if (sessions.state.id) void models.start(); });
  $effect(() => { editor?.setKeys({ send: keys.send, newline: keys.newline }); });
  onMount(() => {
    let mounted = true;
    const unlisten = listenToBackend(sessions, workspaces, settings, models);
    void Promise.all([settings.start(), sessions.start(), workspaces.load()]).then(() => {
      if (!mounted) return;
      workspaces.follow(conversation.workspaceId);
      editor = createEditor(container, conversation.draft, onChange, {
        keys: { send: keys.send, newline: keys.newline }, onSubmit: () => { void submit(); },
        onImage: (number) => showImage(sessions.current?.attachments[number - 1], number), onAttach: files.image,
        onRef: files.openById, onMention: mentions.update, onMenuKey: files.menuKey,
      });
      void focusInput();
    });
    return () => { mounted = false; sessions.dispose(); editor?.destroy(); unlisten(); };
  });
</script>
<svelte:window onkeydown={keyboard} />
<div class={["app", { hidden: !!doc }]}>
  <main bind:this={scroller} use:links={openLink}>
    {#if shown}
      {#key sessions.state.id}<History messages={conversation.messages} executions={conversation.executions} busy={conversation.busy}
        started={conversation.started} usage={conversation.usage} baseline={conversation.baseline} window={models.state.info.contextWindow}
        showCode={settings.state.showCode} foldCode={settings.state.foldCode} stack={settings.state.stackSteps} onimage={showImage} onref={files.openRef} ondocument={(next) => { doc = next; }} />{/key}
    {/if}
    <div class="input-wrap">
      {#if conversation.steering}<div class="steering"><PendingDots /></div>{/if}
      <div class={["editor", { empty: !conversation.messages.length }]} bind:this={container}></div>
    </div>
    <Notices {conversation} error={sessions.state.error} onretry={() => { void contextAction(sessions, 'compact'); }} />
    {#if conversation.messages.length}<div class="scroll-space" aria-hidden="true"></div>{/if}
  </main>
  <StatusBar workspaceIcon={workspaces.state.items.find((item) => item.id === workspaces.state.selected)?.icon ?? ''} workspace={workspaces.state.items.find((item) => item.id === workspaces.state.selected)?.name ?? ''} working={sessions.activity.working} done={sessions.activity.done} context={conversation.usage?.context ?? null} info={models.state.info} commands={keys.commands} stop={keys.stop} busy={conversation.busy} onstop={() => chat.stop()} onopen={(name) => { modals.clear(); open(name); }} />
</div>
{#if doc}
  {#key doc}<DocumentEditor bind:this={documentEditor} title={doc.title} load={doc.load} save={doc.save} saveKey={keys.save} onclose={stepBack} />{/key}
{/if}
<FileOverlays {files} {mentions} workspace={workspaces.state.selected} hidden={!!doc} />
<ImagePicker bind:this={picker} onpick={(url) => { editor?.insert(files.image(url)); editor?.focus(); }} />
<Panels {modal} back={modals.canBack} {image} {settings} {sessions} {models} {conversation} onclose={stepBack} ondismiss={closeModal} {open}
  {workspaces} onagents={async () => { await workspaces.refreshAgents(); if (workspaces.state.agents) { modals.clear(); doc = agentsDocument(workspaces.state.agents); } }}
  target={actions.state.target} onaction={(action) => { void actions.run(action); }} onrename={actions.rename}
  onnew={() => { void selectSession(); }} onselect={(id) => { void selectSession(id); }} ondocument={(name) => { doc = namedDocument(name); }} />
<Toasts items={toasts.items} onopen={(toast) => { toasts.dismiss(toast.id); if (toast.session) void selectSession(toast.session); }} />

<style>
  .app { width: 100%; height: 100%; min-width: 0; display: grid; grid-template-columns: minmax(0, 1fr); grid-template-rows: minmax(0, 1fr) auto; }
  .hidden { display: none; }
  main { min-height: 0; overflow-x: hidden; overflow-y: auto; }
  .input-wrap { position: relative; }
  .steering { position: absolute; left: calc(clamp(48px, 6vw, 80px) - 30px); top: 38px; z-index: 1; }
  .editor { min-height: 180px; }
  .empty { min-height: calc(100vh - 44px); }
  .scroll-space { height: 50vh; }
</style>
