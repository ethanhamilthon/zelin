<script lang="ts">
 import Dialog from './lib/ui/dialog.svelte';
 import type { Snippet } from 'svelte';
 let { title, hint, back = false, onclose, ondismiss, top, children, actions }: {
  title: string; hint?: string; back?: boolean; onclose: () => void; ondismiss?: () => void;
  top?: Snippet; children: Snippet; actions?: Snippet;
 } = $props();
</script>
<Dialog {title} {onclose} {ondismiss} {top}>
 {#snippet header()}
  <div class="title"><span class="heading">{title}</span>{#if actions}{@render actions()}{:else if hint}<span>{hint}</span>{/if}</div>
 {/snippet}
 {@render children()}
</Dialog>
<style>
 .title { color: var(--muted); display: flex; justify-content: space-between; align-items: center; gap: 16px; }
 .title { flex-wrap: wrap; }
 .heading { font-weight: 600; color: var(--foreground); }
 .title > span:last-child { font-size: .8em; color: var(--dim); }
</style>
