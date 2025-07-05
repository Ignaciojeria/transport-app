package request

type GenerateTokenRequest struct {
	Sub      string            `json:"sub" example:"1234567890"`
	Scopes   []string          `json:"scopes"`
	Context  map[string]string `json:"context"`
	Audience string            `json:"aud" example:"https://api.miapp.com"`
}
