package domain

import "github.com/biter777/countries"

type Organization struct {
	ID                    int64
	OrganizationCountryID int64
	Country               countries.CountryCode `json:"country"`
	Name                  string                `json:"name"`
	Email                 string                `json:"email"`
	Key                   string                `json:"key"`
}
