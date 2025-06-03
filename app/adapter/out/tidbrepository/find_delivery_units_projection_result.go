package tidbrepository

import (
	"context"
	"time"
	"transport-app/app/adapter/out/tidbrepository/projectionresult"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/projection/deliveryunits"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/doug-martin/goqu/v9"
)

type FindDeliveryUnitsProjectionResult func(
	ctx context.Context,
	filters domain.DeliveryUnitsFilter) (projectionresult.DeliveryUnitsProjectionResults, bool, error)

func init() {
	ioc.Registry(
		NewFindDeliveryUnitsProjectionResult,
		database.NewConnectionFactory,
		deliveryunits.NewProjection,
	)
}

func NewFindDeliveryUnitsProjectionResult(
	conn database.ConnectionFactory,
	projection deliveryunits.Projection) FindDeliveryUnitsProjectionResult {
	const (
		duh  = "duh"  // delivery_units_status_histories
		o    = "o"    // orders
		or   = "or"   // order_references
		dadi = "dadi" // destination_address_infos
		oadi = "oadi" // origin_address_infos
		dd   = "dd"   // destination_districts
		dp   = "dp"   // destination_provinces
		dst  = "dst"  // destination_states
		od   = "od"   // origin_districts
		op   = "op"   // origin_provinces
		ost  = "ost"  // origin_states
		du   = "du"   // delivery_units
		dul  = "dul"  // delivery_unit_labels
		oh   = "oh"   // order_headers
		ot   = "ot"   // order_types
		s    = "s"    // status
	)

	return func(ctx context.Context, filters domain.DeliveryUnitsFilter) (projectionresult.DeliveryUnitsProjectionResults, bool, error) {
		var results projectionresult.DeliveryUnitsProjectionResults
		hasMoreResults := false

		// Dataset base
		baseQuery := goqu.From(goqu.T("delivery_units_status_histories").As(duh)).
			Select(
				goqu.I(duh+".id").As("id"),
				goqu.I(duh+".updated_at").As("updated_at"),
			).
			Where(goqu.Ex{
				duh + ".tenant_id": sharedcontext.TenantIDFromContext(ctx),
			}).
			Order(goqu.I(duh + ".document_id").Asc())

		ds := goqu.From(baseQuery.As("base")).
			Select(
				goqu.I("base.id").As("id"),
				goqu.I("base.updated_at").As("updated_at"),
			).
			InnerJoin(
				goqu.T("delivery_units_status_histories").As(duh),
				goqu.On(goqu.I(duh+".id").Eq(goqu.I("base.id"))),
			).
			InnerJoin(
				goqu.T("orders").As(o),
				goqu.On(goqu.I(o+".document_id").Eq(goqu.I(duh+".order_doc"))),
			)

		// Agregar join con delivery_units si se solicita cualquier campo relacionado
		if projection.DeliveryUnit().Has(filters.RequestedFields) ||
			len(filters.Lpns) > 0 || len(filters.Labels) > 0 {
			ds = ds.InnerJoin(
				goqu.T("delivery_units").As(du),
				goqu.On(goqu.I(du+".document_id").Eq(goqu.I(duh+".delivery_unit_doc"))),
			)
		}

		// Agregar filtro por reference_id si existe
		if len(filters.ReferenceIds) > 0 {
			ds = ds.Where(goqu.I(o + ".reference_id").In(filters.ReferenceIds))
		}

		// Agregar filtro por LPNs si existen
		if len(filters.Lpns) > 0 {
			ds = ds.Where(goqu.I(du + ".lpn").In(filters.Lpns))
		}

		// Agregar filtro por originNodeReferences si existen
		if len(filters.OriginNodeReferences) > 0 {
			nodeDocs := []string{}
			for _, ref := range filters.OriginNodeReferences {
				ni := domain.NodeInfo{
					ReferenceID: domain.ReferenceID(ref),
				}
				nodeDocs = append(nodeDocs, string(ni.DocID(ctx)))
			}
			ds = ds.Where(goqu.I(o + ".origin_node_info_doc").In(nodeDocs))
		}

		// Join address_infos si se requiere algún campo de addressInfo
		if projection.DestinationAddressInfo().Has(filters.RequestedFields) || filters.CoordinatesConfidenceLevel != nil {
			ds = ds.InnerJoin(
				goqu.T("address_infos").As(dadi),
				goqu.On(goqu.I(dadi+".document_id").Eq(goqu.I(o+".destination_address_info_doc"))),
			)
		}

		// Agregar filtro por nivel de confianza de coordenadas si existe
		if filters.CoordinatesConfidenceLevel != nil {
			// Aplicar filtros de nivel de confianza
			if filters.CoordinatesConfidenceLevel.Min != nil {
				ds = ds.Where(goqu.I(dadi + ".coordinate_confidence").Gte(*filters.CoordinatesConfidenceLevel.Min))
			}
			if filters.CoordinatesConfidenceLevel.Max != nil {
				ds = ds.Where(goqu.I(dadi + ".coordinate_confidence").Lte(*filters.CoordinatesConfidenceLevel.Max))
			}
		}

		// Agregar ordenamiento por reference_id
		if filters.Pagination.IsForward() {
			if filters.Pagination.HasAfter() {
				// Si hay cursor, ordenar por updated_at e id
				ds = ds.Order(goqu.I(duh+".updated_at").Asc(), goqu.I(duh+".id").Asc())
			} else {
				// Si no hay cursor, ordenar por reference_id
				ds = ds.Order(goqu.I(o + ".reference_id").Asc())
			}
		} else if filters.Pagination.IsBackward() {
			if filters.Pagination.HasBefore() {
				// Si hay cursor, ordenar por updated_at e id
				ds = ds.Order(goqu.I(duh+".updated_at").Desc(), goqu.I(duh+".id").Desc())
			} else {
				// Si no hay cursor, ordenar por reference_id
				ds = ds.Order(goqu.I(o + ".reference_id").Desc())
			}
		} else {
			// Si no hay paginación, ordenar por reference_id ascendente
			ds = ds.Order(goqu.I(o + ".reference_id").Asc())
		}

		if filters.Pagination.IsForward() {
			ds = ds.Order(goqu.I(duh+".updated_at").Asc(), goqu.I(duh+".id").Asc())

			if afterID, err := filters.Pagination.AfterID(); err != nil {
				return nil, false, err
			} else if afterID != nil {
				// Obtener el updated_at del registro con ese ID
				var updatedAt time.Time
				err := conn.WithContext(ctx).
					Table("delivery_units_status_histories").
					Select("updated_at").
					Where("id = ?", *afterID).
					Scan(&updatedAt).Error
				if err != nil {
					return nil, false, err
				}

				ds = ds.Where(
					goqu.Or(
						goqu.I(duh+".updated_at").Gt(updatedAt),
						goqu.And(
							goqu.I(duh+".updated_at").Eq(updatedAt),
							goqu.I(duh+".id").Gt(*afterID),
						),
					),
				)
			}

			limit := *filters.Pagination.First + 1 // pedir uno extra para saber si hay más
			ds = ds.Limit(uint(limit))
		}

		if filters.Pagination.IsBackward() {
			ds = ds.Order(goqu.I(duh+".updated_at").Desc(), goqu.I(duh+".id").Desc())

			if beforeID, err := filters.Pagination.BeforeID(); err != nil {
				return nil, false, err
			} else if beforeID != nil {
				// Obtener el updated_at del registro con ese ID
				var updatedAt time.Time
				err := conn.WithContext(ctx).
					Table("delivery_units_status_histories").
					Select("updated_at").
					Where("id = ?", *beforeID).
					Scan(&updatedAt).Error
				if err != nil {
					return nil, false, err
				}

				ds = ds.Where(
					goqu.Or(
						goqu.I(duh+".updated_at").Lt(updatedAt),
						goqu.And(
							goqu.I(duh+".updated_at").Eq(updatedAt),
							goqu.I(duh+".id").Lt(*beforeID),
						),
					),
				)
			}

			limit := *filters.Pagination.Last + 1
			ds = ds.Limit(uint(limit))
		}

		// Add order references using WITH clause if either requested or filtered
		if projection.References().Has(filters.RequestedFields) || len(filters.References) > 0 {
			ds = ds.With("order_refs", goqu.From(goqu.T("order_references").As(or)).
				Select(
					goqu.I(or+".order_doc"),
					goqu.L("jsonb_agg(jsonb_build_object('type', type, 'value', value))").As("references"),
				).
				GroupBy(goqu.I(or+".order_doc")),
			).
				InnerJoin(
					goqu.T("order_refs").As(or),
					goqu.On(goqu.I(or+".order_doc").Eq(goqu.I(o+".document_id"))),
				)

			// Only append the references field if it was requested
			if projection.References().Has(filters.RequestedFields) {
				ds = ds.SelectAppend(goqu.Cast(goqu.I(or+".references"), "jsonb").As("order_references"))
			}

			// Add filter conditions for references if provided
			if len(filters.References) > 0 {
				const orf = "orf" // alias exclusivo para evitar colisión con la CTE `order_refs`

				inRefs := []string{}
				for _, ref := range filters.References {
					ref := domain.Reference{
						Type:  ref.Type,
						Value: ref.Value,
					}
					inRefs = append(inRefs, string(ref.DocID(ctx)))
				}

				// Subconsulta simple para obtener IDs únicos
				ds = ds.Where(goqu.I(duh + ".order_doc").In(
					goqu.From(goqu.T("order_references").As(orf)).
						Select(goqu.I(orf + ".order_doc")).
						Where(goqu.I(orf + ".document_id").In(inRefs)),
				))
			}

		}

		// Add delivery unit labels if requested or filtered
		if projection.DeliveryUnitLabels().Has(filters.RequestedFields) || len(filters.Labels) > 0 {

			ds = ds.With("delivery_unit_labels", goqu.From(goqu.T("delivery_units_labels").As(dul)).
				Select(
					goqu.I(dul+".delivery_unit_doc"),
					goqu.L("jsonb_agg(jsonb_build_object('type', type, 'value', value))").As("delivery_unit_labels"),
				).
				GroupBy(goqu.I(dul+".delivery_unit_doc")),
			).
				InnerJoin(
					goqu.T("delivery_unit_labels").As(dul),
					goqu.On(goqu.I(dul+".delivery_unit_doc").Eq(goqu.I(duh+".delivery_unit_doc"))),
				)

			// Only append the labels field if it was requested
			if projection.DeliveryUnitLabels().Has(filters.RequestedFields) {
				ds = ds.SelectAppend(goqu.Cast(goqu.I(dul+".delivery_unit_labels"), "jsonb").As("delivery_unit_labels"))
			}

			// Add filter conditions for labels if provided
			if len(filters.Labels) > 0 {
				const dulf = "dulf" // alias exclusivo para evitar colisión con la CTE `delivery_unit_labels`

				docIds := []string{}
				for _, label := range filters.Labels {
					docIds = append(docIds, string(domain.Reference(label).DocID(ctx)))
				}

				// Subconsulta simple para obtener IDs únicos
				ds = ds.Where(goqu.I(duh + ".delivery_unit_doc").In(
					goqu.From(goqu.T("delivery_units_labels").As(dulf)).
						Select(goqu.I(dulf + ".delivery_unit_doc")).
						Where(goqu.I(dulf + ".document_id").In(docIds)),
				))
			}
		}

		// Campos de delivery_units_status_histories
		if projection.Channel().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(duh + ".channel").As("channel"))
		}

		if projection.Status().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("statuses").As(s),
				goqu.On(goqu.I(s+".document_id").Eq(goqu.I(duh+".delivery_unit_status_doc"))),
			).
				SelectAppend(goqu.I(s + ".status").As("status"))
		}

		if projection.DeliveryUnitLPN().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".lpn").As("lpn"))
		}

		if projection.DeliveryUnitDimensions().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".json_dimensions").As("json_dimensions"))
		}
		if projection.DeliveryUnitWeight().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".json_weight").As("json_weight"))
		}
		if projection.DeliveryUnitInsurance().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".json_insurance").As("json_insurance"))
		}
		if projection.DeliveryUnitItems().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".json_items").As("json_items"))
		}

		// Campos de orders
		if projection.ReferenceID().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".reference_id").As("order_reference_id"))
		}

		if projection.GroupByType().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".group_by_type").As("order_group_by_type"))
		}

		if projection.GroupByValue().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".group_by_value").As("order_group_by_value"))
		}

		if projection.CollectAvailabilityDate().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".collect_availability_date").As("order_collect_availability_date"))
		}

		if projection.DeliveryInstructions().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".delivery_instructions").As("order_delivery_instructions"))
		}

		if projection.CollectAvailabilityDateStartTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".collect_availability_time_range_start").As("order_collect_availability_date_start_time"))
		}
		if projection.CollectAvailabilityDateEndTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".collect_availability_time_range_end").As("order_collect_availability_date_end_time"))
		}

		if projection.PromisedDateDateRangeStartDate().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".promised_date_range_start").As("order_promised_date_start_date"))
		}
		if projection.PromisedDateDateRangeEndDate().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".promised_date_range_end").As("order_promised_date_end_date"))
		}
		if projection.PromisedDateTimeRangeStartTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".promised_time_range_start").As("order_promised_date_start_time"))
		}
		if projection.PromisedDateTimeRangeEndTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".promised_time_range_end").As("order_promised_date_end_time"))
		}
		if projection.PromisedDateServiceCategory().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".service_category").As("order_promised_date_service_category"))
		}

		// Campos de orderType
		if projection.OrderType().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("order_types").As(ot),
				goqu.On(goqu.I(ot+".document_id").Eq(goqu.I(o+".order_type_doc"))),
			)
		}

		if projection.OrderTypeType().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ot + ".type").As("order_type"))
		}

		if projection.OrderTypeDescription().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(ot + ".description").As("order_type_description"))
		}

		// Campos de address_infos
		if projection.DestinationAddressLine2().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".address_line2").As("destination_address_line2"))
		}

		if projection.Commerce().Has(filters.RequestedFields) || projection.Consumer().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("order_headers").As(oh),
				goqu.On(goqu.I(oh+".document_id").Eq(goqu.I(o+".order_headers_doc"))),
			)
		}

		if projection.Commerce().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oh + ".commerce").As("commerce"))
		}

		if projection.Consumer().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oh + ".consumer").As("consumer"))
		}

		// Campos de address_infos
		if projection.DestinationAddressLine1().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".address_line1").As("destination_address_line1"))
		}

		// Join con districts
		if projection.DestinationDistrict().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("districts").As(dd),
				goqu.On(goqu.I(dd+".document_id").Eq(goqu.I(dadi+".district_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(dd + ".name").As("destination_district"))
		}

		// Join con provinces
		if projection.DestinationProvince().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("provinces").As(dp),
				goqu.On(goqu.I(dp+".document_id").Eq(goqu.I(dadi+".province_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(dp + ".name").As("destination_province"))
		}

		// Join con states
		if projection.DestinationState().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("states").As(dst),
				goqu.On(goqu.I(dst+".document_id").Eq(goqu.I(dadi+".state_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(dst + ".name").As("destination_state"))
		}

		if projection.DestinationZipCode().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".zip_code").As("destination_zip_code"))
		}

		// Campos de coordenadas
		if projection.DestinationCoordinatesLatitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".latitude").As("destination_coordinates_latitude"))
		}

		if projection.DestinationCoordinatesLongitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".longitude").As("destination_coordinates_longitude"))
		}

		if projection.DestinationCoordinatesSource().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".coordinate_source").As("destination_coordinates_source"))
		}

		if projection.DestinationCoordinatesConfidenceLevel().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".coordinate_confidence").As("destination_coordinates_confidence_level"))
		}

		if projection.DestinationCoordinatesConfidenceMessage().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".coordinate_message").As("destination_coordinates_confidence_message"))
		}

		if projection.DestinationCoordinatesConfidenceReason().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".coordinate_reason").As("destination_coordinates_confidence_reason"))
		}

		if projection.DestinationTimeZone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".time_zone").As("destination_time_zone"))
		}

		// Campos de contacto del destino
		if projection.DestinationContact().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("contacts").As("dc"),
				goqu.On(goqu.I("dc.document_id").Eq(goqu.I(o+".destination_contact_doc"))),
			)
		}

		if projection.DestinationContactEmail().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.email").As("destination_contact_email"))
		}

		if projection.DestinationContactFullName().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.full_name").As("destination_contact_full_name"))
		}

		if projection.DestinationContactNationalID().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.national_id").As("destination_contact_national_id"))
		}

		if projection.DestinationContactPhone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.phone").As("destination_contact_phone"))
		}

		if projection.DestinationAdditionalContactMethods().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.additional_contact_methods").As("destination_additional_contact_methods"))
		}

		if projection.DestinationContactDocuments().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dc.documents").As("destination_contact_documents"))
		}

		// Join address_infos para origen si se requiere algún campo de addressInfo
		if projection.OriginAddressInfo().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("address_infos").As(oadi),
				goqu.On(goqu.I(oadi+".document_id").Eq(goqu.I(o+".origin_address_info_doc"))),
			)
		}

		// Campos de address_infos para origen
		if projection.OriginAddressLine1().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".address_line1").As("origin_address_line1"))
		}

		if projection.OriginAddressLine2().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".address_line2").As("origin_address_line2"))
		}

		// Join con districts para origen
		if projection.OriginDistrict().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("districts").As(od),
				goqu.On(goqu.I(od+".document_id").Eq(goqu.I(oadi+".district_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(od + ".name").As("origin_district"))
		}

		// Join con provinces para origen
		if projection.OriginProvince().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("provinces").As(op),
				goqu.On(goqu.I(op+".document_id").Eq(goqu.I(oadi+".province_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(op + ".name").As("origin_province"))
		}

		// Join con states para origen
		if projection.OriginState().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("states").As(ost),
				goqu.On(goqu.I(ost+".document_id").Eq(goqu.I(oadi+".state_doc"))),
			)
			ds = ds.SelectAppend(goqu.I(ost + ".name").As("origin_state"))
		}

		if projection.OriginZipCode().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".zip_code").As("origin_zip_code"))
		}

		// Campos de coordenadas para origen
		if projection.OriginCoordinatesLatitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".latitude").As("origin_coordinates_latitude"))
		}

		if projection.OriginCoordinatesLongitude().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".longitude").As("origin_coordinates_longitude"))
		}

		if projection.OriginCoordinatesSource().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".coordinate_source").As("origin_coordinates_source"))
		}

		if projection.OriginCoordinatesConfidenceLevel().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".coordinate_confidence").As("origin_coordinates_confidence_level"))
		}

		if projection.OriginCoordinatesConfidenceMessage().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".coordinate_message").As("origin_coordinates_confidence_message"))
		}

		if projection.OriginCoordinatesConfidenceReason().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".coordinate_reason").As("origin_coordinates_confidence_reason"))
		}

		if projection.OriginTimeZone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".time_zone").As("origin_time_zone"))
		}

		// Campos de contacto del origen
		if projection.OriginContact().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("contacts").As("oc"),
				goqu.On(goqu.I("oc.document_id").Eq(goqu.I(o+".origin_contact_doc"))),
			)
		}

		if projection.OriginContactEmail().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.email").As("origin_contact_email"))
		}

		if projection.OriginContactFullName().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.full_name").As("origin_contact_full_name"))
		}

		if projection.OriginContactNationalID().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.national_id").As("origin_contact_national_id"))
		}

		if projection.OriginContactPhone().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.phone").As("origin_contact_phone"))
		}

		if projection.OriginContactMethods().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.additional_contact_methods").As("origin_additional_contact_methods"))
		}

		if projection.OriginDocuments().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("oc.documents").As("origin_contact_documents"))
		}

		if projection.ExtraFields().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".extra_fields").As("extra_fields"))
		}

		sql, args, err := ds.Prepared(true).ToSQL()
		if err != nil {
			return nil, false, err
		}

		err = conn.WithContext(ctx).Raw(sql, args...).Scan(&results).Error
		if err != nil {
			return nil, false, err
		}

		// Si hay más resultados que el límite solicitado, eliminar el último resultado
		if filters.Pagination.IsForward() && len(results) > *filters.Pagination.First {
			results = results[:*filters.Pagination.First]
			hasMoreResults = true
		} else if filters.Pagination.IsBackward() && len(results) > *filters.Pagination.Last {
			results = results[:*filters.Pagination.Last]
			hasMoreResults = true
		}

		if filters.Pagination.IsBackward() {
			results = results.Reversed()
		}

		return results, hasMoreResults, nil
	}
}
