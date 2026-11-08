export type Action = 'commands' | 'models' | 'sessions' | 'actions' | 'files' | 'newSession' | 'stop' | 'send' | 'newline' | 'save' | 'zoomIn' | 'zoomOut';

export const actions: { id: Action; label: string }[] = [
  { id: 'commands', label: 'Open commands' },
  { id: 'models', label: 'Select model' },
  { id: 'sessions', label: 'Switch session' },
  { id: 'actions', label: 'Session actions' },
  { id: 'files', label: 'Browse files' },
  { id: 'newSession', label: 'New session' },
  { id: 'stop', label: 'Stop the current run' },
  { id: 'send', label: 'Send message' },
  { id: 'newline', label: 'Insert new line' },
  { id: 'save', label: 'Save in editor' },
  { id: 'zoomIn', label: 'Zoom in chat' },
  { id: 'zoomOut', label: 'Zoom out chat' },
];

export const defaultHotkeys: Record<Action, string> = {
  commands: 'Ctrl+K', models: 'Ctrl+M', sessions: 'Ctrl+S', actions: 'Ctrl+E', files: 'Ctrl+F', newSession: 'Ctrl+N', stop: 'Ctrl+C',
  send: 'Enter', newline: 'Shift+Enter', save: 'Ctrl+S', zoomIn: 'Ctrl+Equal', zoomOut: 'Ctrl+Minus',
};

const modifiers = ['Control', 'Shift', 'Alt', 'Meta'];

function keyName(event: KeyboardEvent) {
  const code = event.code;
  if (/^Key[A-Z]$/.test(code)) return code.slice(3);
  if (/^Digit\d$/.test(code)) return code.slice(5);
  if (code) return code.replace(/^Numpad/, 'Num');
  return event.key.length === 1 ? event.key.toUpperCase() : event.key;
}

export function combo(event: KeyboardEvent): string {
  if (modifiers.includes(event.key)) return '';
  const parts = [event.ctrlKey && 'Ctrl', event.altKey && 'Alt', event.shiftKey && 'Shift', event.metaKey && 'Meta'];
  return [...parts.filter(Boolean), keyName(event)].join('+');
}

export const matches = (event: KeyboardEvent, hotkey: string) => Boolean(hotkey) && combo(event) === hotkey;

export function codeMirrorKey(hotkey: string) {
  const parts = hotkey.split('+');
  const key = parts.pop() ?? '';
  return [...parts, key.length === 1 ? key.toLowerCase() : key].join('-');
}
