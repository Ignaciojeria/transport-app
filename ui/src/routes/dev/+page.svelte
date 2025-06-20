<script lang="ts">
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import Map from '$lib/components/Map.svelte';

	let geoJsonData: any = null;

	let mapConfig = {
		center: [-33.52245, -70.575] as [number, number],
		zoom: 12, // Zoom out a bit to see the whole route
		lineColor: '#3388ff',
		lineWeight: 5,
		lineOpacity: 0.7,
		showMarkers: true
	};

	async function loadGeoJson() {
		if (!browser) return;
		try {
			const response = await fetch('/dev/geojson.json');
			if (response.ok) {
				geoJsonData = await response.json();
			} else {
				console.error('Error cargando GeoJSON:', response.statusText);
			}
		} catch (error) {
			console.error('Error cargando GeoJSON:', error);
		}
	}

	onMount(() => {
		loadGeoJson();
	});
</script>

<svelte:head>
	<title>Ruta Optimizada - La Florida, Santiago</title>
</svelte:head>

<div class="w-full h-screen">
	{#if browser}
		<Map 
			geoJson={geoJsonData}
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