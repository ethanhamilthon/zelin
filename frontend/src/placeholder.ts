export function inputPlaceholder(busy: boolean, messages: number) {
  if (busy) return 'Agent is Working...';
  return messages ? 'Continue to Build' : "Let's Build Something Amazing";
}
