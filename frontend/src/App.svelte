<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { createEditor } from './editor';
  import { sendMessage } from './conversation';
  import type { Conversation } from './conversation';
  import { runAgent } from './agent';
  import { createModelController } from './model-controller.svelte';
  import { createSessionController } from './session-controller.svelte';
  import { readToEnd } from './reading';
  import History from './History.svelte';
  import StatusBar from './StatusBar.svelte';
  import Commands from './Commands.svelte';
  import SessionList from './SessionList.svelte';
  import ModelList from './ModelList.svelte';
  import EffortList from './EffortList.svelte';
  import Modal from './Modal.svelte';
  import UsageView from './UsageView.svelte';

  let container: HTMLDivElement;
  let scroller: HTMLElement;
  let editor: ReturnType<typeof createEditor> | undefined;
  type ModalName = 'commands' | 'sessions' | 'models' | 'efforts' | 'usage';
  let modal = $state<ModalName | null>(null);
  let unread = $state(false);
  const conversation = $state<Conversation>({ messages: [], draft: '', busy: false, error: '', activity: 'Ready', executions: [], usage: null });
  const sessions = createSessionController(conversation);
  const models = createModelController();

  function checkRead() { if (unread && readToEnd(scroller)) unread = false; }
  function scrollLatest() { scroller.querySelector('article:last-of-type')?.scrollIntoView({ block: 'center', behavior: 'smooth' }); }
  function closeModal() { modal = null; editor?.focus(); }
  function open(name: ModalName) {
    modal = name;
    if (name === 'sessions') void sessions.refresh();
    if (name === 'models') void models.load();
  }
  function showSessions() { open('sessions'); }
  async function selectModel(id: string) { if (await models.select(id)) closeModal(); }
  async function selectEffort(effort: string) { if (await models.effort(effort)) closeModal(); }

  async function selectSession(id?: string) {
    if (await sessions.select(id)) {
      unread = false;
      editor?.setText(conversation.draft);
      closeModal();
      await tick();
      scroller.scrollTop = 0;
    }
  }

  function keyboard(event: KeyboardEvent) {
    if (event.key === 'Escape' && modal) { event.preventDefault(); closeModal(); return; }
    if (!event.ctrlKey || !['k', 's'].includes(event.key.toLowerCase())) return;
    event.preventDefault();
    if (event.key.toLowerCase() === 's') {
      if (modal === 'sessions') closeModal(); else showSessions();
    } else if (modal === 'commands') closeModal(); else modal = 'commands';
  }

  async function submit() {
    if (!editor || sessions.state.loading || !sessions.state.id || conversation.busy || !conversation.draft.trim()) return;
    unread = false;
    const pending = sendMessage(conversation, async (messages, observe) => {
      await sessions.save();
      return runAgent(messages, observe, sessions.state.id);
    });
    editor.setText('');
    await tick();
    scrollLatest();
    const succeeded = await pending;
    if (editor.view.state.doc.toString() !== conversation.draft) editor.setText(conversation.draft);
    await tick();
    if (succeeded) unread = !readToEnd(scroller);
    if (!modal) editor.focus();
    await sessions.save().catch(() => {});
  }

  onMount(() => {
    let mounted = true;
    void sessions.start().then(() => {
      if (!mounted) return;
      editor = createEditor(container, conversation.draft, (text) => {
        conversation.draft = text;
        if (!sessions.state.loading) sessions.schedule();
      }, () => { void submit(); });
      editor.focus();
    });
    void models.start();
    return () => { mounted = false; sessions.dispose(); editor?.destroy(); };
  });
</script>

<svelte:window onkeydown={keyboard} onresize={checkRead} onfocus={checkRead} />

<div class="app">
  <main bind:this={scroller} onscroll={checkRead}>
    {#if conversation.messages.length}
      <History messages={conversation.messages} executions={conversation.executions} busy={conversation.busy} />
    {/if}
    <div class="editor" class:empty={!conversation.messages.length} bind:this={container}></div>
    {#if conversation.error}<p class="error" role="alert">{conversation.error}</p>{/if}
    {#if sessions.state.error}<p class="error" role="alert">{sessions.state.error}</p>{/if}
    {#if conversation.messages.length}<div class="scroll-space" aria-hidden="true"></div>{/if}
  </main>
  <StatusBar busy={conversation.busy} activity={conversation.activity} info={models.state.info} usage={conversation.usage} {unread} onopen={open} />
</div>
{#if modal === 'commands'}
  <Commands onclose={closeModal} onfocus={closeModal} onlatest={() => { closeModal(); scrollLatest(); }} onnew={() => { void selectSession(); }} onopen={open} reasoning={models.state.info.efforts.length > 0} busy={conversation.busy || sessions.state.loading} />
{:else if modal === 'sessions'}
  <SessionList sessions={sessions.state.items} active={sessions.state.id} busy={conversation.busy || sessions.state.loading} onclose={closeModal} onselect={(id) => { void selectSession(id); }} onnew={() => { void selectSession(); }} />
{:else if modal === 'models'}
  <ModelList models={models.state.models} active={models.state.info.model} loading={models.state.loading} error={models.state.error} onclose={closeModal} onselect={(id) => { void selectModel(id); }} />
{:else if modal === 'efforts'}
  <EffortList efforts={models.state.info.efforts} active={models.state.info.effort} fallback={models.state.info.defaultEffort} error={models.state.error} onclose={closeModal} onselect={(effort) => { void selectEffort(effort); }} />
{:else if modal === 'usage'}
  <Modal title="Context and usage" onclose={closeModal}><UsageView usage={conversation.usage} contextWindow={models.state.info.contextWindow} /></Modal>
{/if}

<style>
  .app { height: 100%; display: grid; grid-template-rows: minmax(0, 1fr) auto; }
  main { min-height: 0; overflow-y: auto; }
  .editor { min-height: 180px; }
  .empty { min-height: calc(100vh - 44px); }
  .scroll-space { height: 50vh; }
  .error { margin: 0; padding: 12px clamp(48px, 6vw, 80px); font: var(--font-size)/1.5 var(--mono); color: #aaa; }
</style>
