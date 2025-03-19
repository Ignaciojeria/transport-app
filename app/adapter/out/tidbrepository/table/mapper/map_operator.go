package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOperator(e domain.Operator) table.Account {
	return table.Account{
		ID:       e.ID,
		IsActive: true,
	}
}
