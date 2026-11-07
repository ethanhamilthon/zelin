<script lang="ts">
  import type { AgentInfo, Usage } from './agent-types';
  import ContextButton from './ContextButton.svelte';
  let { busy, activity, info, usage = null, unread = false, onopen }: {
    busy: boolean; activity: string; info: AgentInfo; usage?: Usage | null; unread?: boolean;
    onopen: (name: 'commands' | 'sessions' | 'models' | 'efforts') => void;
  } = $props();
  const status = $derived(busy ? activity : unread ? 'Unread response' : 'Ready');
</script>

<footer aria-label="Agent status and commands">
  <span class="dot" class:busy class:unread={!busy && unread} role="status" aria-label={status} title={status}></span>
  <ContextButton {usage} contextWindow={info.contextWindow} />
  <div class="metadata">
    <button class="model" onclick={() => onopen('models')} title={`Select model (${info.model || info.name})`}>{info.name}</button>
    {#if info.efforts.length}<button class="effort" onclick={() => onopen('efforts')} title="Select reasoning effort">{info.effort || 'default'}</button>{/if}
  </div>
  <button onclick={() => onopen('sessions')} aria-keyshortcuts="Control+s" title="Sessions">Ctrl + S</button>
  <button onclick={() => onopen('commands')} aria-keyshortcuts="Control+k" title="Commands">Ctrl + K</button>
</footer>

<style>
  footer { display: flex; align-items: center; gap: 18px; height: 44px; padding: 0 20px; border-top: 1px solid #2b2b2b; background: var(--background); color: #888; font: var(--font-size)/1.5 var(--mono); white-space: nowrap; }
  .dot { width: 7px; height: 7px; flex: 0 0 7px; border-radius: 50%; background: #737373; transition: background 180ms; }
  .dot.busy { background: #69b881; animation: blink 1.2s ease-in-out infinite; }
  .dot.unread { background: #679bd4; }
  .metadata { display: flex; gap: 18px; min-width: 0; flex: 1; overflow: hidden; }
  .model { overflow: hidden; text-overflow: ellipsis; min-width: 0; flex-shrink: 1; }
  button { color: inherit; font: inherit; border: 0; background: none; padding: 0; cursor: pointer; flex-shrink: 0; transition: color 140ms; }
  button:hover, button:focus-visible { color: #ddd; }
  @keyframes blink { 50% { opacity: 0.3; } }
  @media (max-width: 480px) { footer { gap: 10px; padding: 0 12px; } .metadata { gap: 10px; } }
</style>
