package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"transport-app/app/domain"
)

type Reference struct {
	Type  string `gorm:"not null" json:"type"`
	Value string `gorm:"not null" json:"value"`
}

type JSONReference []Reference

func (j JSONReference) Map() []domain.Reference {
	mappedReferences := make([]domain.Reference, len(j))
	for i, ref := range j {
		mappedReferences[i] = domain.Reference{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return mappedReferences
}

func (j JSONReference) MapDocuments() []domain.Document {
	mappedReferences := make([]domain.Document, len(j))
	for i, ref := range j {
		mappedReferences[i] = domain.Document{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return mappedReferences
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONReference) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONReference) Value() (driver.Value, error) {
	return json.Marshal(j)
}
