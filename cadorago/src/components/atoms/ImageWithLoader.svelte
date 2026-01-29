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
  let lastSrc = $state(''); // Para detectar cambios en la URL
  let lastImageSrc = $state(''); // Para detectar cambios en imageSrc (por retryCount)
  const maxRetries = 30; // Máximo 30 reintentos (~90 segundos)
  const retryInterval = 3000; // 3 segundos entre reintentos
  const initialPollDelay = 2000; // Esperar 2 segundos antes del primer polling (dar tiempo a que se suba)

  // Normalizar URL: corregir "https.storage" → "https://storage"
  const normalizedSrc = $derived.by(() => {
    if (!src) return '';
    // Corregir URLs mal formateadas
    return src
      .replace(/https\.storage\.googleapis\.com/g, 'https://storage.googleapis.com')
      .replace(/http\.storage\.googleapis\.com/g, 'http://storage.googleapis.com');
  });

  // Determinar si es una imagen de GCS que podría estar siendo generada
  const isGCSImage = $derived(normalizedSrc.includes('storage.googleapis.com'));

  // URL de la imagen con timestamp para evitar cacheo del navegador
  // Para imágenes de GCS, siempre agregamos timestamp inicial + retryCount
  const imageSrc = $derived.by(() => {
    if (!normalizedSrc) return '';
    if (!isGCSImage) return normalizedSrc;
    
    // Para imágenes de GCS, agregar timestamp inicial + retryCount para forzar recarga
    const separator = normalizedSrc.includes('?') ? '&' : '?';
    return `${normalizedSrc}${separator}_t=${initialTimestamp}${retryCount > 0 ? `_r${retryCount}` : ''}`;
  });

  // Resetear estado cuando cambia la URL original (usar normalizedSrc)
  $effect(() => {
    if (!normalizedSrc) {
      imageLoaded = false;
      imageError = false;
      retryCount = 0;
      showLoader = false;
      lastSrc = '';
      return;
    }

    // Si la URL cambió, resetear todo y generar nuevo timestamp
    if (normalizedSrc !== lastSrc) {
      imageLoaded = false;
      imageError = false;
      retryCount = 0;
      showLoader = false;
      initialTimestamp = Date.now(); // Nuevo timestamp para forzar recarga sin cacheo
      lastSrc = normalizedSrc;
    }
  });

  function handleImageLoad() {
    // SIEMPRE ocultar el loader cuando la imagen carga, incluso si hubo errores previos
    imageLoaded = true;
    imageError = false;
    showLoader = false;
    retryCount = 0;
    // Log para debugging (puedes removerlo después)
    console.log('[ImageWithLoader] Image loaded successfully:', normalizedSrc);
  }

  function handleImageError() {
    // Solo marcar error si la imagen aún no ha cargado
    if (!imageLoaded) {
      imageError = true;
      imageLoaded = false;
      
      // Solo mostrar loader si es una imagen de GCS (podría estar siendo generada)
      if (isGCSImage) {
        // No establecer showLoader aquí - lo hará el $effect si es necesario
        // Log para debugging (puedes removerlo después)
        console.log('[ImageWithLoader] Image error, will retry:', normalizedSrc, 'retryCount:', retryCount);
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

  // Cuando cambia imageSrc (por retryCount), resetear imageError para permitir nuevo intento
  // El elemento <img> se recrea por el {#key}, pero necesitamos resetear el estado de error
  $effect(() => {
    if (imageSrc && imageSrc !== lastImageSrc && !imageLoaded) {
      // Cuando cambia imageSrc (nuevo retry), resetear imageError para permitir nuevo intento
      // Esto asegura que si la imagen ahora está disponible, el navegador pueda cargarla
      imageError = false;
      lastImageSrc = imageSrc;
    } else if (!imageSrc) {
      lastImageSrc = '';
    }
  });

  // Polling para imágenes de GCS: solo cuando hubo error (reintentos).
  // Si la imagen existe y está cargando, NUNCA mostramos el loader.
  $effect(() => {
    // SIEMPRE salir si la imagen cargó
    if (imageLoaded) {
      return;
    }

    if (!normalizedSrc || !isGCSImage) {
      return;
    }

    let interval = null;
    let timeout = null;
    let errorTimeout = null; // Timeout para ocultar loader después de errores sin éxito

    const startPolling = (onlyAfterError) => {
      // Verificar nuevamente antes de iniciar
      if (imageLoaded || retryCount >= maxRetries) return;
      
      // Solo mostrar "Generando imagen..." cuando ya hubo un error (404, etc.)
      if (onlyAfterError && !imageLoaded) {
        showLoader = true;
        
        // Si después de varios reintentos la imagen no carga, ocultar el loader
        // (puede ser un problema de red/CORS, no que la imagen no exista)
        errorTimeout = setTimeout(() => {
          if (!imageLoaded && retryCount >= 5) { // Después de 5 reintentos (~15 segundos)
            showLoader = false;
          }
        }, 15000); // 15 segundos
      }

      interval = setInterval(() => {
        // Verificar SIEMPRE si la imagen cargó (puede haber cargado mientras esperábamos)
        if (imageLoaded) {
          clearInterval(interval);
          interval = null;
          if (errorTimeout) clearTimeout(errorTimeout);
          showLoader = false;
          imageError = false;
          return;
        }
        
        // Incrementar retryCount para cambiar la URL y forzar nueva carga
        if (retryCount < maxRetries) {
          retryCount++;
          // Cuando cambia retryCount, imageSrc cambia, lo que fuerza al navegador a intentar cargar de nuevo
          // El evento onload debería dispararse si la imagen ahora está disponible
        } else {
          clearInterval(interval);
          interval = null;
          if (errorTimeout) clearTimeout(errorTimeout);
          showLoader = false;
        }
      }, retryInterval);
    };

    // Si hubo error inicialmente, mostrar loader y hacer polling
    // Si no hubo error pero la imagen no carga, hacer polling en segundo plano después del delay
    // También continuar polling si retryCount > 0 (ya estamos en modo retry)
    if ((imageError || retryCount > 0) && !imageLoaded) {
      // Hubo error o ya estamos en modo retry → mostrar loader solo si hubo error inicial
      startPolling(imageError);
    } else if (!imageError && retryCount === 0 && !imageLoaded) {
      // Aún no hubo error y no hemos empezado a hacer retry: hacer polling en segundo plano sin mostrar loader
      timeout = setTimeout(() => {
        if (!imageLoaded && !imageError && retryCount < maxRetries) {
          startPolling(false); // false = no mostrar loader
        }
      }, initialPollDelay);
    }

    return () => {
      if (timeout) clearTimeout(timeout);
      if (interval) clearInterval(interval);
      if (errorTimeout) clearTimeout(errorTimeout);
    };
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

    {#key imageSrc}
      <img 
        src={imageSrc} 
        alt={alt}
        class={`${className} ${imageLoaded ? 'opacity-100' : showLoader ? 'opacity-0' : 'opacity-100'} transition-opacity duration-300`}
        style={className.includes('h-auto') ? 'width: 100%; height: auto; display: block; object-fit: contain; object-position: center top;' : ''}
        onload={handleImageLoad}
        onerror={handleImageError}
        loading={loadingAttr}
      />
    {/key}

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
