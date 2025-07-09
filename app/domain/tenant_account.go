package domain

import "time"

type TenantAccount struct {
	Tenant   Tenant
	Account  Account
	Role     string
	Status   string
	Invited  bool
	JoinedAt time.Time
}
