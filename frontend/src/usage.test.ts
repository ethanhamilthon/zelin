import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import UsageView from './UsageView.svelte';
import ContextButton from './ContextButton.svelte';
import type { Usage } from './agent-types';

vi.mock('svelte/transition', () => ({ fly: () => ({ duration: 0 }) }));
const empty: Usage = { requests: 0, unpriced: 0, context: null, input: null, output: null, cacheRead: null, cacheWrite: null, cost: null, estimated: false };
const used: Usage = { ...empty, requests: 3, unpriced: 1, context: 32000, input: 50000, output: 1200, cacheRead: 20000, cost: 0.0042, estimated: true };
let app: ReturnType<typeof mount> | undefined;
afterEach(async () => { if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });

function rows() {
  return Object.fromEntries([...document.querySelectorAll('dt')].map((term) => [term.textContent, term.nextElementSibling?.textContent]));
}

describe('usage view', () => {
  it('shows zero usage before the first request', () => {
    app = mount(UsageView, { target: document.body, props: { usage: empty, contextWindow: 128000 } }); flushSync();
    expect(rows()).toMatchObject({ Context: '0 / 128,000', 'Context used': '0.0%', 'Input tokens': '0', Spending: '$0.00' });
    expect(document.body.textContent).toContain('No requests yet');
  });

  it('labels estimates and keeps unreported metrics unknown', () => {
    app = mount(UsageView, { target: document.body, props: { usage: used, contextWindow: 128000 } }); flushSync();
    expect(rows()).toMatchObject({
      Context: '32,000 / 128,000', 'Context used': '25.0%', 'Input tokens': '50,000', 'Output tokens': '1,200',
      'Cache reads': '20,000', 'Cache writes': 'Unknown', 'Cache hit rate': '40.0%', Spending: '$0.0042 (estimated)',
    });
    expect(document.body.textContent).toContain('3 requests · includes pricing-based estimates · 1 without cost data');
  });

  it('shows an unknown limit when the model metadata is missing', () => {
    app = mount(UsageView, { target: document.body, props: { usage: { ...used, context: null }, contextWindow: 0 } }); flushSync();
    expect(rows()).toMatchObject({ Context: 'Unknown / Unknown', 'Context used': 'Unknown' });
  });

  it('opens the dropdown from the context icon', () => {
    app = mount(ContextButton, { target: document.body, props: { usage: used, contextWindow: 128000 } }); flushSync();
    const button = document.querySelector('button')!;
    expect(button.getAttribute('aria-label')).toBe('Context and usage: 25.0% of context used');
    button.click(); flushSync();
    expect(document.querySelector('[role="dialog"]')?.textContent).toContain('50,000');
    document.querySelector('.popover')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true })); flushSync();
    expect(document.querySelector('[role="dialog"]')).toBeNull();
  });
});
