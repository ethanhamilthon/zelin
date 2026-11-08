import { entryRef, newId, refToken, type Entry, type Ref } from './refs';
import { openPath } from './browser';
import { fileDocument, type Document } from './documents';
import type { Editor } from './editor';
import type { createSessionController } from './session-controller.svelte';
import type { createWorkspaces } from './workspace-controller.svelte';
import type { createMentions } from './mention-controller.svelte';
import type { createModalStack } from './modal-stack.svelte';
import type { createToasts } from './toasts.svelte';

type Deps = {
  sessions: ReturnType<typeof createSessionController>; workspaces: ReturnType<typeof createWorkspaces>;
  mentions: ReturnType<typeof createMentions>; modals: ReturnType<typeof createModalStack>; toasts: ReturnType<typeof createToasts>;
  editor: () => Editor | undefined; openDoc: (doc: Document) => void; showImage: (url: string, number: number) => void;
};

export function createFileRefs({ sessions, workspaces, mentions, modals, toasts, editor, openDoc, showImage }: Deps) {
  const state = $state({ browsing: false, from: '' });
  const failed = (error: unknown) => toasts.push('', String(error));

  function add(make: (conversation: NonNullable<typeof sessions.current>) => Ref) {
    const conversation = sessions.current;
    if (!conversation) return '';
    const ref = make(conversation);
    conversation.refs = [...conversation.refs, ref];
    sessions.schedule();
    return refToken(ref);
  }
  const fromEntry = (entry: Entry) => add(() => entryRef(entry, workspaces.state.selected, sessions.state.id));

  function openEntry(entry: Entry, kind: string) {
    if (kind === 'dir') { modals.clear(); state.from = entry.path; state.browsing = true; return; }
    if (/\.(md|markdown)$/i.test(entry.name)) { modals.clear(); openDoc(fileDocument(entry.path)); return; }
    void openPath(workspaces.state.selected, entry.path).catch(failed);
  }

  function openRef(ref: Ref) {
    if (ref.data) {
      showImage(ref.data, Math.max(1, (sessions.current?.refs ?? []).filter((item) => item.data).indexOf(ref) + 1));
      return;
    }
    if (ref.path) openEntry({ name: ref.label, path: ref.path, kind: ref.kind === 'dir' ? 'dir' : 'file', size: 0 }, ref.kind);
  }

  const root = () => workspaces.state.items.find((item) => item.id === workspaces.state.selected)?.dir ?? '';

  function pick(entry: Entry) {
    const mention = mentions.state.mention;
    if (!mention || !editor()) return;
    if (entry.kind === 'dir') {
      const prefix = mention.query.includes('/') ? mention.query.slice(0, mention.query.lastIndexOf('/') + 1) : '';
      editor()!.replace(mention.from, mention.to, '@' + prefix + entry.name + '/');
      editor()!.focus();
      return;
    }
    editor()!.replace(mention.from, mention.to, fromEntry(entry) + ' ');
    mentions.close();
    editor()!.focus();
  }

  function menuKey(key: string) {
    if (!mentions.state.open) return false;
    if (key === 'ArrowDown') mentions.move(1);
    else if (key === 'ArrowUp') mentions.move(-1);
    else if (key === 'Escape') mentions.close();
    else if (mentions.chosen) pick(mentions.chosen);
    else return false;
    return true;
  }

  function toggle() {
    if (!root()) return;
    modals.clear(); state.from = ''; state.browsing = !state.browsing;
  }

  function tag(entry: Entry) {
    editor()?.insert(fromEntry(entry) + ' ');
    state.browsing = false; modals.clear(); editor()?.focus();
  }

  return {
    state,
    get root() { return root(); },
    image: (url: string) => add(() => ({ id: newId(), kind: 'image', label: 'image', data: url, workspaceId: workspaces.state.selected, sessionId: sessions.state.id })),
    openEntry,
    openRef,
    openById(id: string) { const ref = sessions.current?.refs.find((item) => item.id === id); if (ref) openRef(ref); },
    openApp: (entry: Entry) => { void openPath(workspaces.state.selected, entry.path).catch(failed); },
    toggle,
    tag,
    pick,
    menuKey,
  };
}
