package response

type CreateTenantResponse struct {
	ID      string `json:"id" example:"ed469f7b-1f18-47b5-aca6-6a3df543333b"`
	Country string `json:"country" example:"CL"`
	Tenant  string `json:"tenant" example:"ed469f7b-1f18-47b5-aca6-6a3df543333b-CL"`
	Message string `json:"message" example:"Tenant created successfully."`
}
