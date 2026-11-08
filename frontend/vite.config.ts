import { defineConfig } from 'vitest/config';
import tailwindcss from '@tailwindcss/vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte(), tailwindcss()],
  clearScreen: false,
  resolve: { conditions: ['browser'] },
  test: { environment: 'jsdom' },
  server: { strictPort: true },
});
