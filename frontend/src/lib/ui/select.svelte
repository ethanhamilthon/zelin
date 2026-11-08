<script lang="ts">
 import { Select } from 'bits-ui';
 import { ChevronDown, Check } from '@lucide/svelte';
 let { value, label, options, onchange }: { value: string; label: string; options: { value: string; label: string }[]; onchange: (value: string) => void } = $props();
 const selected = $derived(options.find((item) => item.value === value)?.label ?? value);
</script>
<Select.Root type="single" {value} onValueChange={onchange}>
 <Select.Trigger aria-label={label} class="inline-flex h-9 min-w-28 items-center justify-between gap-3 rounded-md border border-border bg-background px-3 text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring">
  {selected}<ChevronDown size={14} />
 </Select.Trigger>
 <Select.Portal>
  <Select.Content class="z-[70] min-w-32 rounded-lg border border-border bg-background p-1 shadow-xl outline-none" sideOffset={4}>
   <Select.Viewport>
    {#each options as item}
     <Select.Item value={item.value} label={item.label} class="flex cursor-pointer items-center justify-between gap-5 rounded-md px-3 py-2 text-sm outline-none data-highlighted:bg-secondary">
      {item.label}{#if item.value === value}<Check size={14} />{/if}
     </Select.Item>
    {/each}
   </Select.Viewport>
  </Select.Content>
 </Select.Portal>
</Select.Root>
