package request

import (
	"math"
	"testing"

	"github.com/paulmach/orb"
	"github.com/stretchr/testify/assert"
)

func TestArePointsClose(t *testing.T) {
	testCases := []struct {
		name     string
		p1       orb.Point
		p2       orb.Point
		tol      float64
		expected bool
	}{
		{
			name:     "Identical points",
			p1:       orb.Point{10.0, 20.0},
			p2:       orb.Point{10.0, 20.0},
			tol:      50,
			expected: true,
		},
		{
			name:     "Points within tolerance",
			p1:       orb.Point{10.0, 20.0},
			p2:       orb.Point{10.0001, 20.0001},
			tol:      50,
			expected: true,
		},
		{
			name:     "Points outside tolerance",
			p1:       orb.Point{10.0, 20.0},
			p2:       orb.Point{10.01, 20.01},
			tol:      50,
			expected: false,
		},
		{
			name:     "Points with different longitudes",
			p1:       orb.Point{10.0, 20.0},
			p2:       orb.Point{10.0002, 20.0},
			tol:      50,
			expected: true,
		},
		{
			name:     "Points with different latitudes",
			p1:       orb.Point{10.0, 20.0},
			p2:       orb.Point{10.0, 20.0002},
			tol:      50,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := arePointsClose(tc.p1, tc.p2, tc.tol)
			assert.Equal(t, tc.expected, result,
				"arePointsClose failed for points %v and %v with tol %f", tc.p1, tc.p2, tc.tol)
		})
	}
}

func TestPointToString(t *testing.T) {
	testCases := []struct {
		name     string
		point    orb.Point
		expected string
	}{
		{
			name:     "Positive coordinates",
			point:    orb.Point{10.123456, 20.654321},
			expected: "10.123456,20.654321",
		},
		{
			name:     "Negative coordinates",
			point:    orb.Point{-10.123456, -20.654321},
			expected: "-10.123456,-20.654321",
		},
		{
			name:     "Zero coordinates",
			point:    orb.Point{0.0, 0.0},
			expected: "0.000000,0.000000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := pointToString(tc.point)
			assert.Equal(t, tc.expected, result,
				"pointToString failed for point %v", tc.point)
		})
	}
}

func TestArePointsEqual(t *testing.T) {
	testCases := []struct {
		name     string
		p1       orb.Point
		p2       orb.Point
		expected bool
	}{
		{
			name:     "Identical points",
			p1:       orb.Point{10.0, 20.0},
			p2:       orb.Point{10.0, 20.0},
			expected: true,
		},
		{
			name:     "Points very close (within 1 meter)",
			p1:       orb.Point{10.0, 20.0},
			p2:       orb.Point{10.0000090, 20.0000080}, // Aproximadamente 1 metro de diferencia
			expected: true,
		},
		{
			name:     "Points slightly further apart",
			p1:       orb.Point{10.0, 20.0},
			p2:       orb.Point{10.0001, 20.0001}, // Más de 1 metro de diferencia
			expected: false,
		},
		{
			name:     "Points at different locations",
			p1:       orb.Point{10.0, 20.0},
			p2:       orb.Point{11.0, 21.0},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := arePointsEqual(tc.p1, tc.p2)
			assert.Equal(t, tc.expected, result,
				"arePointsEqual failed for points %v and %v", tc.p1, tc.p2)
		})
	}
}

// TestDistanceBetweenPoints simula el cálculo de distancias entre puntos originales y "snapped"
// utilizando los datos de referencia del endpoint de LocationIQ.
func TestDistanceBetweenPoints(t *testing.T) {
	// Nota:
	// Los puntos se definen como orb.Point{lon, lat} para ser consistentes con la librería.
	// Los siguientes casos usan puntos extraídos (o inferidos) del ejemplo del endpoint.
	//
	// Por ejemplo:
	//   - Downing Street:
	//       Punto original (request): [-0.127627, 51.503355]
	//       Punto "snapped" (response): [-0.126406, 51.503174]
	//       Distancia esperada según response: ~87.10 metros
	//   - King William Street:
	//       Original: [-0.087199, 51.509562]
	//       Snapped:  [-0.0872, 51.509562]
	//       Distancia esperada: ~0.07 metros
	//   - Tower Bridge Approach:
	//       Original: [-0.076134, 51.508037]
	//       Snapped:  [-0.074146, 51.507884]
	//       Distancia esperada: ~139.02 metros
	testCases := []struct {
		name              string
		original, snapped orb.Point
		expectedDistance  float64
	}{
		{
			name:             "Downing Street",
			original:         orb.Point{-0.127627, 51.503355},
			snapped:          orb.Point{-0.126406, 51.503174},
			expectedDistance: 87.100385469,
		},
		{
			name:             "King William Street",
			original:         orb.Point{-0.087199, 51.509562},
			snapped:          orb.Point{-0.0872, 51.509562},
			expectedDistance: 0.069402517,
		},
		{
			name:             "Tower Bridge Approach",
			original:         orb.Point{-0.076134, 51.508037},
			snapped:          orb.Point{-0.074146, 51.507884},
			expectedDistance: 139.018328935,
		},
	}

	// Tolerancia en metros para la comparación de distancias.
	tolerance := 1.0

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calculatedDistance := distanceBetweenPoints(tc.original, tc.snapped)
			diff := math.Abs(calculatedDistance - tc.expectedDistance)
			assert.LessOrEqual(t, diff, tolerance,
				"Para %s, la distancia calculada %.6f difiere de la esperada %.6f por más de %.2f metros",
				tc.name, calculatedDistance, tc.expectedDistance, tolerance)
		})
	}
}
