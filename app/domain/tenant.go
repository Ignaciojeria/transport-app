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

/*
func (o Organization) GetOrgKey() string {
	return strconv.FormatInt(o.ID, 10) + "-" + o.Country.Alpha2()
}

func (o *Organization) SetKey(key string) error {
	parts := strings.Split(key, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid key format")
	}
	// Parse ID
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID in key: %v", err)
	}
	o.ID = id
	// Set Country
	countryCode := parts[1]
	o.Country = countries.ByName(countryCode)
	return nil
}*/
