export const imageToken = /\[Image #(\d+)\]/g;
const anyToken = /\[Image #(\d+)\]|\[@([A-Za-z0-9_-]{4,40})(?:\|([^\]\n]*))?\]/g;

export const token = (number: number) => `[Image #${number}]`;

export function imageFiles(list: FileList | null | undefined): File[] {
  return [...(list ?? [])].filter((file) => file.type.startsWith('image/'));
}

export function readImage(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(String(reader.result));
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });
}

export type Part = { text: string; image?: number; ref?: { id: string; label: string } };

export function splitTokens(text: string): Part[] {
  const parts: Part[] = [];
  let last = 0;
  for (const match of text.matchAll(anyToken)) {
    if (match.index > last) parts.push({ text: text.slice(last, match.index) });
    if (match[1]) parts.push({ text: match[0], image: Number(match[1]) });
    else parts.push({ text: '@' + (match[3] || match[2]), ref: { id: match[2], label: match[3] || match[2] } });
    last = match.index + match[0].length;
  }
  if (last < text.length) parts.push({ text: text.slice(last) });
  return parts;
}
