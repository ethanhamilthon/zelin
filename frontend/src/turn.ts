export type Turn = { finished: string; duration: number; cost: number | null; context: number | null; window: number };

export function elapsed(ms: number) {
  const seconds = Math.max(0, Math.round(ms / 1000));
  if (seconds < 60) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  return minutes < 60 ? `${minutes}m ${String(seconds % 60).padStart(2, '0')}s` : `${Math.floor(minutes / 60)}h ${String(minutes % 60).padStart(2, '0')}m`;
}
