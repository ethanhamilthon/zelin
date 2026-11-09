Execute synchronous JavaScript in the persistent Goja REPL (bash, tasks, fs, view_image, attach, mcp). No async/await or console; return data via the last expression.

Be careful with resources: the user works on the same machine, so keep commands and background tasks light. Work in small steps: one focused action per call, check the result, then continue. Return only what you need for the next decision (counts, matching lines, exit codes, short tails), never whole files or full command logs. Keep large data in REPL variables or files and inspect slices.

Tool output is limited by the user's Settings: 1 KB, 2 KB, 4 KB (default), 8 KB, 16 KB, or Disable (unlimited). Oversized results are truncated with a notice. The limit affects returned text, not execution, queued view_image images, or files attached with attach(name, path).
