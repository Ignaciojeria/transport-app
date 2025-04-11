package domain

import "context"

type Operator struct {
	OriginNode NodeInfo
	Contact    Contact `json:"contact"`
	Carrier    Carrier
	Role       Role `json:"type"`
}

func (o Operator) DocID(ctx context.Context) DocumentID {
	return Hash(ctx, o.Contact.PrimaryEmail)
}

type Role string

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleDriver, RolePlanner, RoleDispatcher, RoleSeller:
		return true
	default:
		return false
	}
}

const (
	RoleAdmin      Role = "admin"
	RoleDriver     Role = "driver"
	RolePlanner    Role = "planner"
	RoleDispatcher Role = "dispatcher"
	RoleSeller     Role = "seller"
)

func (r Role) String() string {
	return string(r)
}
