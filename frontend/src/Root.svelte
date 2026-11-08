<script lang="ts">
  import { onMount } from 'svelte';
  import { setupRequired } from './setup';
  import App from './App.svelte';
  import Onboarding from './Onboarding.svelte';
  let mode = $state<'loading' | 'setup' | 'app'>('loading');
  onMount(() => { void setupRequired().then((required) => { mode = required ? 'setup' : 'app'; }); });
</script>

{#if mode === 'setup'}
  <Onboarding ondone={() => { mode = 'app'; }} />
{:else if mode === 'app'}
  <App />
{/if}
