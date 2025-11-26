<script>
  import Home from './components/pages/Home.svelte';
  import AtomsShowcase from './components/pages/AtomsShowcase.svelte';
  
  let currentPage = 'home'; // 'home' or 'showcase'
</script>

<!-- Navigation -->
<nav class="fixed top-0 left-0 right-0 bg-white shadow-md z-50">
  <div class="container mx-auto px-4 py-3 flex gap-4 justify-center">
    <button 
      on:click={() => currentPage = 'home'}
      class="px-6 py-2 rounded-lg font-semibold transition-colors {currentPage === 'home' ? 'bg-gray-800 text-white' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'}"
    >
      Home
    </button>
    <button 
      on:click={() => currentPage = 'showcase'}
      class="px-6 py-2 rounded-lg font-semibold transition-colors {currentPage === 'showcase' ? 'bg-gray-800 text-white' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'}"
    >
      Showcase
    </button>
  </div>
</nav>

<!-- Page Content -->
<div class="pt-16 overflow-x-hidden">
  {#if currentPage === 'home'}
    <Home />
  {:else if currentPage === 'showcase'}
    <AtomsShowcase />
  {/if}
</div>
