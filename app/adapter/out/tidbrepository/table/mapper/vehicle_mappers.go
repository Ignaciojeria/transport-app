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

	return table.Vehicle{
		Plate:             d.Plate,
		IsActive:          d.IsActive,
		CertificateDate:   d.CertificateDate,
		VehicleCategoryID: vehicleCategoryID,
		//	OrganizationID:    d.Organization.ID,
		Weight:          table.JSONB(weight),
		Insurance:       table.JSONB(insurance),
		TechnicalReview: table.JSONB(technicalReview),
		Dimensions:      table.JSONB(dimensions),
	}
}
