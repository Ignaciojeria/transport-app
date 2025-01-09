package model

type SearchOrdersRequest struct {
	Pagination struct {
		Page int `json:"page"`
		Size int `json:"size"`
	} `json:"pagination"`
	Commerces    []string `json:"commerces"`
	ReferenceIDs []string `json:"referenceIDs"`
	PackageLpns  []string `json:"packageLpns"`
}
