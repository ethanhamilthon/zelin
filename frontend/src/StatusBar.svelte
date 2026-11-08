<script lang="ts">
  import type { AgentInfo } from './agent-types';
  import ContextIndicator from './ContextIndicator.svelte';
  import { modelName } from './models';
  import BottomBar from './BottomBar.svelte';
  import Hotkey from './Hotkey.svelte';
  let { info, context = null, commands, stop = 'Ctrl+C', busy = false, onstop = () => {}, onopen }: {
    info: AgentInfo; context?: number | null; commands: string; stop?: string; busy?: boolean; onstop?: () => void;
    onopen: (name: 'models' | 'efforts' | 'usage' | 'commands') => void;
  } = $props();
</script>

<BottomBar label="Model and commands">
  <ContextIndicator used={context} limit={info.contextWindow} onclick={() => onopen('usage')} />
  <button class="model" onclick={() => onopen('models')} title={'Select model (' + info.model + ')'}>{modelName(info.name)}</button>
  {#if info.efforts.length}<button class="effort" onclick={() => onopen('efforts')} title="Select reasoning effort">{info.effort}</button>{/if}
  {#snippet actions()}
    {#if busy}<button onclick={onstop}>Stop <Hotkey value={stop} /></button>{/if}
    <button onclick={() => onopen('commands')}>Commands <Hotkey value={commands} /></button>
  {/snippet}
</BottomBar>
