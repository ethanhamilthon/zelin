<script lang="ts">
  import { summary } from './turn';
  import type { Turn } from './turn';
  let { turn, running = false }: { turn: Turn; running?: boolean } = $props();
</script>

<div class="turn" class:running role="separator" aria-label={running ? 'Current turn' : 'Turn summary'}>
  <span>{summary(turn)}</span>
</div>

<style>
  .turn { display: flex; align-items: center; gap: 16px; margin: 16px calc(-1 * var(--gutter, 0px)) 0; color: var(--faint); white-space: nowrap; }
  .turn::before, .turn::after { content: ''; flex: 1; height: 1px; background: var(--border); }
  .running span { animation: pulse 1.6s ease-in-out infinite; }
  @keyframes pulse { 50% { opacity: 0.55; } }
  @media (prefers-reduced-motion: reduce) { .running span { animation: none; } }
</style>
