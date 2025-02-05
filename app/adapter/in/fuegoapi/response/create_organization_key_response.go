package response

type CreateOrganizationKeyResponse struct {
	OrganizationKey string `json:"organizationKey" example:"f59d9c17-19ab-465c-9dc9-4a48214797fd"`
	Country         string `json:"country" example:"CL"`
	Message         string `json:"message" example:"Organization key created successfully."`
}
