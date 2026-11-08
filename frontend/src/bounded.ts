export const PAGE = 30;
export const ROWS = 200;
export const first = (total: number) => Math.max(0, total - PAGE);
export const earlier = (start: number) => Math.max(0, start - PAGE);
export const more = (hidden: number) => Math.min(hidden, PAGE);
