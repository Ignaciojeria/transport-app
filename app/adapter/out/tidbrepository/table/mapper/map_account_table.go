package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapAccountTable(e domain.Account) table.Account {
	return table.Account{
		Email:      e.Email,
		DocumentID: e.DocID().String(),
		IsActive:   true,
	}
}
