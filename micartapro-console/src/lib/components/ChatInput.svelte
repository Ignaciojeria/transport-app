<script lang="ts">
  import { t as tStore } from '../useLanguage'
  
  let { onSend, disabled = false, onFocus, onBlur, initialValue = '', onTakePhoto, onUploadPhoto, onPastePhoto, hasPendingPhoto = false }: { 
    onSend: (message: string) => void, 
    disabled?: boolean, 
    onFocus?: () => void,
    onBlur?: () => void,
    initialValue?: string,
    onTakePhoto?: () => void,
    onUploadPhoto?: () => void,
    onPastePhoto?: (file: File) => void,
    hasPendingPhoto?: boolean
  } = $props()

  let inputValue = $state(initialValue)
  let textareaRef: HTMLTextAreaElement
  let showPhotoMenu = $state(false)
  let photoButtonRef: HTMLButtonElement | null = $state(null)
  
  // Exponer función para establecer el valor desde el componente padre
  export function setValue(value: string) {
    inputValue = value
    if (textareaRef) {
      adjustTextareaHeight()
      textareaRef.focus()
    }
  }
  
  // Exponer función para hacer focus en el input
  export function focus() {
    if (textareaRef) {
      textareaRef.focus()
    }
  }
  
  // Exponer referencia al textarea para que el padre pueda acceder al valor
  export { textareaRef }
  
  
  // Actualizar cuando cambie initialValue
  $effect(() => {
    if (initialValue) {
      inputValue = initialValue
    }
  })
  
  function handleFocus() {
    if (onFocus) {
      onFocus()
    }
  }

  function handleBlur(e: FocusEvent) {
    // Si hay texto, no llamar onBlur (el usuario está escribiendo)
    if (inputValue.trim().length > 0) {
      return
    }
    
    // En móvil, el blur puede ser temporal cuando aparece el teclado
    // Esperar más tiempo para verificar si realmente perdió el focus
    const isMobile = window.innerWidth <= 768
    const delay = isMobile ? 600 : 200
    
    setTimeout(() => {
      // Verificar si el textarea realmente perdió el focus
      const hasFocus = document.activeElement === textareaRef
      
      // Solo llamar onBlur si realmente perdió el focus y no hay texto
      // En móvil, también verificar que no sea un blur temporal del teclado
      if (!hasFocus && !inputValue.trim()) {
        // En móvil, verificar una vez más después de un pequeño delay adicional
        if (isMobile) {
          setTimeout(() => {
            if (document.activeElement !== textareaRef && !inputValue.trim()) {
              if (onBlur) {
                onBlur()
              }
            }
          }, 200)
        } else {
          if (onBlur) {
            onBlur()
          }
        }
      }
    }, delay)
  }

  function handleSubmit() {
    // Solo permitir enviar si hay texto escrito
    if (inputValue.trim() && !disabled) {
      onSend(inputValue)
      inputValue = ''
      adjustTextareaHeight()
    }
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSubmit()
    }
  }

  function adjustTextareaHeight() {
    if (textareaRef) {
      textareaRef.style.height = 'auto'
      textareaRef.style.height = `${Math.min(textareaRef.scrollHeight, 200)}px`
    }
  }

  $effect(() => {
    if (textareaRef) {
      adjustTextareaHeight()
    }
  })

  // Manejar pegado de imágenes
  function handlePaste(event: ClipboardEvent) {
    const items = event.clipboardData?.items
    if (!items) return

    for (let i = 0; i < items.length; i++) {
      const item = items[i]
      
      // Verificar si es una imagen
      if (item.type.startsWith('image/')) {
        event.preventDefault()
        
        const file = item.getAsFile()
        if (file && onPastePhoto) {
          // Llamar directamente a la función callback con el archivo
          onPastePhoto(file)
        } else if (file) {
          // Fallback: usar el input file si no hay callback directo
          setTimeout(() => {
            const input = document.getElementById('photo-file-input') as HTMLInputElement
            if (input) {
              const dataTransfer = new DataTransfer()
              dataTransfer.items.add(file)
              
              Object.defineProperty(input, 'files', {
                value: dataTransfer.files,
                writable: false,
              })
              
              const changeEvent = new Event('change', { bubbles: true })
              input.dispatchEvent(changeEvent)
            }
          }, 50)
        }
        break
      }
    }
  }

  function handlePhotoClick() {
    if (onTakePhoto && onUploadPhoto) {
      showPhotoMenu = !showPhotoMenu
    } else if (onTakePhoto) {
      onTakePhoto()
    } else if (onUploadPhoto) {
      onUploadPhoto()
    }
  }

  function handleTakePhoto() {
    showPhotoMenu = false
    if (onTakePhoto) {
      onTakePhoto()
    }
  }

  function handleUploadPhoto() {
    showPhotoMenu = false
    if (onUploadPhoto) {
      onUploadPhoto()
    }
  }

  // Cerrar menú al hacer clic fuera
  $effect(() => {
    if (!showPhotoMenu) return

    function handleClickOutside(event: MouseEvent) {
      if (photoButtonRef && !photoButtonRef.contains(event.target as Node)) {
        const menu = document.getElementById('photo-menu')
        if (menu && !menu.contains(event.target as Node)) {
          showPhotoMenu = false
        }
      }
    }

    document.addEventListener('click', handleClickOutside)
    return () => {
      document.removeEventListener('click', handleClickOutside)
    }
  })
</script>

<div class="relative">
  <!-- Efecto de brillo con gradiente (estilo AnimatedChat) -->
  <div class="absolute inset-0 bg-gradient-to-r from-pink-500 via-purple-500 to-blue-500 rounded-2xl blur-sm opacity-30"></div>
  <div class="relative bg-white rounded-2xl border-2 border-transparent p-1">
    <div class="flex items-center gap-2 px-4 py-3 bg-white rounded-xl">
      <!-- Botón de foto con menú -->
      {#if onTakePhoto || onUploadPhoto}
        <div class="relative flex-shrink-0">
          <button
            bind:this={photoButtonRef}
            type="button"
            onclick={handlePhotoClick}
            disabled={disabled}
            class="p-2 rounded-lg transition-colors text-gray-600 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
            title="Tomar o subir foto"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>

          <!-- Menú desplegable -->
          {#if showPhotoMenu && (onTakePhoto && onUploadPhoto)}
            <div 
              id="photo-menu"
              class="absolute bottom-full left-0 mb-2 bg-white rounded-lg shadow-lg border border-gray-200 py-2 min-w-[180px] z-50"
            >
              <button
                type="button"
                onclick={handleTakePhoto}
                disabled={disabled}
                class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 flex items-center gap-3 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                <span>Tomar Foto</span>
              </button>
              <button
                type="button"
                onclick={handleUploadPhoto}
                disabled={disabled}
                class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 flex items-center gap-3 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                </svg>
                <span>Subir Archivo</span>
              </button>
            </div>
          {/if}
        </div>
      {/if}

      <!-- Textarea -->
      <textarea
        bind:this={textareaRef}
        bind:value={inputValue}
        onkeydown={handleKeyDown}
        oninput={adjustTextareaHeight}
        onfocus={handleFocus}
        onblur={handleBlur}
        onpaste={handlePaste}
        placeholder={$tStore.chat.placeholder}
        disabled={disabled}
        class="flex-1 resize-none border-0 focus:ring-0 focus:outline-none text-gray-900 placeholder-gray-400 bg-transparent max-h-[200px] overflow-y-auto text-base"
        rows="1"
      ></textarea>

      <!-- Botón de enviar (triángulo) - solo visible cuando hay texto -->
      {#if inputValue.trim()}
        <button
          type="button"
          onclick={handleSubmit}
          disabled={disabled}
          class="p-2 rounded-lg transition-colors flex-shrink-0 bg-gray-800 hover:bg-gray-700 text-white"
          title="Enviar"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z" />
          </svg>
        </button>
      {/if}
    </div>
  </div>
</div>

