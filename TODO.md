# TODO

1. [x] Move root-level application services and tests into internal/app, leaving only the entry point and embedded assets at the root. Update Wails binding paths.
   - [ ] Further split application services by responsibility instead of treating internal/app as the final package layout.
2. [x] Share one backend across multiple native windows using Wails single-instance support.
   - [x] Synchronize persisted sessions and drafts, workspace changes, settings, model selection, and active-run snapshots.
   - [x] Preserve unsaved local input during refresh and reject duplicate runs of the same session.
   - [x] Verify native Wayland windows and second-launch forwarding on Hyprland.
   - [ ] Complete manual multi-window testing, including concurrent edits and closing the submitting window during a run.
   - [ ] Add revision-based conflict handling for simultaneous edits; persisted changes currently use last-write-wins.
3. [x] Display GPT in uppercase in model names without changing provider model IDs.
4. [x] Add host function attach(name, path) to share files in agent responses with buttons for Markdown, images, and external apps.
5. [x] Migrate UI controls to shadcn-svelte components (bits-ui), use SF Pro Display font, and remove the font picker.
