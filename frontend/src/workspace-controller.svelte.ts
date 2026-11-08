import { agentsFile, createWorkspace, deleteWorkspace, listWorkspaces, updateWorkspace } from './workspaces';
import type { AgentsFile, Workspace } from './workspaces';

const reason = (error: unknown) => String(error).replace(/^Error: /, '');

export function createWorkspaces() {
  const state = $state({ items: [] as Workspace[], selected: '', editing: null as string | null, agents: null as AgentsFile | null });

  async function load() {
    try { state.items = await listWorkspaces(); } catch { state.items = []; }
    if (!state.items.some((item) => item.id === state.selected)) state.selected = state.items[0]?.id ?? '';
  }

  async function save(name: string, icon: string, dir: string): Promise<string | null> {
    if (!name.trim()) return 'Name is required';
    if (!dir.trim()) return 'Directory is required';
    try {
      const item = state.editing ? await updateWorkspace(state.editing, name, icon, dir) : await createWorkspace(name, icon, dir);
      await load();
      state.selected = item.id;
      return null;
    } catch (error) { return reason(error); }
  }

  async function remove(): Promise<string | null> {
    if (!state.editing) return null;
    try { await deleteWorkspace(state.editing); state.editing = null; await load(); return null; }
    catch (error) { return reason(error); }
  }

  return {
    state,
    load,
    save,
    remove,
    get current() { return state.items.find((item) => item.id === state.editing) ?? null; },
    follow(id?: string) { if (id && state.items.some((item) => item.id === id)) state.selected = id; },
    async edit(id: string | null) {
      state.editing = id;
      state.agents = null;
      if (id) state.agents = await agentsFile(id).catch(() => null);
    },
    async refreshAgents() { if (state.editing) state.agents = await agentsFile(state.editing).catch(() => null); },
  };
}
