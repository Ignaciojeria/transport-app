<script>
  import { onMount } from 'svelte';
  import Home from './components/pages/Home.svelte';
  import TrackingView from './components/pages/TrackingView.svelte';
  import { initLanguage } from './lib/useLanguage';

  const isTrackingRoute = $derived(
    typeof window !== 'undefined' && window.location.pathname.startsWith('/track')
  );

  onMount(() => {
    initLanguage();
  });
</script>

{#if isTrackingRoute}
  <TrackingView />
{:else}
  <Home />
{/if}
