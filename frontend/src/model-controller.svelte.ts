import { agentInfo, configureModel, listModels, setTitleModel } from './agent';
import { placeholderInfo } from './agent-types';
import type { ModelOption } from './agent-types';

export function createModelController(session: () => string = () => '') {
  const state = $state({ info: placeholderInfo('Loading model…'), models: [] as ModelOption[], loading: false, error: '' });
  let request = 0;

  async function apply(change: Promise<typeof state.info>) {
    try { state.info = await change; state.error = ''; return true; }
    catch (error) { state.error = `Could not change model: ${String(error)}`; return false; }
  }

  return {
    state,
    async start() {
      const id = session(), current = ++request;
      try { const info = await agentInfo(id); if (current === request) state.info = info; }
      catch { if (current === request) state.info = placeholderInfo('Model unavailable'); }
    },
    async load() {
      if (state.models.length || state.loading) return;
      state.loading = true;
      try { state.models = await listModels(); state.error = ''; }
      catch (error) { state.error = `Could not load models: ${String(error)}`; }
      finally { state.loading = false; }
    },
    select: (id: string) => apply(configureModel(session(), id, state.info.effort)),
    effort: (effort: string) => apply(configureModel(session(), state.info.model, effort)),
    titleModel: (id: string) => apply(setTitleModel(id).then(({ titleModel }) => ({ ...state.info, titleModel }))),
  };
}
