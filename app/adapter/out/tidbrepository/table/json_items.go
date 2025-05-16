package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"transport-app/app/domain"
)

type Items struct {
	Sku    string `gorm:"not null" json:"sku"`
	Skills []struct {
		Type        string `json:"type"`
		Value       string `json:"value"`
		Description string `json:"description"`
	} `gorm:"default:null" json:"skills"`
	QuantityNumber int            `gorm:"not null" json:"quantity_number"`
	QuantityUnit   string         `gorm:"not null" json:"quantity_unit"`
	JSONInsurance  JSONInsurance  `gorm:"type:json" json:"insurance"`
	Description    string         `gorm:"type:text" json:"description"`
	JSONDimensions JSONDimensions `gorm:"type:json" json:"dimensions"`
	JSONWeight     JSONWeight     `gorm:"type:json" json:"weight"`
}

func (i Items) Map() domain.Item {
	return domain.Item{
		Sku:    i.Sku,
		Skills: i.Skills,
		Quantity: domain.Quantity{
			QuantityNumber: i.QuantityNumber,
			QuantityUnit:   i.QuantityUnit,
		},
		Insurance:   i.JSONInsurance.Map(),
		Description: i.Description,
		Dimensions:  i.JSONDimensions.Map(),
		Weight:      i.JSONWeight.Map(),
	}
}

type JSONItems []Items

func (j JSONItems) Map() []domain.Item {
	items := make([]domain.Item, len(j))
	for i, item := range j {
		items[i] = item.Map()
	}
	return items
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONItems) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONItems value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir datos a JSON al guardarlos en la base de datos
func (j JSONItems) Value() (driver.Value, error) {
	return json.Marshal(j)
}
