package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type RouteEndLocation struct {
	Longitude float64
	Latitude  float64
}

type JSONRouteEndLocation RouteEndLocation

// Implementamos los métodos necesarios para el manejo de JSON
func (j *JSONRouteEndLocation) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONPlanLocation value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON
func (j JSONRouteEndLocation) Value() (driver.Value, error) {
	return json.Marshal(j)
}
