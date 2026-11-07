<script lang="ts">
  import type { Message } from './conversation';
  import type { Execution } from './agent-types';
  import { renderMarkdown } from './markdown';
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
        <div class="markdown">{@html renderMarkdown(message.content)}</div>
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
  .markdown { font-size: 16px; line-height: 1.8; }
  .markdown :global(p) { margin: 0 0 12px; }
  .markdown :global(pre) { overflow-x: auto; padding: 0; background: none; }
  .markdown :global(code) { font-family: inherit; color: inherit; }
  .markdown :global(a) { color: #ccc; text-underline-offset: 3px; }
  .markdown :global(blockquote) { margin-left: 0; padding-left: 20px; border-left: 2px solid #404040; color: #888; }
  .markdown :global(table) { display: block; overflow-x: auto; border-collapse: collapse; }
  .markdown :global(th), .markdown :global(td) { border: 1px solid #404040; padding: 6px 12px; }
  @keyframes appear { from { opacity: 0; transform: translateY(3px); } to { opacity: 1; transform: none; } }
</style>
