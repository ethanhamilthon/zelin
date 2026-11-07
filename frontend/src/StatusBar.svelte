<script lang="ts">
  import type { AgentInfo } from './agent-types';
  import ContextIndicator from './ContextIndicator.svelte';
  import { modelName } from './models';
  import BottomBar from './BottomBar.svelte';
  let { info, context = null, commands, onopen }: { info: AgentInfo; context?: number | null; commands: string; onopen: (name: 'models' | 'efforts' | 'usage') => void } = $props();
</script>

<BottomBar hint={`${commands} for Commands`} label="Model and commands">
  <ContextIndicator used={context} limit={info.contextWindow} onclick={() => onopen('usage')} />
  <button class="model" onclick={() => onopen('models')} title={`Select model (${info.model || info.name})`}>{modelName(info.name)}</button>
  {#if info.efforts.length}<button class="effort" onclick={() => onopen('efforts')} title="Select reasoning effort">{info.effort}</button>{/if}
</BottomBar>
