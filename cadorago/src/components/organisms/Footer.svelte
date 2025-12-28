<script>
  import { getRestaurantData } from '../../services/restaurantData.js';
  
  const { className = '' } = $props();
  
  const restaurantData = getRestaurantData();
  const footerImage = restaurantData.footerImage;
</script>

<footer class={`bg-[#E8E4D9] pt-0 pb-8 sm:pb-12 lg:pb-16 ${className}`}>
  <div class="container mx-auto px-4 sm:px-6 lg:px-12 max-w-[1600px]">
    <div class="flex flex-col items-center -mt-8 sm:-mt-10 lg:-mt-12">
      <img 
        src={footerImage} 
        alt="Cadorago Logo" 
        class="h-64 w-auto"
      />
      <p class="text-2xl sm:text-3xl lg:text-4xl font-semibold italic text-gray-800 -mt-4 sm:-mt-6">
        Buen provecho 😊
      </p>
    </div>
  </div>
</footer>

