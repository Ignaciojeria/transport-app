package response

type CreateOrganizationResponse struct {
	OrganizationKey string `json:"organizationKey"`
	Message         string `json:"message" example:"Organization created successfully."`
}
