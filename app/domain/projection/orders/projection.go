package orders

import "strings"

// Projection representa un campo solicitado por el cliente en la query GraphQL.
type Projection string

// Conjunto de constantes para facilitar el uso y evitar errores de escritura.
const (
	// Campos de Order
	ProjectionReferenceID Projection = "edges.node.referenceID"

	// Campos de CollectAvailabilityDate
	ProjectionCollectAvailabilityDate          Projection = "edges.node.collectAvailabilityDate"
	ProjectionCollectAvailabilityDateDate      Projection = "edges.node.collectAvailabilityDate.date"
	ProjectionCollectAvailabilityDateTimeRange Projection = "edges.node.collectAvailabilityDate.timeRange"
	ProjectionCollectAvailabilityDateStartTime Projection = "edges.node.collectAvailabilityDate.timeRange.startTime"
	ProjectionCollectAvailabilityDateEndTime   Projection = "edges.node.collectAvailabilityDate.timeRange.endTime"

	// Campos de Destination Location
	ProjectionDestination                     Projection = "edges.node.destination"
	ProjectionDestinationDeliveryInstructions Projection = "edges.node.destination.deliveryInstructions"

	// Campos de NodeInfo en Destination
	ProjectionDestinationNodeInfo            Projection = "edges.node.destination.nodeInfo"
	ProjectionDestinationNodeInfoReferenceId Projection = "edges.node.destination.nodeInfo.referenceId"
	ProjectionDestinationNodeInfoName        Projection = "edges.node.destination.nodeInfo.name"

	// Campos de AddressInfo en Destination
	ProjectionDestinationAddressInfo  Projection = "edges.node.destination.addressInfo"
	ProjectionDestinationAddressLine1 Projection = "edges.node.destination.addressInfo.addressLine1"
	ProjectionDestinationAddressLine2 Projection = "edges.node.destination.addressInfo.addressLine2"
	ProjectionDestinationDistrict     Projection = "edges.node.destination.addressInfo.district"
	ProjectionDestinationLatitude     Projection = "edges.node.destination.addressInfo.latitude"
	ProjectionDestinationLongitude    Projection = "edges.node.destination.addressInfo.longitude"
	ProjectionDestinationProvince     Projection = "edges.node.destination.addressInfo.province"
	ProjectionDestinationState        Projection = "edges.node.destination.addressInfo.state"
	ProjectionDestinationTimeZone     Projection = "edges.node.destination.addressInfo.timeZone"
	ProjectionDestinationZipCode      Projection = "edges.node.destination.addressInfo.zipCode"

	// Campos de Contact en Destination
	ProjectionDestinationContact           Projection = "edges.node.destination.addressInfo.contact"
	ProjectionDestinationContactEmail      Projection = "edges.node.destination.addressInfo.contact.email"
	ProjectionDestinationContactFullName   Projection = "edges.node.destination.addressInfo.contact.fullName"
	ProjectionDestinationContactNationalID Projection = "edges.node.destination.addressInfo.contact.nationalID"
	ProjectionDestinationContactPhone      Projection = "edges.node.destination.addressInfo.contact.phone"

	// Campos de ContactMethods en Destination
	ProjectionDestinationContactMethods      Projection = "edges.node.destination.addressInfo.contact.additionalContactMethods"
	ProjectionDestinationContactMethodsType  Projection = "edges.node.destination.addressInfo.contact.additionalContactMethods.type"
	ProjectionDestinationContactMethodsValue Projection = "edges.node.destination.addressInfo.contact.additionalContactMethods.value"

	// Campos de Documents en Destination
	ProjectionDestinationDocuments      Projection = "edges.node.destination.addressInfo.contact.documents"
	ProjectionDestinationDocumentsType  Projection = "edges.node.destination.addressInfo.contact.documents.type"
	ProjectionDestinationDocumentsValue Projection = "edges.node.destination.addressInfo.contact.documents.value"

	// Campos de Origin Location
	ProjectionOrigin                     Projection = "edges.node.origin"
	ProjectionOriginDeliveryInstructions Projection = "edges.node.origin.deliveryInstructions"

	// Campos de NodeInfo en Origin
	ProjectionOriginNodeInfo            Projection = "edges.node.origin.nodeInfo"
	ProjectionOriginNodeInfoReferenceId Projection = "edges.node.origin.nodeInfo.referenceId"
	ProjectionOriginNodeInfoName        Projection = "edges.node.origin.nodeInfo.name"

	// Campos de AddressInfo en Origin
	ProjectionOriginAddressInfo  Projection = "edges.node.origin.addressInfo"
	ProjectionOriginAddressLine1 Projection = "edges.node.origin.addressInfo.addressLine1"
	ProjectionOriginAddressLine2 Projection = "edges.node.origin.addressInfo.addressLine2"
	ProjectionOriginDistrict     Projection = "edges.node.origin.addressInfo.district"
	ProjectionOriginLatitude     Projection = "edges.node.origin.addressInfo.latitude"
	ProjectionOriginLongitude    Projection = "edges.node.origin.addressInfo.longitude"
	ProjectionOriginProvince     Projection = "edges.node.origin.addressInfo.province"
	ProjectionOriginState        Projection = "edges.node.origin.addressInfo.state"
	ProjectionOriginTimeZone     Projection = "edges.node.origin.addressInfo.timeZone"
	ProjectionOriginZipCode      Projection = "edges.node.origin.addressInfo.zipCode"

	// Campos de Contact en Origin
	ProjectionOriginContact           Projection = "edges.node.origin.addressInfo.contact"
	ProjectionOriginContactEmail      Projection = "edges.node.origin.addressInfo.contact.email"
	ProjectionOriginContactFullName   Projection = "edges.node.origin.addressInfo.contact.fullName"
	ProjectionOriginContactNationalID Projection = "edges.node.origin.addressInfo.contact.nationalID"
	ProjectionOriginContactPhone      Projection = "edges.node.origin.addressInfo.contact.phone"

	// Campos de ContactMethods en Origin
	ProjectionOriginContactMethods      Projection = "edges.node.origin.addressInfo.contact.additionalContactMethods"
	ProjectionOriginContactMethodsType  Projection = "edges.node.origin.addressInfo.contact.additionalContactMethods.type"
	ProjectionOriginContactMethodsValue Projection = "edges.node.origin.addressInfo.contact.additionalContactMethods.value"

	// Campos de Documents en Origin
	ProjectionOriginDocuments      Projection = "edges.node.origin.addressInfo.contact.documents"
	ProjectionOriginDocumentsType  Projection = "edges.node.origin.addressInfo.contact.documents.type"
	ProjectionOriginDocumentsValue Projection = "edges.node.origin.addressInfo.contact.documents.value"

	// Campos de OrderType
	ProjectionOrderType            Projection = "edges.node.orderType"
	ProjectionOrderTypeType        Projection = "edges.node.orderType.type"
	ProjectionOrderTypeDescription Projection = "edges.node.orderType.description"

	// Campos de Package
	ProjectionPackages    Projection = "edges.node.packages"
	ProjectionPackagesLPN Projection = "edges.node.packages.lpn"

	// Campos de Weight en Package
	ProjectionPackagesWeight      Projection = "edges.node.packages.weight"
	ProjectionPackagesWeightUnit  Projection = "edges.node.packages.weight.unit"
	ProjectionPackagesWeightValue Projection = "edges.node.packages.weight.value"

	// Campos de Dimensions en Package
	ProjectionPackagesDimensions       Projection = "edges.node.packages.dimensions"
	ProjectionPackagesDimensionsLength Projection = "edges.node.packages.dimensions.length"
	ProjectionPackagesDimensionsHeight Projection = "edges.node.packages.dimensions.height"
	ProjectionPackagesDimensionsWidth  Projection = "edges.node.packages.dimensions.width"
	ProjectionPackagesDimensionsUnit   Projection = "edges.node.packages.dimensions.unit"

	// Campos de Insurance en Package
	ProjectionPackagesInsurance          Projection = "edges.node.packages.insurance"
	ProjectionPackagesInsuranceCurrency  Projection = "edges.node.packages.insurance.currency"
	ProjectionPackagesInsuranceUnitValue Projection = "edges.node.packages.insurance.unitValue"

	// Campos de Label en Package
	ProjectionPackagesLabels      Projection = "edges.node.packages.labels"
	ProjectionPackagesLabelsType  Projection = "edges.node.packages.labels.type"
	ProjectionPackagesLabelsValue Projection = "edges.node.packages.labels.value"

	// Campos de Item en Package
	ProjectionPackagesItems            Projection = "edges.node.packages.items"
	ProjectionPackagesItemsSKU         Projection = "edges.node.packages.items.sku"
	ProjectionPackagesItemsDescription Projection = "edges.node.packages.items.description"

	// Campos de Dimensions en Item
	ProjectionPackagesItemsDimensions       Projection = "edges.node.packages.items.dimensions"
	ProjectionPackagesItemsDimensionsLength Projection = "edges.node.packages.items.dimensions.length"
	ProjectionPackagesItemsDimensionsHeight Projection = "edges.node.packages.items.dimensions.height"
	ProjectionPackagesItemsDimensionsWidth  Projection = "edges.node.packages.items.dimensions.width"
	ProjectionPackagesItemsDimensionsUnit   Projection = "edges.node.packages.items.dimensions.unit"

	// Campos de Insurance en Item
	ProjectionPackagesItemsInsurance          Projection = "edges.node.packages.items.insurance"
	ProjectionPackagesItemsInsuranceCurrency  Projection = "edges.node.packages.items.insurance.currency"
	ProjectionPackagesItemsInsuranceUnitValue Projection = "edges.node.packages.items.insurance.unitValue"

	// Campos de Skill en Item
	ProjectionPackagesItemsSkills            Projection = "edges.node.packages.items.skills"
	ProjectionPackagesItemsSkillsType        Projection = "edges.node.packages.items.skills.type"
	ProjectionPackagesItemsSkillsValue       Projection = "edges.node.packages.items.skills.value"
	ProjectionPackagesItemsSkillsDescription Projection = "edges.node.packages.items.skills.description"

	// Campos de Quantity en Item
	ProjectionPackagesItemsQuantity       Projection = "edges.node.packages.items.quantity"
	ProjectionPackagesItemsQuantityNumber Projection = "edges.node.packages.items.quantity.quantityNumber"
	ProjectionPackagesItemsQuantityUnit   Projection = "edges.node.packages.items.quantity.quantityUnit"

	// Campos de Weight en Item
	ProjectionPackagesItemsWeight      Projection = "edges.node.packages.items.weight"
	ProjectionPackagesItemsWeightUnit  Projection = "edges.node.packages.items.weight.unit"
	ProjectionPackagesItemsWeightValue Projection = "edges.node.packages.items.weight.value"

	// Campos de PromisedDate
	ProjectionPromisedDate                Projection = "edges.node.promisedDate"
	ProjectionPromisedDateServiceCategory Projection = "edges.node.promisedDate.serviceCategory"

	// Campos de TimeRange en PromisedDate
	ProjectionPromisedDateTimeRange          Projection = "edges.node.promisedDate.timeRange"
	ProjectionPromisedDateTimeRangeStartTime Projection = "edges.node.promisedDate.timeRange.startTime"
	ProjectionPromisedDateTimeRangeEndTime   Projection = "edges.node.promisedDate.timeRange.endTime"

	// Campos de DateRange en PromisedDate
	ProjectionPromisedDateDateRange          Projection = "edges.node.promisedDate.dateRange"
	ProjectionPromisedDateDateRangeStartDate Projection = "edges.node.promisedDate.dateRange.startDate"
	ProjectionPromisedDateDateRangeEndDate   Projection = "edges.node.promisedDate.dateRange.endDate"

	// Campos de Reference
	ProjectionReferences      Projection = "edges.node.references"
	ProjectionReferencesType  Projection = "edges.node.references.type"
	ProjectionReferencesValue Projection = "edges.node.references.value"

	// Campos de ExtraFields
	ProjectionExtraFields                Projection = "edges.node.extraFields"
	ProjectionExtraFieldsPoliticalAreaId Projection = "edges.node.extraFields.destinationPoliticalAreaId"

	// Campos de paginación
	ProjectionPageInfo                Projection = "pageInfo"
	ProjectionPageInfoHasNextPage     Projection = "pageInfo.hasNextPage"
	ProjectionPageInfoHasPreviousPage Projection = "pageInfo.hasPreviousPage"
	ProjectionPageInfoStartCursor     Projection = "pageInfo.startCursor"
	ProjectionPageInfoEndCursor       Projection = "pageInfo.endCursor"
	ProjectionTotalCount              Projection = "totalCount"
)

// Grupos de proyecciones para casos de uso comunes
var (
	// Proyecciones básicas para información principal de la orden
	ProjectionGroupBasic = []Projection{
		ProjectionReferenceID,
		ProjectionPageInfo,
	}

	// Proyecciones para paquetes y contenido
	ProjectionGroupPackages = []Projection{
		ProjectionPackages,
		ProjectionPackagesLPN,
		ProjectionPackagesWeight,
		ProjectionPackagesDimensions,
	}

	// Proyecciones para información de ubicación de destino
	ProjectionGroupDestination = []Projection{
		ProjectionDestination,
		ProjectionDestinationDeliveryInstructions,
		ProjectionDestinationAddressInfo,
		ProjectionDestinationAddressLine1,
		ProjectionDestinationContact,
	}

	// Proyecciones para información de ubicación de origen
	ProjectionGroupOrigin = []Projection{
		ProjectionOrigin,
		ProjectionOriginAddressInfo,
		ProjectionOriginAddressLine1,
		ProjectionOriginContact,
	}

	// Proyecciones para información detallada de items
	ProjectionGroupItems = []Projection{
		ProjectionPackagesItems,
		ProjectionPackagesItemsSKU,
		ProjectionPackagesItemsDescription,
		ProjectionPackagesItemsWeight,
		ProjectionPackagesItemsDimensions,
	}

	// Proyecciones para información de fechas
	ProjectionGroupDates = []Projection{
		ProjectionPromisedDate,
		ProjectionPromisedDateDateRange,
		ProjectionCollectAvailabilityDate,
	}
)

type ProjectionSet map[string]struct{}

// HasProjection verifica si una proyección está presente en el set.
func (ps ProjectionSet) HasProjection(p Projection) bool {
	_, ok := ps[string(p)]
	return ok
}

// HasAnyProjection verifica si al menos una proyección del grupo está presente
func (ps ProjectionSet) HasAnyProjection(group []Projection) bool {
	for _, p := range group {
		if ps.HasProjection(p) {
			return true
		}
	}
	return false
}

// HasAllProjections verifica si todas las proyecciones del grupo están presentes
func (ps ProjectionSet) HasAllProjections(group []Projection) bool {
	for _, p := range group {
		if !ps.HasProjection(p) {
			return false
		}
	}
	return true
}

// HasAnyWithPrefix verifica si hay alguna proyección que comience con el prefijo dado
func (ps ProjectionSet) HasAnyWithPrefix(prefix string) bool {
	for field := range ps {
		if strings.HasPrefix(field, prefix) {
			return true
		}
	}
	return false
}

// BuildProjectionSet construye un map[string]struct{} a partir del slice de campos solicitados.
func BuildProjectionSet(fields []string) ProjectionSet {
	set := make(ProjectionSet, len(fields))
	for _, f := range fields {
		set[f] = struct{}{}
	}
	return set
}
