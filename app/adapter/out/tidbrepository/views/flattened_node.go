package views

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

type FlattenedNodeView struct {
	NodeID             int64               `gorm:"column:node_id"`
	ReferenceID        string              `gorm:"column:reference_id"`
	NodeName           string              `gorm:"column:node_name"`
	NodeType           string              `gorm:"column:node_type"`
	OperatorID         int64               `gorm:"column:operator_id"`
	OperatorName       string              `gorm:"column:operator_name"`
	OperatorEmail      string              `gorm:"column:operator_email"`
	OperatorPhone      string              `gorm:"column:operator_phone"`
	OperatorNationalID string              `gorm:"column:operator_national_id"`
	OperatorDocuments  table.JSONReference `gorm:"column:operator_documents"`
	NodeReferences     table.JSONReference `gorm:"column:node_references"`
	AddressID          int64               `gorm:"column:address_id"`
	AddressLine1       string              `gorm:"column:address_line1"`
	AddressLine2       string              `gorm:"column:address_line2"`
	AddressLine3       string              `gorm:"column:address_line3"`
	County             string              `gorm:"column:county"`
	District           string              `gorm:"column:district"`
	Latitude           float32             `gorm:"column:latitude"`
	Longitude          float32             `gorm:"column:longitude"`
	Province           string              `gorm:"column:province"`
	State              string              `gorm:"column:state"`
	ZipCode            string              `gorm:"column:zip_code"`
	TimeZone           string              `gorm:"column:time_zone"`
}

func (n FlattenedNodeView) ToNodeInfo() domain.NodeInfo {
	return domain.NodeInfo{
		ID:          n.NodeID,
		ReferenceID: domain.ReferenceID(n.ReferenceID),
		Name:        n.NodeName,
		Type:        n.NodeType,
		References:  mapReferences(n.NodeReferences),
		Operator: domain.Operator{
			ID: n.OperatorID,
			Contact: domain.Contact{
				ID:         n.OperatorID,
				FullName:   n.OperatorName,
				Email:      n.OperatorEmail,
				Phone:      n.OperatorPhone,
				NationalID: n.OperatorNationalID,
				Documents:  mapDocuments(n.OperatorDocuments),
			},
		},
		AddressInfo: domain.AddressInfo{
			ID:           n.AddressID,
			AddressLine1: n.AddressLine1,
			AddressLine2: n.AddressLine2,
			AddressLine3: n.AddressLine3,
			County:       n.County,
			District:     n.District,
			Latitude:     n.Latitude,
			Longitude:    n.Longitude,
			Province:     n.Province,
			State:        n.State,
			ZipCode:      n.ZipCode,
			TimeZone:     n.TimeZone,
		},
	}
}

func mapReferences(references table.JSONReference) []domain.Reference {
	result := make([]domain.Reference, len(references))
	for i, ref := range references {
		result[i] = domain.Reference{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return result
}
