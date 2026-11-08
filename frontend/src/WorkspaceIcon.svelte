<script lang="ts">
  import { workspaceIcon } from './workspace-icons';
  let { icon = '', size = 18 }: { icon?: string; size?: number } = $props();
</script>

<svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" data-workspace-icon={workspaceIcon(icon).id}><path d={workspaceIcon(icon).path} /></svg>

<style>
  svg { flex-shrink: 0; }
</style>
