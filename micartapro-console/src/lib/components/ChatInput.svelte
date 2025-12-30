<script lang="ts">
  import { t as tStore } from '../useLanguage'
  
  let { onSend, disabled = false, onFocus, onBlur, initialValue = '' }: { 
    onSend: (message: string) => void, 
    disabled?: boolean, 
    onFocus?: () => void,
    onBlur?: () => void,
    initialValue?: string
  } = $props()

  let inputValue = $state(initialValue)
  let textareaRef: HTMLTextAreaElement
  
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

  function handleBlur() {
    if (onBlur) {
      onBlur()
    }
  }

  function handleSubmit() {
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
</script>

<div class="relative">
  <!-- Efecto de brillo con gradiente (estilo AnimatedChat) -->
  <div class="absolute inset-0 bg-gradient-to-r from-pink-500 via-purple-500 to-blue-500 rounded-2xl blur-sm opacity-30"></div>
  <div class="relative bg-white rounded-2xl border-2 border-transparent p-1">
    <div class="flex items-center gap-2 px-4 py-3 bg-white rounded-xl">
      <!-- Textarea -->
      <textarea
        bind:this={textareaRef}
        bind:value={inputValue}
        onkeydown={handleKeyDown}
        oninput={adjustTextareaHeight}
        onfocus={handleFocus}
        onblur={handleBlur}
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

