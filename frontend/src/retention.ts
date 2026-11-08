import type { Conversation } from './conversation';

export const KEEP = 12;

export function evictable(conversation: Conversation, saved?: { draft: string; attachments: string[]; refs: string }) {
  return !conversation.busy && !conversation.unread && !conversation.steering && !conversation.error && !!saved
    && conversation.draft === saved.draft && JSON.stringify(conversation.refs) === saved.refs
    && conversation.attachments.length === saved.attachments.length
    && conversation.attachments.every((item, index) => item === saved.attachments[index]);
}

export function stale(order: string[], current: string, keep: (id: string) => boolean) {
  const idle = order.filter((id) => id !== current && keep(id));
  return idle.slice(0, Math.max(0, order.length - KEEP));
}
