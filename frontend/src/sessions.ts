import type { Message } from './conversation';
import type { Usage } from './agent-types';
import { call } from './bridge';
import type { Ref } from './refs';

export type Session = { id: string; title: string; messages: Message[]; draft: string; attachments?: string[]; refs?: Ref[]; updated: string; usage?: Usage; workspaceId?: string };

export const currentSession = () => call<Session>('SessionService.Current');
export const createSession = () => call<Session>('SessionService.Create');
export const openSession = (id: string) => call<Session>('SessionService.Open', id);
export const listSessions = () => call<Session[]>('SessionService.List');
export const saveDraft = (id: string, draft: string) => call<void>('SessionService.SaveDraft', id, draft);
export const saveAttachments = (id: string, attachments: string[]) => call<void>('SessionService.SaveAttachments', id, attachments);
export const saveRefs = (id: string, refs: Ref[]) => call<void>('SessionService.SaveRefs', id, refs);
export const saveSession = (id: string, messages: Message[], draft: string, attachments: string[]) =>
  call<void>('SessionService.Save', id, messages, draft, attachments);
