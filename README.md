# Zelin

A minimal desktop AI agent built with Go, Wails 3.1, Svelte 5, and CodeMirror 6.

## Requirements

- Go 1.26+
- Node.js 22+ and npm
- Bash
- Linux: GTK4 and WebKitGTK 6.0

Arch / Asahi ALARM:

```sh
sudo pacman -S --needed nodejs npm gtk4 webkitgtk-6.0 base-devel pkgconf
```

## Run

Configure an OpenAI-compatible Chat Completions provider in the environment:

```sh
export OPENAI_API_KEY="your-key"
export OPENAI_MODEL="gpt-4.1-mini"
export OPENAI_BASE_URL="https://api.openai.com/v1"
npm --prefix frontend ci
npm --prefix frontend run build
go run .
```

Only the API key is required. The model and base URL above are defaults.
The provider must support function calling. Keys stay in the Go backend.
For a local provider without authentication, use a nonempty placeholder key.

Enter sends the current message; Shift+Enter inserts a newline. All hotkeys
match physical keys, so they work on any keyboard layout, and can be changed in
Commands -> Hotkeys. Defaults: Ctrl+K commands, Ctrl+S sessions, Ctrl+N new
session, Ctrl+C stops the current run when no text is selected.

User messages look like the input. Agent answers render only `**bold**`,
`*italic*`, and `` `code` `` without their markers; everything else stays literal.
Paste or drop images into the input, or use Commands -> Attach image. Each image
becomes a `[Image #1]` token that is deleted as one unit and opens on click.

Sessions, messages, images, execution output, drafts, and provider context are
stored in SQLite at `$XDG_DATA_HOME/zelin/sessions.db` (default
`~/.local/share/zelin/sessions.db`; `ZELIN_HOME` (or the legacy `ZASH_HOME`) overrides the whole data
directory). Several sessions can run at once. When a background session
finishes, a toast at the top center (and a sound) links back to it; both can be
turned off in Commands -> Settings, together with the theme and code display.
The sessions list shows running, finished, and failed sessions. Each session
gets a generated title from the title model chosen in Commands -> Models.

The bottom bar shows the model and reasoning effort on the left and the command
hotkey on the right. Models come from the provider's `/models` endpoint. Context
window and pricing come only from OpenRouter and LiteLLM metadata, cached for a
day in `metadata.json`; they are never used as providers. `OPENAI_CONTEXT_WINDOW`
overrides the context size. Every model offers the same efforts: none, low,
medium, high, xhigh, max, sent as `reasoning_effort`. Hidden models in
Commands -> Models do not appear in the bottom bar picker.

Commands -> Context and usage shows context tokens, session tokens, cache reads
and writes, and spending. Provider-reported costs are used when available;
otherwise the cost is estimated from metadata pricing and labeled as estimated.

Provider reasoning (`reasoning_content`, `reasoning`, `reasoning_details`) is
shown as thinking blocks, stored with the session, and sent back unchanged in
later requests. Responses are not streamed, so a thinking block appears when its
provider response completes. Executed JavaScript appears in code blocks; the
arrow folds a block and clicking the code toggles its output.

The system prompt lives in `SYSTEM.md` in the data directory. Edit it in
Commands -> Instructions -> System Prompt and save with Ctrl+S. The JavaScript
runtime contract is part of the tool description, so prompt edits never break
the tool.
Additional space
below the editor lets the latest message scroll to the center of the viewport.

## Agent runtime

The agent's only tool is `js_repl`, a persistent Goja JavaScript runtime.
Its only host function is synchronous `bash(cmd, timeout)`. No async/await,
console, Node.js APIs, or separate filesystem tools are provided.
The last expression is returned to the model as the tool result.

```js
var result = bash("pwd", 10);
result.stdout
```

Timeout is in seconds: default 30, maximum 120. Bash returns
`{stdout, stderr, exitCode, timedOut}`. Output is capped at 64 KB per stream.
Each tool evaluation is limited to 125 seconds; each agent turn to 5 minutes
and 12 model steps. JavaScript variables and tool history persist across turns.

**Commands run with your user permissions in the application's working directory.
This is not a sandbox or an approval system. Only use a provider and prompts you trust.**

## Build and checks

```sh
npm --prefix frontend test
npm --prefix frontend run build
go test ./...
go build -tags production -o bin/zelin .
```

Optional Wails CLI:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.1.0
wails3 task dev
wails3 task build
wails3 task test
```

## Production install

Install a separate production copy next to the development checkout:

```sh
wails3 task install:prod
```

This builds `~/.local/bin/zelin-prod` and a `Zelin` launcher entry. The production
copy keeps its own data, sessions, system prompt, and web storage in
`~/.local/share/zelin-prod`, so the development build in `~/.local/share/zelin`
can be changed and restarted freely while you work with the production one.
Run `wails3 task install:prod` again to update it, and
`wails3 task uninstall:prod` to remove the binary and launcher (data is kept).
On startup Zelin reads `.env` from its data directory (for the production copy
`~/.local/share/zelin-prod/.env`, mode 600) with `KEY=value` lines such as
`OPENAI_API_KEY` and `OPENAI_BASE_URL`. Variables already set in the environment
take precedence.

Restart the desktop app after code changes. `npm --prefix frontend run dev`
previews the UI in a browser; agent calls require the desktop bridge.
