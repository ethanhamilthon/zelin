<script lang="ts">
  import type { Message } from './conversation';
  import type { Execution } from './agent-types';
  import Steps from './Steps.svelte';
  import MessageIcon from './MessageIcon.svelte';
  import MessageText from './MessageText.svelte';
  let { messages, executions = [], busy = false, showCode = true, foldCode = false, onimage }: {
    messages: Message[]; executions?: Execution[]; busy?: boolean; showCode?: boolean; foldCode?: boolean;
    onimage?: (url: string, number: number) => void;
  } = $props();

  function icon(message: Message, index: number) {
    if (message.role === 'assistant') return 'assistant';
    if (busy && index === messages.length - 1) return 'loading';
    return messages[index + 1]?.role === 'assistant' ? 'done' : 'user';
  }
</script>

<section aria-label="Conversation" aria-live="polite">
  {#each messages as message, index}
    <article class:assistant={message.role === 'assistant'} aria-label={message.role === 'user' ? 'User message' : 'Assistant message'}>
      <Steps executions={message.executions ?? []} {showCode} {foldCode} />
      <div class="message">
        <MessageIcon kind={icon(message, index)} />
        <MessageText content={message.content} markdown={message.role === 'assistant'} images={message.images} {onimage} />
      </div>
    </article>
  {/each}
  <Steps {executions} {showCode} {foldCode} />
</section>

<style>
  section { padding: 32px clamp(48px, 6vw, 80px) 0; }
  article:last-child { margin-bottom: 0; }
  article { margin-bottom: 32px; overflow-wrap: anywhere; color: var(--foreground); animation: appear 180ms ease-out; }
  .message { position: relative; }
  .assistant { color: var(--muted); }
  @keyframes appear { from { opacity: 0; transform: translateY(3px); } to { opacity: 1; transform: none; } }
</style>
