export type Toast = { id: number; session: string; text: string };

export function createToasts(timeout = 6000) {
  const items = $state<Toast[]>([]);
  let next = 0;
  const dismiss = (id: number) => {
    const index = items.findIndex((item) => item.id === id);
    if (index >= 0) items.splice(index, 1);
  };
  return {
    items,
    dismiss,
    push(session: string, text: string) {
      const id = ++next;
      items.push({ id, session, text });
      setTimeout(() => dismiss(id), timeout);
    },
  };
}
