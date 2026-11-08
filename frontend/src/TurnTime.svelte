<script lang="ts">
  import Elapsed from './Elapsed.svelte';
  import { elapsed } from './turn';
  let { duration, started }: { duration?: number; started?: number } = $props();
</script>

<p class="turn-time" aria-label="Turn time">{#if started}<Elapsed {started} />{:else}{elapsed(duration ?? 0)}{/if}</p>

<style>
  .turn-time { margin: 8px 0 0; color: var(--faint); font-size: var(--font-size); line-height: 1.8; }
</style>
