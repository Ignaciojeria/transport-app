package request

import (
	"transport-app/app/domain"

	"github.com/biter777/countries"
	"github.com/google/uuid"
)

type CreateTenantRequest struct {
	ID      string
	Name    string `json:"name" validate:"required" example:"org-name"`
	Email   string `json:"email" validate:"required" example:"org-email@gmail.com"`
	Country string `json:"country" validate:"required" example:"CL"`
}

func (r CreateTenantRequest) Map() domain.TenantAccount {
	uid, _ := uuid.Parse(r.ID)
	return domain.TenantAccount{
		Tenant: domain.Tenant{
			ID:      uid,
			Name:    r.Name,
			Country: countries.ByName(r.Country),
		},
		Account: domain.Account{
			Email: r.Email,
		},
		Role:   "owner",
		Status: "active",
	}
}
