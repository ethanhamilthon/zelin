import { expect, it } from 'vitest';
import manager from './ModelManager.svelte?raw';
import icon from './MessageIcon.svelte?raw';
it('keeps model button styles isolated from conversation icons', () => {
 expect(manager).not.toContain(':global(.icon)');
 expect(manager).toContain('button[data-ui="ui-modelmanager"].icon');
 expect(icon).toContain('.message-icon {');
});
