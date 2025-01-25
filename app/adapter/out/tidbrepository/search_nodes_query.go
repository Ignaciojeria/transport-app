package tidbrepository

import (
	"context"
	views "transport-app/app/adapter/out/tidbrepository/views"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchNodesQuery func(context.Context, domain.NodeSearchFilters) ([]domain.NodeInfo, error)

func init() {
	ioc.Registry(NewSearchNodesQuery, tidb.NewTIDBConnection)
}
func NewSearchNodesQuery(conn tidb.TIDBConnection) SearchNodesQuery {
	return func(ctx context.Context, p domain.NodeSearchFilters) ([]domain.NodeInfo, error) {
		var nodes views.SearchNodesView

		query := `
            SELECT ni.name AS node_name, ni.reference_id
            FROM node_infos ni
            JOIN organization_countries oc ON ni.organization_country_id = oc.id
            JOIN organizations org ON oc.organization_id = org.id
            JOIN api_keys ak ON org.id = ak.organization_id
            WHERE ak.key = ? AND oc.country = ?
            LIMIT ? OFFSET ?
        `
		params := []interface{}{p.Organization.Key, p.Organization.Country.Alpha2(), p.Size, p.Page}

		if err := conn.DB.WithContext(ctx).Raw(query, params...).Scan(&nodes).Error; err != nil {
			return nil, err
		}

		return nodes.Map(), nil
	}
}
