import { expect, it, vi } from 'vitest';
import { mount, unmount, tick } from 'svelte';
import { EditorView } from '@codemirror/view';
import App from './App.svelte';
import { currentSession, listSessions, saveDraft } from './sessions';
import { call, listen } from './bridge';
import { agentInfo } from './agent';
import { placeholderInfo } from './agent-types';
vi.mock('./agent', () => ({ agentInfo: vi.fn(), runAgent: vi.fn(), listModels: vi.fn(), configureModel: vi.fn(), setTitleModel: vi.fn(), stopAgent: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));
vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => false, listen: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
const wait = async (ms = 0) => { await new Promise((resolve) => setTimeout(resolve, ms)); await tick(); };
it('synchronizes actual editor contents between two mounted windows', async () => {
 Element.prototype.scrollIntoView = vi.fn();
 Object.defineProperties(Range.prototype, { getClientRects: { configurable: true, value: () => [] }, getBoundingClientRect: { configurable: true, value: () => ({ left: 0, right: 0, top: 0, bottom: 0 }) } });
 let draft = '';
 const callbacks = new Map<string, Array<(data: unknown) => void>>();
 vi.mocked(listen).mockImplementation(async (name, callback) => { callbacks.set(name, [...(callbacks.get(name) ?? []), callback]); return () => {}; });
 const record = () => ({ id: 'shared', title: 'Shared', draft, messages: [], attachments: [], refs: [], updated: '' });
 vi.mocked(currentSession).mockImplementation(async () => record());
 vi.mocked(listSessions).mockResolvedValue([]);
 vi.mocked(agentInfo).mockResolvedValue(placeholderInfo('GPT'));
 vi.mocked(saveDraft).mockImplementation(async (_id, value) => { draft = value; for (const callback of callbacks.get('session:changed') ?? []) callback('shared'); });
 vi.mocked(call).mockImplementation(async (method) => {
  if (method === 'SessionService.Get') return record();
  if (method === 'AgentService.Activity') return {};
  if (method === 'AgentService.Active') return null;
  if (method === 'SettingsService.Load') return '{}';
  if (method === 'WorkspaceService.List') return [];
  return '';
 });
 const firstHost = document.createElement('div'), secondHost = document.createElement('div');
 document.body.append(firstHost, secondHost);
 const first = mount(App, { target: firstHost }), second = mount(App, { target: secondHost });
 try {
  await wait(); await wait();
  const firstEditor = EditorView.findFromDOM(firstHost.querySelector('.cm-editor')!)!;
  const secondEditor = EditorView.findFromDOM(secondHost.querySelector('.cm-editor')!)!;
  firstEditor.dispatch({ changes: { from: 0, insert: 'from first window' } });
  await wait(400); await wait();
  expect(secondEditor.state.doc.toString()).toBe('from first window');
  secondEditor.dispatch({ changes: { from: 0, to: secondEditor.state.doc.length, insert: 'from second window' } });
  await wait(400); await wait();
  expect(firstEditor.state.doc.toString()).toBe('from second window');
 } finally { await unmount(first); await unmount(second); firstHost.remove(); secondHost.remove(); }
});
