package request

type SearchAccountsRequest struct {
	Pagination struct {
		Page int `json:"page"`
		Size int `json:"size"`
	} `json:"pagination"`
	Emails []string `json:"emails"`
}
