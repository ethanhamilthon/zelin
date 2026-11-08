import { call } from './bridge';

export const readFile = (path: string) => call<string>('FileService.Read', path);
export const writeFile = (path: string, text: string) => call<void>('FileService.Write', path, text);
export const resolveFile = (path: string) => call<string>('FileService.Resolve', path);
export const openFile = (path: string) => call<void>('FileService.Open', path);
export const openUrl = (url: string) => call<void>('FileService.OpenURL', url);
export const fileName = (path: string) => path.split('/').filter(Boolean).pop() ?? path;
export const isMarkdown = (path: string) => /\.(md|markdown)$/i.test(path);
