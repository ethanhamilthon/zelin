import { afterEach, expect, it, vi } from 'vitest';
import { call } from './bridge';
const byName = vi.fn().mockResolvedValue([]);
vi.mock('@wailsio/runtime', () => ({ Call: { ByName: (...args: unknown[]) => byName(...args) } }));
afterEach(() => { delete (window as Window & { wails?: unknown }).wails; byName.mockClear(); });
it('calls moved services using the full Go import path', async () => {
 (window as Window & { wails?: unknown }).wails = { invoke: true };
 await call('SessionService.List');
 expect(byName).toHaveBeenCalledWith('zelin/internal/app.SessionService.List');
});
