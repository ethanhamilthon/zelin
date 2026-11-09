import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ModelList from './ModelList.svelte';
import ModelManager from './ModelManager.svelte';
import EffortList from './EffortList.svelte';
import { createModelController } from './model-controller.svelte';
import { configureModel, setTitleModel } from './agent';
import { placeholderInfo } from './agent-types';
import type { ModelOption } from './agent-types';

vi.mock('./agent', () => ({ agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn(), setTitleModel: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
const pricing = { prompt: 0.000001, completion: 0.000004, cacheRead: null, cacheWrite: null };
const models: ModelOption[] = [
  { id: 'openai/gpt-5', name: 'GPT-5', contextWindow: 400000, pricing },
  { id: 'deepseek/chat', name: 'DeepSeek', contextWindow: 0, pricing: { ...pricing, prompt: null } },
];
let app: ReturnType<typeof mount> | undefined;
afterEach(async () => { if (app) await unmount(app); app = undefined; document.body.replaceChildren(); });
const buttons = () => [...document.querySelectorAll<HTMLButtonElement>('button')];

describe('model selection', () => {
  it('filters models, hides disabled ones, and selects the first match with Enter', () => {
    const onselect = vi.fn();
    app = mount(ModelList, { target: document.body, props: { models, active: '', hidden: ['openai/gpt-5'], loading: false, error: '', onclose: vi.fn(), onselect } });
    flushSync();
    expect(buttons().map((item) => item.querySelector('strong')?.textContent)).toEqual(['DeepSeek']);
    expect(document.activeElement).toBe(document.querySelector('input'));
    document.querySelector('input')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true }));
    expect(onselect).toHaveBeenCalledWith('deepseek/chat');
  });

  it('toggles visibility and chooses the title model', () => {
    const ontoggle = vi.fn();
    const ontitle = vi.fn();
    app = mount(ModelManager, { target: document.body, props: { models, hidden: ['deepseek/chat'], titleModel: 'openai/gpt-5', loading: false, error: '', onclose: vi.fn(), ontoggle, ontitle } });
    flushSync();
    expect(document.querySelector('[aria-label="Hide GPT 5"]')?.getAttribute('aria-pressed')).toBe('true');
    expect(document.querySelector('[aria-label="Show DeepSeek"]')?.getAttribute('aria-pressed')).toBe('false');
    document.querySelector<HTMLButtonElement>('[aria-label="Show DeepSeek"]')!.click();
    expect(ontoggle).toHaveBeenCalledWith('deepseek/chat');
    expect(document.querySelector('[aria-label="Use GPT 5 for background tasks"]')?.getAttribute('aria-pressed')).toBe('true');
    document.querySelector<HTMLButtonElement>('[aria-label="Use DeepSeek for background tasks"]')!.click();
    expect(ontitle).toHaveBeenCalledWith('deepseek/chat');
  });

  it('offers the same efforts for every model', () => {
    const efforts = ['none', 'low', 'medium', 'high', 'xhigh', 'max'];
    app = mount(EffortList, { target: document.body, props: { efforts, active: 'high', onclose: vi.fn(), onselect: vi.fn() } });
    flushSync();
    expect(buttons().map((item) => item.firstChild?.textContent?.trim())).toEqual(['None', 'Low', 'Medium', 'High', 'Xhigh', 'Max']);
  });

  it('keeps the current effort when switching models and sets the title model', async () => {
    const info = { ...placeholderInfo('GPT-5'), model: 'openai/gpt-5', effort: 'xhigh', efforts: ['xhigh'] };
    vi.mocked(configureModel).mockResolvedValue(info);
    vi.mocked(setTitleModel).mockResolvedValue({ ...info, titleModel: 'deepseek/chat' });
    const controller = createModelController(() => 'session-1');
    controller.state.info = info;
    expect(await controller.select('deepseek/chat')).toBe(true);
    expect(configureModel).toHaveBeenCalledWith('session-1', 'deepseek/chat', 'xhigh');
    expect(await controller.titleModel('deepseek/chat')).toBe(true);
    expect(controller.state.info.titleModel).toBe('deepseek/chat');
  });
});
