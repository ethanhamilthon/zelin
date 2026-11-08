<script lang="ts">
  const paths = {
    plus: 'M12 5v14M5 12h14',
    dots: 'M5 12h.01M12 12h.01M19 12h.01',
    down: 'm6 9 6 6 6-6',
  };
  let { name, size = 16 }: { name: keyof typeof paths; size?: number } = $props();
</script>

<svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d={paths[name]} /></svg>
