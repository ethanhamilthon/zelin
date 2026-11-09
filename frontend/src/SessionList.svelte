<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import Modal from './Modal.svelte';
  import Icon from './Icon.svelte';
  import ShowMore from './ShowMore.svelte';
  import SearchInput from './SearchInput.svelte';
  import WorkspaceHeader from './WorkspaceHeader.svelte';
  import type { Session } from './sessions';
  import type { SessionStatus } from './session-controller.svelte';
  import { DEFAULT_WORKSPACE, type Workspace } from './workspaces';
  import { ROWS } from './bounded';
  let { sessions, active, status, steps = () => 0, back = false, workspaces = [], selected = '', onclose, ondismiss, onselect, onworkspacenew = () => {}, onsettings = () => {}, oncreate = () => {} }: {
    sessions: Session[]; active: string; back?: boolean; workspaces?: Workspace[]; selected?: string;
    status: (id: string) => SessionStatus; steps?: (id: string) => number;
    onclose: () => void; ondismiss?: () => void; onselect: (id: string) => void; onworkspacenew?: (id: string) => void;
    onsettings?: (id: string) => void; oncreate?: () => void;
  } = $props();
  let query = $state('');
  let folded = $state<Record<string, boolean>>({});
  let expanded = $state<Record<string, boolean>>({});
  let limit = $state(ROWS);
  const needle = $derived(query.trim().toLowerCase());
  const matches = $derived(needle ? sessions.filter((session) => session.title.toLowerCase().includes(needle)) : sessions);
  const rows = $derived(matches.slice(0, limit));
  const known = $derived(new Set(workspaces.map((item) => item.id)));
  const groups = $derived(workspaces.map((workspace) => ({
    workspace,
    items: rows.filter((session) => (known.has(session.workspaceId ?? '') ? session.workspaceId : DEFAULT_WORKSPACE) === workspace.id),
  })).filter((group) => !needle || group.items.length));
  $effect(() => { needle; limit = ROWS; });

  function formatWhen(iso: string) {
    if (!iso) return '';
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) return '';
    const now = new Date();
    if (date.toDateString() === now.toDateString()) {
      return date.toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' });
    }
    return date.getFullYear() === now.getFullYear()
      ? date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
      : date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
  }

  function meta(session: Session, state: SessionStatus) {
    const when = formatWhen(session.updated);
    const running = state === 'running' ? `(${steps(session.id)} Steps)` : '';
    return [when, running].filter(Boolean).join(' ');
  }
</script>

<Modal title="Workspaces" {back} {onclose} {ondismiss}>
  {#snippet actions()}<span class="actions"><Button data-ui="ui-sessionlist" variant="icon-text" onclick={oncreate}><Icon name="plus" />New workspace</Button></span>{/snippet}
  {#snippet top()}
    <SearchInput bind:value={query} label="Search sessions" onkeydown={(event) => {
      if (event.key === 'Enter' && matches.length) { event.preventDefault(); onselect(matches[0].id); }
    }} />
  {/snippet}
  {#if workspaces.length}
    {#each groups as group (group.workspace.id)}
      <WorkspaceHeader workspace={group.workspace} selected={group.workspace.id === selected} folded={!!folded[group.workspace.id] && !needle} onselect={() => { folded[group.workspace.id] = !folded[group.workspace.id]; }} onsettings={() => onsettings(group.workspace.id)} onnew={() => onworkspacenew(group.workspace.id)} />
      {#if !folded[group.workspace.id] || needle}
      {#each (expanded[group.workspace.id] || needle ? group.items : group.items.slice(0, 3)) as session (session.id)}
        {@const state = status(session.id)}
        <Button data-ui="ui-sessionlist" onclick={() => onselect(session.id)} aria-current={session.id === active ? 'true' : undefined}>
          <strong><i class={state} aria-hidden="true"></i>{session.title}</strong>
          <span>{meta(session, state)}</span>
        </Button>
      {/each}
      {#if !needle && !expanded[group.workspace.id] && group.items.length > 3}
        <ShowMore count={group.items.length - 3} onclick={() => expanded[group.workspace.id] = true} />
      {/if}
      {/if}
    {/each}
    {#if !rows.length}<p>No sessions found</p>{/if}
  {:else}
    {#each rows as session (session.id)}
      {@const state = status(session.id)}
      <Button data-ui="ui-sessionlist" onclick={() => onselect(session.id)} aria-current={session.id === active ? 'true' : undefined}>
        <strong><i class={state} aria-hidden="true"></i>{session.title}</strong>
        <span>{meta(session, state)}</span>
      </Button>
    {:else}
      <p>No sessions found</p>
    {/each}
  {/if}
  {#if matches.length > rows.length}<ShowMore count={matches.length - rows.length} onclick={() => { limit += ROWS; }} />{/if}
</Modal>

<style>
  p { margin: 0; padding: 12px; color: var(--dim); }
  :global(.actions) { display: inline-flex; gap: 12px; align-items: center; }
  strong { display: flex; align-items: center; gap: 10px; min-width: 0; font-weight: 400; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
  span { flex-shrink: 0; white-space: nowrap; }
  i { width: 7px; height: 7px; flex: 0 0 7px; border-radius: 50%; background: var(--border); }
  i.running { background: var(--warning); animation: pulse 1.2s ease-in-out infinite; }
  i.unread { background: var(--success); }
  i.error { background: #d16969; }
  @keyframes pulse { 50% { opacity: 0.3; } }
</style>
