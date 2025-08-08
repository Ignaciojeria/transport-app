package request

import (
	"fmt"
	"testing"
)

func TestAgentOptimizationRequest_IterateVisitsInBatches(t *testing.T) {
	// Preparar datos de prueba
	request := &AgentOptimizationRequest{
		Fleet: []map[string]interface{}{
			{"id": "vehicle1", "capacity": 100},
			{"id": "vehicle2", "capacity": 150},
		},
		Visits: []map[string]interface{}{
			{"id": "visit1", "location": "A", "demand": 10},
			{"id": "visit2", "location": "B", "demand": 15},
			{"id": "visit3", "location": "C", "demand": 20},
			{"id": "visit4", "location": "D", "demand": 25},
			{"id": "visit5", "location": "E", "demand": 30},
			{"id": "visit6", "location": "F", "demand": 35},
		},
	}

	// Test 1: Procesar en lotes de 2 visitas
	t.Run("Procesar en lotes de 2 visitas", func(t *testing.T) {
		var processedBatches [][]map[string]interface{}
		var totalVisits int

		err := request.IterateVisitsInBatches(2, func(batch []map[string]interface{}) error {
			processedBatches = append(processedBatches, batch)
			totalVisits += len(batch)
			return nil
		})

		if err != nil {
			t.Errorf("Error inesperado: %v", err)
		}

		// Verificar que se procesaron 3 lotes (6 visitas / 2 por lote)
		expectedBatches := 3
		if len(processedBatches) != expectedBatches {
			t.Errorf("Se esperaban %d lotes, se obtuvieron %d", expectedBatches, len(processedBatches))
		}

		// Verificar que se procesaron todas las visitas
		if totalVisits != 6 {
			t.Errorf("Se esperaban 6 visitas procesadas, se obtuvieron %d", totalVisits)
		}

		// Verificar el contenido de cada lote
		expectedBatchSizes := []int{2, 2, 2}
		for i, batch := range processedBatches {
			if len(batch) != expectedBatchSizes[i] {
				t.Errorf("Lote %d: se esperaban %d visitas, se obtuvieron %d", i, expectedBatchSizes[i], len(batch))
			}
		}
	})

	// Test 2: Procesar en lotes de 3 visitas
	t.Run("Procesar en lotes de 3 visitas", func(t *testing.T) {
		var processedBatches [][]map[string]interface{}

		err := request.IterateVisitsInBatches(3, func(batch []map[string]interface{}) error {
			processedBatches = append(processedBatches, batch)
			return nil
		})

		if err != nil {
			t.Errorf("Error inesperado: %v", err)
		}

		// Verificar que se procesaron 2 lotes (6 visitas / 3 por lote)
		expectedBatches := 2
		if len(processedBatches) != expectedBatches {
			t.Errorf("Se esperaban %d lotes, se obtuvieron %d", expectedBatches, len(processedBatches))
		}

		// Verificar el tamaño de cada lote
		expectedBatchSizes := []int{3, 3}
		for i, batch := range processedBatches {
			if len(batch) != expectedBatchSizes[i] {
				t.Errorf("Lote %d: se esperaban %d visitas, se obtuvieron %d", i, expectedBatchSizes[i], len(batch))
			}
		}
	})

	// Test 3: Procesar con error en callback
	t.Run("Procesar con error en callback", func(t *testing.T) {
		expectedError := fmt.Errorf("error de procesamiento")

		err := request.IterateVisitsInBatches(2, func(batch []map[string]interface{}) error {
			return expectedError
		})

		if err != expectedError {
			t.Errorf("Se esperaba error %v, se obtuvo %v", expectedError, err)
		}
	})

	// Test 4: Procesar con tamaño de lote inválido (usa valor por defecto)
	t.Run("Procesar con tamaño de lote inválido", func(t *testing.T) {
		var processedBatches [][]map[string]interface{}

		err := request.IterateVisitsInBatches(0, func(batch []map[string]interface{}) error {
			processedBatches = append(processedBatches, batch)
			return nil
		})

		if err != nil {
			t.Errorf("Error inesperado: %v", err)
		}

		// Con 6 visitas y lote por defecto de 100, debería haber 1 lote
		if len(processedBatches) != 1 {
			t.Errorf("Se esperaba 1 lote, se obtuvieron %d", len(processedBatches))
		}
	})
}

func TestAgentOptimizationRequest_GetVisitsBatch(t *testing.T) {
	request := &AgentOptimizationRequest{
		Visits: []map[string]interface{}{
			{"id": "visit1", "location": "A"},
			{"id": "visit2", "location": "B"},
			{"id": "visit3", "location": "C"},
			{"id": "visit4", "location": "D"},
			{"id": "visit5", "location": "E"},
		},
	}

	tests := []struct {
		name       string
		startIndex int
		batchSize  int
		expected   int
	}{
		{"Obtener lote válido", 0, 2, 2},
		{"Obtener lote parcial al final", 3, 3, 2},
		{"Índice fuera de rango", 10, 2, 0},
		{"Tamaño de lote inválido", 0, 0, 0},
		{"Índice negativo", -1, 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batch := request.GetVisitsBatch(tt.startIndex, tt.batchSize)
			if len(batch) != tt.expected {
				t.Errorf("Se esperaban %d visitas, se obtuvieron %d", tt.expected, len(batch))
			}
		})
	}
}

func TestAgentOptimizationRequest_GetTotalVisitsCount(t *testing.T) {
	request := &AgentOptimizationRequest{
		Visits: []map[string]interface{}{
			{"id": "visit1"},
			{"id": "visit2"},
			{"id": "visit3"},
		},
	}

	count := request.GetTotalVisitsCount()
	if count != 3 {
		t.Errorf("Se esperaban 3 visitas, se obtuvieron %d", count)
	}

	// Test con visitas vacías
	emptyRequest := &AgentOptimizationRequest{
		Visits: []map[string]interface{}{},
	}

	emptyCount := emptyRequest.GetTotalVisitsCount()
	if emptyCount != 0 {
		t.Errorf("Se esperaban 0 visitas, se obtuvieron %d", emptyCount)
	}
}

func TestAgentOptimizationRequest_GetTotalBatches(t *testing.T) {
	request := &AgentOptimizationRequest{
		Visits: []map[string]interface{}{
			{"id": "visit1"},
			{"id": "visit2"},
			{"id": "visit3"},
			{"id": "visit4"},
			{"id": "visit5"},
		},
	}

	tests := []struct {
		name      string
		batchSize int
		expected  int
	}{
		{"Lotes de tamaño 2", 2, 3},
		{"Lotes de tamaño 3", 3, 2},
		{"Lotes de tamaño 5", 5, 1},
		{"Lotes de tamaño 10", 10, 1},
		{"Tamaño inválido (usa valor por defecto)", 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batches := request.GetTotalBatches(tt.batchSize)
			if batches != tt.expected {
				t.Errorf("Se esperaban %d lotes, se obtuvieron %d", tt.expected, batches)
			}
		})
	}

	// Test con visitas vacías
	emptyRequest := &AgentOptimizationRequest{
		Visits: []map[string]interface{}{},
	}

	emptyBatches := emptyRequest.GetTotalBatches(10)
	if emptyBatches != 0 {
		t.Errorf("Se esperaban 0 lotes, se obtuvieron %d", emptyBatches)
	}
}

// Ejemplo de uso completo
func ExampleAgentOptimizationRequest_IterateVisitsInBatches() {
	// Crear una solicitud de optimización con datos de ejemplo
	request := &AgentOptimizationRequest{
		Fleet: []map[string]interface{}{
			{"id": "truck1", "capacity": 1000, "type": "delivery"},
			{"id": "truck2", "capacity": 1500, "type": "delivery"},
		},
		Visits: []map[string]interface{}{
			{"id": "visit1", "location": "Madrid", "demand": 100, "priority": "high"},
			{"id": "visit2", "location": "Barcelona", "demand": 150, "priority": "medium"},
			{"id": "visit3", "location": "Valencia", "demand": 200, "priority": "low"},
			{"id": "visit4", "location": "Sevilla", "demand": 120, "priority": "high"},
			{"id": "visit5", "location": "Bilbao", "demand": 80, "priority": "medium"},
			{"id": "visit6", "location": "Málaga", "demand": 90, "priority": "low"},
		},
	}

	// Procesar visitas en lotes de 2
	fmt.Printf("Total de visitas: %d\n", request.GetTotalVisitsCount())
	fmt.Printf("Total de lotes necesarios: %d\n", request.GetTotalBatches(2))

	var processedVisits int
	err := request.IterateVisitsInBatches(2, func(batch []map[string]interface{}) error {
		fmt.Printf("Procesando lote con %d visitas:\n", len(batch))
		for _, visit := range batch {
			fmt.Printf("  - Visit ID: %s, Location: %s, Demand: %v\n",
				visit["id"], visit["location"], visit["demand"])
			processedVisits++
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Printf("Error durante el procesamiento: %v\n", err)
		return
	}

	fmt.Printf("Total de visitas procesadas: %d\n", processedVisits)

	// Output:
	// Total de visitas: 6
	// Total de lotes necesarios: 3
	// Procesando lote con 2 visitas:
	//   - Visit ID: visit1, Location: Madrid, Demand: 100
	//   - Visit ID: visit2, Location: Barcelona, Demand: 150
	//
	// Procesando lote con 2 visitas:
	//   - Visit ID: visit3, Location: Valencia, Demand: 200
	//   - Visit ID: visit4, Location: Sevilla, Demand: 120
	//
	// Procesando lote con 2 visitas:
	//   - Visit ID: visit5, Location: Bilbao, Demand: 80
	//   - Visit ID: visit6, Location: Málaga, Demand: 90
	//
	// Total de visitas procesadas: 6
}
