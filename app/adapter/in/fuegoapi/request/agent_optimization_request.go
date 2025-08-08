package request

type AgentOptimizationRequest struct {
	Fleet  []map[string]interface{} `json:"fleet"`
	Visits []map[string]interface{} `json:"visits"`
}

// IterateVisitsInBatches itera sobre las visitas en lotes del tamaño especificado
// y ejecuta la función callback para cada lote
func (r *AgentOptimizationRequest) IterateVisitsInBatches(batchSize int, callback func([]map[string]interface{}) error) error {
	if batchSize <= 0 {
		batchSize = 100 // tamaño por defecto
	}

	totalVisits := len(r.Visits)
	if totalVisits == 0 {
		return nil
	}

	for i := 0; i < totalVisits; i += batchSize {
		end := i + batchSize
		if end > totalVisits {
			end = totalVisits
		}

		batch := r.Visits[i:end]
		if err := callback(batch); err != nil {
			return err
		}
	}

	return nil
}

// GetVisitsBatch devuelve un lote específico de visitas
func (r *AgentOptimizationRequest) GetVisitsBatch(startIndex, batchSize int) []map[string]interface{} {
	if startIndex < 0 || batchSize <= 0 {
		return nil
	}

	totalVisits := len(r.Visits)
	if startIndex >= totalVisits {
		return nil
	}

	endIndex := startIndex + batchSize
	if endIndex > totalVisits {
		endIndex = totalVisits
	}

	return r.Visits[startIndex:endIndex]
}

// GetTotalVisitsCount devuelve el número total de visitas
func (r *AgentOptimizationRequest) GetTotalVisitsCount() int {
	return len(r.Visits)
}

// GetTotalBatches calcula el número total de lotes necesarios para procesar todas las visitas
func (r *AgentOptimizationRequest) GetTotalBatches(batchSize int) int {
	if batchSize <= 0 {
		batchSize = 100
	}

	totalVisits := len(r.Visits)
	if totalVisits == 0 {
		return 0
	}

	return (totalVisits + batchSize - 1) / batchSize
}
