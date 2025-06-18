package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"transport-app/app/domain"
)

type Dimensions struct {
	Height int64  `gorm:"not null" json:"height"`
	Width  int64  `gorm:"not null" json:"width"`
	Length int64  `gorm:"not null" json:"length"`
	Unit   string `gorm:"not null" json:"unit"`
}

type JSONDimensions Dimensions

func (j JSONDimensions) Map() domain.Dimensions {
	return domain.Dimensions{
		Height: j.Height,
		Width:  j.Width,
		Length: j.Length,
		Unit:   j.Unit,
	}
}

// Scan implementa la interfaz sql.Scanner para convertir datos JSON desde la base de datos
func (j *JSONDimensions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONDocuments value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON al guardar en la base de datos
func (j JSONDimensions) Value() (driver.Value, error) {
	return json.Marshal(j)
}
