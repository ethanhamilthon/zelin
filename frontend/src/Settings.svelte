<script lang="ts">
  import { fonts } from './appearance';
  import Modal from './Modal.svelte';
  import { themes } from './themes';
  import type { Settings } from './settings.svelte';
  let { settings, back = false, onclose, ondismiss, onchange }: { settings: Settings; back?: boolean; onclose: () => void; ondismiss?: () => void; onchange: (change: Partial<Settings>) => void } = $props();
  const state = (value: boolean) => value ? 'On' : 'Off';
</script>

<Modal title="Settings" {back} {onclose} {ondismiss}>
  <button onclick={() => onchange({ toast: !settings.toast })} aria-pressed={settings.toast}>Toast when a session finishes <span>{state(settings.toast)}</span></button>
  <button onclick={() => onchange({ sound: !settings.sound })} aria-pressed={settings.sound}>Sound when a session finishes <span>{state(settings.sound)}</span></button>
  <button onclick={() => onchange({ showCode: !settings.showCode })} aria-pressed={settings.showCode}>Show code <span>{state(settings.showCode)}</span></button>
  {#if settings.showCode}
    <button onclick={() => onchange({ foldCode: !settings.foldCode })} aria-pressed={settings.foldCode}>Fold code blocks by default <span>{state(settings.foldCode)}</span></button>
  {/if}
  <button onclick={() => onchange({ stackSteps: !settings.stackSteps })} aria-pressed={settings.stackSteps}>Stack consecutive steps <span>{state(settings.stackSteps)}</span></button>
  <h3>Compaction</h3>
  <button onclick={() => onchange({ autoCompact: !settings.autoCompact })} aria-pressed={settings.autoCompact}>Auto compaction <span>{state(settings.autoCompact)}</span></button>
  <label>Context threshold (%) <input aria-label="Compaction threshold" type="number" min="1" max="100" value={settings.compactThreshold}
    onchange={(event) => { const value = Number(event.currentTarget.value); if (value >= 1 && value <= 100) onchange({ compactThreshold: value }); }} /></label>
  <h3>Appearance</h3>
  <label>Font<select aria-label="Application font" value={settings.font} onchange={(event) => onchange({ font: event.currentTarget.value })}>{#each fonts as font}<option value={font.id}>{font.name}</option>{/each}</select></label>
  <h3>JS tool</h3>
  <label>Output limit<select aria-label="JS tool output limit" value={settings.toolOutputLimit} onchange={(event) => onchange({ toolOutputLimit: Number(event.currentTarget.value) })}>{#each [1024, 2048, 4096, 8192, 16384, 0] as limit}<option value={limit}>{limit ? limit / 1024 + ' KB' : 'Disable'}</option>{/each}</select></label>
  <h3>Theme</h3>
  {#each themes as theme (theme.id)}
    <button onclick={() => onchange({ theme: theme.id })} aria-current={theme.id === settings.theme ? 'true' : undefined}>
      <i style:background={theme.colors[0]} style:border-color={theme.colors[2]}></i>{theme.name}
      <span>{theme.id === settings.theme ? 'Current' : theme.dark ? 'Dark' : 'Light'}</span>
    </button>
  {/each}
</Modal>

<style>
  label { display: flex; align-items: center; justify-content: space-between; padding: 12px; color: var(--muted); }
  select { max-width: 60%; padding: 6px; color: var(--foreground); background: var(--background); border: 1px solid var(--border); font: inherit; }
  input { width: 72px; padding: 6px; color: var(--foreground); background: var(--background); border: 1px solid var(--border); font: inherit; }
  h3 { margin: 12px 0 4px; padding: 0 12px; font: inherit; color: var(--dim); }
  i { display: inline-block; width: 14px; height: 14px; margin: 0 12px 0 0; border: 3px solid; border-radius: 50%; flex-shrink: 0; align-self: center; }
  button { justify-content: flex-start !important; }
  button span { margin-left: auto; }
</style>
