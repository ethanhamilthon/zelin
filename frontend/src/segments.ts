import type { Execution } from './agent-types';

export type Segment = { id: string; note?: Execution; steps: Execution[] };

export function segments(executions: Execution[]): Segment[] {
  const result: Segment[] = [];
  let steps: Execution[] = [];
  const flush = () => { if (steps.length) result.push({ id: steps[0].id, steps }); steps = []; };
  for (const execution of executions) {
    if (execution.kind !== 'message') { steps.push(execution); continue; }
    flush();
    result.push({ id: execution.id, note: execution, steps: [] });
  }
  flush();
  return result;
}
