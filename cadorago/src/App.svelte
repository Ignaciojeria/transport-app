<script>
  import { onMount } from 'svelte';
  import Home from './components/pages/Home.svelte';
  import { initLanguage } from './lib/useLanguage';
  
  onMount(() => {
    initLanguage();
  });
</script>

<Home />
