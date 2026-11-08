import { expect, it } from 'vitest';
import { updateRemote, localRuns } from './remote-session';
import { blank } from './session-controller.svelte';
import type { Session } from './sessions';
const session: Session = { id: 'one', title: 'Chat', messages: [{ role: 'user', content: 'Hello' }], draft: 'shared', updated: '' };
it('restores running state and executions in a second window', () => {
 const target = blank();
 updateRemote(target, session, { started: 123, events: [{ kind: 'execution', execution: { id: 'tool', code: '1', output: '', state: 'running' } }] }, false);
 expect(target.busy).toBe(true);
 expect(target.started).toBe(123);
 expect(target.executions).toHaveLength(1);
 expect(target.messages).toEqual(session.messages);
 expect(target.draft).toBe('shared');
 updateRemote(target, session, null, false);
 expect(target.busy).toBe(false);
 expect(target.started).toBeUndefined();
});
it('preserves unsaved local input and does not overwrite the submitting window', () => {
 const target = { ...blank(), draft: 'typing' };
 updateRemote(target, session, null, true);
 expect(target.draft).toBe('typing');
 localRuns.add('one');
 updateRemote(target, { ...session, messages: [] }, null, false);
 expect(target.messages).toHaveLength(1);
 expect(target.draft).toBe('typing');
 localRuns.clear();
});
