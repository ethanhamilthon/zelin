# Zelin

A desktop AI agent with a Go backend, Svelte UI, and a JavaScript tool runtime powered by Goja. Built with Wails v3.

## Development

Requires Go 1.26+, Node.js with npm, Task, and the native dependencies for Wails v3 on your platform.

```sh
task dev
task test
task build
```

These tasks install frontend dependencies automatically. The build output is written to `bin/`. Restart the desktop app to apply changes when using `task dev`.

## Structure

- `frontend/src/`: Svelte UI and frontend tests
- `internal/agent/`: agent execution and context management
- `internal/app/`: desktop application services
- `internal/provider/`: model provider integration
- `internal/files/`, `internal/host/`, `internal/jsruntime/`: filesystem, shell, and JavaScript tools
- `internal/prompts/`: built-in agent prompts

## Local installation

`task install:prod` installs `zelin-prod` and a desktop launcher, with separate data in `~/.local/share/zelin-prod`.

`task uninstall:prod` removes the installed binary and launcher, keeping its data.
