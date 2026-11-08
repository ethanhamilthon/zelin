<script lang="ts">
  import type { Message } from './conversation';
  import type { Execution } from './agent-types';
  import Attachments from './Attachments.svelte';
  import { openAttachment } from './attachments';
  import type { Document } from './documents';
  import type { FileAttachment } from './agent-types';
  import Timeline from './Timeline.svelte';
  import MessageIcon from './MessageIcon.svelte';
  import MessageText from './MessageText.svelte';
  import TurnTime from './TurnTime.svelte';
  import type { Ref } from './refs';
  import MoreButton from './MoreButton.svelte';
  import { earlier, first, more as count } from './bounded';
  import { keepScroll } from './scroll-anchor';
  let { messages, executions = [], busy = false, showCode = true, foldCode = false, stack = false, started, activity = '', streaming = '', onimage, onref, ondocument }: {
    messages: Message[]; executions?: Execution[]; busy?: boolean; showCode?: boolean; foldCode?: boolean; stack?: boolean;
    started?: number; activity?: string; streaming?: string;
    onimage?: (url: string, number: number) => void; onref?: (ref: Ref) => void; ondocument?: (doc: Document) => void;
  } = $props();

  let error = $state('');
  const open = (file: FileAttachment) => { error = ''; void openAttachment(file, (doc) => ondocument?.(doc), onimage).catch((reason) => { error = String(reason); }); };
  let start = $state<number | null>(null);
  const from = $derived(Math.min(start ?? first(messages.length), messages.length));
  let anchor: HTMLElement;
  const reveal = (event: MouseEvent) => {
    (event.currentTarget as HTMLElement).blur();
    void keepScroll(anchor, () => { start = earlier(from); });
  };

  function icon(message: Message, index: number) {
    if (message.error) return 'error';
    if (message.role === 'assistant') return 'assistant';
    if (busy && index === messages.length - 1) return 'loading';
    return 'done';
  }
</script>

<section aria-label="Conversation" aria-live="polite" bind:this={anchor}>
  {#if from > 0}<MoreButton count={count(from)} label="messages" onclick={reveal} />{/if}
  {#each messages.slice(from) as message, offset (from + offset)}
    {@const index = from + offset}
    <article class:assistant={message.role === 'assistant'} aria-label={message.role === 'user' ? 'User message' : 'Assistant message'}>
      <Timeline executions={message.executions ?? []} {showCode} {foldCode} {stack} />
      <div class="message">
        <MessageIcon kind={icon(message, index)} />
        <MessageText content={message.content} markdown={message.role === 'assistant'} images={message.images} refs={message.refs} {onimage} {onref} />
      </div>
      {#if message.files?.length}<Attachments files={message.files} onopen={open} />{/if}
      {#if message.turn}<TurnTime duration={message.turn.duration} />{/if}
    </article>
  {/each}
  {#if error}<p role="alert">{error}</p>{/if}
  {#if busy}<Timeline {executions} running {showCode} {foldCode} {stack} {started} {activity} />{/if}
  {#if busy && streaming}<div class="message streaming" aria-label="Streaming answer"><MessageText content={streaming} markdown /></div>{/if}
</section>

<style>
  section { --gutter: clamp(48px, 6vw, 80px); padding: 20px var(--gutter) 0; }
  article:last-child { margin-bottom: 0; }
  article { min-width: 0; margin-bottom: 16px; overflow-wrap: anywhere; color: var(--foreground); animation: appear 180ms ease-out; }
  .message { position: relative; }
  .assistant, .streaming { color: var(--muted); }
  .streaming { margin-top: 8px; }
  @keyframes appear { from { opacity: 0; transform: translateY(3px); } to { opacity: 1; transform: none; } }
</style>
