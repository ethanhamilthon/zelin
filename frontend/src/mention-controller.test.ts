import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { createMentions } from './mention-controller.svelte';
import { suggest } from './browser';

vi.mock('./browser', () => ({ suggest: vi.fn() }));
const entry = (name: string) => ({ name, path: '/w/' + name, kind: 'file' as const, size: 1 });
const at = (query: string) => ({ from: 0, to: query.length + 1, query, x: 1, y: 2 });

beforeEach(() => { vi.useFakeTimers(); vi.resetAllMocks(); });
afterEach(() => { vi.useRealTimers(); });

it('debounces lookups and retains every matching result', async () => {
  vi.mocked(suggest).mockResolvedValue(Array.from({ length: 9 }, (_, i) => entry('f' + i)));
  const mentions = createMentions(() => 'w');
  mentions.update(at('a')); mentions.update(at('ab')); mentions.update(at('abc'));
  expect(suggest).not.toHaveBeenCalled();
  await vi.advanceTimersByTimeAsync(200);
  expect(suggest).toHaveBeenCalledTimes(1);
  expect(suggest).toHaveBeenCalledWith('w', '', 'abc');
  expect(mentions.state.items).toHaveLength(9);
  expect(mentions.state.loading).toBe(false);
});

it('ignores stale responses', async () => {
  let release!: (items: ReturnType<typeof entry>[]) => void;
  vi.mocked(suggest).mockImplementationOnce(() => new Promise((resolve) => { release = resolve; })).mockResolvedValueOnce([entry('fresh')]);
  const mentions = createMentions(() => 'w');
  mentions.update(at('a')); await vi.advanceTimersByTimeAsync(200);
  mentions.update(at('ab')); await vi.advanceTimersByTimeAsync(200);
  release([entry('stale')]); await vi.advanceTimersByTimeAsync(0);
  expect(mentions.state.items.map((item) => item.name)).toEqual(['fresh']);
});

it('moves with wraparound, closes, and cancels pending lookups', async () => {
  vi.mocked(suggest).mockResolvedValue([entry('a'), entry('b'), entry('c')]);
  const mentions = createMentions(() => 'w');
  mentions.update(at('')); await vi.advanceTimersByTimeAsync(200);
  mentions.move(-1);
  expect(mentions.chosen.name).toBe('c');
  mentions.move(1);
  expect(mentions.chosen.name).toBe('a');
  mentions.update(at('x')); mentions.close(); await vi.advanceTimersByTimeAsync(300);
  expect(mentions.state.open).toBe(false);
  expect(suggest).toHaveBeenCalledTimes(1);
});

it('stays closed without a workspace and survives lookup errors', async () => {
  const closed = createMentions(() => '');
  closed.update(at('a'));
  expect(closed.state.open).toBe(false);
  vi.mocked(suggest).mockRejectedValue(new Error('denied'));
  const mentions = createMentions(() => 'w');
  mentions.update(at('a')); await vi.advanceTimersByTimeAsync(200);
  expect(mentions.state.items).toEqual([]);
  expect(mentions.state.loading).toBe(false);
});
