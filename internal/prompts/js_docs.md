# JavaScript Runtime (js_repl)

Your only tool is `js_repl`, which evaluates synchronous JavaScript in a persistent Goja VM.

## Host APIs
- `bash(cmd, timeout)`: runs a shell command and returns `{stdout, stderr, exitCode, timedOut}`. Timeout is in seconds (default 30, maximum 120).
- `fs.read(path, {offset, limit})`: reads a text file with line numbers (default 2000 lines from line 1).
- `fs.write(path, text)`: creates or overwrites a file, creating parent directories.
- `fs.edit(path, oldText, newText)`: replaces exactly one occurrence of `oldText`.
- `fs.patch(path, diff)`: applies a unified diff with `@@` hunks; nothing is written if any hunk fails.
- `fs.grep(pattern, path, {glob})`: regex search returning `["file:line:text"]`, skipping `.git` and `node_modules`.
- `cd(path)`: changes the working directory of this session only and returns `{ok: true, path}` or `{error}`. Accepts absolute, relative, `..` and `~/` paths; `cd()` returns the current directory. It never changes the application's own directory.
- `view_image(path)`: loads an image file from disk (PNG, JPEG, GIF, WebP) and attaches it to the next model step so you can inspect it visually.

## Runtime Rules
- `fs` functions and `view_image` return `{error}` on failure; `fs.write`, `fs.edit`, and `fs.patch` return `{ok: true}` on success; `view_image` returns `{ok: true, path, mimeType}`.
- Paths may start with `~`. Relative paths in `fs`, `view_image` and `bash` resolve against the current `cd()` directory, which is saved with the session.
- A `cd` inside a `bash` command only affects that one command; use the JS `cd()` to change directories for later calls.
- No `async`/`await`, `Promise`, `console`, `require`, or Node.js/browser APIs.
- The value of the last expression is serialized and returned as the tool output.
- Top-level variables persist across turns; reuse variables or declare with `var` to avoid redeclaration errors.
- Commands run with the user's permissions in the current working directory.
