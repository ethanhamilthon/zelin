<script lang="ts" module>
  export type ModalName = 'commands' | 'sessions' | 'about' | 'workspace' | 'actions' | 'rename' | 'models' | 'efforts' | 'usage' | 'manager' | 'settings' | 'hotkeys' | 'instructions';
</script>

<script lang="ts">
  import type { Conversation } from './conversation';
  import type { createSessionController } from './session-controller.svelte';
  import type { createModelController } from './model-controller.svelte';
  import type { createSettings } from './settings.svelte';
  import Commands from './Commands.svelte';
  import SessionActions from './SessionActions.svelte';
  import RenameSession from './RenameSession.svelte';
  import WorkspaceForm from './WorkspaceForm.svelte';
  import About from './About.svelte';
  import type { createWorkspaces } from './workspace-controller.svelte';
  import type { SessionAction } from './session-actions';
  import SessionList from './SessionList.svelte';
  import ModelList from './ModelList.svelte';
  import ModelManager from './ModelManager.svelte';
  import EffortList from './EffortList.svelte';
  import Settings from './Settings.svelte';
  import Hotkeys from './Hotkeys.svelte';
  import Instructions from './Instructions.svelte';
  import Modal from './Modal.svelte';
  import UsageView from './UsageView.svelte';
  let { modal, back = false, image, settings, sessions, models, conversation, target = '', workspaces, onclose, ondismiss, open, onnew, onselect, ondocument, onaction, onrename, onagents }: {
    modal: ModalName | null; back?: boolean; image: { url: string; number: number } | null; conversation: Conversation;
    settings: ReturnType<typeof createSettings>; sessions: ReturnType<typeof createSessionController>; models: ReturnType<typeof createModelController>;
    onclose: () => void; ondismiss: () => void; open: (name: ModalName) => void; onnew: () => void; onselect: (id: string) => void;
    onaction: (action: SessionAction) => void; ondocument: (name: string) => void;
    target?: string; onrename: (title: string) => Promise<string | null>;
    workspaces: ReturnType<typeof createWorkspaces>; onagents: () => void;
  } = $props();
  const info = $derived(models.state.info);
  const subject = $derived(sessions.conversations[target]);
  async function apply(change: Promise<boolean>) { if (await change) ondismiss(); }
</script>

{#if image}
  <Modal title={`Image #${image.number}`} {onclose} {ondismiss}><img src={image.url} alt={`Attachment ${image.number}`} /></Modal>
{:else if modal === 'commands'}
  <Commands hotkeys={settings.state.hotkeys} {back} {onclose} {ondismiss} onopen={open} />
{:else if modal === 'actions'}
  <SessionActions busy={subject?.busy} hasMessages={!!subject?.messages.length} {back} {onclose} {ondismiss} {onaction} />
{:else if modal === 'rename'}
  <RenameSession title={sessions.state.items.find((item) => item.id === target)?.title ?? ''} {back} {onclose} {ondismiss} onsave={onrename} />
{:else if modal === 'sessions'}
  <SessionList hotkey={settings.state.hotkeys.newSession} sessions={sessions.state.items} active={sessions.state.id} status={sessions.status} steps={sessions.steps} {back} {onclose} {ondismiss} {onselect} {onnew}
    workspaces={workspaces.state.items} selected={workspaces.state.selected}
    onworkspace={(id) => workspaces.follow(id)} onsettings={(id) => { void workspaces.edit(id).then(() => open('workspace')); }}
    oncreate={() => { void workspaces.edit(null).then(() => open('workspace')); }} />
{:else if modal === 'about'}
  <About {back} {onclose} {ondismiss} />
{:else if modal === 'workspace'}
  <WorkspaceForm workspace={workspaces.current} agents={workspaces.state.agents} {back} {onclose} {ondismiss} onsave={async (name, icon, dir) => { const error = await workspaces.save(name, icon, dir); if (!error) onclose(); return error; }} {onagents} />
{:else if modal === 'models'}
  <ModelList models={models.state.models} active={info.model} hidden={settings.state.hiddenModels} loading={models.state.loading} error={models.state.error} {back} {onclose} {ondismiss} onselect={(id) => { void apply(models.select(id)); }} />
{:else if modal === 'manager'}
  <ModelManager models={models.state.models} hidden={settings.state.hiddenModels} titleModel={info.titleModel || info.model} loading={models.state.loading} error={models.state.error}
    {back} {onclose} {ondismiss} ontoggle={settings.toggleModel} ontitle={(id) => { void models.titleModel(id); }} />
{:else if modal === 'efforts'}
  <EffortList efforts={info.efforts} active={info.effort} error={models.state.error} {back} {onclose} {ondismiss} onselect={(effort) => { void apply(models.effort(effort)); }} />
{:else if modal === 'settings'}
  <Settings settings={settings.state} {back} {onclose} {ondismiss} onchange={settings.update} />
{:else if modal === 'hotkeys'}
  <Hotkeys hotkeys={settings.state.hotkeys} {back} {onclose} {ondismiss} onbind={settings.bind} />
{:else if modal === 'instructions'}
  <Instructions {back} {onclose} {ondismiss} onopen={ondocument} />
{:else if modal === 'usage'}
  <Modal title="Context and usage" {back} {onclose} {ondismiss}><UsageView usage={conversation.usage} contextWindow={info.contextWindow} /></Modal>
{/if}

<style>
  img { display: block; max-width: 100%; max-height: 55vh; margin: 0 auto 8px; border-radius: 4px; }
</style>
