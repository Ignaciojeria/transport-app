package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type EvidencePhoto struct {
	URL     string
	Type    string
	TakenAt *time.Time
}

type JSONEvidencePhotos []EvidencePhoto

// Implementamos los métodos necesarios para el manejo de JSON
func (j *JSONEvidencePhotos) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONEvidencePhotos value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON
func (j JSONEvidencePhotos) Value() (driver.Value, error) {
	return json.Marshal(j)
}
