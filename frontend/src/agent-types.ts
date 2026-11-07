export type Execution = {
  id: string;
  code: string;
  output: string;
  state: 'running' | 'complete';
};

export type AgentEvent =
  | { kind: 'activity'; activity: string }
  | { kind: 'execution'; execution: Execution };

export type AgentResult = { content: string; executions: Execution[]; error?: string };
export type AgentInfo = { model: string; contextWindow: number };
