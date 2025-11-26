<script>
  import Logo from '../atoms/Logo.svelte';
  import LogoText from '../atoms/LogoText.svelte';
  import Birds from '../atoms/Birds.svelte';
  import Slogan from '../atoms/Slogan.svelte';
  import BrandLogo from '../molecules/BrandLogo.svelte';
</script>

<div class="min-h-screen bg-gray-50 p-8">
  <div class="max-w-4xl mx-auto">
    <h1 class="text-3xl font-bold text-gray-900 mb-8">Atoms Showcase</h1>

    <section class="bg-white rounded-lg shadow p-6 mb-6">
      <h2 class="text-xl font-semibold text-gray-800 mb-4 border-b pb-2">Logo Atom</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8 items-end bg-gray-100 p-6 rounded-lg">
        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Micro (h-6)</span>
          <Logo className="h-6 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Navbar (h-12)</span>
          <Logo className="h-12 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Footer (h-24)</span>
          <Logo className="h-24 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Hero (h-64)</span>
          <Logo className="h-64 w-auto" />
        </div>
      </div>
    </section>

    <section class="bg-white rounded-lg shadow p-6 mb-6">
      <h2 class="text-xl font-semibold text-gray-800 mb-4 border-b pb-2">LogoText Atom</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8 items-end bg-gray-100 p-6 rounded-lg">
        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Micro (h-6)</span>
          <LogoText className="h-6 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Navbar (h-12)</span>
          <LogoText className="h-12 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Footer (h-24)</span>
          <LogoText className="h-24 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Hero (h-64)</span>
          <LogoText className="h-64 w-auto" />
        </div>
      </div>
    </section>

    <section class="bg-white rounded-lg shadow p-6 mb-6">
      <h2 class="text-xl font-semibold text-gray-800 mb-4 border-b pb-2">Birds Atom</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8 items-end bg-gray-100 p-6 rounded-lg">
        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Micro (h-6)</span>
          <Birds className="h-6 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Navbar (h-12)</span>
          <Birds className="h-12 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Footer (h-24)</span>
          <Birds className="h-24 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Hero (h-64)</span>
          <Birds className="h-64 w-auto" />
        </div>
      </div>
    </section>

    <section class="bg-white rounded-lg shadow p-6 mb-6">
      <h2 class="text-xl font-semibold text-gray-800 mb-4 border-b pb-2">Slogan Atom</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8 items-end bg-gray-100 p-6 rounded-lg">
        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Micro (h-6)</span>
          <Slogan className="h-6 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Navbar (h-12)</span>
          <Slogan className="h-12 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Footer (h-24)</span>
          <Slogan className="h-24 w-auto" />
        </div>

        <div class="flex flex-col items-center space-y-2">
          <span class="text-xs text-gray-500 uppercase tracking-wide">Hero (h-64)</span>
          <Slogan className="h-64 w-auto" />
        </div>
      </div>
    </section>

    <!-- Molecules Section -->
    <h1 class="text-3xl font-bold text-gray-900 mb-8 mt-12">Molecules Showcase</h1>

    <section class="bg-white rounded-lg shadow p-6 mb-6">
      <h2 class="text-xl font-semibold text-gray-800 mb-4 border-b pb-2">Brand Logo Molecule</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 items-end bg-[#E8E4D9] p-12 rounded-lg">
        <div class="flex flex-col items-center space-y-4">
          <span class="text-sm text-gray-600 uppercase tracking-wide font-semibold">Footer Size</span>
          <BrandLogo size="footer" />
        </div>

        <div class="flex flex-col items-center space-y-4">
          <span class="text-sm text-gray-600 uppercase tracking-wide font-semibold">Hero Size</span>
          <BrandLogo size="hero" />
        </div>
      </div>
    </section>
  </div>
</div>
