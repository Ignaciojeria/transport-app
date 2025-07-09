package request

type RegisterRequest struct {
	Email   string `json:"email"`
	Country string `json:"country"`
}
