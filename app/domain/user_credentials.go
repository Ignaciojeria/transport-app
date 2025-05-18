package domain

type UserCredentials struct {
	Email    string
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
