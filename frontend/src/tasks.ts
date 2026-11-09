import { call, listen } from './bridge';

export type TaskInfo = { id: number; name: string; cmd: string; level: string; state: string; exitCode?: number; elapsed: number };

export const listTasks = (session: string) => call<TaskInfo[]>('TaskService.List', session);
export const taskOutput = (session: string, id: number) => call<string>('TaskService.Output', session, id);
export const killTask = (session: string, id: number) => call<void>('TaskService.Kill', session, id);
export const onTasks = (handler: (session: string) => void) => listen('tasks:changed', (data) => handler(String(data)));

export function elapsed(seconds: number) {
  if (seconds < 60) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  return minutes < 60 ? `${minutes}m ${seconds % 60}s` : `${Math.floor(minutes / 60)}h ${minutes % 60}m`;
}
