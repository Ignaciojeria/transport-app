package domain

type Entitlement struct {
	V        int    `json:"v" example:"1"`
	UserID   string `json:"user_id" example:"763a590a-9b8e-4a91-b8ee-47f2a64d003d"`
	Plan     string `json:"plan" example:"pro"`
	Status   string `json:"status" example:"trialing"`
	Access   bool   `json:"access" example:"true"`
	StartsAt string `json:"starts_at" example:"2026-01-03T01:36:48.000Z"`
	EndsAt   string `json:"ends_at" example:"2026-01-17T01:36:48.000Z"`
}
