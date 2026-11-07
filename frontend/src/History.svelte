<script lang="ts">
  import type { Message } from './conversation';
  import type { Execution } from './agent-types';
  import { highlightMarkdown } from './highlight';
  import ExecutionBlock from './ExecutionBlock.svelte';
  import MessageIcon from './MessageIcon.svelte';

  let { messages, executions = [], busy = false }: { messages: Message[]; executions?: Execution[]; busy?: boolean } = $props();

  function icon(message: Message, index: number) {
    if (message.role === 'assistant') return 'assistant';
    if (busy && index === messages.length - 1) return 'loading';
    return messages[index + 1]?.role === 'assistant' ? 'done' : 'user';
  }
</script>

<section aria-label="Conversation" aria-live="polite">
  {#each messages as message, index}
    <article class:assistant={message.role === 'assistant'} aria-label={message.role === 'user' ? 'User message' : 'Assistant message'}>
      {#each message.executions ?? [] as execution}
        <ExecutionBlock {execution} />
      {/each}
      <div class="message">
        <MessageIcon kind={icon(message, index)} />
        <div class="text">{@html highlightMarkdown(message.content)}</div>
      </div>
    </article>
  {/each}
  {#each executions as execution}
    <ExecutionBlock {execution} />
  {/each}
</section>

<style>
  section { padding: 32px clamp(48px, 6vw, 80px) 0; }
  article:last-child { margin-bottom: 0; }
  article { margin-bottom: 32px; overflow-wrap: anywhere; color: #e0e0e0; animation: appear 180ms ease-out; }
  .message { position: relative; }
  .assistant { color: #b0b0b0; }
  .text { font-size: var(--font-size); line-height: 1.8; white-space: pre-wrap; }
  @keyframes appear { from { opacity: 0; transform: translateY(3px); } to { opacity: 1; transform: none; } }
</style>
