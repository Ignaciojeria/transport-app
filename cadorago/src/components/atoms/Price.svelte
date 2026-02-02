<script>
  import { restaurantDataStore } from '../../stores/restaurantDataStore.svelte.js';
  import { getEffectiveCurrency, formatPrice } from '../../lib/currency';

  const props = $props();
  const price = typeof props.price === 'number' ? props.price : 0;
  const className = props.className ?? '';
  const currency = props.currency ?? '';

  const effectiveCurrency = $derived(
    currency || getEffectiveCurrency(restaurantDataStore.value)
  );
  const formattedPrice = $derived(formatPrice(price, effectiveCurrency));
</script>

<span class={`font-bold text-gray-900 text-base sm:text-lg ${className}`}>
  {formattedPrice}
</span>
