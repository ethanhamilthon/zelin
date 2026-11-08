import { expect, it } from 'vitest';
import { inputPlaceholder } from './placeholder';

it('picks the input placeholder by session state', () => {
  expect(inputPlaceholder(false, 0)).toBe("Let's Build Something Amazing");
  expect(inputPlaceholder(true, 1)).toBe('Agent is Working...');
  expect(inputPlaceholder(false, 2)).toBe('Continue to Build');
});
