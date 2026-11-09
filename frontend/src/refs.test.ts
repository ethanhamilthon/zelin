import { describe, expect, it } from 'vitest';
import { entryRef, refToken, tagged, unresolved, usedRefs, newId, imageName } from './refs';
import { splitTokens } from './images';
import { findMention } from './mention';
import { crumbs } from './browser';

describe('references', () => {
  it('builds tokens with safe labels and unique ids', () => {
    const ref = entryRef({ name: 'my file.md', path: '/w/my file.md', kind: 'file', size: 1 }, 'w', 's');
    expect(refToken(ref)).toBe('[@' + ref.id + '|my file.md]');
    expect(refToken({ id: 'abcd1234', label: 'a]b|c\nd' })).toBe('[@abcd1234|a b c d]');
    expect(newId()).not.toBe(newId());
    expect(entryRef({ name: 'p.PNG', path: '/w/p.PNG', kind: 'file', size: 1 }, 'w', 's').kind).toBe('image');
    expect(entryRef({ name: 'src', path: '/w/src', kind: 'dir', size: 0 }, 'w', 's').kind).toBe('dir');
  });

  it('selects only the references used by the text and flags unknown ones', () => {
    const a = { id: 'aaaa1111', kind: 'file' as const, label: 'a' }, b = { id: 'bbbb2222', kind: 'file' as const, label: 'b' };
    expect(usedRefs('see [@aaaa1111|a] and me@example.com', [a, b])).toEqual([a]);
    expect(unresolved('[@cccc3333|c] [@aaaa1111|a]', [a])).toEqual(['c']);
    expect(unresolved('plain @mention and [@x]', [])).toEqual([]);
  });

  it('shows references as plain @ labels', () => {
    expect(tagged('open [@aaaa1111|my file.md] now')).toBe('open @my file.md now');
  });

  it('splits message text into image and reference parts without touching literal @', () => {
    const parts = splitTokens('a [Image #2] b [@aaaa1111|x y.md] c me@example.com');
    expect(parts.map((part) => part.image ?? part.ref?.id ?? part.text)).toEqual(['a ', 2, ' b ', 'aaaa1111', ' c me@example.com']);
    expect(parts[3].text).toBe('@x y.md');
  });
});

describe('mention detection', () => {
  it('opens only for an @ that starts a word', () => {
    expect(findMention('hello @sr')).toEqual({ query: 'sr', length: 3 });
    expect(findMention('@')).toEqual({ query: '', length: 1 });
    expect(findMention('(@src/ma')).toEqual({ query: 'src/ma', length: 7 });
    expect(findMention('mail me@example')).toBeNull();
    expect(findMention('hello @a b')).toBeNull();
    expect(findMention('@my docs/main')).toEqual({ query: 'my docs/main', length: 13 });
    expect(findMention('no at here')).toBeNull();
  });
});

describe('breadcrumbs', () => {
  it('builds a trail from the workspace root', () => {
    expect(crumbs('/w/proj', '/w/proj/src/deep').map((c) => c.name)).toEqual(['proj', 'src', 'deep']);
    expect(crumbs('/w/proj', '/w/proj')[0].path).toBe('/w/proj');
  });
});

it('names pasted images uniquely and keeps dropped names', () => {
  const png = 'data:image/png;base64,QQ==';
  const one = { id: 'a1a1a1', kind: 'image' as const, label: 'image-1.png', data: png };
  expect(imageName([], png)).toBe('image-1.png');
  expect(imageName([one], png)).toBe('image-2.png');
  expect(imageName([], 'data:image/jpeg;base64,QQ==')).toBe('image-1.jpg');
  expect(imageName([one], png, 'shot 1.png')).toBe('shot_1.png');
  expect(imageName([{ ...one, label: 'shot.png' }], png, 'shot.png')).toBe('image-2.png');
});
