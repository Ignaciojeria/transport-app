package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertOrderPackages func(context.Context, domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertOrderPackages,
		database.NewConnectionFactory)
}
func NewUpsertOrderPackages(conn database.ConnectionFactory) UpsertOrderPackages {
	return func(ctx context.Context, order domain.Order) error {
		orderPackages := mapper.MapOrderPackages(ctx, order)
		return conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("order_doc = ?", order.DocID(ctx)).
				Delete(&table.OrderPackage{}).Error; err != nil {
				return err
			}
			if len(orderPackages) > 0 {
				if err := tx.Save(&orderPackages).Error; err != nil {
					return err
				}
			}

			if len(orderPackages) == 0 {
				pkg := domain.Package{}
				if err := tx.Save(&table.OrderPackage{
					OrderDoc:   order.DocID(ctx).String(),
					PackageDoc: pkg.DocID(ctx, order.ReferenceID.String()).String()}).Error; err != nil {
					return err
				}
			}

			return nil
		})
	}
}
