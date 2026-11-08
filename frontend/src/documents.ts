import { call } from './bridge';
import { readFile, writeFile } from './files';

export type Document = { title: string; load: () => Promise<string>; save: (text: string) => Promise<void> };

export const documents = [{ name: 'system', title: 'System Prompt', file: 'SYSTEM.md' }];

export const readDocument = (name: string) => call<string>('SettingsService.Read', name);
export const writeDocument = (name: string, text: string) => call<void>('SettingsService.Write', name, text);
export const documentTitle = (name: string) => documents.find((item) => item.name === name)?.title ?? name;

export const namedDocument = (name: string): Document => ({
  title: documentTitle(name), load: () => readDocument(name), save: (text) => writeDocument(name, text),
});

export const fileDocument = (path: string): Document => ({
  title: path, load: () => readFile(path), save: (text) => writeFile(path, text),
});
