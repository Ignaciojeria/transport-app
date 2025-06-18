package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"transport-app/app/domain"
)

type Weight struct {
	WeightValue int64  `gorm:"not null" json:"weight_value"`
	WeightUnit  string `gorm:"not null" json:"weight_unit"`
}

type JSONWeight Weight

func (j JSONWeight) Map() domain.Weight {
	return domain.Weight{
		Value: j.WeightValue,
		Unit:  j.WeightUnit,
	}
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONWeight) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONWeight) Value() (driver.Value, error) {
	return json.Marshal(j)
}
