<script lang="ts">
  import { FileArchive, FileCode, FileBraces, FileText, File, FileImage, FileType } from '@lucide/svelte';
  import type { FileKind } from './file-kind';
  let { kind }: { kind: FileKind } = $props();
  const icons = { code: FileCode, text: FileText, data: FileBraces, archive: FileArchive, pdf: FileType, image: FileImage, file: File };
  const Glyph = $derived(icons[kind]);
</script>

<Glyph size={16} strokeWidth={2} aria-hidden="true" />
