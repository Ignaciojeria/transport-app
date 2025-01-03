package model

type CreateOrganizationKeyRequest struct {
	Name  string `json:"name" validate:"required" example:"org-name"`
	Email string `json:"email" validate:"required" example:"org-email@gmail.com"`
}
