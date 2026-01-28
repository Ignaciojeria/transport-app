<script>
  import Price from '../atoms/Price.svelte';
  
  const {
    pricing = null,
    min = 0,
    max = 1000,
    step = 1,
    value = 0,
    unit = 'g',
    onConfirm = () => {},
    onCancel = () => {}
  } = $props();
  
  let quantity = $state(value || min);
  let inputValue = $state(String(value || min));
  
  // Detectar si necesitamos mostrar decimales (para KILOGRAM, LITER, etc.)
  const needsDecimals = $derived(
    pricing?.unit === 'KILOGRAM' || 
    pricing?.unit === 'LITER' || 
    pricing?.unit === 'METER' ||
    pricing?.unit === 'SQUARE_METER' ||
    step < 1
  );
  
  // Detectar si es venta por peso (ocultar slider)
  const isWeightSale = $derived(pricing?.mode === 'WEIGHT');
  
  // Función para formatear cantidad con decimales si es necesario
  function formatQuantity(qty) {
    if (needsDecimals) {
      return parseFloat(qty.toFixed(1));
    }
    return Math.round(qty);
  }
  
  // Función para formatear cantidad como string para mostrar
  function formatQuantityString(qty) {
    if (needsDecimals) {
      return qty.toFixed(1);
    }
    return String(Math.round(qty));
  }
  
  // Calcular precio basado en la cantidad
  const calculatedPrice = $derived.by(() => {
    if (!pricing) return 0;
    
    if (pricing.mode === 'UNIT') {
      return (pricing.pricePerUnit || 0) * quantity;
    }
    
    // Para WEIGHT, VOLUME, etc: precio = (cantidad / baseUnit) * pricePerUnit
    const baseUnit = pricing.baseUnit || 1;
    const pricePerUnit = pricing.pricePerUnit || 0;
    if (baseUnit === 0) return 0;
    
    return (quantity / baseUnit) * pricePerUnit;
  });
  
  // Sincronizar input con slider
  function handleSliderChange(event) {
    const newValue = parseFloat(event.target.value) || min;
    quantity = Math.max(min, Math.min(max, newValue));
    inputValue = formatQuantityString(quantity);
  }
  
  function handleInputChange(event) {
    const newValue = parseFloat(event.target.value) || min;
    const clampedValue = Math.max(min, Math.min(max, newValue));
    quantity = clampedValue;
    inputValue = formatQuantityString(clampedValue);
  }
  
  function handleInputBlur() {
    // Asegurar que el input tenga un valor válido
    if (!inputValue || isNaN(parseFloat(inputValue))) {
      quantity = min;
      inputValue = String(min);
    }
  }
  
  function handleConfirmClick() {
    onConfirm(quantity);
  }
  
  function handleCancelClick() {
    onCancel();
  }
  
  // Formatear la unidad para mostrar
  const unitLabel = $derived(
    unit && unit !== 'g' ? unit :
    !pricing || !pricing.unit ? (unit || 'g') :
    pricing.unit === 'GRAM' ? 'g' :
    pricing.unit === 'KILOGRAM' ? 'kg' :
    pricing.unit === 'MILLILITER' ? 'ml' :
    pricing.unit === 'LITER' ? 'L' :
    pricing.unit === 'METER' ? 'm' :
    pricing.unit === 'SQUARE_METER' ? 'm²' :
    (unit || 'g')
  );
</script>

<div class="space-y-4">
  <!-- Slider (oculto para venta por peso) -->
  {#if !isWeightSale}
    <div class="space-y-2">
      <div class="flex justify-between items-center">
        <label class="text-sm font-medium text-gray-700">
          Cantidad
        </label>
        <span class="text-sm text-gray-600">
          {formatQuantityString(quantity)} {unitLabel}
        </span>
      </div>
      <input
        type="range"
        min={min}
        max={max}
        step={step}
        value={quantity}
        oninput={handleSliderChange}
        class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer accent-gray-600"
      />
      <div class="flex justify-between text-xs text-gray-500">
        <span>{needsDecimals ? min.toFixed(1) : min} {unitLabel}</span>
        <span>{needsDecimals ? max.toFixed(1) : max} {unitLabel}</span>
      </div>
    </div>
  {/if}
  
  <!-- Input numérico -->
  <div class="space-y-2">
    <label class="text-sm font-medium text-gray-700">
      {isWeightSale ? 'Cantidad' : 'Cantidad exacta'}
    </label>
    <div class="flex items-center gap-2">
      <input
        type="number"
        min={min}
        max={max}
        step={needsDecimals ? "0.1" : step}
        value={inputValue}
        oninput={(e) => {
          inputValue = e.target.value;
          handleInputChange(e);
        }}
        onblur={handleInputBlur}
        class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-500 focus:border-transparent text-base"
        placeholder="0"
      />
      <span class="text-sm text-gray-600 font-medium min-w-[2rem]">
        {unitLabel}
      </span>
    </div>
  </div>
  
  <!-- Precio calculado -->
  <div class="pt-2 border-t border-gray-200">
    <div class="flex justify-between items-center">
      <span class="text-sm font-medium text-gray-700">
        Precio:
      </span>
      <Price price={calculatedPrice} className="text-lg font-bold" />
    </div>
  </div>
  
  <!-- Botones -->
  <div class="flex gap-3 pt-2">
    <button
      type="button"
      onclick={handleCancelClick}
      class="flex-1 px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg transition-colors font-medium text-sm"
    >
      Cancelar
    </button>
    <button
      type="button"
      onclick={handleConfirmClick}
      class="flex-1 px-4 py-2 bg-gray-800 hover:bg-gray-900 text-white rounded-lg transition-colors font-medium text-sm"
    >
      Agregar
    </button>
  </div>
</div>
