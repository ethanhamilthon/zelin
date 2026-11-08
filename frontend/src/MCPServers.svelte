<script lang="ts">
  import Button from "./lib/ui/button/button.svelte";
  import Input from "./lib/ui/input/input.svelte";
  import Modal from './Modal.svelte';
  import { addServer, listServers, parseServer, removeServer, testServer, type MCPServer } from './mcp';
  let { back = false, onclose, ondismiss }: { back?: boolean; onclose: () => void; ondismiss?: () => void } = $props();
  let servers = $state<MCPServer[]>([]), status = $state<Record<string, string>>({});
  let name = $state(''), command = $state(''), env = $state(''), error = $state(''), adding = $state(false);
  const fail = (reason: unknown) => String(reason).replace(/^Error: /, '');
  $effect(() => { void listServers().then((list) => { servers = list; }).catch((reason) => { error = fail(reason); }); });
  async function add() {
    error = '';
    const server = parseServer(name, command, env);
    if (!server.name || !server.command) { error = 'Name and command are required'; return; }
    adding = true;
    try { servers = await addServer(server); name = command = env = ''; void check(server.name); }
    catch (reason) { error = fail(reason); }
    finally { adding = false; }
  }
  async function check(server: string) {
    status[server] = 'Connecting…';
    try { const tools = await testServer(server); status[server] = `${tools.length} ${tools.length === 1 ? 'tool' : 'tools'}`; }
    catch (reason) { status[server] = fail(reason); }
  }
  async function remove(server: string) {
    try { servers = await removeServer(server); delete status[server]; } catch (reason) { error = fail(reason); }
  }
</script>

<Modal title="MCP servers" {back} {onclose} {ondismiss}>
  {#each servers as server (server.name)}
    <div class="server">
      <div class="info"><strong>{server.name}</strong><code>{[server.command, ...server.args].join(' ')}</code>{#if status[server.name]}<span class="status">{status[server.name]}</span>{/if}</div>
      <Button data-ui="ui-mcp" onclick={() => { void check(server.name); }}>Test</Button>
      <Button data-ui="ui-mcp" class="delete" onclick={() => { void remove(server.name); }}>Remove</Button>
    </div>
  {:else}
    <p class="empty">No servers yet. The agent reaches them through mcp.tools() and mcp.call().</p>
  {/each}
  <h3>Add server</h3>
  <label>Name<Input bind:value={name} aria-label="Server name" placeholder="filesystem" /></label>
  <label>Command<Input bind:value={command} aria-label="Server command" placeholder="npx -y @modelcontextprotocol/server-filesystem ~/Projects" /></label>
  <label>Environment<textarea bind:value={env} aria-label="Server environment" rows="2" placeholder="API_KEY=value"></textarea></label>
  {#if error}<p role="alert">{error}</p>{/if}
  <div class="form-actions"><Button data-ui="ui-mcp" class="primary" onclick={() => { void add(); }} disabled={adding}>Add</Button></div>
</Modal>

<style>
  .server { display: flex; align-items: center; gap: 8px; padding: 10px 0; border-bottom: 1px solid var(--border); }
  .info { display: grid; gap: 2px; flex: 1; min-width: 0; }
  code { color: var(--dim); font: inherit; font-size: .85em; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .status { color: var(--faint); font-size: .85em; overflow-wrap: anywhere; }
  .empty, p { margin: 0 0 10px; color: var(--dim); }
  label { display: block; margin: 0 0 10px; color: var(--muted); }
  textarea { display: block; width: 100%; margin-top: 4px; padding: 10px 12px; border: 1px solid var(--border); border-radius: 4px; background: var(--background); color: var(--foreground); font: inherit; outline: none; resize: vertical; }
  textarea:focus { border-color: var(--dim); }
  .form-actions { display: flex; justify-content: flex-end; }
  :global(button[data-ui="ui-mcp"]) { width: auto !important; min-height: 0; padding: 6px 12px; border: 0; color: var(--muted); }
  :global(button[data-ui="ui-mcp"].primary) { color: var(--accent); }
  :global(button[data-ui="ui-mcp"].delete) { color: var(--error); }
</style>
