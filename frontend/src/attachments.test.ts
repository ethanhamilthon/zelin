import { afterEach, expect, it, vi } from 'vitest';
import { mount, unmount } from 'svelte';
import { openAttachment } from './attachments';
import { call } from './bridge';
import { openFile } from './files';
import History from './History.svelte';
vi.mock('./bridge', () => ({ call: vi.fn() }));
vi.mock('./files', () => ({ isMarkdown: (path: string) => /.(md|markdown)$/i.test(path), openFile: vi.fn(), fileName: (path: string) => path.split('/').pop() }));
afterEach(() => { vi.clearAllMocks(); document.body.replaceChildren(); });
it('routes Markdown, images, and other files to their respective viewers', async () => {
 const edit = vi.fn(), image = vi.fn();
 await openAttachment({ name: 'Report', path: '/tmp/report.md' }, edit, image);
 expect(edit).toHaveBeenCalledWith(expect.objectContaining({ title: '/tmp/report.md' }));
 vi.mocked(call).mockResolvedValue('data:image/png;base64,a');
 await openAttachment({ name: 'Chart', path: '/tmp/chart.PNG' }, edit, image);
 expect(call).toHaveBeenCalledWith('FileService.Image', '/tmp/chart.PNG');
 expect(image).toHaveBeenCalledWith('data:image/png;base64,a', 1);
 await openAttachment({ name: 'PDF', path: '/tmp/report.pdf' }, edit, image);
 expect(openFile).toHaveBeenCalledWith('/tmp/report.pdf');
});
it('shows image previews and file type icons', async () => {
 vi.mocked(call).mockResolvedValue('data:image/png;base64,a');
 const { tick } = await import('svelte');
 const app = mount(History, { target: document.body, props: { messages: [{ role: 'assistant', content: 'Files', files: [{ name: 'Shot', path: '/tmp/a.png' }, { name: 'Code', path: '/tmp/a.go' }] }] } });
 try {
  await new Promise((resolve) => setTimeout(resolve, 0)); await tick();
  const buttons = document.querySelectorAll('[aria-label="Attached files"] button');
  expect(buttons[0].querySelector('img')?.getAttribute('src')).toBe('data:image/png;base64,a');
  expect(buttons[1].querySelector('svg')).not.toBeNull();
  expect(call).toHaveBeenCalledTimes(1);
 } finally { await unmount(app); }
});
it('renders clickable named attachments after an assistant answer', async () => {
 const edit = vi.fn();
 const app = mount(History, { target: document.body, props: { messages: [{ role: 'assistant', content: 'Your report', files: [{ name: 'Download report', path: '/tmp/report.md' }] }], ondocument: edit } });
 try {
  const button = document.querySelector<HTMLButtonElement>('[aria-label="Attached files"] button')!;
  expect(button.textContent?.trim()).toBe('Download report');
  expect(button.querySelector('svg')).not.toBeNull();
  button.click();
  expect(edit).toHaveBeenCalledWith(expect.objectContaining({ title: '/tmp/report.md' }));
 } finally { await unmount(app); }
});
