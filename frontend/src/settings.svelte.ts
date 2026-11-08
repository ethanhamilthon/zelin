import { call } from './bridge';
import { applyTheme } from './themes';
import { defaultHotkeys } from './hotkeys';
import type { Action } from './hotkeys';

export type Settings = {
  theme: string; toast: boolean; sound: boolean; showCode: boolean; foldCode: boolean; stackSteps: boolean;
  autoCompact: boolean; compactThreshold: number;
  hiddenModels: string[]; hotkeys: Record<Action, string>;
};

export const defaultSettings = (): Settings => ({
  theme: 'zen', toast: true, sound: true, showCode: true, foldCode: false, stackSteps: false,
  autoCompact: true, compactThreshold: 60,
  hiddenModels: [], hotkeys: { ...defaultHotkeys },
});

export const loadSettings = () => call<string>('SettingsService.Load');
export const saveSettings = (value: string) => call<void>('SettingsService.Save', value);

export function createSettings() {
  const state = $state<Settings>(defaultSettings());

  function persist() {
    applyTheme(state.theme);
    void saveSettings(JSON.stringify(state)).catch(() => {});
  }

  return {
    state,
    async start() {
      try {
        const saved = JSON.parse((await loadSettings()) || '{}') as Partial<Settings>;
        Object.assign(state, saved, { hotkeys: { ...defaultHotkeys, ...saved.hotkeys } });
      } catch { }
      applyTheme(state.theme);
    },
    update(change: Partial<Settings>) { Object.assign(state, change); persist(); },
    bind(action: Action, hotkey: string) { state.hotkeys = { ...state.hotkeys, [action]: hotkey }; persist(); },
    toggleModel(id: string) {
      const hidden = state.hiddenModels.includes(id);
      state.hiddenModels = hidden ? state.hiddenModels.filter((item) => item !== id) : [...state.hiddenModels, id];
      persist();
    },
  };
}
