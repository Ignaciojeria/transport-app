package orders

import "reflect"

// Field representa un campo específico con métodos para verificación
type Field struct {
	path string
}

// Has verifica si este campo está en el mapa de campos solicitados
func (f Field) Has(requestedFields map[string]struct{}) bool {
	_, exists := requestedFields[f.path]
	return exists
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
func (p Projection) Packages() Field {
	return Field{path: "packages"}
}

func (p Projection) PackagesLPN() Field {
	return Field{path: "packages.lpn"}
}

// Métodos para campos de Weight en Package
func (p Projection) PackagesWeight() Field {
	return Field{path: "packages.weight"}
}

func (p Projection) PackagesWeightUnit() Field {
	return Field{path: "packages.weight.unit"}
}

func (p Projection) PackagesWeightValue() Field {
	return Field{path: "packages.weight.value"}
}

// Métodos para campos de Dimensions en Package
func (p Projection) PackagesDimensions() Field {
	return Field{path: "packages.dimensions"}
}

func (p Projection) PackagesDimensionsLength() Field {
	return Field{path: "packages.dimensions.length"}
}

func (p Projection) PackagesDimensionsHeight() Field {
	return Field{path: "packages.dimensions.height"}
}

func (p Projection) PackagesDimensionsWidth() Field {
	return Field{path: "packages.dimensions.width"}
}

func (p Projection) PackagesDimensionsUnit() Field {
	return Field{path: "packages.dimensions.unit"}
}

// Métodos para campos de Insurance en Package
func (p Projection) PackagesInsurance() Field {
	return Field{path: "packages.insurance"}
}

func (p Projection) PackagesInsuranceCurrency() Field {
	return Field{path: "packages.insurance.currency"}
}

func (p Projection) PackagesInsuranceUnitValue() Field {
	return Field{path: "packages.insurance.unitValue"}
}

// Métodos para campos de Label en Package
func (p Projection) PackagesLabels() Field {
	return Field{path: "packages.labels"}
}

func (p Projection) PackagesLabelsType() Field {
	return Field{path: "packages.labels.type"}
}

func (p Projection) PackagesLabelsValue() Field {
	return Field{path: "packages.labels.value"}
}

// Métodos para campos de Item en Package
func (p Projection) PackagesItems() Field {
	return Field{path: "packages.items"}
}

func (p Projection) PackagesItemsSKU() Field {
	return Field{path: "packages.items.sku"}
}

func (p Projection) PackagesItemsDescription() Field {
	return Field{path: "packages.items.description"}
}

// Métodos para campos de Dimensions en Item
func (p Projection) PackagesItemsDimensions() Field {
	return Field{path: "packages.items.dimensions"}
}

func (p Projection) PackagesItemsDimensionsLength() Field {
	return Field{path: "packages.items.dimensions.length"}
}

func (p Projection) PackagesItemsDimensionsHeight() Field {
	return Field{path: "packages.items.dimensions.height"}
}

func (p Projection) PackagesItemsDimensionsWidth() Field {
	return Field{path: "packages.items.dimensions.width"}
}

func (p Projection) PackagesItemsDimensionsUnit() Field {
	return Field{path: "packages.items.dimensions.unit"}
}

// Métodos para campos de Insurance en Item
func (p Projection) PackagesItemsInsurance() Field {
	return Field{path: "packages.items.insurance"}
}

func (p Projection) PackagesItemsInsuranceCurrency() Field {
	return Field{path: "packages.items.insurance.currency"}
}

func (p Projection) PackagesItemsInsuranceUnitValue() Field {
	return Field{path: "packages.items.insurance.unitValue"}
}

// Métodos para campos de Skill en Item
func (p Projection) PackagesItemsSkills() Field {
	return Field{path: "packages.items.skills"}
}

func (p Projection) PackagesItemsSkillsType() Field {
	return Field{path: "packages.items.skills.type"}
}

func (p Projection) PackagesItemsSkillsValue() Field {
	return Field{path: "packages.items.skills.value"}
}

func (p Projection) PackagesItemsSkillsDescription() Field {
	return Field{path: "packages.items.skills.description"}
}

// Métodos para campos de Quantity en Item
func (p Projection) PackagesItemsQuantity() Field {
	return Field{path: "packages.items.quantity"}
}

func (p Projection) PackagesItemsQuantityNumber() Field {
	return Field{path: "packages.items.quantity.quantityNumber"}
}

func (p Projection) PackagesItemsQuantityUnit() Field {
	return Field{path: "packages.items.quantity.quantityUnit"}
}

// Métodos para campos de Weight en Item
func (p Projection) PackagesItemsWeight() Field {
	return Field{path: "packages.items.weight"}
}

func (p Projection) PackagesItemsWeightUnit() Field {
	return Field{path: "packages.items.weight.unit"}
}

func (p Projection) PackagesItemsWeightValue() Field {
	return Field{path: "packages.items.weight.value"}
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

func (p Projection) ExtraFieldsPoliticalAreaId() Field {
	return Field{path: "extraFields.destinationPoliticalAreaId"}
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
