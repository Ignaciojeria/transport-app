<script lang="ts">
  import { t as tStore } from '../useLanguage'
  
  let { onSend, disabled = false, onFocus }: { onSend: (message: string) => void, disabled?: boolean, onFocus?: () => void } = $props()

  let inputValue = $state('')
  let textareaRef: HTMLTextAreaElement
  
  function handleFocus() {
    if (onFocus) {
      onFocus()
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

<div class="relative flex items-end gap-2 bg-white rounded-2xl border border-gray-300 focus-within:border-blue-500 focus-within:ring-2 focus-within:ring-blue-200 transition-all">
  <textarea
    bind:this={textareaRef}
    bind:value={inputValue}
    onkeydown={handleKeyDown}
    oninput={adjustTextareaHeight}
    onfocus={handleFocus}
    placeholder={$tStore.chat.placeholder}
    disabled={disabled}
    class="flex-1 resize-none border-0 focus:ring-0 focus:outline-none py-3 px-4 text-gray-900 placeholder-gray-400 bg-transparent max-h-[200px] overflow-y-auto"
    rows="1"
  ></textarea>

  <div class="flex items-center pr-2 pb-2">
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

