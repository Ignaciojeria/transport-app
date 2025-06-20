export interface MapConfig {
	lineString: number[][];
	center: [number, number];
	zoom: number;
	height: string;
	showMarkers: boolean;
	lineColor: string;
	lineWeight: number;
	lineOpacity: number;
}

export interface RoutePoint {
	lat: number;
	lng: number;
	name?: string;
	description?: string;
}

export interface Route {
	id: string;
	name: string;
	points: RoutePoint[];
	color: string;
	weight: number;
	opacity: number;
}

export interface MapEvent {
	type: 'click' | 'route-click' | 'marker-click';
	coordinates?: [number, number];
	routeId?: string;
	pointIndex?: number;
} 