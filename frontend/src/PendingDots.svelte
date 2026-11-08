<script lang="ts">
  let { label = 'Steering queued' }: { label?: string } = $props();
</script>

<span class="dots" role="status" aria-label={label}><i></i><i></i><i></i></span>

<style>
  .dots { display: inline-flex; align-items: center; justify-content: center; gap: 3px; width: 16px; height: 16px; color: var(--dim); }
  i { width: 3px; height: 3px; border-radius: 50%; background: currentColor; animation: pulse 1.2s ease-in-out infinite; }
  i:nth-child(2) { animation-delay: 150ms; }
  i:nth-child(3) { animation-delay: 300ms; }
  @keyframes pulse { 0%, 70%, 100% { opacity: 0.25; } 35% { opacity: 1; } }
</style>
