import { compact, cost } from './format';
import type { Usage } from './agent-types';

export type Turn = { finished: string; duration: number; cost: number | null; context: number | null; window: number };

export function elapsed(ms: number) {
  const seconds = Math.max(0, Math.round(ms / 1000));
  if (seconds < 60) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  return minutes < 60 ? `${minutes}m ${String(seconds % 60).padStart(2, '0')}s` : `${Math.floor(minutes / 60)}h ${String(minutes % 60).padStart(2, '0')}m`;
}

export function live(started: number, now: number, usage: Usage | null, baseline: Usage | null | undefined, window: number): Turn {
  const cost = usage?.cost != null ? Math.max(0, usage.cost - (baseline?.cost ?? 0)) : null;
  return { finished: new Date(now).toISOString(), duration: now - started, cost, context: usage?.context ?? null, window };
}

export function summary(turn: Turn) {
  const time = new Date(turn.finished).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' });
  const parts = [time, elapsed(turn.duration)];
  if (turn.cost != null) parts.push(cost(turn.cost));
  if (turn.context != null) parts.push(turn.window > 0 ? `${compact(turn.context)} / ${compact(turn.window)}` : compact(turn.context));
  return parts.join(' · ');
}
