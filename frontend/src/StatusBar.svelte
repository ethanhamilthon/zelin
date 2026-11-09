<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import WorkspaceIcon from './WorkspaceIcon.svelte';
  import type { AgentInfo } from './agent-types';
  import ContextIndicator from './ContextIndicator.svelte';
  import { effortName, modelName } from './models';
  import BottomBar from './BottomBar.svelte';
  import ActivitySummary from './ActivitySummary.svelte';
  let { info, context = null, working = 0, done = 0, workspace = '', workspaceIcon = '', title = '', onopen }: {
    workspaceIcon?: string; workspace?: string; title?: string; info: AgentInfo; context?: number | null; working?: number; done?: number;
    onopen: (name: 'models' | 'efforts' | 'usage' | 'sessions') => void;
  } = $props();
  const short = (text: string) => text.length > 24 ? text.slice(0, 24).trimEnd() + '…' : text;
</script>

<BottomBar label="Model and workspace">
  <ContextIndicator used={context} limit={info.contextWindow} onclick={() => onopen('usage')} />
  <Button data-ui="ui-statusbar" class="model" onclick={() => onopen('models')} title={'Select model (' + info.model + ')'}>{modelName(info.name)}</Button>
  {#if info.efforts.length}<Button data-ui="ui-statusbar" class="effort" onclick={() => onopen('efforts')} title="Select reasoning effort">{effortName(info.effort)}</Button>{/if}
  {#snippet actions()}
    <div class="summary">
      {#if workspace}<Button data-ui="ui-statusbar" class="workspace" onclick={() => onopen('sessions')} title={workspace}><WorkspaceIcon icon={workspaceIcon} /><span>{short(title) || workspace}</span></Button>{/if}
      <ActivitySummary {working} {done} />
    </div>
  {/snippet}
</BottomBar>

<style>
  :global(.summary) { display: flex; align-items: center; gap: 12px; min-width: 0; color: var(--dim); }
  :global(.workspace) { display: inline-flex; align-items: center; gap: 6px; max-width: 18vw; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  :global(.summary) :global(button.workspace) { flex-shrink: 1; min-width: 0; overflow: hidden; }
  :global(.summary) :global(button.workspace > span) { overflow: hidden; text-overflow: ellipsis; }
</style>
