<script lang="ts">
  import type { Usage } from './agent-types';
  import { cost, percent, tokens } from './format';
  let { usage, contextWindow }: { usage: Usage | null; contextWindow: number } = $props();
  const idle = $derived(!usage || usage.requests === 0);
  const count = (value: number | null | undefined) => tokens(idle ? 0 : value);
  const limit = $derived(contextWindow > 0 ? tokens(contextWindow) : 'Unknown');
  const hitRate = $derived(usage?.cacheRead != null && usage.input ? percent(usage.cacheRead, usage.input) : 'Unknown');
  const spent = $derived(idle ? cost(0) : cost(usage?.cost));
  const note = $derived.by(() => {
    if (!usage || idle) return 'No requests yet';
    const notes = [`${usage.requests} ${usage.requests === 1 ? 'request' : 'requests'}`];
    if (usage.estimated) notes.push('includes pricing-based estimates');
    if (usage.unpriced) notes.push(`${usage.unpriced} without cost data`);
    return notes.join(' · ');
  });
</script>

<dl>
  <dt>Context</dt><dd>{count(usage?.context)} / {limit}</dd>
  <dt>Context used</dt><dd>{idle ? percent(0, contextWindow) : percent(usage?.context, contextWindow)}</dd>
  <dt>Input tokens</dt><dd>{count(usage?.input)}</dd>
  <dt>Output tokens</dt><dd>{count(usage?.output)}</dd>
  <dt>Cache reads</dt><dd>{count(usage?.cacheRead)}</dd>
  <dt>Cache writes</dt><dd>{count(usage?.cacheWrite)}</dd>
  <dt>Cache hit rate</dt><dd>{idle ? 'Unknown' : hitRate}</dd>
  <dt>Spending</dt><dd>{spent}{usage?.estimated ? ' (estimated)' : ''}</dd>
</dl>
<p>{note}</p>

<style>
  dl { display: grid; grid-template-columns: auto 1fr; gap: 8px 24px; margin: 0; padding: 4px 12px; }
  dt { color: #777; }
  dd { margin: 0; color: #ddd; text-align: right; white-space: nowrap; }
  p { margin: 12px 0 0; padding: 0 12px 8px; color: #666; white-space: normal; }
</style>
