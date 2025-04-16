package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertOperator func(context.Context, domain.Operator) (domain.Operator, error)

func init() {
	ioc.Registry(NewUpsertOperator, database.NewConnectionFactory)
}
func NewUpsertOperator(conn database.ConnectionFactory) UpsertOperator {
	return func(ctx context.Context, o domain.Operator) (domain.Operator, error) {
		/*
			operatorTBL := mapper.MapOperator(o)
			err := conn.DB.WithContext(ctx).
				Table("accounts").
				Where("reference_id = ? AND organization_id = ?",
					//o.ReferenceID,
					"TODO").
				First(&operatorTBL).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.Operator{}, err
			}
			operatorToUpdate := operatorTBL.MapOperator(o.Organization).UpdateIfChanged(o)
			DBOperatorToUpdate := mapper.MapOperator(operatorToUpdate)
			DBOperatorToUpdate.CreatedAt = operatorTBL.CreatedAt
			return DBOperatorToUpdate.MapOperator(o.Organization), conn.
				Omit("Contact").
				Omit("AddressInfo").
				Omit("OriginNodeInfo").
				Omit("Organization").
				Save(&DBOperatorToUpdate).Error*/
		return domain.Operator{}, nil
	}
}
