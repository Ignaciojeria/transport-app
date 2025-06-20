<script lang="ts">
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import Map from '$lib/components/Map.svelte';

	let routeData: any = null;
	let lineString: number[][] = [];
	let markers: any[] = [];

	let mapConfig = {
		center: [-33.52245, -70.575] as [number, number],
		zoom: 12, // Zoom out a bit to see the whole route
		lineColor: '#3388ff',
		lineWeight: 5,
		lineOpacity: 0.7,
		showMarkers: true
	};

	async function loadRouteData() {
		if (!browser) return;
		try {
			const response = await fetch('/dev/polyline.json');
			if (response.ok) {
				routeData = await response.json();
				
				// Procesar la primera ruta (asumiendo que solo hay una)
				if (routeData && routeData.length > 0) {
					const route = routeData[0];
					
					// Usar las coordenadas decodificadas del polyline
					lineString = route.route || [];
					
					// Procesar los steps para marcadores
					markers = route.steps || [];
					
					console.log('Ruta cargada:', {
						vehicle: route.vehicle,
						cost: route.cost,
						duration: route.duration,
						routePoints: lineString.length,
						steps: markers.length
					});
				}
			} else {
				console.error('Error cargando datos de ruta:', response.statusText);
			}
		} catch (error) {
			console.error('Error cargando datos de ruta:', error);
		}
	}

	onMount(() => {
		loadRouteData();
	});
</script>

<svelte:head>
	<title>Ruta Optimizada - La Florida, Santiago</title>
</svelte:head>

<div class="w-full h-screen">
	{#if browser}
		<Map 
			lineString={lineString}
			customMarkers={markers}
			center={mapConfig.center}
			zoom={mapConfig.zoom}
			height="100vh"
			lineColor={mapConfig.lineColor}
			lineWeight={mapConfig.lineWeight}
			lineOpacity={mapConfig.lineOpacity}
			showMarkers={mapConfig.showMarkers}
		/>
	{:else}
		<div class="w-full h-screen bg-gray-200 flex items-center justify-center">
			<p class="text-gray-500 text-xl">Cargando mapa...</p>
		</div>
	{/if}
</div> 