<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import WorkspaceIcon from './WorkspaceIcon.svelte';
  import type { AgentInfo } from './agent-types';
  import ContextIndicator from './ContextIndicator.svelte';
  import { modelName } from './models';
  import BottomBar from './BottomBar.svelte';
  import Hotkey from './Hotkey.svelte';
  import ActivitySummary from './ActivitySummary.svelte';
  let { info, context = null, commands, stop = 'Ctrl+C', busy = false, working = 0, done = 0, workspace = '', workspaceIcon = '', onstop = () => {}, onopen }: {
    workspaceIcon?: string; workspace?: string; info: AgentInfo; context?: number | null; commands: string; stop?: string; busy?: boolean; working?: number; done?: number; onstop?: () => void;
    onopen: (name: 'models' | 'efforts' | 'usage' | 'commands') => void;
  } = $props();
</script>

<BottomBar label="Model and commands">
  {#snippet center()}<div class="summary">{#if workspace}<span class="workspace" title={workspace}><WorkspaceIcon icon={workspaceIcon} /><span>{workspace}</span></span>{/if}<ActivitySummary {working} {done} /></div>{/snippet}
  <ContextIndicator used={context} limit={info.contextWindow} onclick={() => onopen('usage')} />
  <Button data-ui="ui-statusbar" class="model" onclick={() => onopen('models')} title={'Select model (' + info.model + ')'}>{modelName(info.name)}</Button>
  {#if info.efforts.length}<Button data-ui="ui-statusbar" class="effort" onclick={() => onopen('efforts')} title="Select reasoning effort">{info.effort}</Button>{/if}
  {#snippet actions()}
    {#if busy}<Button data-ui="ui-statusbar" onclick={onstop}>Stop <Hotkey value={stop} /></Button>{/if}
    <Button data-ui="ui-statusbar" onclick={() => onopen('commands')}>Commands <Hotkey value={commands} /></Button>
  {/snippet}
</BottomBar>

<style>
  :global(.workspace) > span { overflow: hidden; text-overflow: ellipsis; }
  :global(.summary) { display: flex; align-items: center; gap: 12px; min-width: 0; color: var(--dim); }
  :global(.workspace) { display: inline-flex; align-items: center; gap: 6px; max-width: 18vw; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
