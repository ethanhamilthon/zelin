export const MIN_ZOOM = 12, MAX_ZOOM = 26;

export const zoomed = (size: number, step: number) => Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, (size || 16) + step));
