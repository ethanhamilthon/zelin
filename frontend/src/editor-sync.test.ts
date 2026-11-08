import { expect, it, vi } from 'vitest';
import { createEditor } from './editor';
it('applies remote drafts without reporting them as local edits', () => {
 const host = document.createElement('div'); document.body.append(host);
 const onChange = vi.fn(); const editor = createEditor(host, '', onChange);
 editor.setText('remote draft', false);
 expect(editor.text()).toBe('remote draft');
 expect(onChange).not.toHaveBeenCalled();
 editor.view.dispatch({ changes: { from: 12, insert: '!' } });
 expect(onChange).toHaveBeenCalledWith('remote draft!');
 editor.destroy(); host.remove();
});
