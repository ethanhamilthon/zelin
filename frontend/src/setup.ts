import { call } from './bridge';
import type { AgentInfo, ModelOption } from './agent-types';

export type SetupCheck = { model: string; reply: string; error: string; millis: number };

export const setupRequired = () => call<boolean>('SetupService.Required').catch(() => false);
export const connectProvider = (base: string, key: string) => call<ModelOption[]>('SetupService.Connect', base, key);
export const checkModel = (model: string) => call<SetupCheck>('SetupService.Check', model);
export const finishSetup = (model: string, titleModel: string) => call<AgentInfo>('SetupService.Finish', model, titleModel);

export const openRouter = 'https://openrouter.ai/api/v1';
export const mask = (key: string) => key.length > 10 ? `${key.slice(0, 6)}…${key.slice(-4)}` : '•'.repeat(key.length);
export const reason = (error: unknown) => error instanceof Error ? error.message : String(error);
