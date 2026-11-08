<script lang="ts">
 import { Switch } from 'bits-ui';
 let { checked = false, label, disabled = false }: { checked?: boolean; label: string; disabled?: boolean } = $props();
</script>
<Switch.Root {checked} {disabled} aria-label={label} class="switch pointer-events-none inline-flex h-5 w-9 shrink-0 items-center rounded-full border border-border bg-secondary transition-colors data-[state=checked]:bg-primary" tabindex={-1}>
 {#snippet child({ props })}<span {...props}><span class="block size-4 rounded-full bg-foreground transition-transform" style:transform={checked ? "translateX(16px)" : "translateX(2px)"}></span></span>{/snippet}
</Switch.Root>
