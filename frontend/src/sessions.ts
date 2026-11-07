import type { Message } from './conversation';
import type { Usage } from './agent-types';
import { call } from './bridge';

export type Session = { id: string; title: string; messages: Message[]; draft: string; attachments?: string[]; updated: string; usage?: Usage };

export const currentSession = () => call<Session>('SessionService.Current');
export const createSession = () => call<Session>('SessionService.Create');
export const openSession = (id: string) => call<Session>('SessionService.Open', id);
export const listSessions = () => call<Session[]>('SessionService.List');
export const saveSession = (id: string, messages: Message[], draft: string, attachments: string[]) =>
  call<void>('SessionService.Save', id, messages, draft, attachments);
