<script lang="ts">
 import { onDestroy, tick } from 'svelte';
 import { Dialog } from 'bits-ui';
 import type { Snippet } from 'svelte';
 let { title, onclose, ondismiss, children, header, top }: { title: string; onclose: () => void; ondismiss?: () => void; children: Snippet; header: Snippet; top?: Snippet } = $props();
 onDestroy(() => { setTimeout(() => { if (!document.querySelector('[role="dialog"]')) document.querySelector<HTMLElement>('.cm-content')?.focus(); }, 0); });
 let panel = $state<HTMLDivElement>(null!);
 function navigate(event: KeyboardEvent) {
  if (event.defaultPrevented || !['ArrowDown', 'ArrowUp'].includes(event.key) || event.target instanceof HTMLSelectElement) return;
  const items = [...panel.querySelectorAll<HTMLElement>('input:not(:disabled), textarea:not(:disabled), select:not(:disabled), button:not(:disabled)')];
  if (!items.length) return;
  event.preventDefault();
  event.stopPropagation();
  const index = items.indexOf(document.activeElement as HTMLElement);
  items[(index + (event.key === 'ArrowUp' ? -1 : 1) + items.length) % items.length].focus();
 }
</script>
<Dialog.Root open={true} onOpenChange={(open) => { if (!open) onclose(); }}>
 <Dialog.Portal>
  <Dialog.Overlay class="overlay fixed inset-0 z-40 bg-black/55 backdrop-blur-sm" onclick={() => (ondismiss ?? onclose)()} />
  <Dialog.Content bind:ref={panel} aria-label={title} aria-describedby={undefined} class="panel dialog-content fixed left-1/2 top-[min(20vh,120px)] z-50 flex max-h-[70vh] w-[calc(100%-40px)] max-w-[640px] -translate-x-1/2 flex-col rounded-xl border border-border bg-background p-3 shadow-2xl outline-none" onkeydowncapture={navigate}
   onCloseAutoFocus={(event) => { event.preventDefault(); void tick().then(() => { if (!document.querySelector('[role="dialog"]')) document.querySelector<HTMLElement>('.cm-content')?.focus(); }); }} onEscapeKeydown={(event) => { event.preventDefault(); event.stopPropagation(); onclose(); }} onInteractOutside={(event) => { event.preventDefault(); (ondismiss ?? onclose)(); }}>
   <div class="dialog-heading"><Dialog.Title class="sr-only">{title}</Dialog.Title>{@render header()}</div>
   {#if top}<div class="top dialog-top">{@render top()}</div>{/if}
   <div class="body dialog-body">{@render children()}</div>
  </Dialog.Content>
 </Dialog.Portal>
</Dialog.Root>
<style>
 .dialog-heading { padding: 10px 12px 16px; }
 :global(.dialog-top input) { width: 100%; }
 .dialog-top { flex-shrink: 0; padding: 0 12px; }
 .dialog-body { min-height: 0; overflow-y: auto; padding: 0 12px; }
 :global(.dialog-content .dialog-body > button) { width: 100%; justify-content: space-between; text-align: left; }
 :global(.dialog-content .dialog-body > h3) { color: var(--dim); font-size: .8em; margin: 20px 0 10px; font-weight: 600; }
</style>
