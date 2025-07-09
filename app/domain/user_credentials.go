package domain

import "github.com/biter777/countries"

type UserCredentials struct {
	Email    string
	Country  countries.CountryCode
	Password string
}

func (u UserCredentials) DocID() DocumentID {
	return HashInputs(u.Email)
}

type ProviderToken struct {
	TokenType    string
	TokenValue   string
	RefreshToken string
	ExpiresIn    int64
	Provider     string
}
