<script lang="ts">
  import { t as tStore } from '../useLanguage'
  
  interface Message {
    id: string
    role: 'user' | 'assistant'
    content: string
    timestamp: Date
    showExploreButton?: boolean
    imageUrl?: string
    isPreview?: boolean
    pendingVersionId?: string
  }

  let { message, onExploreOptions }: { 
    message: Message, 
    onExploreOptions?: () => void
  } = $props()
</script>

<div class="flex items-start gap-3 {message.role === 'user' ? 'flex-row-reverse' : ''}">
  {#if message.role === 'assistant'}
    <div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center flex-shrink-0">
      <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
    </div>
  {:else}
    <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center flex-shrink-0">
      <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
      </svg>
    </div>
  {/if}

  <div class="flex-1 {message.role === 'user' ? 'flex justify-end' : ''}">
    <div class="
      {message.role === 'user' 
        ? message.isPreview
          ? 'bg-blue-100 border-2 border-blue-300 border-dashed text-gray-700 rounded-2xl rounded-tr-sm' 
          : 'bg-blue-600 text-white rounded-2xl rounded-tr-sm'
        : 'bg-gray-50 text-gray-900 rounded-2xl rounded-tl-sm'
      } 
      px-4 py-3 max-w-[85%] inline-block
    ">
      {#if message.imageUrl}
        <div class="mb-2">
          <img 
            src={message.imageUrl} 
            alt={message.isPreview ? "Foto pendiente" : "Foto enviada"} 
            class="max-w-full max-h-64 rounded-lg object-contain"
          />
        </div>
      {/if}
      <p class="text-sm leading-relaxed whitespace-pre-wrap">{message.content}</p>
      {#if message.showExploreButton && onExploreOptions}
        <button
          onclick={onExploreOptions}
          class="mt-3 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-colors"
        >
          {$tStore.chat.exploreOptions}
        </button>
      {/if}
    </div>
    <p class="text-xs text-gray-500 mt-1 {message.role === 'user' ? 'text-right' : 'text-left'}">
      {new Date(message.timestamp).toLocaleTimeString('es-ES', { 
        hour: '2-digit', 
        minute: '2-digit' 
      })}
    </p>
  </div>
</div>

