/**
 * Store para la plantilla seleccionada en modo previsualización (preview-first).
 * Solo tiene efecto cuando la URL tiene ?station=true (preview desde consola).
 * Cambio al instante sin guardar; "Usar este diseño" confirma (postMessage al padre).
 */

class PreviewTemplateStore {
  constructor() {
    /** @type {'hero' | 'modern' | null} - override en preview; null = usar el del menú */
    this.override = $state(null);
  }

  setOverride(value) {
    this.override = value === 'hero' || value === 'modern' ? value : null;
  }

  clearOverride() {
    this.override = null;
  }
}

export const previewTemplateStore = new PreviewTemplateStore();
