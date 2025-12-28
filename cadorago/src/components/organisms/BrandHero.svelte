<script>
  import { getRestaurantData } from '../../services/restaurantData.js';
  
  /** @type {string} - Additional custom classes */
  export let className = '';
  
  const restaurantData = getRestaurantData();
  const coverImage = restaurantData.coverImage;
</script>

<div class={`w-full overflow-hidden ${className}`}>
  <!-- Cover Image como foto de portada ocupando todo el ancho sin recortarse -->
  <img 
    src={coverImage} 
    alt="Portada de Cadorago" 
    class="w-full h-auto"
    style="width: 100%; height: auto; display: block; object-fit: contain; object-position: center top;"
  />
</div>
