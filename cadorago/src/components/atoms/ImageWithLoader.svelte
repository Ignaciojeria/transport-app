<script>
  const { 
    src = '',
    alt = '',
    className = '',
    fallbackIcon = true,
    /** 'eager' para imágenes above-the-fold (portada); 'lazy' por defecto */
    loading = 'lazy'
  } = $props();
  const loadingAttr = loading === 'eager' ? 'eager' : 'lazy';

  let imageLoaded = $state(false);
  let imageError = $state(false);
  let showLoader = $state(false);
  let lastSrc = $state(''); // Para detectar cambios en la URL

  // Normalizar URL: corregir "https.storage" → "https://storage"
  // Es idempotente: puede aplicarse múltiples veces sin causar efectos secundarios
  const normalizedSrc = $derived.by(() => {
    if (!src) return '';
    
    // Verificar si la URL ya está correctamente formateada
    if (src.startsWith('https://storage.googleapis.com') || src.startsWith('http://storage.googleapis.com')) {
      // Aún así, verificar si hay duplicaciones de "https"
      return src
        .replace(/httpshttps:\/\//g, 'https://')
        .replace(/httphttp:\/\//g, 'http://');
    }
    
    // Corregir URLs mal formateadas
    let normalized = src
      .replace(/https\.storage\.googleapis\.com/g, 'https://storage.googleapis.com')
      .replace(/http\.storage\.googleapis\.com/g, 'http://storage.googleapis.com');
    
    // Manejar casos donde se duplicó "https" (httpshttps://storage...)
    normalized = normalized
      .replace(/httpshttps:\/\//g, 'https://')
      .replace(/httphttp:\/\//g, 'http://');
    
    return normalized;
  });

  // Determinar si es una imagen de GCS que podría estar siendo generada
  const isGCSImage = $derived(normalizedSrc.includes('storage.googleapis.com'));

  // URL de la imagen: usar siempre la URL normalizada (sin polling complejo)
  // El navegador manejará el cacheo naturalmente
  const imageSrc = $derived(normalizedSrc);

  // Resetear estado cuando cambia la URL original
  $effect(() => {
    if (!normalizedSrc) {
      imageLoaded = false;
      imageError = false;
      showLoader = false;
      lastSrc = '';
      return;
    }

    // Si la URL cambió, resetear todo
    if (normalizedSrc !== lastSrc) {
      imageLoaded = false;
      imageError = false;
      showLoader = false;
      lastSrc = normalizedSrc;
    }
  });

  function handleImageLoad() {
    imageLoaded = true;
    imageError = false;
    showLoader = false;
  }

  function handleImageError() {
    if (!imageLoaded) {
      imageError = true;
      imageLoaded = false;
      
      // Solo mostrar loader si es una imagen de GCS (podría estar siendo generada)
      // El polling del backend se encargará de refrescar el iframe cuando esté lista
      if (isGCSImage) {
        showLoader = true;
      } else {
        showLoader = false;
      }
    }
  }

  // Asegurar que el loader se oculte cuando la imagen carga
  $effect(() => {
    if (imageLoaded) {
      showLoader = false;
      imageError = false;
    }
  });
</script>

<div class="relative {className}">
  {#if normalizedSrc}
    {#if showLoader && !imageLoaded}
      <!-- Loader solo cuando la imagen falla y podría estar siendo generada -->
      <div class="w-full flex items-center justify-center bg-gray-100" style={className.includes('h-auto') ? 'min-height: 200px; padding: 40px 0;' : 'height: 100%;'}>
        <div class="flex flex-col items-center gap-2">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-400"></div>
          <p class="text-xs text-gray-500">Generando imagen...</p>
        </div>
      </div>
    {/if}

    <img 
      src={imageSrc} 
      alt={alt}
      class={`${className} ${imageLoaded ? 'opacity-100' : showLoader ? 'opacity-0' : 'opacity-100'} transition-opacity duration-300`}
      style={className.includes('h-auto') ? 'width: 100%; height: auto; display: block; object-fit: contain; object-position: center top;' : ''}
      onload={handleImageLoad}
      onerror={handleImageError}
      loading={loadingAttr}
    />

    {#if imageError && !imageLoaded}
      <!-- Fallback cuando la imagen no existe después de todos los reintentos -->
      {#if fallbackIcon}
        <div class="w-full flex items-center justify-center bg-gray-100" style={className.includes('h-auto') ? 'min-height: 200px; padding: 40px 0;' : 'height: 100%;'}>
          <svg class="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
      {/if}
    {/if}
  {:else if fallbackIcon}
    <!-- Fallback cuando no hay URL -->
    <div class="w-full h-full flex items-center justify-center bg-gray-100">
      <svg class="w-12 h-12 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
    </div>
  {/if}
</div>
