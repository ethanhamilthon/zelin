import type { Message } from './conversation';

export type Session = { id: string; title: string; messages: Message[]; draft: string; updated: string };

async function call<T>(method: string, ...args: unknown[]): Promise<T> {
  const { Call } = await import('@wailsio/runtime');
  return Call.ByName(`main.SessionService.${method}`, ...args);
}

export const currentSession = () => call<Session>('Current');
export const createSession = () => call<Session>('Create');
export const openSession = (id: string) => call<Session>('Open', id);
export const listSessions = () => call<Session[]>('List');
export const saveSession = (id: string, messages: Message[], draft: string) => call<void>('Save', id, messages, draft);
