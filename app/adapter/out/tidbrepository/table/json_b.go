package table

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB json.RawMessage

// Scan implementa la interfaz `sql.Scanner` para deserializar valores JSON
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = JSONB([]byte("null"))
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed for JSONB")
	}
	*j = JSONB(bytes)
	return nil
}

// Value implementa la interfaz `driver.Valuer` para serializar valores JSON
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
