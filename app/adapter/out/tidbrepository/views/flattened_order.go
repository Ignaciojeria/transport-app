package views

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	"github.com/biter777/countries"
)

type FlattenedOrderView struct {
	OrderID                      int64               `gorm:"column:order_id"`
	ReferenceID                  string              `gorm:"column:reference_id"`
	OrganizationCountryID        int64               `gorm:"column:organization_country_id"`
	OrganizationCountry          string              `gorm:"column:organization_country"`
	CommerceID                   int64               `gorm:"column:commerce_id"`
	CommerceName                 string              `gorm:"column:commerce_name"`
	ConsumerID                   int64               `gorm:"column:consumer_id"`
	ConsumerName                 string              `gorm:"column:consumer_name"`
	OrderStatusID                int64               `gorm:"column:order_status_id"`
	OrderStatus                  string              `gorm:"column:order_status"`
	OrderTypeID                  int64               `gorm:"column:order_type_id"`
	OrderType                    string              `gorm:"column:order_type"`
	OrderTypeDescription         string              `gorm:"column:order_type_description"`
	DeliveryInstructions         string              `gorm:"column:delivery_instructions"`
	OriginContactID              int64               `gorm:"column:origin_contact_id"`
	OriginContactName            string              `gorm:"column:origin_contact_name"`
	OriginContactPhone           string              `gorm:"column:origin_contact_phone"`
	OriginContactDocuments       table.JSONReference `gorm:"column:origin_contact_documents"`
	DestinationContactID         int64               `gorm:"column:destination_contact_id"`
	DestinationContactName       string              `gorm:"column:destination_contact_name"`
	DestinationContactPhone      string              `gorm:"column:destination_contact_phone"`
	DestinationContactDocuments  table.JSONReference `gorm:"column:destination_contact_documents"`
	OriginAddressInfoID          int64               `gorm:"column:origin_address_info_id"`
	OriginAddressLine1           string              `gorm:"column:origin_address_line1"`
	OriginAddressLine2           string              `gorm:"column:origin_address_line2"`
	OriginAddressLine3           string              `gorm:"column:origin_address_line3"`
	OriginContactEmail           string              `gorm:"column:origin_contact_email"`
	OriginContactNationalID      string              `gorm:"column:origin_contact_national_id"`
	OriginState                  string              `gorm:"column:origin_state"`
	OriginProvince               string              `gorm:"column:origin_province"`
	OriginCounty                 string              `gorm:"column:origin_county"`
	OriginDistrict               string              `gorm:"column:origin_district"`
	OriginZipCode                string              `gorm:"column:origin_zipcode"`
	OriginLatitude               float32             `gorm:"column:origin_latitude"`
	OriginLongitude              float32             `gorm:"column:origin_longitude"`
	OriginTimeZone               string              `gorm:"column:origin_timezone"`
	OriginNodeInfoID             int64               `gorm:"column:origin_node_info_id"`
	OriginNodeReferenceID        string              `gorm:"column:origin_node_reference_id"`
	OriginNodeName               string              `gorm:"column:origin_node_name"`
	OriginNodeType               string              `gorm:"column:origin_node_type"`
	OriginNodeOperatorID         int64               `gorm:"column:origin_node_operator_id"`
	OriginNodeOperatorName       string              `gorm:"column:origin_node_operator_name"`
	DestinationAddressInfoID     int64               `gorm:"column:destination_address_info_id"`
	DestinationAddressLine1      string              `gorm:"column:destination_address_line1"`
	DestinationAddressLine2      string              `gorm:"column:destination_address_line2"`
	DestinationAddressLine3      string              `gorm:"column:destination_address_line3"`
	DestinationContactEmail      string              `gorm:"column:destination_contact_email"`
	DestinationContactNationalID string              `gorm:"column:destination_contact_national_id"`
	DestinationState             string              `gorm:"column:destination_state"`
	DestinationProvince          string              `gorm:"column:destination_province"`
	DestinationCounty            string              `gorm:"column:destination_county"`
	DestinationDistrict          string              `gorm:"column:destination_district"`
	DestinationZipCode           string              `gorm:"column:destination_zipcode"`
	DestinationLatitude          float32             `gorm:"column:destination_latitude"`
	DestinationLongitude         float32             `gorm:"column:destination_longitude"`
	DestinationTimeZone          string              `gorm:"column:destination_timezone"`
	DestinationNodeInfoID        int64               `gorm:"column:destination_node_info_id"`
	DestinationNodeReferenceID   string              `gorm:"column:destination_node_reference_id"`
	DestinationNodeName          string              `gorm:"column:destination_node_name"`
	DestinationNodeType          string              `gorm:"column:destination_node_type"`
	DestinationNodeOperatorID    int64               `gorm:"column:destination_node_operator_id"`
	DestinationNodeOperatorName  string              `gorm:"column:destination_node_operator_name"`
	Items                        table.JSONItems     `gorm:"column:items"`
	CollectAvailabilityDate      string              `gorm:"column:collect_availability_date"`
	CollectStartTime             string              `gorm:"column:collect_start_time"`
	CollectEndTime               string              `gorm:"column:collect_end_time"`
	PromisedStartDate            string              `gorm:"column:promised_start_date"`
	PromisedEndDate              string              `gorm:"column:promised_end_date"`
	PromisedStartTime            string              `gorm:"column:promised_start_time"`
	PromisedEndTime              string              `gorm:"column:promised_end_time"`
	TransportRequirements        table.JSONReference `gorm:"column:transport_requirements"`
}

type FlattenedPackageView struct {
	PackageID             int64                    `gorm:"column:package_id"`
	OrderID               int64                    `gorm:"column:order_id"`
	Items                 table.JSONItemReferences `gorm:"column:items_references"`
	OrganizationCountryID int64                    `gorm:"column:organization_country_id"`
	Lpn                   string                   `gorm:"column:lpn"`
	Height                float64                  `gorm:"column:height"`
	Width                 float64                  `gorm:"column:width"`
	Depth                 float64                  `gorm:"column:depth"`
	Unit                  string                   `gorm:"column:unit"`
	WeightValue           float64                  `gorm:"column:weight_value"`
	WeightUnit            string                   `gorm:"column:weight_unit"`
	UnitValue             float64                  `gorm:"column:unit_value"`
	Currency              string                   `gorm:"column:currency"`
}

type FlattenedOrderReferenceView struct {
	ReferenceID int64  `gorm:"column:reference_id"`
	OrderID     int64  `gorm:"column:order_id"`
	Type        string `gorm:"column:type"`
	Value       string `gorm:"column:value"`
}

type FlattenedVisitView struct {
	VisitID        int64  `gorm:"column:visit_id"`
	OrderID        int64  `gorm:"column:order_id"`
	Date           string `gorm:"column:date"`
	TimeRangeStart string `gorm:"column:time_range_start"`
	TimeRangeEnd   string `gorm:"column:time_range_end"`
}

func (o FlattenedOrderView) ToOrder(packages []FlattenedPackageView, refs []FlattenedOrderReferenceView, visits []FlattenedVisitView) domain.Order {
	references := make([]domain.Reference, len(refs))
	for i, ref := range refs {
		references[i] = domain.Reference{
			ID:    ref.ReferenceID,
			Type:  ref.Type,
			Value: ref.Value,
		}
	}

	var transportReqs []domain.Reference
	if o.TransportRequirements != nil {
		transportReqs = make([]domain.Reference, len(o.TransportRequirements))
		for i, req := range o.TransportRequirements {
			transportReqs[i] = domain.Reference{
				Type:  req.Type,
				Value: req.Value,
			}
		}
	}

	return domain.Order{
		ID:          o.OrderID,
		ReferenceID: domain.ReferenceID(o.ReferenceID),
		Organization: domain.Organization{
			OrganizationCountryID: o.OrganizationCountryID,
			Country:               countries.ByName(o.OrganizationCountry),
		},
		Commerce: domain.Commerce{
			ID:    o.CommerceID,
			Value: o.CommerceName,
		},
		Consumer: domain.Consumer{
			ID:    o.ConsumerID,
			Value: o.ConsumerName,
		},
		OrderStatus: domain.OrderStatus{
			ID:     o.OrderStatusID,
			Status: o.OrderStatus,
		},
		OrderType: domain.OrderType{
			ID:          o.OrderTypeID,
			Type:        o.OrderType,
			Description: o.OrderTypeDescription,
		},
		Origin: domain.NodeInfo{
			ID:          o.OriginNodeInfoID,
			ReferenceID: domain.ReferenceID(o.OriginNodeReferenceID),
			Name:        o.OriginNodeName,
			Type:        o.OriginNodeType,

			AddressInfo: domain.AddressInfo{
				ID: o.OriginAddressInfoID,
				Contact: domain.Contact{
					ID:         o.OriginContactID,
					FullName:   o.OriginContactName,
					Phone:      o.OriginContactPhone,
					Email:      o.OriginContactEmail,
					NationalID: o.OriginContactNationalID,
					Documents:  mapDocuments(o.OriginContactDocuments),
				},
				AddressLine1: o.OriginAddressLine1,
				AddressLine2: o.OriginAddressLine2,
				AddressLine3: o.OriginAddressLine3,
				State:        o.OriginState,
				Province:     o.OriginProvince,
				ZipCode:      o.OriginZipCode,
				County:       o.OriginCounty,
				District:     o.OriginDistrict,
				Latitude:     o.OriginLatitude,
				Longitude:    o.OriginLongitude,
				TimeZone:     o.OriginTimeZone,
			},
		},
		Destination: domain.NodeInfo{
			ID:          o.DestinationNodeInfoID,
			ReferenceID: domain.ReferenceID(o.DestinationNodeReferenceID),
			Name:        o.DestinationNodeName,
			Type:        o.DestinationNodeType,
			AddressInfo: domain.AddressInfo{
				ID: o.DestinationAddressInfoID,
				Contact: domain.Contact{
					ID:         o.DestinationContactID,
					FullName:   o.DestinationContactName,
					Phone:      o.DestinationContactPhone,
					Email:      o.DestinationContactEmail,
					NationalID: o.DestinationContactNationalID,
					Documents:  mapDocuments(o.DestinationContactDocuments),
				},
				AddressLine1: o.DestinationAddressLine1,
				AddressLine2: o.DestinationAddressLine2,
				AddressLine3: o.DestinationAddressLine3,
				State:        o.DestinationState,
				Province:     o.DestinationProvince,
				ZipCode:      o.DestinationZipCode,
				County:       o.DestinationCounty,
				District:     o.DestinationDistrict,
				Latitude:     o.DestinationLatitude,
				Longitude:    o.DestinationLongitude,
				TimeZone:     o.DestinationTimeZone,
			},
		},
		Packages:              mapPackages(packages),
		Items:                 mapJSONItems(o.Items),
		References:            references,
		TransportRequirements: transportReqs,
		CollectAvailabilityDate: domain.CollectAvailabilityDate{
			Date: o.CollectAvailabilityDate,
			TimeRange: domain.TimeRange{
				StartTime: o.CollectStartTime,
				EndTime:   o.CollectEndTime,
			},
		},
		PromisedDate: domain.PromisedDate{
			DateRange: domain.DateRange{
				StartDate: o.PromisedStartDate,
				EndDate:   o.PromisedEndDate,
			},
			TimeRange: domain.TimeRange{
				StartTime: o.PromisedStartTime,
				EndTime:   o.PromisedEndTime,
			},
		},
		Visits: mapVisits(visits),
	}
}

func mapPackages(packages []FlattenedPackageView) []domain.Package {
	result := make([]domain.Package, len(packages))
	for i, p := range packages {
		result[i] = domain.Package{
			ID:  p.PackageID,
			Lpn: p.Lpn,
			Dimensions: domain.Dimensions{
				Height: p.Height,
				Width:  p.Width,
				Depth:  p.Depth,
				Unit:   p.Unit,
			},
			Weight: domain.Weight{
				Value: p.WeightValue,
				Unit:  p.WeightUnit,
			},
			Insurance: domain.Insurance{
				UnitValue: p.UnitValue,
				Currency:  p.Currency,
			},
			ItemReferences: mapItemReferences(p.Items),
		}
	}
	return result
}

func mapItemReferences(items table.JSONItemReferences) []domain.ItemReference {
	result := make([]domain.ItemReference, len(items))
	for i, item := range items {
		result[i] = domain.ItemReference{
			ReferenceID: domain.ReferenceID(item.ReferenceID),
			Quantity: domain.Quantity{
				QuantityNumber: item.QuantityNumber,
				QuantityUnit:   item.QuantityUnit,
			},
		}
	}
	return result
}

func mapVisits(visits []FlattenedVisitView) []domain.Visit {
	result := make([]domain.Visit, len(visits))
	for i, v := range visits {
		result[i] = domain.Visit{
			Date: v.Date,
			TimeRange: domain.TimeRange{
				StartTime: v.TimeRangeStart,
				EndTime:   v.TimeRangeEnd,
			},
		}
	}
	return result
}

func mapDocuments(docs table.JSONReference) []domain.Document {
	result := make([]domain.Document, len(docs))
	for i, doc := range docs {
		result[i] = domain.Document{
			Value: doc.Value,
			Type:  doc.Type,
		}
	}
	return result
}

func mapJSONItems(items table.JSONItems) []domain.Item {
	result := make([]domain.Item, len(items))
	for i, item := range items {
		result[i] = domain.Item{
			ReferenceID:       domain.ReferenceID(item.ReferenceID),
			LogisticCondition: item.LogisticCondition,
			Quantity: domain.Quantity{
				QuantityNumber: item.QuantityNumber,
				QuantityUnit:   item.QuantityUnit,
			},
			Insurance: domain.Insurance{
				UnitValue: item.JSONInsurance.UnitValue,
				Currency:  item.JSONInsurance.Currency,
			},
			Description: item.Description,
			Dimensions: domain.Dimensions{
				Height: item.JSONDimensions.Height,
				Width:  item.JSONDimensions.Width,
				Depth:  item.JSONDimensions.Depth,
				Unit:   item.JSONDimensions.Unit,
			},
			Weight: domain.Weight{
				Value: item.JSONWeight.WeightValue,
				Unit:  item.JSONWeight.WeightUnit,
			},
		}
	}
	return result
}
