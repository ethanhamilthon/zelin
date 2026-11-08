<script lang="ts">
  import type { Message } from './conversation';
  import type { Execution } from './agent-types';
  import Steps from './Steps.svelte';
  import MessageIcon from './MessageIcon.svelte';
  import MessageText from './MessageText.svelte';
  import TurnSummary from './TurnSummary.svelte';
  import LiveTurn from './LiveTurn.svelte';
  import type { Usage } from './agent-types';
  let { messages, executions = [], busy = false, showCode = true, foldCode = false, stack = false, started, usage = null, baseline = null, window = 0, onimage }: {
    messages: Message[]; executions?: Execution[]; busy?: boolean; showCode?: boolean; foldCode?: boolean; stack?: boolean;
    started?: number; usage?: Usage | null; baseline?: Usage | null; window?: number;
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
      <Steps executions={message.executions ?? []} {showCode} {foldCode} {stack} />
      <div class="message">
        <MessageIcon kind={icon(message, index)} />
        <MessageText content={message.content} markdown={message.role === 'assistant'} images={message.images} {onimage} />
      </div>
      {#if message.turn}<TurnSummary turn={message.turn} />{/if}
    </article>
  {/each}
  <Steps {executions} {showCode} {foldCode} {stack} />
  {#if busy && started}<LiveTurn {started} {usage} {baseline} {window} />{/if}
</section>

<style>
  section { --gutter: clamp(48px, 6vw, 80px); padding: 20px var(--gutter) 0; }
  article:last-child { margin-bottom: 0; }
  article { min-width: 0; margin-bottom: 16px; overflow-wrap: anywhere; color: var(--foreground); animation: appear 180ms ease-out; }
  .message { position: relative; }
  .assistant { color: var(--muted); }
  @keyframes appear { from { opacity: 0; transform: translateY(3px); } to { opacity: 1; transform: none; } }
</style>
