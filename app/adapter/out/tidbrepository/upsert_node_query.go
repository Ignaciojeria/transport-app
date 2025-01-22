package tidbrepository

import (
	"context"
	"fmt"
	views "transport-app/app/adapter/out/tidbrepository/views"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertNodeQuery func(
	ctx context.Context,
	origin domain.NodeInfo) (domain.NodeInfo, error)

func init() {
	ioc.Registry(
		NewUpsertNodeQuery,
		tidb.NewTIDBConnection)
}

func NewUpsertNodeQuery(conn tidb.TIDBConnection) UpsertNodeQuery {
	return func(ctx context.Context, origin domain.NodeInfo) (domain.NodeInfo, error) {
		var flattenedNode views.FlattenedNodeView

		// Realiza la consulta combinando las tablas originales
		result := conn.DB.WithContext(ctx).
			Table("node_infos AS n").
			Select(`
                n.id AS node_id,
                n.reference_id,
                n.name AS node_name,
                n.type AS node_type,
                n.address_id,
                a.address_line1,
                a.address_line2,
                a.address_line3,
                a.county,
                a.district,
                a.latitude,
                a.longitude,
                a.province,
                a.state,
                a.zip_code,
                a.time_zone,
                o.id AS operator_id,
                o.contact_id,
                c.full_name AS operator_name,
                c.email AS operator_email,
                c.phone AS operator_phone,
                c.national_id AS operator_national_id,
                c.documents AS operator_documents,
                n.node_references
            `).
			Joins("LEFT JOIN operators AS o ON n.operator_id = o.id").
			Joins("LEFT JOIN contacts AS c ON o.contact_id = c.id").
			Joins("LEFT JOIN address_infos AS a ON n.address_id = a.id").
			Where("n.reference_id = ?", origin.ReferenceID).
			Scan(&flattenedNode)

		// Manejo de errores
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return domain.NodeInfo{}, fmt.Errorf("node with ReferenceID '%s' not found", origin.ReferenceID)
			}
			return domain.NodeInfo{}, fmt.Errorf("failed to query node: %w", result.Error)
		}

		// Mapea el resultado al dominio y retorna
		return flattenedNode.ToNodeInfo(), nil
	}
}
