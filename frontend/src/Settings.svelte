<script lang="ts">
  import Select from "./lib/ui/select.svelte";
  import Switch from "./lib/ui/switch.svelte";
  import Input from "./lib/ui/input/input.svelte";
  import Separator from "./lib/ui/separator.svelte";
  import Button from "./lib/ui/button/button.svelte";
  import Modal from './Modal.svelte';
  import { themes } from './themes';
  import type { Settings } from './settings.svelte';
  let { settings, back = false, onclose, ondismiss, onchange }: { settings: Settings; back?: boolean; onclose: () => void; ondismiss?: () => void; onchange: (change: Partial<Settings>) => void } = $props();
  const state = (value: boolean) => value ? 'On' : 'Off';
</script>

<Modal title="Settings" {back} {onclose} {ondismiss}>
  <Button data-ui="ui-settings" onclick={() => onchange({ toast: !settings.toast })} aria-pressed={settings.toast}>Toast when a session finishes <span class="setting-state">{state(settings.toast)}<Switch checked={settings.toast} label="toast" /></span></Button>
  <Button data-ui="ui-settings" onclick={() => onchange({ sound: !settings.sound })} aria-pressed={settings.sound}>Sound when a session finishes <span class="setting-state">{state(settings.sound)}<Switch checked={settings.sound} label="sound" /></span></Button>
  <Button data-ui="ui-settings" onclick={() => onchange({ showCode: !settings.showCode })} aria-pressed={settings.showCode}>Show code <span class="setting-state">{state(settings.showCode)}<Switch checked={settings.showCode} label="showCode" /></span></Button>
  {#if settings.showCode}
    <Button data-ui="ui-settings" onclick={() => onchange({ foldCode: !settings.foldCode })} aria-pressed={settings.foldCode}>Fold code blocks by default <span class="setting-state">{state(settings.foldCode)}<Switch checked={settings.foldCode} label="foldCode" /></span></Button>
  {/if}
  <Button data-ui="ui-settings" onclick={() => onchange({ stackSteps: !settings.stackSteps })} aria-pressed={settings.stackSteps}>Stack consecutive steps <span class="setting-state">{state(settings.stackSteps)}<Switch checked={settings.stackSteps} label="stackSteps" /></span></Button>
  <Separator /><h3>Compaction</h3>
  <Button data-ui="ui-settings" onclick={() => onchange({ autoCompact: !settings.autoCompact })} aria-pressed={settings.autoCompact}>Auto compaction <span class="setting-state">{state(settings.autoCompact)}<Switch checked={settings.autoCompact} label="autoCompact" /></span></Button>
  <label>Context threshold (%) <Input class="threshold-input" aria-label="Compaction threshold" type="number" min="1" max="100" value={settings.compactThreshold}
    onchange={(event) => { const value = Number(event.currentTarget.value); if (value >= 1 && value <= 100) onchange({ compactThreshold: value }); }} /></label>
  <h3>JS tool</h3>
  <label>Output limit<Select label="JS tool output limit" value={String(settings.toolOutputLimit)} options={[1024, 2048, 4096, 8192, 16384, 0].map((limit) => ({ value: String(limit), label: limit ? limit / 1024 + " KB" : "Disable" }))} onchange={(value) => onchange({ toolOutputLimit: Number(value) })} /></label>
  <h3>Theme</h3>
  {#each themes as theme (theme.id)}
    <Button data-ui="ui-settings" onclick={() => onchange({ theme: theme.id })} aria-current={theme.id === settings.theme ? 'true' : undefined}>
      <i style:background={theme.colors[0]} style:border-color={theme.colors[2]}></i>{theme.name}
      <span>{theme.id === settings.theme ? 'Current' : theme.dark ? 'Dark' : 'Light'}</span>
    </Button>
  {/each}
</Modal>

<style>
 :global(.setting-state) { display: inline-flex; align-items: center; gap: 12px; }
  label { display: flex; align-items: center; justify-content: space-between; padding: 12px; color: var(--muted); }
  :global(input.threshold-input) { width: 72px; padding: 6px; color: var(--foreground); background: var(--background); border: 1px solid var(--border); font: inherit; }
  h3 { margin: 12px 0 4px; padding: 0 12px; font: inherit; color: var(--dim); }
  i { display: inline-block; width: 14px; height: 14px; margin: 0 12px 0 0; border: 3px solid; border-radius: 50%; flex-shrink: 0; align-self: center; }
  :global(button[data-ui="ui-settings"]) { justify-content: flex-start !important; }
  :global(button[data-ui="ui-settings"]) span { margin-left: auto; }
</style>
