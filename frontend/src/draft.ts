const key = 'zash:draft:v1';

export function loadDraft(storage: Pick<Storage, 'getItem'>): string {
  try {
    return storage.getItem(key) ?? '';
  } catch {
    return '';
  }
}

export function saveDraft(storage: Pick<Storage, 'setItem'>, text: string): boolean {
  try {
    storage.setItem(key, text);
    return true;
  } catch {
    return false;
  }
}
