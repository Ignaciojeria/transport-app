<script>
  import BrandLogo from '../molecules/BrandLogo.svelte';
  
  const { className = '' } = $props();
</script>

<footer class={`bg-[#E8E4D9] py-8 sm:py-12 lg:py-16 ${className}`}>
  <div class="container mx-auto px-4 sm:px-6 lg:px-12 max-w-[1600px]">
    <div class="flex justify-center">
      <BrandLogo size="footer" />
    </div>
  </div>
</footer>

