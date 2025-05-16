package domain

import (
	"github.com/biter777/countries"
	"github.com/google/uuid"
)

type Tenant struct {
	ID       uuid.UUID
	Operator Operator
	Country  countries.CountryCode `json:"country"`
	Name     string                `json:"name"`
}
