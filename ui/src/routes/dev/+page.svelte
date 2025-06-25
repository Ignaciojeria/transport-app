<script lang="ts">
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import Map from '$lib/components/Map.svelte';

	let routeData: any = null;
	let allLineStrings: number[][][] = [];
	let allMarkers: any[] = [];
	let multipleRoutes: { coordinates: number[][], color: string }[] = [];
	let vehicleColors: string[] = [
		'#3388ff', // Azul
		'#ff4444', // Rojo
		'#44ff44', // Verde
		'#ffaa00', // Naranja
		'#aa44ff', // Púrpura
		'#ff44aa', // Rosa
		'#44aaff', // Azul claro
		'#aaff44', // Verde claro
		'#ffaa44', // Naranja claro
		'#aa44aa'  // Púrpura oscuro
	];

	let mapConfig = {
		center: [-33.52245, -70.575] as [number, number],
		zoom: 12,
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
				
				// Procesar todas las rutas
				if (routeData && routeData.length > 0) {
					routeData.forEach((route: any, routeIndex: number) => {
						// Usar las coordenadas decodificadas del polyline
						if (route.route && route.route.length > 0) {
							allLineStrings.push(route.route);
							
							// Agregar a múltiples rutas con color
							const vehicleColor = vehicleColors[routeIndex % vehicleColors.length];
							multipleRoutes.push({
								coordinates: route.route,
								color: vehicleColor
							});
						}
						
						// Procesar los steps para marcadores con información del vehículo
						if (route.steps && route.steps.length > 0) {
							const vehicleColor = vehicleColors[routeIndex % vehicleColors.length];
							route.steps.forEach((step: any) => {
								allMarkers.push({
									...step,
									vehicle: route.vehicle,
									vehiclePlate: route.vehicle_plate,
									vehicleColor: vehicleColor,
									routeIndex: routeIndex
								});
							});
						}
					});
					
					console.log('Rutas cargadas:', {
						totalRoutes: routeData.length,
						totalLineStrings: allLineStrings.length,
						totalMarkers: allMarkers.length,
						vehicles: routeData.map((r: any) => ({ 
							id: r.vehicle, 
							plate: r.vehicle_plate 
						}))
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
	<title>Rutas Optimizadas - Múltiples Vehículos</title>
</svelte:head>

<div class="w-full h-screen">
	{#if browser}
		<Map 
			multipleRoutes={multipleRoutes}
			customMarkers={allMarkers}
			center={mapConfig.center}
			zoom={mapConfig.zoom}
			height="100vh"
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