<script lang="ts">
 import type { Snippet } from 'svelte';
 let { children }: { children: Snippet } = $props();
</script>
<span data-slot="badge" class="inline-flex shrink-0 items-center rounded-md border border-border px-2 py-0.5 text-xs font-medium text-muted-foreground">{@render children()}</span>
