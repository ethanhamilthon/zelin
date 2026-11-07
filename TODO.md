# TODO

## Spacing

- [x] Reduce the vertical spacing between chat elements (messages, tool call blocks, thinking blocks). The current gaps are too large.

## Stack folding

- [x] Add stack folding: consecutive tool call blocks and thinking blocks are grouped into one foldable stack.
- [x] Make stack folding a toggle. When enabled, the stack is folded as one block; when expanded, each block inside can still be folded individually.

## Turn summary

- [x] After the final agent answer, show a full-width divider line with a centered summary: current time, turn duration (seconds or minutes), cost of this turn, and context usage / context window.

## Sessions window

- [x] Fix the "New session" button: it stretches across the whole header row. Make it a compact button next to the "Sessions" title.
- [x] Show the Ctrl+N hotkey on the "New session" button.
