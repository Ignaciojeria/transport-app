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
  let retryCount = $state(0);
  let showLoader = $state(false);
  let initialTimestamp = $state(Date.now()); // Timestamp inicial para evitar cacheo
  const maxRetries = 20; // Máximo 20 reintentos (~60 segundos)
  const retryInterval = 3000; // 3 segundos entre reintentos

  // Determinar si es una imagen de GCS que podría estar siendo generada
  const isGCSImage = $derived(src.includes('storage.googleapis.com'));

  // URL de la imagen con timestamp para evitar cacheo del navegador
  // Para imágenes de GCS, siempre agregamos timestamp inicial + retryCount si hay error
  const imageSrc = $derived.by(() => {
    if (!src) return '';
    if (!isGCSImage) return src;
    
    // Para imágenes de GCS, agregar timestamp inicial + retryCount para forzar recarga
    const separator = src.includes('?') ? '&' : '?';
    return `${src}${separator}_t=${initialTimestamp}${retryCount > 0 ? `_r${retryCount}` : ''}`;
  });

  // Resetear estado cuando cambia la URL original
  $effect(() => {
    if (!src) {
      imageLoaded = false;
      imageError = false;
      retryCount = 0;
      showLoader = false;
      return;
    }

    // Resetear estado para nueva imagen y generar nuevo timestamp para evitar cacheo
    imageLoaded = false;
    imageError = false;
    retryCount = 0;
    showLoader = false; // No mostrar loader inicialmente
    initialTimestamp = Date.now(); // Nuevo timestamp para forzar recarga sin cacheo
  });

  function handleImageLoad() {
    imageLoaded = true;
    imageError = false;
    showLoader = false;
    retryCount = 0;
  }

  function handleImageError() {
    imageError = true;
    imageLoaded = false;
    
    // Solo mostrar loader si es una imagen de GCS (podría estar siendo generada)
    if (isGCSImage) {
      showLoader = true;
      // Hacer polling para reintentar cargar la imagen
      if (retryCount < maxRetries) {
        setTimeout(() => {
          retryCount++;
        }, retryInterval);
      }
    } else {
      // Para otras URLs, no mostrar loader
      showLoader = false;
    }
  }

  // Polling automático solo si hay error y es imagen de GCS
  $effect(() => {
    if (!src || !isGCSImage || imageLoaded || !imageError || !showLoader) {
      return;
    }

    // Si la imagen falló y es de GCS, hacer polling para reintentar
    const interval = setInterval(() => {
      if (imageError && !imageLoaded && retryCount < maxRetries) {
        retryCount++;
      } else if (imageLoaded || retryCount >= maxRetries) {
        clearInterval(interval);
      }
    }, retryInterval);

    return () => clearInterval(interval);
  });
</script>

<div class="relative {className}">
  {#if src}
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

    {#if imageError && retryCount >= maxRetries && !imageLoaded}
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
