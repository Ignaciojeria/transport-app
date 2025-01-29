package domain

type Operator struct {
	ID           int64
	Organization Organization
	OriginNode   NodeInfo
	Contact      Contact `json:"contact"`
	Type         string  `json:"type"`
}
