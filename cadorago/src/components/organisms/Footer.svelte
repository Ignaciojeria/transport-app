<script>
  import Logo from '../atoms/Logo.svelte';
  
  const { className = '' } = $props();
</script>

<footer class={`bg-[#E8E4D9] py-8 sm:py-12 lg:py-16 ${className}`}>
  <div class="container mx-auto px-4 sm:px-6 lg:px-12 max-w-[1600px]">
    <div class="flex justify-center">
      <Logo className="h-24 sm:h-32 lg:h-40 w-auto" alt="Cadorago Logo" />
    </div>
  </div>
</footer>

