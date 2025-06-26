package model

import (
	"fmt"
	"testing"
	"transport-app/app/domain/optimization"
)

func TestGroupReferencesByCoordinates(t *testing.T) {
	// Crear datos de prueba
	req := optimization.FleetOptimization{
		Visits: []optimization.Visit{
			{
				Orders: []optimization.Order{
					{ReferenceID: "REF001"},
					{ReferenceID: "REF002"},
				},
			},
			{
				Orders: []optimization.Order{
					{ReferenceID: "REF003"},
				},
			},
		},
	}

	// Crear steps con coordenadas similares
	steps := []Step{
		{
			Type:     "delivery",
			Job:      1,
			Location: [2]float64{-70.123456, -33.123456}, // lon, lat
		},
		{
			Type:     "delivery",
			Job:      2,
			Location: [2]float64{-70.123457, -33.123457}, // coordenadas muy similares
		},
		{
			Type:     "pickup",
			Shipment: 1,
			Location: [2]float64{-70.500000, -33.500000}, // coordenadas diferentes
		},
	}

	// Probar la agrupación
	coordinateGroups := groupReferencesByCoordinates(steps, req)

	fmt.Println("Grupos de coordenadas:")
	for coordKey, refs := range coordinateGroups {
		fmt.Printf("Coordenadas: %s -> Referencias: %v\n", coordKey, refs)
	}

	// Verificar que las coordenadas similares se agrupen
	if len(coordinateGroups) < 2 {
		t.Errorf("Se esperaban al menos 2 grupos de coordenadas, se obtuvieron %d", len(coordinateGroups))
	}

	// Verificar que las referencias se agrupen correctamente
	for coordKey, refs := range coordinateGroups {
		if len(refs) == 0 {
			t.Errorf("El grupo %s no tiene referencias", coordKey)
		}
		fmt.Printf("Grupo %s tiene %d referencias: %v\n", coordKey, len(refs), refs)
	}
}

func TestGetGroupedReferencesForStep(t *testing.T) {
	// Crear datos de prueba
	req := optimization.FleetOptimization{
		Visits: []optimization.Visit{
			{
				Orders: []optimization.Order{
					{ReferenceID: "REF001"},
					{ReferenceID: "REF002"},
				},
			},
			{
				Orders: []optimization.Order{
					{ReferenceID: "REF003"},
				},
			},
		},
	}

	// Crear steps con coordenadas similares
	allSteps := []Step{
		{
			Type:     "delivery",
			Job:      1,
			Location: [2]float64{-70.123456, -33.123456},
		},
		{
			Type:     "delivery",
			Job:      2,
			Location: [2]float64{-70.123457, -33.123457}, // coordenadas muy similares
		},
	}

	// Probar obtener referencias agrupadas para el primer step
	step := allSteps[0]
	groupedRefs := getGroupedReferencesForStep(step, req, allSteps)

	fmt.Printf("Referencias agrupadas para step en %v: %v\n", step.Location, groupedRefs)

	// Verificar que se obtengan todas las referencias de ambos steps
	expectedRefs := []string{"REF001", "REF002", "REF003"}
	if len(groupedRefs) != len(expectedRefs) {
		t.Errorf("Se esperaban %d referencias, se obtuvieron %d", len(expectedRefs), len(groupedRefs))
	}

	// Verificar que no haya duplicados
	seen := make(map[string]bool)
	for _, ref := range groupedRefs {
		if seen[ref] {
			t.Errorf("Referencia duplicada encontrada: %s", ref)
		}
		seen[ref] = true
	}
}

func ExampleGroupReferencesByCoordinates() {
	// Este ejemplo muestra cómo se agrupan las referencias por coordenadas
	req := optimization.FleetOptimization{
		Visits: []optimization.Visit{
			{
				Orders: []optimization.Order{
					{ReferenceID: "REF001"},
					{ReferenceID: "REF002"},
				},
			},
			{
				Orders: []optimization.Order{
					{ReferenceID: "REF003"},
				},
			},
		},
	}

	steps := []Step{
		{
			Type:     "delivery",
			Job:      1,
			Location: [2]float64{-70.123456, -33.123456},
		},
		{
			Type:     "delivery",
			Job:      2,
			Location: [2]float64{-70.123457, -33.123457}, // coordenadas muy similares
		},
	}

	coordinateGroups := groupReferencesByCoordinates(steps, req)

	fmt.Println("Agrupación de referencias por coordenadas:")
	for coordKey, refs := range coordinateGroups {
		fmt.Printf("Coordenadas: %s -> Referencias: %v\n", coordKey, refs)
	}

	// Output:
	// Agrupación de referencias por coordenadas:
	// Coordenadas: -33.123456,-70.123456 -> Referencias: [REF001 REF002 REF003]
}
