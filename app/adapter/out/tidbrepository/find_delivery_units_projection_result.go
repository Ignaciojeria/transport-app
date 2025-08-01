package tidbrepository

import (
	"context"
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
		dus  = "dus"  // delivery_units_skills
		sk   = "sk"   // skills
	)

	return func(ctx context.Context, filters domain.DeliveryUnitsFilter) (projectionresult.DeliveryUnitsProjectionResults, bool, error) {
		var results projectionresult.DeliveryUnitsProjectionResults
		hasMoreResults := false

		duh := "duh"

		var baseQuery *goqu.SelectDataset

		if filters.OnlyLatestStatus {
			// Subquery que agrupa y obtiene el último id por combinación
			latestIDsSubquery := goqu.From(goqu.T("delivery_units_status_histories").As(duh)).
				Select(
					goqu.MAX(goqu.I(duh+".id")).As("id"),
				).
				Where(goqu.Ex{
					duh + ".tenant_id": sharedcontext.TenantIDFromContext(ctx),
				}).
				GroupBy(
					goqu.I(duh+".delivery_unit_doc"),
					goqu.I(duh+".order_doc"),
				).As("latest_ids")

			baseQuery = goqu.From("delivery_units_status_histories").
				Join(latestIDsSubquery, goqu.On(
					goqu.I("delivery_units_status_histories.id").Eq(goqu.I("latest_ids.id")),
				)).
				Select(goqu.I("delivery_units_status_histories.id"))
		}

		if !filters.OnlyLatestStatus {
			// Opción sin filtrar por último estado
			baseQuery = goqu.From("delivery_units_status_histories").
				Where(goqu.Ex{
					"tenant_id": sharedcontext.TenantIDFromContext(ctx),
				}).
				Select(goqu.I("id"))
		}

		ds := goqu.From(baseQuery.As("base")).
			Select(
				goqu.I("base.id").As("id"),
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
			(filters.DeliveryUnit != nil && (len(filters.DeliveryUnit.Lpns) > 0 ||
				len(filters.DeliveryUnit.Labels) > 0 ||
				len(filters.DeliveryUnit.SizeCategories) > 0)) {
			ds = ds.InnerJoin(
				goqu.T("delivery_units").As(du),
				goqu.On(goqu.I(du+".document_id").Eq(goqu.I(duh+".delivery_unit_doc"))),
			)
		}

		// Agregar filtro por reference_id si existe
		if filters.Order != nil && len(filters.Order.ReferenceIds) > 0 {
			ds = ds.Where(goqu.I(o + ".reference_id").In(filters.Order.ReferenceIds))
		}

		// Agregar filtro por LPNs si existen
		if filters.DeliveryUnit != nil && len(filters.DeliveryUnit.Lpns) > 0 {
			ds = ds.Where(goqu.I(du + ".lpn").In(filters.DeliveryUnit.Lpns))
		}

		// Agregar filtro por SizeCategories si existen
		if filters.DeliveryUnit != nil && len(filters.DeliveryUnit.SizeCategories) > 0 {
			sizeCategoriesDocs := []string{}
			for _, sizeCategory := range filters.DeliveryUnit.SizeCategories {
				sizeCategoriesDocs = append(sizeCategoriesDocs, string(domain.SizeCategory{Code: sizeCategory}.DocumentID(ctx)))
			}
			ds = ds.Where(goqu.I(du + ".size_category_doc").In(sizeCategoriesDocs))
		}

		// Agregar filtro por originNodeReferences si existen
		if filters.Origin != nil && len(filters.Origin.NodeReferences) > 0 {
			nodeDocs := []string{}
			for _, ref := range filters.Origin.NodeReferences {
				ni := domain.NodeInfo{
					ReferenceID: domain.ReferenceID(ref),
				}
				nodeDocs = append(nodeDocs, string(ni.DocID(ctx)))
			}
			ds = ds.Where(goqu.I(o + ".origin_node_info_doc").In(nodeDocs))
		}

		// Join address_infos si se requiere algún campo de addressInfo
		if projection.DestinationAddressInfo().Has(filters.RequestedFields) ||
			(filters.Destination != nil && filters.Destination.CoordinatesConfidence != nil) {
			ds = ds.InnerJoin(
				goqu.T("address_infos").As(dadi),
				goqu.On(goqu.I(dadi+".document_id").Eq(goqu.I(o+".destination_address_info_doc"))),
			)
		}

		// Agregar filtro por nivel de confianza de coordenadas si existe
		if filters.Destination != nil && filters.Destination.CoordinatesConfidence != nil {
			// Aplicar filtros de nivel de confianza
			if filters.Destination.CoordinatesConfidence.Min != nil {
				ds = ds.Where(goqu.I(dadi + ".coordinate_confidence").Gte(*filters.Destination.CoordinatesConfidence.Min))
			}
			if filters.Destination.CoordinatesConfidence.Max != nil {
				ds = ds.Where(goqu.I(dadi + ".coordinate_confidence").Lte(*filters.Destination.CoordinatesConfidence.Max))
			}
		}

		// Agregar filtro por rango de fecha prometida si existe
		if filters.PromisedDate != nil && filters.PromisedDate.DateRange != nil {
			if filters.PromisedDate.DateRange.StartDate != nil {
				ds = ds.Where(goqu.I(o + ".promised_date_range_start").Gte(*filters.PromisedDate.DateRange.StartDate))
			}
			if filters.PromisedDate.DateRange.EndDate != nil {
				ds = ds.Where(goqu.I(o + ".promised_date_range_end").Lte(*filters.PromisedDate.DateRange.EndDate))
			}
		}

		// Agregar filtro por CollectAvailabilityDates si existen
		if filters.CollectAvailability != nil && len(filters.CollectAvailability.Dates) > 0 {
			ds = ds.Where(goqu.I(o + ".collect_availability_date").In(filters.CollectAvailability.Dates))
		}

		// Agregar ordenamiento por reference_id
		ds = ds.Order(goqu.I(o + ".reference_id").Asc())

		if filters.Pagination.IsForward() {
			ds = ds.Order(goqu.I(duh + ".id").Asc())
		}

		if filters.Pagination.IsBackward() {
			ds = ds.Order(goqu.I(duh + ".id").Desc())
		}

		if filters.Pagination.IsForward() {
			afterID, err := filters.Pagination.AfterID()
			if err != nil {
				return nil, false, err
			}

			if afterID != nil {
				ds = ds.Where(goqu.I(duh + ".id").Gt(*afterID))
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
				ds = ds.Where(goqu.I(duh + ".id").Lt(*beforeID))
			}

			limit := *filters.Pagination.Last + 1
			ds = ds.Limit(uint(limit))
		}

		if projection.DeliveryUnitSkills().Has(filters.RequestedFields) {
			ds = ds.With("delivery_unit_skills", goqu.From(goqu.T("delivery_units_skills").As(dus)).
				Select(
					goqu.I(dus+".delivery_unit_doc"),
					goqu.L("jsonb_agg(skill)").As("skills"),
				).
				GroupBy(goqu.I(dus+".delivery_unit_doc")),
			).
				LeftJoin(
					goqu.T("delivery_unit_skills").As(dus),
					goqu.On(goqu.I(dus+".delivery_unit_doc").Eq(goqu.I(duh+".delivery_unit_doc"))),
				)

			ds = ds.SelectAppend(goqu.Cast(goqu.I(dus+".skills"), "jsonb").As("delivery_unit_skills"))
		}

		// Add order references using WITH clause if either requested or filtered
		if projection.References().Has(filters.RequestedFields) ||
			(filters.Order != nil && len(filters.Order.References) > 0) {
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
			if filters.Order != nil && len(filters.Order.References) > 0 {
				const orf = "orf" // alias exclusivo para evitar colisión con la CTE `order_refs`

				inRefs := []string{}
				for _, ref := range filters.Order.References {
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
		if projection.DeliveryUnitLabels().Has(filters.RequestedFields) ||
			(filters.DeliveryUnit != nil && len(filters.DeliveryUnit.Labels) > 0) {
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
			if filters.DeliveryUnit != nil && len(filters.DeliveryUnit.Labels) > 0 {
				const dulf = "dulf" // alias exclusivo para evitar colisión con la CTE `delivery_unit_labels`

				docIds := []string{}
				for _, label := range filters.DeliveryUnit.Labels {
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

		if projection.ManualChangePerformedBy().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(duh + ".manual_change_performed_by").As("manual_change_performed_by"))
		}

		if projection.ManualChangeReason().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(duh + ".manual_change_reason").As("manual_change_reason"))
		}

		// Campos de delivery failure
		if projection.DeliveryFailureReferenceID().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(duh + ".non_delivery_reason_reference_id").As("non_delivery_reason_reference_id"))
		}

		if projection.DeliveryFailureReason().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(duh + ".non_delivery_reason").As("non_delivery_reason"))
		}

		if projection.DeliveryFailureDetail().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(duh + ".non_delivery_detail").As("non_delivery_detail"))
		}

		if projection.DeliveryEvidencePhotos().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(duh + ".evidence_photos").As("evidence_photos"))
		}

		if projection.DeliveryUnitLPN().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".lpn").As("lpn"))
		}

		if projection.DeliveryUnitSizeCategory().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("size_categories").As("sc"),
				goqu.On(goqu.I("sc.document_id").Eq(goqu.I(du+".size_category_doc"))),
			)
			ds = ds.SelectAppend(goqu.I("sc.code").As("size_category"))
		}

		if projection.DeliveryUnitVolume().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".volume").As("volume"))
		}
		if projection.DeliveryUnitWeight().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".weight").As("weight"))
		}
		if projection.DeliveryUnitInsurance().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(du + ".insurance").As("insurance"))
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
			ds = ds.SelectAppend(goqu.L("to_char(" + o + ".collect_availability_time_range_start, 'HH24:MI')").As("order_collect_availability_date_start_time"))
		}
		if projection.CollectAvailabilityDateEndTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.L("to_char(" + o + ".collect_availability_time_range_end, 'HH24:MI')").As("order_collect_availability_date_end_time"))
		}

		if projection.PromisedDateDateRangeStartDate().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".promised_date_range_start").As("order_promised_date_start_date"))
		}
		if projection.PromisedDateDateRangeEndDate().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(o + ".promised_date_range_end").As("order_promised_date_end_date"))
		}
		if projection.PromisedDateTimeRangeStartTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.L("to_char(" + o + ".promised_time_range_start, 'HH24:MI')").As("order_promised_date_start_time"))
		}
		if projection.PromisedDateTimeRangeEndTime().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.L("to_char(" + o + ".promised_time_range_end, 'HH24:MI')").As("order_promised_date_end_time"))
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

		if projection.DestinationPoliticalArea().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("political_areas").As("dpa"),
				goqu.On(goqu.I("dpa.document_id").Eq(goqu.I(dadi+".political_area_doc"))),
			)
			ds = ds.SelectAppend(goqu.I("dpa.code").As("destination_political_area_code"))
		}

		if projection.DestinationPoliticalAreaConfidenceLevel().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".political_area_confidence").As("destination_political_area_confidence_level"))
		}

		if projection.DestinationPoliticalAreaConfidenceMessage().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".political_area_message").As("destination_political_area_confidence_message"))
		}

		if projection.DestinationPoliticalAreaConfidenceReason().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(dadi + ".political_area_reason").As("destination_political_area_confidence_reason"))
		}

		// Join con admin area levels
		if projection.DestinationAdminAreaLevel1().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dpa.admin_area_level1").As("destination_admin_area_level1"))
		}

		// Join con admin area levels
		if projection.DestinationAdminAreaLevel2().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dpa.admin_area_level2").As("destination_admin_area_level2"))
		}

		// Join con admin area levels
		if projection.DestinationAdminAreaLevel3().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dpa.admin_area_level3").As("destination_admin_area_level3"))
		}

		// Join con admin area levels
		if projection.DestinationAdminAreaLevel4().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("dpa.admin_area_level4").As("destination_admin_area_level4"))
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
			ds = ds.SelectAppend(goqu.I("dpa.time_zone").As("destination_time_zone"))
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

		if projection.OriginPoliticalArea().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("opa.code").As("origin_political_area_code"))
		}

		if projection.OriginPoliticalAreaConfidenceLevel().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".political_area_confidence").As("origin_political_area_confidence_level"))
		}

		if projection.OriginPoliticalAreaConfidenceMessage().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".political_area_message").As("origin_political_area_confidence_message"))
		}

		if projection.OriginPoliticalAreaConfidenceReason().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I(oadi + ".political_area_reason").As("origin_political_area_confidence_reason"))
		}

		// Join con admin area levels para origen
		if projection.OriginAdminAreaLevel1().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("opa.admin_area_level1").As("origin_admin_area_level1"))
		}

		// Join con admin area levels para origen
		if projection.OriginAdminAreaLevel2().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("opa.admin_area_level2").As("origin_admin_area_level2"))
		}

		// Join con admin area levels para origen
		if projection.OriginAdminAreaLevel3().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("opa.admin_area_level3").As("origin_admin_area_level3"))
		}

		// Join con admin area levels para origen
		if projection.OriginAdminAreaLevel4().Has(filters.RequestedFields) {
			ds = ds.SelectAppend(goqu.I("opa.admin_area_level4").As("origin_admin_area_level4"))
		}

		// Agregar campos de political area para origen
		if projection.OriginPoliticalArea().Has(filters.RequestedFields) {
			ds = ds.InnerJoin(
				goqu.T("political_areas").As("opa"),
				goqu.On(goqu.I("opa.document_id").Eq(goqu.I(oadi+".political_area_doc"))),
			)
			ds = ds.SelectAppend(goqu.I("opa.code").As("origin_political_area_code"))
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
			ds = ds.SelectAppend(goqu.I("opa.time_zone").As("origin_time_zone"))
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
