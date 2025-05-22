package orders

import (
	"reflect"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

// Field representa un campo específico con métodos para verificación
type Field struct {
	path string
}

// Has verifica si este campo está en el mapa de campos solicitados
func (f Field) Has(requestedFields map[string]struct{}) bool {
	_, exists := requestedFields[f.path]
	return exists
}

func init() {
	ioc.Registry(NewProjection)
}

func NewProjection() Projection {
	return Projection{}
}

// Projection representa un conjunto de campos solicitados por el cliente en la query GraphQL.
type Projection struct{}

// Métodos para campos de Order
func (p Projection) ReferenceID() Field {
	return Field{path: "referenceID"}
}

// Métodos para campos de CollectAvailabilityDate
func (p Projection) CollectAvailabilityDate() Field {
	return Field{path: "collectAvailabilityDate"}
}

func (p Projection) CollectAvailabilityDateDate() Field {
	return Field{path: "collectAvailabilityDate.date"}
}

func (p Projection) CollectAvailabilityDateTimeRange() Field {
	return Field{path: "collectAvailabilityDate.timeRange"}
}

func (p Projection) CollectAvailabilityDateStartTime() Field {
	return Field{path: "collectAvailabilityDate.timeRange.startTime"}
}

func (p Projection) CollectAvailabilityDateEndTime() Field {
	return Field{path: "collectAvailabilityDate.timeRange.endTime"}
}

// Métodos para campos de Destination Location
func (p Projection) Destination() Field {
	return Field{path: "destination"}
}

func (p Projection) DestinationDeliveryInstructions() Field {
	return Field{path: "destination.deliveryInstructions"}
}

// Métodos para campos de NodeInfo en Destination
func (p Projection) DestinationNodeInfo() Field {
	return Field{path: "destination.nodeInfo"}
}

func (p Projection) DestinationNodeInfoReferenceId() Field {
	return Field{path: "destination.nodeInfo.referenceId"}
}

func (p Projection) DestinationNodeInfoName() Field {
	return Field{path: "destination.nodeInfo.name"}
}

// Métodos para campos de AddressInfo en Destination
func (p Projection) DestinationAddressInfo() Field {
	return Field{path: "destination.addressInfo"}
}

func (p Projection) DestinationAddressLine1() Field {
	return Field{path: "destination.addressInfo.addressLine1"}
}

func (p Projection) DestinationAddressLine2() Field {
	return Field{path: "destination.addressInfo.addressLine2"}
}

func (p Projection) DestinationDistrict() Field {
	return Field{path: "destination.addressInfo.district"}
}

func (p Projection) DestinationLatitude() Field {
	return Field{path: "destination.addressInfo.latitude"}
}

func (p Projection) DestinationLongitude() Field {
	return Field{path: "destination.addressInfo.longitude"}
}

func (p Projection) DestinationProvince() Field {
	return Field{path: "destination.addressInfo.province"}
}

func (p Projection) DestinationState() Field {
	return Field{path: "destination.addressInfo.state"}
}

func (p Projection) DestinationTimeZone() Field {
	return Field{path: "destination.addressInfo.timeZone"}
}

func (p Projection) DestinationZipCode() Field {
	return Field{path: "destination.addressInfo.zipCode"}
}

// Métodos para campos de Contact en Destination
func (p Projection) DestinationContact() Field {
	return Field{path: "destination.addressInfo.contact"}
}

func (p Projection) DestinationContactEmail() Field {
	return Field{path: "destination.addressInfo.contact.email"}
}

func (p Projection) DestinationContactFullName() Field {
	return Field{path: "destination.addressInfo.contact.fullName"}
}

func (p Projection) DestinationContactNationalID() Field {
	return Field{path: "destination.addressInfo.contact.nationalID"}
}

func (p Projection) DestinationContactPhone() Field {
	return Field{path: "destination.addressInfo.contact.phone"}
}

// Métodos para campos de ContactMethods en Destination
func (p Projection) DestinationContactMethods() Field {
	return Field{path: "destination.addressInfo.contact.additionalContactMethods"}
}

func (p Projection) DestinationContactMethodsType() Field {
	return Field{path: "destination.addressInfo.contact.additionalContactMethods.type"}
}

func (p Projection) DestinationContactMethodsValue() Field {
	return Field{path: "destination.addressInfo.contact.additionalContactMethods.value"}
}

// Métodos para campos de Documents en Destination
func (p Projection) DestinationDocuments() Field {
	return Field{path: "destination.addressInfo.contact.documents"}
}

func (p Projection) DestinationDocumentsType() Field {
	return Field{path: "destination.addressInfo.contact.documents.type"}
}

func (p Projection) DestinationDocumentsValue() Field {
	return Field{path: "destination.addressInfo.contact.documents.value"}
}

// Métodos para campos de Origin Location
func (p Projection) Origin() Field {
	return Field{path: "origin"}
}

func (p Projection) OriginDeliveryInstructions() Field {
	return Field{path: "origin.deliveryInstructions"}
}

// Métodos para campos de NodeInfo en Origin
func (p Projection) OriginNodeInfo() Field {
	return Field{path: "origin.nodeInfo"}
}

func (p Projection) OriginNodeInfoReferenceId() Field {
	return Field{path: "origin.nodeInfo.referenceId"}
}

func (p Projection) OriginNodeInfoName() Field {
	return Field{path: "origin.nodeInfo.name"}
}

// Métodos para campos de AddressInfo en Origin
func (p Projection) OriginAddressInfo() Field {
	return Field{path: "origin.addressInfo"}
}

func (p Projection) OriginAddressLine1() Field {
	return Field{path: "origin.addressInfo.addressLine1"}
}

func (p Projection) OriginAddressLine2() Field {
	return Field{path: "origin.addressInfo.addressLine2"}
}

func (p Projection) OriginDistrict() Field {
	return Field{path: "origin.addressInfo.district"}
}

func (p Projection) OriginLatitude() Field {
	return Field{path: "origin.addressInfo.latitude"}
}

func (p Projection) OriginLongitude() Field {
	return Field{path: "origin.addressInfo.longitude"}
}

func (p Projection) OriginProvince() Field {
	return Field{path: "origin.addressInfo.province"}
}

func (p Projection) OriginState() Field {
	return Field{path: "origin.addressInfo.state"}
}

func (p Projection) OriginTimeZone() Field {
	return Field{path: "origin.addressInfo.timeZone"}
}

func (p Projection) OriginZipCode() Field {
	return Field{path: "origin.addressInfo.zipCode"}
}

// Métodos para campos de Contact en Origin
func (p Projection) OriginContact() Field {
	return Field{path: "origin.addressInfo.contact"}
}

func (p Projection) OriginContactEmail() Field {
	return Field{path: "origin.addressInfo.contact.email"}
}

func (p Projection) OriginContactFullName() Field {
	return Field{path: "origin.addressInfo.contact.fullName"}
}

func (p Projection) OriginContactNationalID() Field {
	return Field{path: "origin.addressInfo.contact.nationalID"}
}

func (p Projection) OriginContactPhone() Field {
	return Field{path: "origin.addressInfo.contact.phone"}
}

// Métodos para campos de ContactMethods en Origin
func (p Projection) OriginContactMethods() Field {
	return Field{path: "origin.addressInfo.contact.additionalContactMethods"}
}

func (p Projection) OriginContactMethodsType() Field {
	return Field{path: "origin.addressInfo.contact.additionalContactMethods.type"}
}

func (p Projection) OriginContactMethodsValue() Field {
	return Field{path: "origin.addressInfo.contact.additionalContactMethods.value"}
}

// Métodos para campos de Documents en Origin
func (p Projection) OriginDocuments() Field {
	return Field{path: "origin.addressInfo.contact.documents"}
}

func (p Projection) OriginDocumentsType() Field {
	return Field{path: "origin.addressInfo.contact.documents.type"}
}

func (p Projection) OriginDocumentsValue() Field {
	return Field{path: "origin.addressInfo.contact.documents.value"}
}

// Métodos para campos de OrderType
func (p Projection) OrderType() Field {
	return Field{path: "orderType"}
}

func (p Projection) OrderTypeType() Field {
	return Field{path: "orderType.type"}
}

func (p Projection) OrderTypeDescription() Field {
	return Field{path: "orderType.description"}
}

// Métodos para campos de Package
func (p Projection) DeliveryUnit() Field {
	return Field{path: "deliveryUnit"}
}

func (p Projection) DeliveryUnitLPN() Field {
	return Field{path: "deliveryUnit.lpn"}
}

// Métodos para campos de Weight en Package
func (p Projection) DeliveryUnitWeight() Field {
	return Field{path: "deliveryUnit.weight"}
}

func (p Projection) DeliveryUnitWeightUnit() Field {
	return Field{path: "deliveryUnit.weight.unit"}
}

func (p Projection) DeliveryUnitWeightValue() Field {
	return Field{path: "deliveryUnit.weight.value"}
}

// Métodos para campos de Dimensions en Package
func (p Projection) DeliveryUnitDimensions() Field {
	return Field{path: "deliveryUnit.dimensions"}
}

func (p Projection) DeliveryUnitDimensionsLength() Field {
	return Field{path: "deliveryUnit.dimensions.length"}
}

func (p Projection) DeliveryUnitDimensionsHeight() Field {
	return Field{path: "deliveryUnit.dimensions.height"}
}

func (p Projection) DeliveryUnitDimensionsWidth() Field {
	return Field{path: "deliveryUnit.dimensions.width"}
}

func (p Projection) DeliveryUnitDimensionsUnit() Field {
	return Field{path: "deliveryUnit.dimensions.unit"}
}

// Métodos para campos de Insurance en Package
func (p Projection) DeliveryUnitInsurance() Field {
	return Field{path: "deliveryUnit.insurance"}
}

func (p Projection) DeliveryUnitInsuranceCurrency() Field {
	return Field{path: "deliveryUnit.insurance.currency"}
}

func (p Projection) DeliveryUnitInsuranceUnitValue() Field {
	return Field{path: "deliveryUnit.insurance.unitValue"}
}

// Métodos para campos de Label en Package
func (p Projection) DeliveryUnitLabels() Field {
	return Field{path: "deliveryUnit.labels"}
}

func (p Projection) DeliveryUnitLabelsType() Field {
	return Field{path: "deliveryUnit.labels.type"}
}

func (p Projection) DeliveryUnitLabelsValue() Field {
	return Field{path: "deliveryUnit.labels.value"}
}

// Métodos para campos de Item en Package
func (p Projection) DeliveryUnitItems() Field {
	return Field{path: "deliveryUnit.items"}
}

func (p Projection) DeliveryUnitItemsSKU() Field {
	return Field{path: "deliveryUnit.items.sku"}
}

func (p Projection) DeliveryUnitItemsDescription() Field {
	return Field{path: "deliveryUnit.items.description"}
}

// Métodos para campos de Dimensions en Item
func (p Projection) DeliveryUnitItemsDimensions() Field {
	return Field{path: "deliveryUnit.items.dimensions"}
}

func (p Projection) DeliveryUnitItemsDimensionsLength() Field {
	return Field{path: "deliveryUnit.items.dimensions.length"}
}

func (p Projection) DeliveryUnitItemsDimensionsHeight() Field {
	return Field{path: "deliveryUnit.items.dimensions.height"}
}

func (p Projection) DeliveryUnitItemsDimensionsWidth() Field {
	return Field{path: "deliveryUnit.items.dimensions.width"}
}

func (p Projection) DeliveryUnitItemsDimensionsUnit() Field {
	return Field{path: "deliveryUnit.items.dimensions.unit"}
}

// Métodos para campos de Insurance en Item
func (p Projection) DeliveryUnitItemsInsurance() Field {
	return Field{path: "deliveryUnit.items.insurance"}
}

func (p Projection) DeliveryUnitItemsInsuranceCurrency() Field {
	return Field{path: "deliveryUnit.items.insurance.currency"}
}

func (p Projection) DeliveryUnitItemsInsuranceUnitValue() Field {
	return Field{path: "deliveryUnit.items.insurance.unitValue"}
}

// Métodos para campos de Skill en Item
func (p Projection) DeliveryUnitItemsSkills() Field {
	return Field{path: "deliveryUnit.items.skills"}
}

func (p Projection) DeliveryUnitItemsSkillsType() Field {
	return Field{path: "deliveryUnit.items.skills.type"}
}

func (p Projection) DeliveryUnitItemsSkillsValue() Field {
	return Field{path: "deliveryUnit.items.skills.value"}
}

func (p Projection) DeliveryUnitItemsSkillsDescription() Field {
	return Field{path: "deliveryUnit.items.skills.description"}
}

// Métodos para campos de Quantity en Item
func (p Projection) DeliveryUnitItemsQuantity() Field {
	return Field{path: "deliveryUnit.items.quantity"}
}

func (p Projection) DeliveryUnitItemsQuantityNumber() Field {
	return Field{path: "deliveryUnit.items.quantity.quantityNumber"}
}

func (p Projection) DeliveryUnitItemsQuantityUnit() Field {
	return Field{path: "deliveryUnit.items.quantity.quantityUnit"}
}

// Métodos para campos de Weight en Item
func (p Projection) DeliveryUnitItemsWeight() Field {
	return Field{path: "deliveryUnit.items.weight"}
}

func (p Projection) DeliveryUnitItemsWeightUnit() Field {
	return Field{path: "deliveryUnit.items.weight.unit"}
}

func (p Projection) DeliveryUnitItemsWeightValue() Field {
	return Field{path: "deliveryUnit.items.weight.value"}
}

// Métodos para campos de PromisedDate
func (p Projection) PromisedDate() Field {
	return Field{path: "promisedDate"}
}

func (p Projection) PromisedDateServiceCategory() Field {
	return Field{path: "promisedDate.serviceCategory"}
}

// Métodos para campos de TimeRange en PromisedDate
func (p Projection) PromisedDateTimeRange() Field {
	return Field{path: "promisedDate.timeRange"}
}

func (p Projection) PromisedDateTimeRangeStartTime() Field {
	return Field{path: "promisedDate.timeRange.startTime"}
}

func (p Projection) PromisedDateTimeRangeEndTime() Field {
	return Field{path: "promisedDate.timeRange.endTime"}
}

// Métodos para campos de DateRange en PromisedDate
func (p Projection) PromisedDateDateRange() Field {
	return Field{path: "promisedDate.dateRange"}
}

func (p Projection) PromisedDateDateRangeStartDate() Field {
	return Field{path: "promisedDate.dateRange.startDate"}
}

func (p Projection) PromisedDateDateRangeEndDate() Field {
	return Field{path: "promisedDate.dateRange.endDate"}
}

// Métodos para campos de Reference
func (p Projection) References() Field {
	return Field{path: "references"}
}

func (p Projection) ReferencesType() Field {
	return Field{path: "references.type"}
}

func (p Projection) ReferencesValue() Field {
	return Field{path: "references.value"}
}

// Métodos para campos de ExtraFields
func (p Projection) ExtraFields() Field {
	return Field{path: "extraFields"}
}

func (p Projection) ExtraFieldsKey() Field {
	return Field{path: "extraFields.key"}
}

func (p Projection) ExtraFieldValue() Field {
	return Field{path: "extraFields.value"}
}

// Métodos para campos de paginación
func (p Projection) PageInfo() Field {
	return Field{path: "pageInfo"}
}

func (p Projection) PageInfoHasNextPage() Field {
	return Field{path: "pageInfo.hasNextPage"}
}

func (p Projection) PageInfoHasPreviousPage() Field {
	return Field{path: "pageInfo.hasPreviousPage"}
}

func (p Projection) PageInfoStartCursor() Field {
	return Field{path: "pageInfo.startCursor"}
}

func (p Projection) PageInfoEndCursor() Field {
	return Field{path: "pageInfo.endCursor"}
}

// Métodos para campos de DeliveryUnitsReport
func (p Projection) Commerce() Field {
	return Field{path: "commerce"}
}

func (p Projection) Consumer() Field {
	return Field{path: "consumer"}
}

func (p Projection) Channel() Field {
	return Field{path: "channel"}
}

func (p Projection) GroupBy() Field {
	return Field{path: "groupBy"}
}

func (p Projection) GroupByType() Field {
	return Field{path: "groupBy.type"}
}

func (p Projection) GroupByValue() Field {
	return Field{path: "groupBy.value"}
}

// Métodos para campos de Delivery
func (p Projection) Delivery() Field {
	return Field{path: "delivery"}
}

func (p Projection) DeliveryStatus() Field {
	return Field{path: "delivery.status"}
}

func (p Projection) DeliveryRecipient() Field {
	return Field{path: "delivery.recipient"}
}

func (p Projection) DeliveryRecipientFullName() Field {
	return Field{path: "delivery.recipient.fullName"}
}

func (p Projection) DeliveryRecipientNationalID() Field {
	return Field{path: "delivery.recipient.nationalID"}
}

func (p Projection) DeliveryHandledAt() Field {
	return Field{path: "delivery.handledAt"}
}

func (p Projection) DeliveryFailure() Field {
	return Field{path: "delivery.failure"}
}

func (p Projection) DeliveryFailureDetail() Field {
	return Field{path: "delivery.failure.detail"}
}

func (p Projection) DeliveryFailureReason() Field {
	return Field{path: "delivery.failure.reason"}
}

func (p Projection) DeliveryFailureReferenceID() Field {
	return Field{path: "delivery.failure.referenceID"}
}

func (p Projection) DeliveryLocation() Field {
	return Field{path: "delivery.location"}
}

func (p Projection) DeliveryLocationLatitude() Field {
	return Field{path: "delivery.location.latitude"}
}

func (p Projection) DeliveryLocationLongitude() Field {
	return Field{path: "delivery.location.longitude"}
}

// Métodos para campos de EvidencePhotos (array)
func (p Projection) DeliveryEvidencePhotos() Field {
	return Field{path: "delivery.evidencePhotos"}
}

func (p Projection) DeliveryEvidencePhotosTakenAt() Field {
	return Field{path: "delivery.evidencePhotos.takenAt"}
}

func (p Projection) DeliveryEvidencePhotosType() Field {
	return Field{path: "delivery.evidencePhotos.type"}
}

func (p Projection) DeliveryEvidencePhotosURL() Field {
	return Field{path: "delivery.evidencePhotos.url"}
}

// GetAllProjections devuelve un mapa con todas las proyecciones disponibles
func getAllProjections() map[string]Field {
	var p Projection
	projections := make(map[string]Field)

	// Obtener el tipo de Projection
	pType := reflect.TypeOf(p)
	pValue := reflect.ValueOf(p)

	// Iterar sobre todos los métodos de Projection
	for i := 0; i < pType.NumMethod(); i++ {
		method := pType.Method(i)

		// Solo considerar métodos que no reciben parámetros y devuelven Field
		if method.Type.NumIn() == 1 && method.Type.NumOut() == 1 {
			if method.Type.Out(0).Name() == "Field" {
				// Invocar el método para obtener el Field
				result := pValue.Method(i).Call(nil)[0].Interface().(Field)
				projections[method.Name] = result
			}
		}
	}

	return projections
}

// GetAllProjectionPaths devuelve un slice con todas las rutas de campos disponibles
func GetAllProjectionPaths() []string {
	projections := getAllProjections()
	paths := make([]string, 0, len(projections))

	for _, field := range projections {
		paths = append(paths, field.String())
	}

	return paths
}

// String devuelve la ruta del campo
func (f Field) String() string {
	return f.path
}
