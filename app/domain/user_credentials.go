package domain

type UserCredentials struct {
	Email    string
	Password string
}

type ProviderToken struct {
	TokenType    string
	TokenValue   string
	RefreshToken string
	ExpiresIn    int64
	Provider     string
}
