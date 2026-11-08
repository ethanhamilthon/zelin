import { afterEach, expect, it, vi } from 'vitest';
import { flushSync, mount, tick, unmount } from 'svelte';
import { parseServer } from './mcp';
import { call } from './bridge';
import MCPServers from './MCPServers.svelte';

vi.mock('./bridge', () => ({ call: vi.fn() }));
vi.mock('svelte/transition', () => ({ fade: () => ({ duration: 0 }), fly: () => ({ duration: 0 }) }));
afterEach(() => { vi.clearAllMocks(); document.body.replaceChildren(); });
const wait = async () => { for (let i = 0; i < 4; i++) { await new Promise((resolve) => setTimeout(resolve, 0)); await tick(); } };

it('parses a command line and environment', () => {
  expect(parseServer(' fs ', 'npx -y "@scope/server fs" ~/p', 'KEY=a=b\nbad\n')).toEqual({ name: 'fs', command: 'npx', args: ['-y', '@scope/server fs', '~/p'], env: { KEY: 'a=b' } });
  expect(parseServer('x', '', '')).toEqual({ name: 'x', command: '', args: [] });
});

it('adds, tests and removes servers', async () => {
  const server = { name: 'fs', command: 'npx', args: ['srv'] };
  vi.mocked(call).mockImplementation(async (method) => {
    if (method === 'MCPService.List') return [];
    if (method === 'MCPService.Add') return [server];
    if (method === 'MCPService.Test') return [{ name: 'read' }, { name: 'write' }];
    if (method === 'MCPService.Remove') return [];
    throw new Error(method);
  });
  const app = mount(MCPServers, { target: document.body, props: { onclose: vi.fn() } });
  try {
    await wait();
    const set = (label: string, value: string) => { const input = document.querySelector<HTMLInputElement>(`[aria-label="${label}"]`)!; input.value = value; input.dispatchEvent(new Event('input', { bubbles: true })); };
    set('Server name', 'fs'); set('Server command', 'npx srv'); flushSync();
    [...document.querySelectorAll('button')].find((b) => b.textContent === 'Add')!.click();
    await wait();
    expect(call).toHaveBeenCalledWith('MCPService.Add', { name: 'fs', command: 'npx', args: ['srv'] });
    expect(document.querySelector('.server')?.textContent).toContain('2 tools');
    [...document.querySelectorAll('button')].find((b) => b.textContent === 'Remove')!.click();
    await wait();
    expect(document.querySelector('.server')).toBeNull();
  } finally { await unmount(app); }
});
