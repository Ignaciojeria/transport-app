<script lang="ts">
  let { onSend, disabled = false }: { onSend: (message: string) => void, disabled?: boolean } = $props()

  let inputValue = $state('')
  let textareaRef: HTMLTextAreaElement

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

<div class="relative flex items-end gap-2 bg-white rounded-2xl border border-gray-300 focus-within:border-blue-500 focus-within:ring-2 focus-within:ring-blue-200 transition-all">
  <button 
    class="p-2 text-gray-500 hover:text-gray-700 transition-colors flex-shrink-0"
    type="button"
    aria-label="Adjuntar archivo"
  >
    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
    </svg>
  </button>

  <textarea
    bind:this={textareaRef}
    bind:value={inputValue}
    onkeydown={handleKeyDown}
    oninput={adjustTextareaHeight}
    placeholder="Escribe tu menú y precios aquí..."
    disabled={disabled}
    class="flex-1 resize-none border-0 focus:ring-0 focus:outline-none py-3 px-2 text-gray-900 placeholder-gray-400 bg-transparent max-h-[200px] overflow-y-auto"
    rows="1"
  ></textarea>

  <div class="flex items-center gap-1 pr-2 pb-2">
    <button
      type="button"
      class="p-1.5 text-gray-500 hover:text-gray-700 transition-colors"
      title="Micrófono"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
      </svg>
    </button>

    <button
      type="button"
      onclick={handleSubmit}
      disabled={disabled || !inputValue.trim()}
      class="p-1.5 rounded-lg transition-colors {disabled || !inputValue.trim() 
        ? 'text-gray-300 cursor-not-allowed' 
        : 'text-blue-600 hover:bg-blue-50'}"
      title="Enviar"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
      </svg>
    </button>
  </div>
</div>

