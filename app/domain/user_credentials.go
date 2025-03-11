package domain

type UserCredentials struct {
	PrimaryOrganization  Organization   `json:"primaryOrganization"`            // Organización principal del usuario
	InvitedOrganizations []Organization `json:"invitedOrganizations,omitempty"` // Organizaciones donde ha sido invitado
	Email                string         `json:"email"`
	Password             string         `json:"password"`
}

type ProviderToken struct {
	TokenType    string `json:"tokenType"`
	TokenValue   string `json:"tokenValue"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn,omitempty"`
	Provider     string `json:"provider,omitempty"`
}
