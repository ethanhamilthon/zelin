import { call } from './bridge';

export const documents = [{ name: 'system', title: 'System Prompt', file: 'SYSTEM.md' }];

export const readDocument = (name: string) => call<string>('SettingsService.Read', name);
export const writeDocument = (name: string, text: string) => call<void>('SettingsService.Write', name, text);
export const documentTitle = (name: string) => documents.find((item) => item.name === name)?.title ?? name;
