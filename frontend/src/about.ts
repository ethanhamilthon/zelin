import { call } from './bridge';

export type AboutInfo = {
  name: string; version: string; commit: string; date: string; description: string;
  runtime: string; platform: string; profile: string; daemon: string;
};

export type Notice = { name: string; license: string; url: string };

export const notices: Notice[] = [
  { name: 'Wails', license: 'MIT', url: 'https://github.com/wailsapp/wails' },
  { name: 'Goja', license: 'MIT', url: 'https://github.com/dop251/goja' },
  { name: 'go-sqlite3', license: 'MIT', url: 'https://github.com/mattn/go-sqlite3' },
  { name: 'Svelte', license: 'MIT', url: 'https://github.com/sveltejs/svelte' },
  { name: 'CodeMirror', license: 'MIT', url: 'https://github.com/codemirror/dev' },
  { name: 'marked', license: 'MIT', url: 'https://github.com/markedjs/marked' },
  { name: 'DOMPurify', license: 'MPL-2.0 OR Apache-2.0', url: 'https://github.com/cure53/DOMPurify' },
];

export const loadAbout = () => call<AboutInfo>('AboutService.Info');
