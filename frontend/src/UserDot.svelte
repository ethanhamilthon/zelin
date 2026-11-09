<script lang="ts">
  let { queued = false }: { queued?: boolean } = $props();
</script>

<span class={["user-dot", { queued }]} role={queued ? 'status' : undefined} aria-label={queued ? 'Steering queued' : undefined} aria-hidden={!queued}><i></i></span>

<style>
  .user-dot { display: inline-flex; align-items: center; justify-content: center; width: 16px; height: 16px; color: var(--dim); }
  .queued { color: var(--accent); }
  i { display: block; width: 8px; height: 8px; border-radius: 50%; background: currentColor; }
</style>
