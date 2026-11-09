<script lang="ts">
  import Select from "./lib/ui/select.svelte";
  import Switch from "./lib/ui/switch.svelte";
  import Input from "./lib/ui/input/input.svelte";
  import Separator from "./lib/ui/separator.svelte";
  import Button from "./lib/ui/button/button.svelte";
  import Modal from './Modal.svelte';
  import type { Settings } from './settings.svelte';
  let { settings, back = false, onclose, ondismiss, onchange }: { settings: Settings; back?: boolean; onclose: () => void; ondismiss?: () => void; onchange: (change: Partial<Settings>) => void } = $props();
  const state = (value: boolean) => value ? 'On' : 'Off';
</script>

<Modal title="Settings" {back} {onclose} {ondismiss}>
  <Button data-ui="ui-settings" onclick={() => onchange({ toast: !settings.toast })} aria-pressed={settings.toast}>Toast when a session finishes <span class="setting-state">{state(settings.toast)}<Switch checked={settings.toast} label="toast" /></span></Button>
  <Button data-ui="ui-settings" onclick={() => onchange({ sound: !settings.sound })} aria-pressed={settings.sound}>Sound when a session finishes <span class="setting-state">{state(settings.sound)}<Switch checked={settings.sound} label="sound" /></span></Button>
  <Separator /><h3>Compaction</h3>
  <Button data-ui="ui-settings" onclick={() => onchange({ autoCompact: !settings.autoCompact })} aria-pressed={settings.autoCompact}>Auto compaction <span class="setting-state">{state(settings.autoCompact)}<Switch checked={settings.autoCompact} label="autoCompact" /></span></Button>
  <label>Context threshold (%) <Input class="threshold-input" aria-label="Compaction threshold" type="number" min="1" max="100" value={settings.compactThreshold}
    onchange={(event) => { const value = Number(event.currentTarget.value); if (value >= 1 && value <= 100) onchange({ compactThreshold: value }); }} /></label>
  <h3>JS tool</h3>
  <label>Output limit<Select label="JS tool output limit" value={String(settings.toolOutputLimit)} options={[1024, 2048, 4096, 8192, 16384, 0].map((limit) => ({ value: String(limit), label: limit ? limit / 1024 + " KB" : "Disable" }))} onchange={(value) => onchange({ toolOutputLimit: Number(value) })} /></label>
</Modal>

<style>
 :global(.setting-state) { display: inline-flex; align-items: center; gap: 12px; }
  label { display: flex; align-items: center; justify-content: space-between; padding: 12px; color: var(--muted); }
  :global(input.threshold-input) { width: 72px; padding: 6px; color: var(--foreground); background: var(--background); border: 1px solid var(--border); font: inherit; }
  h3 { margin: 12px 0 4px; padding: 0 12px; font: inherit; color: var(--dim); }
  :global(button[data-ui="ui-settings"]) { justify-content: flex-start !important; }
  :global(button[data-ui="ui-settings"]) span { margin-left: auto; }
</style>
