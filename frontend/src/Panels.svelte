<script lang="ts" module>
  export type ModalName = 'commands' | 'sessions' | 'models' | 'efforts' | 'usage' | 'manager' | 'settings' | 'hotkeys' | 'instructions';
</script>

<script lang="ts">
  import type { Conversation } from './conversation';
  import type { createSessionController } from './session-controller.svelte';
  import type { createModelController } from './model-controller.svelte';
  import type { createSettings } from './settings.svelte';
  import Commands from './Commands.svelte';
  import SessionList from './SessionList.svelte';
  import ModelList from './ModelList.svelte';
  import ModelManager from './ModelManager.svelte';
  import EffortList from './EffortList.svelte';
  import Settings from './Settings.svelte';
  import Hotkeys from './Hotkeys.svelte';
  import Instructions from './Instructions.svelte';
  import Modal from './Modal.svelte';
  import UsageView from './UsageView.svelte';
  let { modal, image, settings, sessions, models, conversation, onclose, open, onnew, onselect, onstop, onattach, onlatest, ondocument }: {
    modal: ModalName | null; image: { url: string; number: number } | null; conversation: Conversation;
    settings: ReturnType<typeof createSettings>; sessions: ReturnType<typeof createSessionController>; models: ReturnType<typeof createModelController>;
    onclose: () => void; open: (name: ModalName) => void; onnew: () => void; onselect: (id: string) => void;
    onstop: () => void; onattach: () => void; onlatest: () => void; ondocument: (name: string) => void;
  } = $props();
  const info = $derived(models.state.info);
  async function apply(change: Promise<boolean>) { if (await change) onclose(); }
</script>

{#if image}
  <Modal title={`Image #${image.number}`} {onclose}><img src={image.url} alt={`Attachment ${image.number}`} /></Modal>
{:else if modal === 'commands'}
  <Commands hotkeys={settings.state.hotkeys} {onclose} onopen={open} />
{:else if modal === 'sessions'}
  <SessionList sessions={sessions.state.items} active={sessions.state.id} status={sessions.status} {onclose} {onselect} {onnew} />
{:else if modal === 'models'}
  <ModelList models={models.state.models} active={info.model} hidden={settings.state.hiddenModels} loading={models.state.loading} error={models.state.error} {onclose} onselect={(id) => { void apply(models.select(id)); }} />
{:else if modal === 'manager'}
  <ModelManager models={models.state.models} hidden={settings.state.hiddenModels} titleModel={info.titleModel || info.model} loading={models.state.loading} error={models.state.error}
    {onclose} ontoggle={settings.toggleModel} ontitle={(id) => { void models.titleModel(id); }} />
{:else if modal === 'efforts'}
  <EffortList efforts={info.efforts} active={info.effort} error={models.state.error} {onclose} onselect={(effort) => { void apply(models.effort(effort)); }} />
{:else if modal === 'settings'}
  <Settings settings={settings.state} {onclose} onchange={settings.update} />
{:else if modal === 'hotkeys'}
  <Hotkeys hotkeys={settings.state.hotkeys} {onclose} onbind={settings.bind} />
{:else if modal === 'instructions'}
  <Instructions {onclose} onopen={ondocument} />
{:else if modal === 'usage'}
  <Modal title="Context and usage" {onclose}><UsageView usage={conversation.usage} contextWindow={info.contextWindow} /></Modal>
{/if}

<style>
  img { display: block; max-width: 100%; max-height: 55vh; margin: 0 auto 8px; border-radius: 4px; }
</style>
