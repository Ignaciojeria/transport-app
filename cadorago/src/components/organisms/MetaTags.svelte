<script>
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  
  const restaurantData = $derived(restaurantDataStore.value);
  const coverImage = $derived(restaurantData?.coverImage || '');
  const businessName = $derived(restaurantData?.businessInfo?.businessName || 'Carta Digital');
  const description = $derived(
    restaurantData?.businessInfo?.description || 
    `Menú digital de ${businessName}. Explora nuestro catálogo y realiza tu pedido.`
  );
  
  // URL actual
  const currentUrl = $derived(() => {
    if (typeof window !== 'undefined') {
      return window.location.href;
    }
    return '';
  });
  
  // Función para actualizar meta tags existentes o crear nuevos
  function updateMetaTag(property, content) {
    if (typeof document === 'undefined') return;
    
    // Buscar por property o name
    let meta = document.querySelector(`meta[property="${property}"]`) || 
               document.querySelector(`meta[name="${property}"]`);
    
    if (meta) {
      meta.setAttribute('content', content);
    } else {
      // Crear nuevo meta tag
      meta = document.createElement('meta');
      if (property.startsWith('og:')) {
        meta.setAttribute('property', property);
      } else {
        meta.setAttribute('name', property);
      }
      meta.setAttribute('content', content);
      document.head.appendChild(meta);
    }
  }
  
  // Actualizar meta tags cuando cambien los datos
  $effect(() => {
    if (restaurantData && typeof document !== 'undefined') {
      // Actualizar título
      document.title = `${businessName} - Carta Digital`;
      
      // Actualizar meta description
      updateMetaTag('description', description);
      
      // Actualizar Open Graph
      updateMetaTag('og:url', currentUrl());
      updateMetaTag('og:title', businessName);
      updateMetaTag('og:description', description);
      
      if (coverImage) {
        updateMetaTag('og:image', coverImage);
        
        // Agregar meta tags adicionales para la imagen
        updateMetaTag('og:image:type', 'image/webp');
        updateMetaTag('og:image:width', '1200');
        updateMetaTag('og:image:height', '630');
      }
      
      // Actualizar Twitter Card
      updateMetaTag('twitter:card', 'summary_large_image');
      updateMetaTag('twitter:title', businessName);
      updateMetaTag('twitter:description', description);
      
      if (coverImage) {
        updateMetaTag('twitter:image', coverImage);
      }
    }
  });
</script>
