<script>
  import { onMount } from 'svelte';
  import ModernHeader from './ModernHeader.svelte';
  import ModernCategoriesNav from './ModernCategoriesNav.svelte';
  import BrandHero from '../organisms/BrandHero.svelte';
  import Footer from '../organisms/Footer.svelte';
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  import './modern-template.css';
  
  const { bgColor = 'bg-white', className = '', children } = $props();
  
  const restaurantData = $derived(restaurantDataStore.value);
  const menu = $derived(restaurantData?.menu || []);
  
  // Extraer categorías del menú para la navegación
  const categories = $derived(() => {
    const cats = new Map();
    menu.forEach(section => {
      if (section?.title && section?.items?.length > 0) {
        const firstItem = section.items[0];
        cats.set(section.title, {
          id: section.title.toLowerCase().replace(/\s+/g, '-'),
          title: section.title,
          image: firstItem?.photoUrl || ''
        });
      }
    });
    return Array.from(cats.values());
  });
  
  let selectedCategory = $state('all');
  
  function handleCategoryChange(category) {
    selectedCategory = category;
  }
  
  // Agregar clase al body para que los componentes puedan detectar el template moderno
  onMount(() => {
    if (typeof document !== 'undefined') {
      document.body.classList.add('modern-template');
      return () => {
        document.body.classList.remove('modern-template');
      };
    }
  });
</script>

<div class={`min-h-screen modern-template ${bgColor} ${className}`}>
  <!-- Header sticky -->
  <ModernHeader />
  
  <!-- Hero Section con BrandHero y overlay oscuro -->
  <section class="hero-section">
    <div class="hero-image-wrapper">
      <BrandHero />
    </div>
    <div class="hero-overlay"></div>
  </section>
  
  <!-- Navegación de categorías sticky -->
  <ModernCategoriesNav categories={categories()} onCategoryChange={handleCategoryChange} />
   
  <!-- Contenedor para el contenido adicional con diseño moderno -->
  <div class="content-container pt-6 sm:pt-8 pb-8 sm:pb-12 lg:pb-20">
    <!-- Render children content -->
    {@render children()}
  </div>
  
  <!-- Footer -->
  <Footer />
</div>
