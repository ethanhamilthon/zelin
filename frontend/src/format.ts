export const unknown = 'Unknown';

export function tokens(value: number | null | undefined) {
  return value == null ? unknown : value.toLocaleString('en-US');
}

export function compact(value: number) {
  return Intl.NumberFormat('en-US', { notation: 'compact', maximumFractionDigits: 1 }).format(value);
}

export function cost(value: number | null | undefined) {
  if (value == null) return unknown;
  return `$${value.toFixed(value > 0 && value < 0.01 ? 4 : 2)}`;
}

export function perMillion(price: number | null) {
  return price == null ? unknown : `$${(price * 1e6).toFixed(2)}/M`;
}

export function percent(used: number | null | undefined, total: number) {
  return used == null || total <= 0 ? unknown : `${((used / total) * 100).toFixed(1)}%`;
}
