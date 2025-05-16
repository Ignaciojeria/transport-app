package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"transport-app/app/domain"
)

type ItemReference struct {
	Sku            string `json:"sku"`
	QuantityNumber int    `json:"quantity_number"`
	QuantityUnit   string `json:"quantity_unit"`
}

type JSONItemReferences []ItemReference

func (j JSONItemReferences) Map() []domain.ItemReference {
	mappedReferences := make([]domain.ItemReference, len(j))
	for i, ref := range j {
		mappedReferences[i] = domain.ItemReference{
			Sku: ref.Sku,
			Quantity: domain.Quantity{
				QuantityNumber: ref.QuantityNumber,
				QuantityUnit:   ref.QuantityUnit,
			},
		}
	}
	return mappedReferences
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONItemReferences) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONItemReferences) Value() (driver.Value, error) {
	return json.Marshal(j)
}
