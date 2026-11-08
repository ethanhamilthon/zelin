type BridgeWindow = Window & {
  chrome?: { webview?: { postMessage?: unknown } };
  webkit?: { messageHandlers?: { external?: { postMessage?: unknown } } };
  wails?: { invoke?: unknown };
};

export function desktop() {
  const bridge = window as BridgeWindow;
  return Boolean(bridge.webkit?.messageHandlers?.external?.postMessage ||
    bridge.chrome?.webview?.postMessage || bridge.wails?.invoke);
}

export async function call<T>(method: string, ...args: unknown[]): Promise<T> {
  if (!desktop()) throw new Error('Available in the desktop app only');
  const { Call } = await import('@wailsio/runtime');
  return Call.ByName(`zelin/internal/app.${method}`, ...args);
}

export async function listen(name: string, handler: (data: unknown) => void) {
  if (!desktop()) return () => {};
  const { Events } = await import('@wailsio/runtime');
  return Events.On(name, (event) => handler(event.data));
}
