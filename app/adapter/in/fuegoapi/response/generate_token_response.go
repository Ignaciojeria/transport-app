package response

type GenerateTokenResponse struct {
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt int64  `json:"expiresAt" example:"1714857600"`
	TokenType string `json:"tokenType" example:"Bearer"`
}
