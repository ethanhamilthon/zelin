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

## Steps alignment

- [x] Steps text is offset 14px to the right of message text: `.execution`, `.thinking` and `.stack .toggle` use `padding: 4px 14px` while `.message` has no horizontal padding. Align the left edge of code, thinking and stack labels with message text.
- [x] Stack toggle chevron sits ~2px above its text line: `StepStack` svg uses `top: 8px` while the line center is at 4px + 28.8px / 2. Execution fold (`top: 10px`) and thinking icon (`top: 10px`) are centered. Use one shared offset for all gutter icons.
- [x] Uneven vertical rhythm inside an open stack: toggle → first step 8px, step → step 16px (4 + 8 margin + 4), last step → answer 12px (margins collapse with `.stack`). Use one gap for all steps.
- [x] Nested steps in an open stack have no visual hierarchy: their chevrons sit in the same gutter column as the stack chevron and message icons, so it is unclear what belongs to the stack. Indent or otherwise mark the stack content.
- [x] Folded code shows two kinds of truncation: a manual ` …` for multiline code and the CSS `text-overflow: ellipsis` for long first lines. Use one.

## File tools

- [x] Add an `fs` host object to the JS runtime next to `bash` (codemode, no new provider tools): `fs.read`, `fs.write`, `fs.edit`, `fs.patch`, `fs.grep`.
- [x] `fs.read(path, {offset, limit})`: read a text file with line numbers, limited output for large files.
- [x] `fs.write(path, text)`: create or overwrite a file, create parent directories.
- [x] `fs.edit(path, oldText, newText)`: replace one exact match, fail on zero or multiple matches.
- [x] `fs.patch(path, patch)`: apply a unified diff, fail without partial writes.
- [x] `fs.grep(pattern, path, {glob})`: regex search with `file:line:text` results, skip `.git` and `node_modules`.
- [x] Keep the Go side in a separate package (`internal/fs`), return `{error}` instead of throwing raw Go errors.
- [x] Describe `fs` in the `js_repl` tool description.

## Clickable links

- [x] Render Markdown links `[text](url)` and bare URLs in messages, open them in the system browser.
- [x] Detect file paths in messages (`~/...`, `/...`, relative paths in inline code) and make them clickable.
- [x] Open `.md` files in the built-in `DocumentEditor`, like the System Prompt. It currently works only with named documents through `SettingsService.Read/Write`, so it needs to accept any file path.
- [x] Open other files with the system default app.

## Live turn summary

- [x] Show the turn summary divider as soon as the user sends a message, not only after the final answer.
- [x] Update it live while the turn runs: elapsed time ticks, cost and context update after each provider response, final values when the turn ends.
- [x] Keep the same divider for finished turns in history.

## Onboarding

- [x] Show onboarding on first launch when no provider is configured, in the same layout as the chat: timeline, input, status bar.
- [x] Status bar: "Welcome to Zash" on the left, "Setup x/n" on the right.
- [x] Step 1: add a provider (base URL and API key) through the input. Store it in settings instead of only reading `OPENAI_*` env.
- [x] Step 2: pick the main model. Only this model is enabled after setup.
- [x] Step 3: pick the title model.
- [x] Step 4: run a short test request for both models, show the result in the timeline, block finishing until both work.
- [x] After a successful check, switch to a normal new session.

## Hotkey labels

- [x] Use one pattern for every button or hint with a hotkey: name first, then hotkey, e.g. `Commands Ctrl+K`, `New session Ctrl+N`.
- [x] Render the hotkey with less emphasis: a noticeably darker color than the name.

## Stop button

- [x] While the agent is running, show a `Stop Ctrl+C` button in the status bar, left of the `Commands Ctrl+K` button, following the hotkey label pattern.
- [x] The button and Ctrl+C stop the current turn. Hide the button when the agent is idle. Ctrl+C copies instead when text is selected.

## Compaction

- [x] Add manual compaction (command in the command menu).
- [x] Add auto compaction: trigger when context usage reaches the threshold, default 60% of the context window.
- [x] Settings: auto compaction toggle and threshold value.

## Handoff

- [x] Add a "Handoff" command to the command menu.
- [x] The title model takes the whole session history and writes a handoff text.
- [x] Create a new session that starts with this handoff text.

## Steering

- [x] Sending a message while a turn is running queues it as a steer message.
- [x] The steer waits until all current tool calls finish, then goes to the agent.
- [x] While it waits, show an animated three-dot icon on the left side of the first input line, same as in chat messages.
- [x] Typing in the input again removes the icon and takes the steer out of the queue.

## Selection

- [x] Use one selection color everywhere (dark gray).
- [x] Allow text selection only for chat text, disable it for the rest of the UI.
