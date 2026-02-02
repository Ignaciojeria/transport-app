<script>
  import { previewTemplateStore } from '../../stores/previewTemplateStore.svelte.js';
  import { t } from '../../lib/useLanguage';

  /** @type {'hero' | 'modern'} - template actualmente guardada en el menú */
  let savedTemplate = $props('hero');

  const currentOverride = $derived(previewTemplateStore.override);
  const effectiveTemplate = $derived(currentOverride || savedTemplate);
  const hasChanged = $derived(currentOverride !== null && currentOverride !== savedTemplate);

  function selectTemplate(template) {
    previewTemplateStore.setOverride(template);
  }

  function useThisDesign() {
    if (!hasChanged) return;
    const style = effectiveTemplate === 'hero' ? 'HERO' : 'MODERN';
    if (typeof window !== 'undefined' && window.opener) {
      window.opener.postMessage(
        { type: 'MICARTAPRO_SAVE_TEMPLATE', presentationStyle: style },
        '*'
      );
    }
    if (typeof window !== 'undefined' && window.parent !== window) {
      window.parent.postMessage(
        { type: 'MICARTAPRO_SAVE_TEMPLATE', presentationStyle: style },
        '*'
      );
    }
    previewTemplateStore.clearOverride();
  }
</script>

<div class="template-selector">
  <span class="template-selector__label">{$t.preview.templateLabel}</span>
  <div class="template-selector__tabs" role="tablist" aria-label="Estilo de carta">
    <button
      type="button"
      role="tab"
      aria-selected={effectiveTemplate === 'hero'}
      class="template-selector__tab"
      class:template-selector__tab_active={effectiveTemplate === 'hero'}
      onclick={() => selectTemplate('hero')}
    >
      {$t.preview.templateHero}
    </button>
    <button
      type="button"
      role="tab"
      aria-selected={effectiveTemplate === 'modern'}
      class="template-selector__tab"
      class:template-selector__tab_active={effectiveTemplate === 'modern'}
      onclick={() => selectTemplate('modern')}
    >
      {$t.preview.templateModern}
    </button>
  </div>
  {#if hasChanged}
    <button
      type="button"
      class="template-selector__confirm"
      onclick={useThisDesign}
    >
      {$t.preview.useThisDesign}
    </button>
  {/if}
</div>

<style>
  .template-selector {
    position: fixed;
    bottom: 1rem;
    left: 50%;
    transform: translateX(-50%);
    z-index: 50;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem 0.75rem;
    background: rgba(255, 255, 255, 0.98);
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
    border: 1px solid rgba(0, 0, 0, 0.08);
    flex-wrap: wrap;
    justify-content: center;
    max-width: calc(100vw - 2rem);
  }

  .template-selector__label {
    font-size: 0.8125rem;
    font-weight: 600;
    color: #374151;
  }

  .template-selector__tabs {
    display: flex;
    gap: 0.25rem;
    background: #f3f4f6;
    padding: 0.2rem;
    border-radius: 8px;
  }

  .template-selector__tab {
    padding: 0.4rem 0.75rem;
    font-size: 0.8125rem;
    font-weight: 500;
    color: #6b7280;
    background: transparent;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
  }

  .template-selector__tab:hover {
    color: #374151;
    background: rgba(255, 255, 255, 0.7);
  }

  .template-selector__tab_active {
    color: #111827;
    background: #fff;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
  }

  .template-selector__confirm {
    padding: 0.4rem 0.75rem;
    font-size: 0.8125rem;
    font-weight: 600;
    color: #fff;
    background: #059669;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.15s;
  }

  .template-selector__confirm:hover {
    background: #047857;
  }
</style>
