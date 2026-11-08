<span class="spinner" role="status" aria-label="Running"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round" aria-hidden="true"><path d="M21 12a9 9 0 1 1-6.2-8.55" /></svg></span>

<style>
  .spinner { display: inline-flex; width: 16px; height: 16px; color: var(--accent); }
  svg { display: block; }
  svg { animation: spin 0.9s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
  @media (prefers-reduced-motion: reduce) { svg { animation: none; } }
</style>
