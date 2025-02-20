package request

import (
	"fmt"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

// arePointsClose retorna true si la distancia entre p1 y p2 es menor que la tolerancia dada (en metros).
func arePointsClose(p1, p2 orb.Point, tolerance float64) bool {
	return geo.Distance(p1, p2) < tolerance
}

// pointToString convierte un Point al formato "longitude,latitude" para LocationIQ.
func pointToString(p orb.Point) string {
	return fmt.Sprintf("%f,%f", p[0], p[1])
}

// arePointsEqual compara si dos puntos están lo suficientemente cerca (dentro de 1 metro de tolerancia).
func arePointsEqual(p1, p2 orb.Point) bool {
	return arePointsClose(p1, p2, 1)
}

// distanceMatches calcula la distancia geodésica entre p1 y p2 y la compara con la distancia esperada,
// dentro de una tolerancia definida (ambos en metros).
func distanceMatches(p1, p2 orb.Point, expectedDistance, tolerance float64) bool {
	actualDistance := geo.Distance(p1, p2)
	return math.Abs(actualDistance-expectedDistance) <= tolerance
}

// distanceBetweenPoints retorna la distancia geodésica en metros entre dos puntos.
func distanceBetweenPoints(p1, p2 orb.Point) float64 {
	return geo.Distance(p1, p2)
}
