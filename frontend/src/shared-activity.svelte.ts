import { call } from './bridge';
type Activity = { running: boolean; unread: boolean };
export function createSharedActivity() {
 const state = $state({ items: {} as Record<string, Activity> });
 return {
 state,
 async sync() { state.items = await call<Record<string, Activity>>('AgentService.Activity'); },
 read(id: string) {
  if (!state.items[id]?.unread) return;
  state.items[id].unread = false;
  void call('AgentService.Read', id).catch(() => {});
 },
 };
}
