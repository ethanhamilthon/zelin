<script lang="ts">
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
  <button class="model" onclick={() => onopen('models')} title={'Select model (' + info.model + ')'}>{modelName(info.name)}</button>
  {#if info.efforts.length}<button class="effort" onclick={() => onopen('efforts')} title="Select reasoning effort">{info.effort}</button>{/if}
  {#snippet actions()}
    {#if busy}<button onclick={onstop}>Stop <Hotkey value={stop} /></button>{/if}
    <button onclick={() => onopen('commands')}>Commands <Hotkey value={commands} /></button>
  {/snippet}
</BottomBar>

<style>
  .workspace > span { overflow: hidden; text-overflow: ellipsis; }
  .summary { display: flex; align-items: center; gap: 12px; min-width: 0; color: var(--dim); }
  .workspace { display: inline-flex; align-items: center; gap: 6px; max-width: 18vw; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
