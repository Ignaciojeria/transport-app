package request

type CreateOrganizationRequest struct {
	Name    string `json:"name" validate:"required" example:"org-name"`
	Email   string `json:"email" validate:"required" example:"org-email@gmail.com"`
	Country string `json:"country" validate:"required" example:"CL"`
}
