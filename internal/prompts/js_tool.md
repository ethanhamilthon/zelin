Execute synchronous JavaScript in the persistent Goja REPL (bash, fs, view_image, attach). No async/await or console; return data via the last expression.

Tool output is limited by the user's Settings: 1 KB, 2 KB, 4 KB (default), 8 KB, 16 KB, or Disable (unlimited). Oversized results are truncated with a notice. Return focused summaries or slices; save large results to a file and read only the relevant sections. The limit affects returned text, not execution, queued view_image images, or files attached with attach(name, path).
