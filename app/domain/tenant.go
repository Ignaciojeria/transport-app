package domain

import (
	"github.com/biter777/countries"
	"github.com/google/uuid"
)

type Tenant struct {
	ID      uuid.UUID
	Name    string
	Country countries.CountryCode
}
