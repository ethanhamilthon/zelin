# TODO

1. [x] Left-align step stack summary and expanded items, and display steps as a tree:
   - [x] Left-align the "n steps" summary row with user/assistant messages instead of centering it.
   - [x] Left-align all content inside the expanded step stack.
   - [x] Render expanded execution and thinking steps as a tree structure with connecting guide lines.

2. [x] Rework steering, history persistence, and step streaming:
   - [x] Fix steering context loss bug where the agent forgets prior actions after a steer message.
   - [x] Make conversation history strictly append-only: grows only at the end and is never mutated, replaced, or truncated prior to compaction.
   - [x] Record errors as tool call outputs or system messages instead of resetting or wiping history.
   - [x] Persist every intermediate step (tool call, output, thinking) immediately to the database to prevent data loss.
