package request

type RegisterRequest struct {
	Email            string `json:"email"`
	OrganizationName string `json:"organizationName"`
	Country          string `json:"country"`
}
