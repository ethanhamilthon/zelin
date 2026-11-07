export const imageToken = /\[Image #(\d+)\]/g;

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

export function splitTokens(text: string): { text: string; image?: number }[] {
  const parts: { text: string; image?: number }[] = [];
  let last = 0;
  for (const match of text.matchAll(imageToken)) {
    if (match.index > last) parts.push({ text: text.slice(last, match.index) });
    parts.push({ text: match[0], image: Number(match[1]) });
    last = match.index + match[0].length;
  }
  if (last < text.length) parts.push({ text: text.slice(last) });
  return parts;
}
