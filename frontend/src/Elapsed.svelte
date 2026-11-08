<script lang="ts">
  import { elapsed } from './turn';
  let { started }: { started: number } = $props();
  let now = $state(Date.now());
  $effect(() => {
    const timer = setInterval(() => { now = Date.now(); }, 1000);
    return () => clearInterval(timer);
  });
</script>

{elapsed(Math.max(0, now - started))}
