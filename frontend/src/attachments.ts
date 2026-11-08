import { call } from './bridge';
import { isMarkdown, openFile } from './files';
import { fileDocument, type Document } from './documents';
import type { FileAttachment } from './agent-types';
import { fileKind } from './file-kind';
export async function openAttachment(file: FileAttachment, edit: (doc: Document) => void, image?: (url: string, number: number) => void) {
 if (isMarkdown(file.path)) { edit(fileDocument(file.path)); return; }
 if (fileKind(file.path) === 'image') {
  const url = await call<string>('FileService.Image', file.path);
  image?.(url, 1);
  return;
 }
 await openFile(file.path);
}
