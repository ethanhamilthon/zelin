import { call } from './bridge';
import type { Session } from './sessions';

export type Workspace = { id: string; name: string; icon: string; dir: string };
export type AgentsFile = { path: string; exists: boolean };

export const DEFAULT_WORKSPACE = 'default';

export async function listWorkspaces(): Promise<Workspace[]> {
  const list = await call<Workspace[]>('WorkspaceService.List');
  return Array.isArray(list) ? list : [];
}
export const createWorkspace = (name: string, icon: string, dir: string) => call<Workspace>('WorkspaceService.Create', name, icon, dir);
export const updateWorkspace = (id: string, name: string, icon: string, dir: string) => call<Workspace>('WorkspaceService.Update', id, name, icon, dir);
export const deleteWorkspace = (id: string) => call<void>('WorkspaceService.Delete', id);
export const workspaceSession = (id: string) => call<Session>('WorkspaceService.NewSession', id);
export const agentsFile = (id: string) => call<AgentsFile>('WorkspaceService.Agents', id);
