<script lang="ts">
  import FileBrowser from './FileBrowser.svelte';
  import MentionMenu from './MentionMenu.svelte';
  import type { createFileRefs } from './file-refs.svelte';
  import type { createMentions } from './mention-controller.svelte';
  let { files, mentions, workspace, hidden = false }: {
    files: ReturnType<typeof createFileRefs>; mentions: ReturnType<typeof createMentions>; workspace: string; hidden?: boolean;
  } = $props();
</script>

{#if files.state.browsing && files.root && !hidden}
  <FileBrowser {workspace} root={files.root} start={files.state.from} onclose={() => { files.state.browsing = false; }} ondismiss={() => { files.state.browsing = false; }}
    onopenfile={(entry) => { files.state.browsing = false; files.openEntry(entry, 'file'); }} onopenapp={files.openApp} ontag={(entry) => files.tag(entry)} />
{/if}
{#if mentions.state.open && mentions.state.mention}
  <MentionMenu items={mentions.state.items} active={mentions.state.active} x={mentions.state.mention.x} y={mentions.state.mention.y} loading={mentions.state.loading}
    onpick={(entry) => files.pick(entry)} onhover={(index) => { mentions.state.active = index; }} />
{/if}
