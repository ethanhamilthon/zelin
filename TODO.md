# TODO

## Workspace settings button

- [ ] Make the three-dot workspace settings button in the Sessions window icon-sized instead of stretching across the row.

## Workspace creation

- [ ] Replace the free-text icon input with a picker offering 10 predefined icons.
- [ ] Replace manual working-directory path entry with a directory chooser.

## Status bar

- [ ] Show the current workspace name in the status bar activity area, before the Working/Done counters.
- [ ] Use the Commands / Ctrl+K color hierarchy throughout the bottom status bar: important content uses the Commands foreground color; secondary content uses the muted Ctrl+K color.

## UI appearance

- [ ] Add a font selector to the UI settings so users can change the application font.
- [ ] Introduce distinct bright accent colors for different UI elements, coordinated with the current theme and preserving readable contrast.

## JS tool output limit

- [ ] Add a configurable js_repl output size limit in Settings with options: 1 KB, 2 KB, 4 KB, 8 KB, 16 KB, and Disable (unlimited output).
- [ ] Apply the selected limit to tool output and clearly indicate when output is truncated.
- [ ] Document the configurable output limit and truncation behavior in the JS tool description.

## Markdown rendering

- [ ] Render unordered-list markers as a plain hyphen (-), matching the input editor, instead of round bullets.

## Files window

- [ ] Remove the leading arrows from folder rows.
- [ ] Replace text action buttons with compact icon-only buttons: @ for Tag, a folder icon for Open folder, and an arrow icon for Open in app / Open in default app.
- [ ] Preserve accessible labels and tooltips for icon-only actions.

## File mention picker

- [ ] Make all relevant matching files and folders accessible instead of limiting results to five.
- [ ] Show at most five items at once, with scrolling and arrow-key navigation through the remaining results.
- [ ] Support Tab to accept the selected suggestion without submitting the message.
- [ ] Support multi-stage selection: enter a folder, browse or filter its children, and continue selecting nested files or folders.
