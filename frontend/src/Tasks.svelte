<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import Modal from './Modal.svelte';
  import { elapsed, killTask, listTasks, onTasks, taskOutput, type TaskInfo } from './tasks';
  let { session, back = false, onclose, ondismiss }: { session: string; back?: boolean; onclose: () => void; ondismiss?: () => void } = $props();
  let items = $state<TaskInfo[]>([]), output = $state<{ id: number; text: string } | null>(null), error = $state('');
  const fail = (reason: unknown) => String(reason).replace(/^Error: /, '');
  const load = () => listTasks(session).then((list) => { items = list; }).catch((reason) => { error = fail(reason); });
  $effect(() => {
    void load();
    const timer = setInterval(() => { void load(); }, 1000);
    const off = onTasks((id) => { if (id === session) void load(); });
    return () => { clearInterval(timer); void off.then((stop) => stop()); };
  });
  async function show(id: number) {
    if (output?.id === id) { output = null; return; }
    try { output = { id, text: await taskOutput(session, id) || 'No output yet' }; } catch (reason) { error = fail(reason); }
  }
  async function kill(id: number) {
    try { await killTask(session, id); if (output?.id === id) output = null; await load(); } catch (reason) { error = fail(reason); }
  }
</script>

<Modal title="Tasks" {back} {onclose} {ondismiss}>
  {#each items as task (task.id)}
    <div class="task">
      <div class="info"><strong>#{task.id} {task.name}</strong><span>{elapsed(task.elapsed)} · {task.level}</span></div>
      <Button data-ui="ui-tasks" onclick={() => { void show(task.id); }}>{output?.id === task.id ? 'Hide' : 'Output'}</Button>
      <Button data-ui="ui-tasks" class="delete" onclick={() => { void kill(task.id); }}>Kill</Button>
    </div>
    {#if output?.id === task.id}<pre>{output.text}</pre>{/if}
  {:else}
    <p>No running tasks</p>
  {/each}
  {#if error}<p role="alert">{error}</p>{/if}
</Modal>

<style>
  .task { display: flex; align-items: center; gap: 8px; padding: 10px 0; border-bottom: 1px solid var(--border); }
  .info { display: grid; gap: 2px; flex: 1; min-width: 0; }
  strong { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  span, p { color: var(--dim); font-size: .85em; }
  p { margin: 0; font-size: 1em; }
  pre { margin: 8px 0; padding: 10px 12px; max-height: 320px; overflow: auto; border-radius: 4px; background: var(--code-bg); color: var(--muted); font-size: .85em; white-space: pre-wrap; overflow-wrap: anywhere; }
  :global(button[data-ui="ui-tasks"]) { width: auto !important; min-height: 0; padding: 6px 12px; border: 0; color: var(--muted); }
  :global(button[data-ui="ui-tasks"].delete) { color: var(--error); }
</style>
