package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/projectionresult"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/projection/nodes"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/doug-martin/goqu/v9"
)

type FindNodesProjectionResult func(
	ctx context.Context,
	filters domain.NodesFilter) (projectionresult.NodesProjectionResults, bool, error)

func init() {
	ioc.Registry(
		NewFindNodesProjectionResult,
		database.NewConnectionFactory,
		nodes.NewProjection,
	)
}

func NewFindNodesProjectionResult(
	conn database.ConnectionFactory,
	projection nodes.Projection) FindNodesProjectionResult {
	const (
		ni = "ni" // node_infos
		nt = "nt" // node_types
		ai = "ai" // address_infos
		c  = "c"  // contacts
		d  = "d"  // districts
		p  = "p"  // provinces
		s  = "s"  // states
	)

	return func(ctx context.Context, filters domain.NodesFilter) (projectionresult.NodesProjectionResults, bool, error) {
		var results projectionresult.NodesProjectionResults
		hasMoreResults := false

		ds := goqu.From(goqu.T("node_infos").As(ni)).
			Select(goqu.I(ni + ".document_id")).
			Where(goqu.Ex{
				ni + ".tenant_id": sharedcontext.TenantIDFromContext(ctx),
			})

		// Agregar filtro por reference_ids si existen
		if len(filters.ReferenceIds) > 0 {
			ds = ds.Where(goqu.I(ni + ".reference_id").In(filters.ReferenceIds))
		}

		// Agregar filtro por nombre si existe
		if filters.Name != nil {
			ds = ds.Where(goqu.I(ni + ".name").Like("%" + *filters.Name + "%"))
		}

		// Agregar filtro por tipo de nodo si existe
		if filters.NodeType != nil {
			ds = ds.InnerJoin(
				goqu.T("node_types").As(nt),
				goqu.On(goqu.I(nt+".document_id").Eq(goqu.I(ni+".node_type_doc"))),
			).
				Where(goqu.I(nt + ".value").Eq(filters.NodeType.Value))
		}

		// Agregar filtro por referencias si existen
		if len(filters.References) > 0 {
			const nrf = "nrf" // alias exclusivo para evitar colisión con la CTE `node_refs`

			inRefs := []string{}
			for _, ref := range filters.References {
				ref := domain.Reference{
					Type:  ref.Type,
					Value: ref.Value,
				}
				inRefs = append(inRefs, string(ref.DocID(ctx)))
			}

			// Subconsulta simple para obtener IDs únicos
			ds = ds.Where(goqu.I(ni + ".document_id").In(
				goqu.From(goqu.T("node_references").As(nrf)).
					Select(goqu.I(nrf + ".node_doc")).
					Where(goqu.I(nrf + ".document_id").In(inRefs)),
			))
		}

		// Add node references using WITH clause if either requested or filtered
		if projection.NodeReferences().Has(filters.RequestedFields) ||
			len(filters.References) > 0 {
			ds = ds.With("node_refs", goqu.From(goqu.T("node_references").As("nr")).
				Select(
					goqu.I("nr.node_doc"),
					goqu.L("jsonb_agg(jsonb_build_object('type', type, 'value', value))").As("references"),
				).
				GroupBy(goqu.I("nr.node_doc")),
			).
				InnerJoin(
					goqu.T("node_refs").As("nr"),
					goqu.On(goqu.I("nr.node_doc").Eq(goqu.I(ni+".document_id"))),
				)

			// Only append the references field if it was requested
			if projection.NodeReferences().Has(filters.RequestedFields) {
				ds = ds.SelectAppend(goqu.Cast(goqu.I("nr.references"), "jsonb").As("node_references"))
			}

			// Add filter conditions for references if provided
			if len(filters.References) > 0 {
				const nrf = "nrf" // alias exclusivo para evitar colisión con la CTE `node_refs`

				inRefs := []string{}
				for _, ref := range filters.References {
					ref := domain.Reference{
						Type:  ref.Type,
						Value: ref.Value,
					}
					inRefs = append(inRefs, string(ref.DocID(ctx)))
				}

				// Subconsulta simple para obtener IDs únicos
				ds = ds.Where(goqu.I(ni + ".document_id").In(
					goqu.From(goqu.T("node_references").As(nrf)).
						Select(goqu.I(nrf + ".node_doc")).
						Where(goqu.I(nrf + ".document_id").In(inRefs)),
				))
			}
		}

		// Join con address_infos si se requiere algún campo de addressInfo
		if projection.AddressInfo().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("address_infos").As(ai),
				goqu.On(goqu.I(ai+".document_id").Eq(goqu.I(ni+".address_info_doc"))),
			)
		}

		// Campos de node_infos
		if projection.NodeInfo().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(
				goqu.I(ni+".id").As("id"),
				goqu.I(ni+".name").As("node_name"),
				goqu.I(ni+".node_references").As("node_references"),
				goqu.I(ni+".document_id").As("document_id"),
			)
		} else {
			// Si no se solicitan campos de node_info, al menos necesitamos el document_id para los joins
			ds = ds.SelectAppend(goqu.I(ni + ".document_id").As("document_id"))
		}

		if projection.NodeType().Has(filters.RequestedFields) {
			// Only add the join if it hasn't been added by the filter
			if filters.NodeType == nil {
				ds = ds.InnerJoin(
					goqu.T("node_types").As(nt),
					goqu.On(goqu.I(nt+".document_id").Eq(goqu.I(ni+".node_type_doc"))),
				)
			}
			ds = ds.SelectAppend(goqu.I(nt + ".value").As("node_type"))
		}

		// Campos de address_infos - solo seleccionar los campos específicos solicitados
		if projection.AddressLine1().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ai + ".address_line1").As("address_line1"))
		}

		if projection.AddressLine2().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ai + ".address_line2").As("address_line2"))
		}

		if projection.ZipCode().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ai + ".zip_code").As("zip_code"))
		}

		// Campos de coordenadas - solo seleccionar los campos específicos solicitados
		if projection.CoordinatesLatitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ai + ".latitude").As("coordinates_latitude"))
		}

		if projection.CoordinatesLongitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ai + ".longitude").As("coordinates_longitude"))
		}

		if projection.CoordinatesSource().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ai + ".coordinate_source").As("coordinates_source"))
		}

		if projection.CoordinatesConfidenceLevel().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ai + ".coordinate_confidence").As("coordinates_confidence_level"))
		}

		if projection.CoordinatesConfidenceMessage().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ai + ".coordinate_message").As("coordinates_confidence_message"))
		}

		if projection.CoordinatesConfidenceReason().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ai + ".coordinate_reason").As("coordinates_confidence_reason"))
		}

		if projection.TimeZone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ai + ".time_zone").As("time_zone"))
		}

		// Campos de contacto
		if projection.Contact().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("contacts").As(c),
				goqu.On(goqu.I(c+".document_id").Eq(goqu.I(ni+".contact_doc"))),
			)
		}

		if projection.ContactEmail().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(c + ".email").As("contact_email"))
		}

		if projection.ContactFullName().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(c + ".full_name").As("contact_full_name"))
		}

		if projection.ContactNationalID().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(c + ".national_id").As("contact_national_id"))
		}

		if projection.ContactPhone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(c + ".phone").As("contact_phone"))
		}

		if projection.ContactDocuments().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(c + ".documents").As("contact_documents"))
		}

		if projection.ContactAdditionalContactMethods().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(c + ".additional_contact_methods").As("additional_contact_methods"))
		}

		// Agregar ordenamiento por id
		if filters.Pagination.IsForward() {
			ds = ds.Order(goqu.I(ni + ".id").Asc())
		}

		if filters.Pagination.IsBackward() {
			ds = ds.Order(goqu.I(ni + ".id").Desc())
		}

		// Apply pagination filters
		if filters.Pagination.IsForward() {
			afterID, err := filters.Pagination.AfterID()
			if err != nil {
				return nil, false, err
			}

			if afterID != nil {
				ds = ds.Where(goqu.I(ni + ".id").Gt(*afterID))
			}

			limit := *filters.Pagination.First + 1
			ds = ds.Limit(uint(limit))
		}

		if filters.Pagination.IsBackward() {
			beforeID, err := filters.Pagination.BeforeID()
			if err != nil {
				return nil, false, err
			}

			if beforeID != nil {
				ds = ds.Where(goqu.I(ni + ".id").Lt(*beforeID))
			}

			limit := *filters.Pagination.Last + 1
			ds = ds.Limit(uint(limit))
		}

		sql, args, err := ds.Prepared(true).ToSQL()
		if err != nil {
			return nil, false, err
		}

		err = conn.WithContext(ctx).Raw(sql, args...).Scan(&results).Error
		if err != nil {
			return nil, false, err
		}

		// Check if there are more results than requested
		if filters.Pagination.IsForward() && len(results) > *filters.Pagination.First {
			results = results[:*filters.Pagination.First]
			hasMoreResults = true
		} else if filters.Pagination.IsBackward() && len(results) > *filters.Pagination.Last {
			results = results[:*filters.Pagination.Last]
			hasMoreResults = true
		}

		// Reverse results for backward pagination
		if filters.Pagination.IsBackward() {
			results = results.Reversed()
		}

		return results, hasMoreResults, nil
	}
}
