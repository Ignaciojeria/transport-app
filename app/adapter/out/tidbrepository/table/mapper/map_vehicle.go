package mapper

import (
	"context"
	"encoding/json"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapVehicle(ctx context.Context, v domain.Vehicle) table.Vehicle {
	// Convertir los campos JSON a bytes
	weightJSON, _ := json.Marshal(v.Weight)
	insuranceJSON, _ := json.Marshal(v.Insurance)
	technicalReviewJSON, _ := json.Marshal(v.TechnicalReview)
	dimensionsJSON, _ := json.Marshal(v.Dimensions)

	return table.Vehicle{
		DocumentID:         string(v.DocID(ctx)),
		TenantID:           sharedcontext.TenantIDFromContext(ctx),
		Plate:              v.Plate,
		CertificateDate:    v.CertificateDate,
		VehicleCategoryDoc: string(v.VehicleCategory.DocID(ctx)),
		VehicleHeadersDoc:  string(v.Headers.DocID(ctx)),
		CarrierDoc:         string(v.Carrier.DocID(ctx)),
		Weight:             weightJSON,
		Insurance:          insuranceJSON,
		TechnicalReview:    technicalReviewJSON,
		Dimensions:         dimensionsJSON,
	}
}
