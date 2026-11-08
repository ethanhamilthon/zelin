import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import App from './App.svelte';
import { agentInfo } from './agent';
import { currentSession, listSessions, createSession } from './sessions';
import { call } from './bridge';
import { placeholderInfo } from './agent-types';

vi.mock('./agent', () => ({ runAgent: vi.fn(), agentInfo: vi.fn(), listModels: vi.fn(), configureModel: vi.fn(), setTitleModel: vi.fn(), steerAgent: vi.fn(), stopAgent: vi.fn() }));
vi.mock('./sessions', () => ({ currentSession: vi.fn(), createSession: vi.fn(), openSession: vi.fn(), listSessions: vi.fn(), saveSession: vi.fn(), saveDraft: vi.fn(), saveAttachments: vi.fn(), saveRefs: vi.fn() }));
vi.mock('./bridge', () => ({ call: vi.fn(), desktop: () => false, listen: vi.fn(async () => () => {}) }));
vi.mock('./sound', () => ({ chime: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));

let app: ReturnType<typeof mount>;
let spaces: { id: string; name: string; icon: string; dir: string }[];
const $ = (selector: string) => document.querySelector<HTMLElement>(selector);
const all = (selector: string) => [...document.querySelectorAll<HTMLElement>(selector)];
const wait = async () => { await new Promise((r) => setTimeout(r, 0)); await tick(); };
const button = (text: string) => [...document.querySelectorAll<HTMLButtonElement>('button')].find((b) => b.textContent?.includes(text))!;
const key = (k: string, code: string, init: KeyboardEventInit = {}) => window.dispatchEvent(new KeyboardEvent('keydown', { key: k, code, cancelable: true, ...init }));
const row = (id: string, title: string, workspaceId: string) => ({ id, title, draft: '', messages: [], updated: '', workspaceId });
const type = (label: string, value: string) => { const input = $('[aria-label="' + label + '"]') as HTMLInputElement; input.value = value; input.dispatchEvent(new Event('input', { bubbles: true })); };
const enter = (label: string) => $('[aria-label="' + label + '"]')!.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true }));

beforeEach(async () => {
  vi.resetAllMocks();
  spaces = [{ id: 'default', name: 'Default', icon: '', dir: '/home' }, { id: 'w2', name: 'Docs', icon: '📚', dir: '/docs' }];
  Element.prototype.scrollIntoView = vi.fn();
  Object.defineProperties(Range.prototype, { getClientRects: { configurable: true, value: () => [] }, getBoundingClientRect: { configurable: true, value: () => new DOMRect() } });
  vi.mocked(currentSession).mockResolvedValue(row('a', 'Alpha', 'default'));
  vi.mocked(listSessions).mockResolvedValue([row('a', 'Alpha', 'default'), row('b', 'Beta', 'w2'), row('c', 'Gamma', '')]);
  vi.mocked(agentInfo).mockResolvedValue(placeholderInfo('Test'));
  vi.mocked(call).mockImplementation(async (method, ...args) => {
    if (method === 'WorkspaceService.List') return spaces;
    if (method === 'WorkspaceService.Create') { const item = { id: 'w3', name: String(args[0]).trim(), icon: String(args[1]), dir: String(args[2]) }; spaces = [...spaces, item]; return item; }
    if (method === 'WorkspaceService.Update') { const item = { id: String(args[0]), name: String(args[1]), icon: String(args[2]), dir: String(args[3]) }; spaces = spaces.map((s) => s.id === item.id ? item : s); return item; }
    if (method === 'WorkspaceService.Agents') return { path: '/docs/AGENTS.md', exists: false };
    if (method === 'WorkspaceService.NewSession') return { ...row('n', 'New session', String(args[0])) };
    return '';
  });
  app = mount(App, { target: document.body }); flushSync(); await wait();
  key('s', 'KeyS', { ctrlKey: true }); await wait();
});
afterEach(async () => { await unmount(app); document.body.replaceChildren(); });

it('groups sessions by workspace and puts legacy sessions in the default one', () => {
  const heads = all('.head .name').map((el) => el.textContent?.trim());
  expect(heads).toEqual(['▸Default', '📚Docs']);
  const groups = all('.body > *').map((el) => el.classList.contains('head') ? 'H:' + el.textContent!.replace(/[⋯\s]+$/, '').trim() : el.textContent!.trim().slice(0, 5));
  expect(groups).toEqual(['H:▸Default', 'Alpha', 'Gamma', 'H:📚Docs', 'Beta']);
});

it('gives every workspace row an independent settings button', async () => {
  expect(all('.head .more').map((el) => el.getAttribute('aria-label'))).toEqual(['Settings for Default', 'Settings for Docs']);
  all('.head .more')[1].click(); await wait();
  expect($('[role="dialog"]')?.getAttribute('aria-label')).toBe('Workspace settings');
  expect(($('[aria-label="Workspace name"]') as HTMLInputElement).value).toBe('Docs');
  expect(($('[aria-label="Working directory"]') as HTMLInputElement).value).toBe('/docs');
});

it('keeps session search working across groups', async () => {
  type('Search sessions', 'bet');
  await wait();
  expect(all('.head .name')).toHaveLength(1);
  expect(all('.body > button:not(.name)').map((el) => el.textContent?.trim().slice(0, 4))).toEqual(['Beta']);
});

it('starts a new session in the selected workspace', async () => {
  all('.head .name')[1].click(); await wait();
  expect(all('.head')[1].classList.contains('selected')).toBe(true);
  button('New session').click(); await wait();
  expect(call).toHaveBeenCalledWith('WorkspaceService.NewSession', 'w2');
  expect(createSession).not.toHaveBeenCalled();
});

it('creates a workspace and validates input', async () => {
  button('New workspace').click(); await wait();
  expect($('[role="dialog"]')?.getAttribute('aria-label')).toBe('New workspace');
  enter('Workspace name'); await wait();
  expect($('[role="alert"]')?.textContent).toBe('Name is required');
  type('Workspace name', '  Fresh '); type('Working directory', ' ');
  enter('Workspace name'); await wait();
  expect($('[role="alert"]')?.textContent).toBe('Directory is required');
  type('Working directory', '~/my project');
  enter('Workspace name'); await wait();
  expect(call).toHaveBeenCalledWith('WorkspaceService.Create', '  Fresh ', '', '~/my project');
  await vi.waitFor(() => expect(all('.head .name').map((el) => el.textContent?.trim())).toContain('▸Fresh'));
});

it('shows the backend error and keeps the form open', async () => {
  all('.head .more')[0].click(); await wait();
  vi.mocked(call).mockImplementation(async (method) => { if (method === 'WorkspaceService.Update') throw new Error('Directory is not accessible'); return ''; });
  type('Working directory', '/nope'); enter('Working directory'); await wait();
  expect($('[role="alert"]')?.textContent).toContain('Directory is not accessible');
  expect($('[aria-label="Workspace name"]')).not.toBeNull();
});

it('explains that existing sessions keep their directory', async () => {
  all('.head .more')[1].click(); await wait();
  expect($('.note')?.textContent).toContain('Existing sessions keep their own directory');
});

it('offers to create AGENTS.md and opens it in the in-app editor', async () => {
  all('.head .more')[1].click(); await wait();
  button('Create AGENTS.md').click(); await wait(); await wait();
  expect($('[aria-label="/docs/AGENTS.md"]')).not.toBeNull();
  expect(call).not.toHaveBeenCalledWith('FileService.Write', expect.anything(), expect.anything());
});

it('opens an existing AGENTS.md and writes only on explicit save', async () => {
  const base = vi.mocked(call).getMockImplementation()!;
  vi.mocked(call).mockImplementation(async (method, ...args) => {
    if (method === 'WorkspaceService.Agents') return { path: '/docs/AGENTS.md', exists: true };
    if (method === 'FileService.Read') return 'Be brief';
    return base(method, ...args);
  });
  all('.head .more')[1].click(); await wait();
  button('Open AGENTS.md').click(); await wait(); await wait();
  expect(($('.cm-content[aria-label="/docs/AGENTS.md"]'))?.textContent).toBe('Be brief');
  expect(call).not.toHaveBeenCalledWith('FileService.Write', expect.anything(), expect.anything());
  key('s', 'KeyS', { ctrlKey: true }); await wait();
  expect(call).toHaveBeenCalledWith('FileService.Write', '/docs/AGENTS.md', 'Be brief');
});
