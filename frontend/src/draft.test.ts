import { describe, expect, it } from 'vitest';
import { loadDraft, saveDraft } from './draft';

describe('draft persistence', () => {
  it('starts empty', () => {
    expect(loadDraft({ getItem: () => null })).toBe('');
  });
  it('round-trips markdown and unicode, including an empty draft', () => {
    const values = new Map<string, string>();
    const storage = {
      getItem: (key: string) => values.get(key) ?? null,
      setItem: (key: string, value: string) => { values.set(key, value); },
    };
    for (const text of ['# привет\n\n**текст** 📝', '']) {
      expect(saveDraft(storage, text)).toBe(true);
      expect(loadDraft(storage)).toBe(text);
    }
  });
  it('survives unavailable storage', () => {
    expect(loadDraft({ getItem: () => { throw new Error('disabled'); } })).toBe('');
    expect(saveDraft({ setItem: () => { throw new Error('quota'); } }, 'draft')).toBe(false);
  });
});
