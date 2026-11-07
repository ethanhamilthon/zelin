import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ModelList from './ModelList.svelte';
import EffortList from './EffortList.svelte';
import StatusBar from './StatusBar.svelte';
import { createModelController } from './model-controller.svelte';
import { configureModel, listModels } from './agent';
import { placeholderInfo } from './agent-types';
import type { ModelOption } from './agent-types';

vi.mock('./agent', () => ({ agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
const pricing = { prompt: 0.000001, completion: 0.000004, cacheRead: null, cacheWrite: null };
const models: ModelOption[] = [
  { id: 'openai/gpt-5', name: 'GPT-5', contextWindow: 400000, pricing, efforts: ['high', 'low'], defaultEffort: 'low' },
  { id: 'deepseek/chat', name: 'DeepSeek', contextWindow: 0, pricing: { ...pricing, prompt: null }, efforts: [], defaultEffort: '' },
];
let app: ReturnType<typeof mount> | undefined;
afterEach(async () => { if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });

describe('model selection', () => {
  it('filters models and selects the first match with Enter', () => {
    const onselect = vi.fn();
    app = mount(ModelList, { target: document.body, props: { models, active: 'openai/gpt-5', loading: false, error: '', onclose: vi.fn(), onselect } });
    flushSync();
    const input = document.querySelector('input')!;
    expect(document.activeElement).toBe(input);
    expect(document.querySelector('[aria-current="true"]')?.textContent).toContain('GPT-5');
    expect(document.body.textContent).toContain('400K · $1.00/M in · $4.00/M out');
    expect(document.body.textContent).toContain('? · Unknown in');
    input.value = 'seek'; input.dispatchEvent(new Event('input', { bubbles: true })); flushSync();
    expect(document.querySelectorAll('button')).toHaveLength(1);
    input.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true, cancelable: true }));
    expect(onselect).toHaveBeenCalledWith('deepseek/chat');
  });

  it('offers only supported efforts', () => {
    const onselect = vi.fn();
    app = mount(EffortList, { target: document.body, props: { efforts: ['high', 'low'], active: 'low', fallback: 'medium', onclose: vi.fn(), onselect } });
    flushSync();
    expect([...document.querySelectorAll('button')].map((button) => button.textContent?.trim())).toEqual(['Default medium', 'high', 'low Current']);
    document.querySelectorAll('button')[0].click();
    expect(onselect).toHaveBeenCalledWith('');
  });

  it('shows the effort control only for reasoning models', async () => {
    const info = { ...placeholderInfo('GPT-5'), model: 'openai/gpt-5', efforts: ['high'], effort: 'high' };
    const props = { busy: false, activity: '', info, onopen: vi.fn() };
    app = mount(StatusBar, { target: document.body, props }); flushSync();
    expect(document.querySelector('.effort')?.textContent).toBe('high');
    await unmount(app);
    app = mount(StatusBar, { target: document.body, props: { ...props, info: { ...info, efforts: [] } } }); flushSync();
    expect(document.querySelector('.effort')).toBeNull();
  });

  it('keeps the current effort when switching models and reports failures', async () => {
    vi.mocked(listModels).mockResolvedValue(models);
    vi.mocked(configureModel).mockImplementation(async (model, effort) => ({ ...placeholderInfo(model), model, effort }));
    const controller = createModelController();
    controller.state.info = { ...placeholderInfo('GPT-5'), model: 'openai/gpt-5', effort: 'high' };
    await controller.load(); await controller.load();
    expect(listModels).toHaveBeenCalledOnce();
    expect(await controller.select('deepseek/chat')).toBe(true);
    expect(configureModel).toHaveBeenCalledWith('deepseek/chat', 'high');
    vi.mocked(configureModel).mockRejectedValueOnce(new Error('Unknown model'));
    expect(await controller.effort('low')).toBe(false);
    expect(controller.state.error).toContain('Unknown model');
    expect(controller.state.info.model).toBe('deepseek/chat');
  });
});
