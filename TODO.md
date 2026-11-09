# TODO

## Pasted images as file paths

Goal: the agent sees the full path of a pasted image in the text, e.g. `@/abs/path/image.png`, instead of `[Image #N]`.

- [ ] Save each pasted/dropped image to disk (session assets dir) when it is attached; keep the image data URL for the model as before
- [ ] Store the absolute path next to the image in the message (`Images` keeps the data, new field keeps the path)
- [ ] In `agent/begin.go` replace `[Image #N]` in the text sent to the model with `@<full path>`, keep the image part in the same message
- [ ] In the UI editor and history show `@nameoftheimage.png` instead of `[Image #N]` (generate a unique name on paste, e.g. `image-1.png`, keep the original name on drop)
- [ ] Make the image token in `images.ts` / `image-tokens.ts` carry the file name, keep it clickable to open the preview
- [ ] Keep old `[Image #N]` messages in saved sessions working
- [ ] Update `provider/message.go` so the token-to-image-part split still works with the path text
- [ ] Update compaction plan to keep image paths in the summary
- [ ] Cleanup: remove saved images together with the session
- [ ] Tests: path substitution, missing file, multiple images, compaction
