const LIMIT = 240;

export function oneLine(text: string, pattern = /\s*\n\s*/g) {
  return text.slice(0, LIMIT * 4).trim().replace(pattern, ' ').slice(0, LIMIT);
}
