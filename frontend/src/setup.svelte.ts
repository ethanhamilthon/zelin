import type { Message } from './conversation';
import type { ModelOption } from './agent-types';
import { modelName, search } from './models';
import { checkModel, connectProvider, finishSetup, mask, openRouter, reason } from './setup';
import type { SetupCheck } from './setup';

export type Step = 'url' | 'key' | 'model' | 'title' | 'check';
export const total = 4;
const numbers: Record<Step, number> = { url: 1, key: 1, model: 2, title: 3, check: 4 };
const placeholders: Record<Step, string> = {
  url: 'Provider URL, or Enter for OpenRouter…', key: 'API key…', model: 'Filter models, ↑↓ to move, Enter to pick…',
  title: 'Filter models, ↑↓ to move, Enter to pick…', check: 'Checking…',
};
const ask = {
  url: `Provider URL? Press Enter to use OpenRouter (\`${openRouter}\`). Any OpenAI-compatible API works.`,
  key: 'API key? It is stored in the Zelin \`.env\` file.',
  model: '**Main model**. Pick it from the list below. Only this model will be enabled, you can enable more later in the Model manager.',
  title: '**Title model** for session names. A small, fast model works best. The main model is preselected.',
};
const result = (name: string, check: SetupCheck) => check.error
  ? `✗ **${name}** failed: ${check.error}` : `✓ **${name}** replied in ${(check.millis / 1000).toFixed(1)}s`;

export function createSetup(done: (model: string, models: ModelOption[]) => void) {
  const state = $state({
    step: 'url' as Step, busy: false, base: '', models: [] as ModelOption[], model: '', query: '', active: 0,
    messages: [{ role: 'assistant', content: `Welcome to **Zelin**. Let's connect a provider and pick your models.\n\n${ask.url}` }] as Message[],
  });
  const say = (content: string) => { state.messages = [...state.messages, { role: 'assistant', content }]; };
  const reply = (content: string) => { state.messages = [...state.messages, { role: 'user', content }]; };
  const name = (id: string) => modelName(state.models.find((model) => model.id === id)?.name ?? id);
  const picking = () => state.step === 'model' || state.step === 'title';
  const ordered = () => state.step === 'title' ? [...state.models].sort((a, b) => Number(b.id === state.model) - Number(a.id === state.model)) : state.models;
  const options = () => picking() ? search(ordered(), state.query, 200) : [];

  async function verify(title: string) {
    state.step = 'check';
    say(`Checking **${name(state.model)}**${title === state.model ? '' : ` and **${name(title)}**`}…`);
    const [main, titles] = await Promise.all([checkModel(state.model), title === state.model ? null : checkModel(title)]);
    const checks = [result(name(state.model), main), ...(titles ? [result(name(title), titles)] : [])].join('\n');
    if (main.error || titles?.error) {
      state.step = main.error ? 'model' : 'title';
      return say(`${checks}\n\n${main.error ? ask.model : ask.title}`);
    }
    await finishSetup(state.model, title);
    say(`${checks}\n\nAll set. Starting Zelin…`);
    done(state.model, state.models);
  }

  async function choose(id: string) {
    reply(name(id));
    state.query = '';
    state.active = 0;
    if (state.step === 'model') {
      state.model = id;
      state.step = 'title';
      return say(`Main model: **${name(id)}**.\n\n${ask.title}`);
    }
    await verify(id);
  }

  async function handle(text: string) {
    if (state.step === 'url') {
      state.base = text || openRouter;
      reply(state.base);
      state.step = 'key';
      return say(ask.key);
    }
    if (state.step === 'key') {
      if (!text) return say('The API key is required. ' + ask.key);
      reply(mask(text));
      try { state.models = await connectProvider(state.base, text); }
      catch (error) { state.step = 'url'; return say(`Could not connect: ${reason(error)}\n\n${ask.url}`); }
      state.step = 'model';
      return say(`Connected, ${state.models.length} models available.\n\n${ask.model}`);
    }
  }

  async function run(task: () => Promise<unknown>) {
    if (state.busy) return;
    state.busy = true;
    try { await task(); }
    catch (error) { if (state.step === 'check') state.step = 'title'; say(`Something went wrong: ${reason(error)}`); }
    finally { state.busy = false; }
  }

  return {
    state,
    get number() { return numbers[state.step]; },
    get placeholder() { return placeholders[state.step]; },
    get picking() { return picking(); },
    get options() { return options(); },
    filter(query: string) { state.query = query; state.active = 0; },
    move(step: number) { const count = options().length; if (count) state.active = (state.active + step + count) % count; },
    pick: (id: string) => run(() => choose(id)),
    submit(text: string) {
      if (!picking()) return run(() => handle(text.trim()));
      const option = options()[state.active];
      return option ? run(() => choose(option.id)) : Promise.resolve();
    },
  };
}
