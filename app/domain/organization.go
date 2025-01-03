package domain

import "github.com/biter777/countries"

type Organization struct {
	Country countries.CountryCode `json:"country"`
	Name    string                `json:"name"`
	Email   string                `json:"email"`
	Key     string                `json:"key"`
}
