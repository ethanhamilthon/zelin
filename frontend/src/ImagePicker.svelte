<script lang="ts">
  import { imageFiles, readImage } from './images';
  let { onpick }: { onpick: (url: string, name?: string) => void } = $props();
  let input: HTMLInputElement;
  export function open() { input.click(); }
  async function pick() {
    for (const file of imageFiles(input.files)) onpick(await readImage(file), file.name);
    input.value = '';
  }
</script>

<input type="file" accept="image/*" multiple bind:this={input} onchange={() => { void pick(); }} />

<style>
  input { display: none; }
</style>
