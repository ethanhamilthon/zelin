import type { ModelOption } from './agent-types';

export function search(models: ModelOption[], query: string, limit = 100) {
  const needle = query.trim().toLowerCase();
  return models.filter((model) => `${model.id} ${model.name} ${modelName(model.name)}`.toLowerCase().includes(needle)).slice(0, limit);
}

export function modelName(name: string) {
  return name.replace(/(\d)[-_](?=\d)/g, '$1.')
    .split(/[-_\s]+/).filter(Boolean)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1)).join(' ').replace(/gpt/gi, 'GPT');
}

export const effortName = (effort: string) => effort.charAt(0).toUpperCase() + effort.slice(1);
