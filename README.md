# Zash

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

Enter inserts a newline. Shift+Enter sends the current message.
Messages and answers become selectable, non-editable rendered Markdown.
The input stays editable during a run, so you can draft the next message.
Sending again is disabled until the current turn finishes. New drafts are never
replaced by a completed answer or an error. Failed input is restored when no
new draft exists and no commands were executed. Commands are never rolled back.
Sessions, messages, execution output, drafts, and provider context are stored in
SQLite at `$XDG_DATA_HOME/zash/sessions.db` (default
`~/.local/share/zash/sessions.db`). The last active session reopens on startup.
Ctrl+S opens saved sessions; the command menu can create a new one. Switching
sessions is disabled during an agent turn. JavaScript bindings are isolated per
session and survive switching, but reset when the application restarts; the
agent is informed of that reset. Completed answers are committed before the
backend returns them to the UI. Drafts are saved shortly after typing stops,
and drafts typed during a run are saved when that turn finishes.

The bottom bar shows activity, the active model, and its context window.
Context size comes from the provider's `/models` metadata or the explicit
`OPENAI_CONTEXT_WINDOW` token count. Unsupported metadata is shown as Unknown,
not an estimated or invented limit. Ctrl+K opens commands; Escape closes them.
Executed JavaScript appears as dimmed, syntax-highlighted plain code; click it
to toggle its output. The bottom status dot is gray when idle, blinking green
while running, and blue for an unread response until you scroll to its end.
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
go build -tags production -o bin/zash .
```

Optional Wails CLI:

```sh
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.1.0
wails3 task dev
wails3 task build
wails3 task test
```

Restart the desktop app after code changes. `npm --prefix frontend run dev`
previews the UI in a browser; agent calls require the desktop bridge.
