export type Execution = {
  id: string;
  code: string;
  output: string;
  state: 'running' | 'complete';
};

export type Usage = {
  requests: number;
  unpriced: number;
  context: number | null;
  input: number | null;
  output: number | null;
  cacheRead: number | null;
  cacheWrite: number | null;
  cost: number | null;
  estimated: boolean;
};

export type AgentEvent =
  | { kind: 'activity'; activity: string }
  | { kind: 'execution'; execution: Execution }
  | { kind: 'usage'; usage: Usage };

export type Pricing = { prompt: number | null; completion: number | null; cacheRead: number | null; cacheWrite: number | null };
export type ModelOption = { id: string; name: string; contextWindow: number; pricing: Pricing; efforts: string[]; defaultEffort: string };
export type AgentInfo = Omit<ModelOption, 'id'> & { model: string; effort: string };
export type AgentResult = { content: string; executions: Execution[]; error?: string };

const noPricing: Pricing = { prompt: null, completion: null, cacheRead: null, cacheWrite: null };

export function placeholderInfo(name: string): AgentInfo {
  return { model: '', name, effort: '', contextWindow: 0, pricing: noPricing, efforts: [], defaultEffort: '' };
}
