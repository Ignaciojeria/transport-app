<script>
  import { onMount, onDestroy } from 'svelte';
  
  let { lat = 0, lon = 0, height = '400px', zoom = 18 } = $props();
  
  let mapContainer;
  let map = null;
  let marker = null;
  let L = null;
  const LOCATIONIQ_API_KEY = 'pk.6fb271168ebb42320bc9248737e83834';
  
  onMount(async () => {
    if (typeof window === 'undefined' || !mapContainer) return;
    
    // Importar Leaflet dinámicamente
    L = (await import('leaflet')).default;
    await import('leaflet/dist/leaflet.css');
    
    // Importar iconos por defecto de Leaflet
    if (L.Icon.Default.prototype && '_getIconUrl' in L.Icon.Default.prototype) {
      delete L.Icon.Default.prototype._getIconUrl;
    }
    L.Icon.Default.mergeOptions({
      iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
      iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
      shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
    });
    
    // Inicializar el mapa
    map = L.map(mapContainer).setView([lat, lon], zoom);
    
    // Usar LocationIQ como proveedor de tiles
    L.tileLayer(`https://{s}-tiles.locationiq.com/v2/obk/r/{z}/{x}/{y}.png?key=${LOCATIONIQ_API_KEY}`, {
      attribution: '© LocationIQ',
      maxZoom: 19,
      subdomains: ['a', 'b', 'c']
    }).addTo(map);
    
    // Agregar marcador
    marker = L.marker([lat, lon], {
      icon: L.icon({
        iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-red.png',
        shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
        iconSize: [25, 41],
        iconAnchor: [12, 41],
        popupAnchor: [1, -34],
        shadowSize: [41, 41]
      })
    }).addTo(map);
    
    // Ajustar la vista al marcador
    map.fitBounds([[lat, lon]], { padding: [20, 20] });
  });
  
  onDestroy(() => {
    if (map) {
      map.remove();
      map = null;
      marker = null;
    }
  });
  
  // Actualizar mapa cuando cambien las coordenadas
  $effect(() => {
    if (map && L && lat && lon) {
      map.setView([lat, lon], zoom);
      if (marker) {
        marker.setLatLng([lat, lon]);
      } else {
        marker = L.marker([lat, lon], {
          icon: L.icon({
            iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-red.png',
            shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
            iconSize: [25, 41],
            iconAnchor: [12, 41],
            popupAnchor: [1, -34],
            shadowSize: [41, 41]
          })
        }).addTo(map);
      }
      map.fitBounds([[lat, lon]], { padding: [20, 20] });
    }
  });
</script>

<div 
  bind:this={mapContainer} 
  class="w-full rounded-lg overflow-hidden border-2 border-green-500 shadow-lg"
  style="height: {height}; z-index: 1;"
></div>

<style>
  :global(.leaflet-container) {
    z-index: 1;
    border-radius: 0.5rem;
  }
</style>
