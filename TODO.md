# TODO

## Layout and appearance

- [x] Remove the gap between the last message and the user input.
- [x] Use filled icons where possible and increase stroke width for line icons.
- [x] Replace the user message icon with a checkmark when the agent finishes the turn.
- [x] Set the application background to `#111111`.
- [x] Choose and apply a monospace font throughout the interface.

## Model and reasoning effort selection

- [x] Add a model selector.
- [x] Add a reasoning effort selector with options supported by the selected model.
- [x] Fetch model metadata from OpenRouter, including context limits, pricing, and supported parameters.
- [x] Investigate how to obtain supported reasoning efforts per model. Do not assume every model supports the same effort levels; define a fallback when metadata is unavailable.

## Context and session usage

- [x] Replace the context text in the status bar with an icon that opens a context and usage dropdown.
- [x] Expose the same context and usage view through the command menu.
- [x] Show current context token usage, the model maximum context size, and the percentage used.
- [x] Show cumulative input and output tokens for the session.
- [x] Show a cache summary, including cache reads and writes when reported by the provider.
- [x] Show cumulative session spending, using provider-reported costs when available and clearly labeling pricing-based estimates.
- [x] Persist session usage totals and avoid double counting requests. Show unavailable metrics as unknown rather than zero.

## Tool output

- [x] Remove horizontal scrolling from tool call output. Render it as regular wrapping text, breaking long lines when necessary.

## Chat text consistency

- [ ] Make chat message rendering match the input text one-to-one. Preserve literal Markdown markers instead of replacing them with formatted elements (for example, keep `-` rather than rendering a bullet).
- [ ] Use the input font size consistently throughout the interface, including chat messages, tool output, menus, and the status bar.
