package domain

type Operator struct {
	Organization Organization
	OriginNode   NodeInfo
	Contact      Contact `json:"contact"`
	Carrier      Carrier
	Role         Role `json:"type"`
}

func (o Operator) DocID() DocumentID {
	return Hash(o.Organization, o.Contact.PrimaryEmail)
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

func (o Operator) UpdateIfChanged(newOperator Operator) Operator {
	if newOperator.Role.String() != "" {
		o.Role = newOperator.Role
	}
	return o
}
