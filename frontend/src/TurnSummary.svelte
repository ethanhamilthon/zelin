<script lang="ts">
  import { summary } from './turn';
  import type { Turn } from './turn';
  let { turn }: { turn: Turn } = $props();
</script>

<div class="turn" role="separator" aria-label="Turn summary">
  <span>{summary(turn)}</span>
</div>

<style>
  .turn { display: flex; align-items: center; gap: 16px; margin: 16px calc(-1 * var(--gutter, 0px)) 0; color: var(--faint); white-space: nowrap; }
  .turn::before, .turn::after { content: ''; flex: 1; height: 1px; background: var(--border); }
</style>
