package request

type TokenRequest struct {
	GrantType    string `json:"grant_type" example:"client_credentials"`
	ClientID     string `json:"client_id" example:"abc123"`
	ClientSecret string `json:"client_secret" example:"secret456"`
	Scope        string `json:"scope" example:"read:orders"`
}
