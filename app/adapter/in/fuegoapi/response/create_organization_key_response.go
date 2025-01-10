package response

type CreateOrganizationKeyResponse struct {
	OrganizationKey string `json:"organizationKey" example:"f59d9c17-19ab-465c-9dc9-4a48214797fd"`
	Message         string `json:"message" example:"Organization key created successfully. Please save it securely as it will be required for API authentication."`
}
