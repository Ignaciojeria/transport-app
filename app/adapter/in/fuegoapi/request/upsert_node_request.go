package request

import "transport-app/app/domain"

type UpsertNodeRequest struct {
	ReferenceID string `json:"referenceId" validate:"required"`
	Name        string `json:"name" validate:"required"`
	NodeAddress struct {
		AddressLine1 string  `json:"addressLine1"`
		AddressLine2 string  `json:"addressLine2"`
		AddressLine3 string  `json:"addressLine3"`
		County       string  `json:"county"`
		District     string  `json:"district"`
		Latitude     float32 `json:"latitude"`
		Longitude    float32 `json:"longitude"`
		Province     string  `json:"province"`
		State        string  `json:"state"`
		TimeZone     string  `json:"timeZone"`
		ZipCode      string  `json:"zipCode"`
	} `json:"nodeAddress"`
	OperatorContact struct {
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		NationalID string `json:"nationalID"`
		Documents  []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"documents"`
		FullName string `json:"fullName"`
	} `json:"operatorContact"`
}

func (req UpsertNodeRequest) Map() domain.Origin {
	nodeInfo := domain.NodeInfo{
		ReferenceID: domain.ReferenceID(req.ReferenceID),
		Name:        &req.Name,
	}
	nodeAddress := domain.AddressInfo{
		AddressLine1: req.NodeAddress.AddressLine1,
		AddressLine2: req.NodeAddress.AddressLine2,
		AddressLine3: req.NodeAddress.AddressLine3,
		County:       req.NodeAddress.County,
		District:     req.NodeAddress.District,
		Latitude:     req.NodeAddress.Latitude,
		Longitude:    req.NodeAddress.Longitude,
		Province:     req.NodeAddress.Province,
		State:        req.NodeAddress.State,
		TimeZone:     req.NodeAddress.TimeZone,
		ZipCode:      req.NodeAddress.ZipCode,
		Contact: domain.Contact{
			Email:      req.OperatorContact.Email,
			Phone:      req.OperatorContact.Phone,
			NationalID: req.OperatorContact.NationalID,
			FullName:   req.OperatorContact.FullName,
			Documents: func() []domain.Document {
				docs := make([]domain.Document, len(req.OperatorContact.Documents))
				for i, doc := range req.OperatorContact.Documents {
					docs[i] = domain.Document{
						Type:  doc.Type,
						Value: doc.Value,
					}
				}
				return docs
			}(),
		},
	}
	return domain.Origin{
		NodeInfo:    nodeInfo,
		AddressInfo: nodeAddress,
	}
}
