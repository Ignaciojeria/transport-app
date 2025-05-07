package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapAccountTable(e domain.Operator) table.Account {
	return table.Account{
		Email:    e.Contact.PrimaryEmail,
		IsActive: true,
	}
}
