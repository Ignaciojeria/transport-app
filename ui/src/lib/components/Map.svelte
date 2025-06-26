<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';

	// Extender Window interface para routeLayers
	declare global {
		interface Window {
			routeLayers?: any[];
		}
	}

	export let lineString: number[][] = [];
	export let geoJson: any = null;
	export let customMarkers: any[] = [];
	export let multipleRoutes: { coordinates: number[][], color: string }[] = [];
	export let center: [number, number] = [40.4168, -3.7038];
	export let zoom: number = 6;
	export let height: string = '400px';
	export let showMarkers: boolean = true;
	export let lineColor: string = 'red';
	export let lineWeight: number = 3;
	export let lineOpacity: number = 0.7;

	// Paleta de colores para 10 rutas diferentes
	const routeColors = [
		'#FF6B6B', // Rojo coral
		'#4ECDC4', // Turquesa
		'#45B7D1', // Azul claro
		'#96CEB4', // Verde menta
		'#FFEAA7', // Amarillo claro
		'#DDA0DD', // Ciruela
		'#98D8C8', // Verde agua
		'#F7DC6F', // Amarillo dorado
		'#BB8FCE', // Lavanda
		'#85C1E9'  // Azul cielo
	];

	let mapContainer: HTMLDivElement;
	let map: any = null;
	let routeLayer: any = null;
	let geoJsonLayer: any = null;
	let markers: any[] = [];
	let L: any = null;

	// Función para obtener el color de una ruta por índice
	function getRouteColor(index: number): string {
		return routeColors[index % routeColors.length];
	}

	onMount(async () => {
		if (!browser) return;

		// Importar Leaflet dinámicamente solo en el navegador
		const leaflet = await import('leaflet');
		L = leaflet.default;
		
		// Importar CSS
		await import('leaflet/dist/leaflet.css');

		// Inicializar el mapa
		map = L.map(mapContainer).setView(center, zoom);

		// Agregar capa de tiles de OpenStreetMap
		L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
			attribution: '© OpenStreetMap contributors'
		}).addTo(map);

		// Dibujar la ruta inicial si existe
		if (lineString.length > 0) {
			drawRoute();
		}

		// Dibujar GeoJSON si existe
		if (geoJson) {
			drawGeoJson();
		}
	});

	onDestroy(() => {
		if (map) {
			map.remove();
		}
	});

	// Función para dibujar la ruta
	function drawRoute() {
		if (!L || !map) return;

		// Limpiar capas existentes
		clearRoute();

		// Dibujar múltiples rutas si están disponibles
		if (multipleRoutes.length > 0) {
			drawMultipleRoutes();
		} else if (lineString.length > 0) {
			drawSingleRoute();
		}

		// Agregar marcadores personalizados si están disponibles
		if (showMarkers && customMarkers.length > 0) {
			drawCustomMarkers();
		}
	}

	// Función para dibujar múltiples rutas
	function drawMultipleRoutes() {
		multipleRoutes.forEach((route, index) => {
			if (route.coordinates && route.coordinates.length > 0) {
				// Usar color específico de la ruta o color de la paleta
				const routeColor = route.color || getRouteColor(index);
				
				const routeLayer = L.polyline(route.coordinates, {
					color: routeColor,
					weight: lineWeight,
					opacity: lineOpacity
				}).addTo(map);
				
				// Guardar referencia para poder limpiar después
				if (!window.routeLayers) {
					window.routeLayers = [];
				}
				window.routeLayers.push(routeLayer);
			}
		});

		// Ajustar la vista para mostrar todas las rutas
		if (window.routeLayers && window.routeLayers.length > 0) {
			const group = new L.featureGroup(window.routeLayers);
			map.fitBounds(group.getBounds());
		}
	}

	// Función para dibujar una sola ruta
	function drawSingleRoute() {
		// Dibujar la línea de la ruta
		routeLayer = L.polyline(lineString, {
			color: lineColor,
			weight: lineWeight,
			opacity: lineOpacity
		}).addTo(map);

		// Ajustar la vista para mostrar toda la ruta
		map.fitBounds(routeLayer.getBounds());
	}

	// Función para dibujar marcadores personalizados
	function drawCustomMarkers() {
		customMarkers.forEach((step, index) => {
			if (step.location && step.location.length === 2) {
				const stepType = step.step_type;
				let markerHtml = '';
				
				// Determinar el color del marcador basado en la ruta
				let bgColor = lineColor; // color por defecto
				
				// Si el marcador tiene información de ruta, usar ese color
				if (step.route_index !== undefined) {
					bgColor = getRouteColor(step.route_index);
				} else if (step.vehicleColor) {
					bgColor = step.vehicleColor;
				}

				if (stepType === 'start') {
					markerHtml = '▶';
					bgColor = '#28a745'; // verde para start
				} else if (stepType === 'end') {
					markerHtml = '⏹️';
					bgColor = '#dc3545'; // rojo para end
				} else if (stepType === 'pickup') {
					markerHtml = '📦';
					// Mantener el color de la ruta para pickup
				} else if (stepType === 'delivery') {
					markerHtml = step.step_number ? step.step_number.toString() : '';
					// Mantener el color de la ruta para delivery
				} else if (step.step_number) {
					markerHtml = step.step_number.toString();
					// Mantener el color de la ruta para jobs
				}

				// No mostrar marcador si no hay nada que mostrar
				if (!markerHtml) return;

				const customIcon = L.divIcon({
					className: 'custom-numbered-marker',
					html: `<div style="
						background-color: ${bgColor};
						color: white;
						border: 2px solid white;
						border-radius: 50%;
						width: 30px;
						height: 30px;
						display: flex;
						align-items: center;
						justify-content: center;
						font-weight: bold;
						font-size: 14px;
						box-shadow: 0 2px 4px rgba(0,0,0,0.3);
					">${markerHtml}</div>`,
					iconSize: [30, 30],
					iconAnchor: [15, 15]
				});
				
				const vehicleInfo = step.vehicle ? `Vehículo ${step.vehicle}<br>` : '';
				const routeInfo = step.route_index !== undefined ? `Ruta ${step.route_index + 1}<br>` : '';
				const stepDescription = stepType === 'pickup' ? 'Recogida' : 
									  stepType === 'delivery' ? 'Entrega' : 
									  stepType === 'job' ? 'Entrega directa' : stepType;
				
				// Añadir información de órdenes si está disponible
				const orderInfo = step.order_refs && step.order_refs.length > 0 
					? `<br><strong>Órdenes:</strong> ${step.order_refs.join(', ')}` 
					: '';
				
				const marker = L.marker(step.location, { icon: customIcon })
					.addTo(map)
					.bindPopup(
						`${routeInfo}${vehicleInfo}Paso ${step.step_number || index + 1}: ${stepDescription}<br>Llegada: ${step.arrival || 'N/A'} seg${orderInfo}`
					);
				markers.push(marker);
			}
		});
	}

	// Nueva función para dibujar GeoJSON
	function drawGeoJson() {
		if (!L || !map || !geoJson) return;

		// Limpiar capas existentes
		clearGeoJson();

		// Crear capa GeoJSON con estilo personalizado
		geoJsonLayer = L.geoJSON(geoJson, {
			style: function(feature) {
				const props = feature.properties || {};
				const routeIndex = props.route_index || 0;
				const routeColor = getRouteColor(routeIndex);
				
				return {
					color: routeColor,
					weight: lineWeight,
					opacity: lineOpacity,
					fillOpacity: 0.1
				};
			},
			pointToLayer: function(feature, latlng) {
				if (showMarkers) {
					const props = feature.properties || {};
					const stepType = props.step_type;
					let markerHtml = '';
					
					// Determinar el color del marcador basado en la ruta
					const routeIndex = props.route_index || 0;
					let bgColor = getRouteColor(routeIndex);

					if (stepType === 'start') {
						markerHtml = '▶';
						bgColor = '#28a745'; // verde para start
					} else if (stepType === 'end') {
						markerHtml = '⏹️';
						bgColor = '#dc3545'; // rojo para end
					} else if (stepType === 'pickup') {
						markerHtml = '📦';
						// Mantener el color de la ruta para pickup
					} else if (stepType === 'delivery') {
						markerHtml = props.step_number ? props.step_number : '';
						// Mantener el color de la ruta para delivery
					} else if (props.step_number) {
						markerHtml = props.step_number;
						// Mantener el color de la ruta para jobs
					}

					// No mostrar marcador si no hay nada que mostrar
					if (!markerHtml) return;

					const customIcon = L.divIcon({
						className: 'custom-numbered-marker',
						html: `<div style="
							background-color: ${bgColor};
							color: white;
							border: 2px solid white;
							border-radius: 50%;
							width: 30px;
							height: 30px;
							display: flex;
							align-items: center;
							justify-content: center;
							font-weight: bold;
							font-size: 14px;
							box-shadow: 0 2px 4px rgba(0,0,0,0.3);
						">${markerHtml}</div>`,
						iconSize: [30, 30],
						iconAnchor: [15, 15]
					});
					
					const routeInfo = props.route_index !== undefined ? `Ruta ${props.route_index + 1}<br>` : '';
					const popupContent = `${routeInfo}${props.popup || props.name || `Punto`}`;
					
					return L.marker(latlng, { icon: customIcon }).bindPopup(popupContent);
				}
				
				const routeIndex = feature.properties?.route_index || 0;
				const routeColor = getRouteColor(routeIndex);
				
				return L.circleMarker(latlng, {
					radius: 6,
					fillColor: routeColor,
					color: routeColor,
					weight: 2,
					opacity: 1,
					fillOpacity: 0.8
				});
			},
			onEachFeature: function(feature, layer) {
				if (feature.properties && feature.properties.name) {
					layer.bindPopup(feature.properties.name);
				}
			}
		}).addTo(map);

		// Ajustar la vista para mostrar todo el GeoJSON
		if (geoJsonLayer.getBounds) {
			map.fitBounds(geoJsonLayer.getBounds());
		}
	}

	// Función para limpiar la ruta
	function clearRoute() {
		if (!map) return;

		if (routeLayer) {
			map.removeLayer(routeLayer);
			routeLayer = null;
		}

		// Limpiar múltiples rutas
		if (window.routeLayers) {
			window.routeLayers.forEach(layer => {
				map.removeLayer(layer);
			});
			window.routeLayers = [];
		}

		markers.forEach(marker => {
			map.removeLayer(marker);
		});
		markers = [];
	}

	// Nueva función para limpiar GeoJSON
	function clearGeoJson() {
		if (!map) return;

		if (geoJsonLayer) {
			map.removeLayer(geoJsonLayer);
			geoJsonLayer = null;
		}
	}

	// Observar cambios en lineString
	$: if (map && lineString && L && lineString.length > 0) {
		drawRoute();
	}

	// Observar cambios en customMarkers
	$: if (map && customMarkers && L && customMarkers.length > 0) {
		drawRoute();
	}

	// Observar cambios en multipleRoutes
	$: if (map && multipleRoutes && L && multipleRoutes.length > 0) {
		drawRoute();
	}

	// Observar cambios en geoJson
	$: if (map && geoJson && L) {
		drawGeoJson();
	}

	// Observar cambios en las propiedades de estilo
	$: if (routeLayer && lineString.length > 0 && L) {
		routeLayer.setStyle({
			color: lineColor,
			weight: lineWeight,
			opacity: lineOpacity
		});
	}

	// Observar cambios en las propiedades de estilo para GeoJSON
	$: if (geoJsonLayer && geoJson && L) {
		geoJsonLayer.setStyle({
			color: lineColor,
			weight: lineWeight,
			opacity: lineOpacity,
			fillOpacity: 0.1
		});
	}
</script>

<div 
	bind:this={mapContainer} 
	class="w-full"
	style="height: {height}; z-index: 1;"
></div>

<style>
	:global(.leaflet-container) {
		z-index: 1;
	}
</style> 