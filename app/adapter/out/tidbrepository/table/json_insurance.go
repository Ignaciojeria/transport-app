package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"transport-app/app/domain"
)

type Insurance struct {
	UnitValue float64 `gorm:"not null" json:"unit_value"`
	Currency  string  `gorm:"not null" json:"currency"`
}

type JSONInsurance Insurance

func (j JSONInsurance) Map() domain.Insurance {
	return domain.Insurance{
		UnitValue: j.UnitValue,
		Currency:  j.Currency,
	}
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONInsurance) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONInsurance) Value() (driver.Value, error) {
	return json.Marshal(j)
}
