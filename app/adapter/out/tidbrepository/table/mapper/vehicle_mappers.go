package mapper

import (
	"encoding/json"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func DomainToTableVehicle(d domain.Vehicle) table.Vehicle {
	weight, _ := json.Marshal(d.Weight)
	insurance, _ := json.Marshal(d.Insurance)
	technicalReview, _ := json.Marshal(d.TechnicalReview)
	dimensions, _ := json.Marshal(d.Dimensions)
	var vehicleCategoryID *int64 = nil
	if d.VehicleCategory.ID != 0 {
		vehicleCategoryID = &d.VehicleCategory.ID
	}
	var carrierID *int64 = nil
	if d.Carrier.ID != 0 {
		carrierID = &d.Carrier.ID
	}
	return table.Vehicle{
		ID:                d.ID,
		ReferenceID:       d.ReferenceID,
		Plate:             d.Plate,
		IsActive:          d.IsActive,
		CertificateDate:   d.CertificateDate,
		VehicleCategoryID: vehicleCategoryID,
		VehicleHeadersID:  d.Headers.ID,
		OrganizationID:    d.Organization.ID,
		CarrierID:         carrierID,
		Weight:            table.JSONB(weight),
		Insurance:         table.JSONB(insurance),
		TechnicalReview:   table.JSONB(technicalReview),
		Dimensions:        table.JSONB(dimensions),
	}
}
