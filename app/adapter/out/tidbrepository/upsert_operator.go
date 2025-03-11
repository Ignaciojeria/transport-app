package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertOperator func(context.Context, domain.Operator) (domain.Operator, error)

func init() {
	ioc.Registry(NewUpsertOperator, tidb.NewTIDBConnection)
}
func NewUpsertOperator(conn tidb.TIDBConnection) UpsertOperator {
	return func(ctx context.Context, o domain.Operator) (domain.Operator, error) {
		operatorTBL := mapper.MapOperator(o)
		err := conn.DB.WithContext(ctx).
			Table("accounts").
			Where("reference_id = ? AND organization_id = ?",
				o.ReferenceID,
				o.Organization.ID).
			First(&operatorTBL).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Operator{}, err
		}
		operatorToUpdate := operatorTBL.MapOperator().UpdateIfChanged(o)
		DBOperatorToUpdate := mapper.MapOperator(operatorToUpdate)
		DBOperatorToUpdate.CreatedAt = operatorTBL.CreatedAt
		return DBOperatorToUpdate.MapOperator(), conn.
			Omit("Contact").
			Omit("AddressInfo").
			Omit("OriginNodeInfo").
			Omit("Organization").
			Save(&DBOperatorToUpdate).Error
	}
}
