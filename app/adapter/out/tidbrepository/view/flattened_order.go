package views

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

type FlattenedOrder struct {
	OrderID                     int64               `gorm:"order_id"`
	ReferenceID                 string              `gorm:"reference_id"`
	OrganizationCountry         string              `gorm:"organization_country"`
	CommerceName                string              `gorm:"commerce_name"`
	ConsumerName                string              `gorm:"consumer_name"`
	OrderStatus                 string              `gorm:"order_status"`
	OrderType                   string              `gorm:"order_type"`
	DeliveryInstructions        string              `gorm:"delivery_instructions"`
	OriginContactName           string              `gorm:"origin_contact_name"`
	OriginContactPhone          string              `gorm:"origin_contact_phone"`
	OriginContactDocuments      table.JSONReference `gorm:"origin_contact_documents"`
	DestinationContactName      string              `gorm:"destination_contact_name"`
	DestinationContactPhone     string              `gorm:"destination_contact_phone"`
	DestinationContactDocuments table.JSONReference `gorm:"destination_contact_documents"`
	OriginAddressLine1          string              `gorm:"origin_address_line1"`
	OriginAddressLine2          string              `gorm:"origin_address_line2"`
	OriginAddressLine3          string              `gorm:"origin_address_line3"`
	OriginState                 string              `gorm:"origin_state"`
	OriginProvince              string              `gorm:"origin_province"`
	OriginCounty                string              `gorm:"origin_county"`
	OriginDistrict              string              `gorm:"origin_district"`
	OriginZipCode               string              `gorm:"origin_zipcode"`
	OriginLatitude              float32             `gorm:"origin_latitude"`
	OriginLongitude             float32             `gorm:"origin_longitude"`
	OriginTimeZone              string              `gorm:"origin_timezone"`
	OriginNodeReferenceID       string              `gorm:"origin_node_reference_id"`
	OriginNodeName              string              `gorm:"origin_node_name"`
	OriginNodeType              string              `gorm:"origin_node_type"`
	OriginNodeOperatorName      string              `gorm:"origin_node_operator_name"`
	DestinationAddressLine1     string              `gorm:"destination_address_line1"`
	DestinationAddressLine2     string              `gorm:"destination_address_line2"`
	DestinationAddressLine3     string              `gorm:"destination_address_line3"`
	DestinationState            string              `gorm:"destination_state"`
	DestinationProvince         string              `gorm:"destination_province"`
	DestinationCounty           string              `gorm:"destination_county"`
	DestinationDistrict         string              `gorm:"destination_district"`
	DestinationZipCode          string              `gorm:"destination_zipcode"`
	DestinationLatitude         float32             `gorm:"destination_latitude"`
	DestinationLongitude        float32             `gorm:"destination_longitude"`
	DestinationTimeZone         string              `gorm:"destination_timezone"`
	DestinationNodeReferenceID  string              `gorm:"destination_node_reference_id"`
	DestinationNodeName         string              `gorm:"destination_node_name"`
	DestinationNodeType         string              `gorm:"destination_node_type"`
	DestinationNodeOperatorName string              `gorm:"destination_node_operator_name"`
	Packages                    []FlattenedPackage  `gorm:"packages"`
	Items                       table.JSONItems     `gorm:"items"`
	CollectAvailabilityDate     string              `gorm:"collect_availability_date"`
	CollectStartTime            string              `gorm:"collect_start_time"`
	CollectEndTime              string              `gorm:"collect_end_time"`
	PromisedStartDate           string              `gorm:"promised_start_date"`
	PromisedEndDate             string              `gorm:"promised_end_date"`
	PromisedStartTime           string              `gorm:"promised_start_time"`
	PromisedEndTime             string              `gorm:"promised_end_time"`
	TransportRequirements       []string            `gorm:"transport_requirements"`
	Visits                      []FlattenedVisit    `gorm:"visits"`
}

type FlattenedPackage struct {
	Lpn         string  `gorm:"lpn"`
	Height      float64 `gorm:"height"`
	Width       float64 `gorm:"width"`
	Depth       float64 `gorm:"depth"`
	Unit        string  `gorm:"unit"`
	WeightValue float64 `gorm:"weight_value"`
	WeightUnit  string  `gorm:"weight_unit"`
	Description string  `gorm:"description"`
	Quantity    int     `gorm:"quantity"`
	UnitValue   float64 `gorm:"unit_value"`
	Currency    string  `gorm:"currency"`
}

type FlattenedVisit struct {
	Date           string `gorm:"date"`
	TimeRangeStart string `gorm:"time_range_start"`
	TimeRangeEnd   string `gorm:"time_range_end"`
}

func (o FlattenedOrder) Map() domain.Order {
	return domain.Order{
		ID:          o.OrderID,
		ReferenceID: domain.ReferenceID(o.ReferenceID),
		BusinessIdentifiers: domain.BusinessIdentifiers{
			Commerce: o.CommerceName,
			Consumer: o.ConsumerName,
		},
		OrderStatus: domain.OrderStatus{
			Status: o.OrderStatus,
		},
		OrderType: domain.OrderType{
			Type: o.OrderType,
		},
		Origin: domain.Origin{
			NodeInfo: domain.NodeInfo{
				ReferenceID: domain.ReferenceID(o.OriginNodeReferenceID),
				Name:        &o.OriginNodeName,
				Type:        o.OriginNodeType,
			},
			AddressInfo: domain.AddressInfo{
				Contact: domain.Contact{
					FullName:  o.OriginContactName,
					Phone:     o.OriginContactPhone,
					Documents: mapDocuments(o.OriginContactDocuments),
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
		Destination: domain.Destination{
			DeliveryInstructions: o.DeliveryInstructions,
			NodeInfo: domain.NodeInfo{
				ReferenceID: domain.ReferenceID(o.DestinationNodeReferenceID),
				Name:        &o.DestinationNodeName,
				Type:        o.DestinationNodeType,
			},
			AddressInfo: domain.AddressInfo{
				Contact: domain.Contact{
					FullName:  o.DestinationContactName,
					Phone:     o.DestinationContactPhone,
					Documents: mapDocuments(o.DestinationContactDocuments),
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
		Packages: mapPackages(o.Packages),
		Items:    mapJSONItems(o.Items), // Assuming domain and flattened items are compatible
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
		TransportRequirements: mapTransportRequirements(o.TransportRequirements),
		Visits:                mapVisits(o.Visits),
	}
}

func mapJSONItems(flattenedItems table.JSONItems) []domain.Item {
	items := make([]domain.Item, len(flattenedItems))
	for i, fi := range flattenedItems {
		items[i] = domain.Item{
			ReferenceID:       domain.ReferenceID(fi.ReferenceID),
			LogisticCondition: fi.LogisticCondition,
			Quantity: domain.Quantity{
				QuantityNumber: fi.Quantity.QuantityNumber,
				QuantityUnit:   fi.Quantity.QuantityUnit,
			},
			Insurance: domain.Insurance{
				UnitValue: fi.Insurance.UnitValue,
				Currency:  fi.Insurance.Currency,
			},
			Description: fi.Description,
			Dimensions: domain.Dimensions{
				Height: fi.Dimensions.Height,
				Width:  fi.Dimensions.Width,
				Depth:  fi.Dimensions.Depth,
				Unit:   fi.Dimensions.Unit,
			},
			Weight: domain.Weight{
				Value: fi.Weight.Value,
				Unit:  fi.Weight.Unit,
			},
		}
	}
	return items
}

func mapPackages(flattenedPackages []FlattenedPackage) []domain.Package {
	packages := make([]domain.Package, len(flattenedPackages))
	for i, fp := range flattenedPackages {
		packages[i] = domain.Package{
			Lpn: fp.Lpn,
			Dimensions: domain.Dimensions{
				Height: fp.Height,
				Width:  fp.Width,
				Depth:  fp.Depth,
				Unit:   fp.Unit,
			},
			Weight: domain.Weight{
				Value: fp.WeightValue,
				Unit:  fp.WeightUnit,
			},
			ItemReferences: []domain.ItemReference{ // Assuming the item references are in flattened form
				{
					Quantity: domain.Quantity{
						QuantityNumber: fp.Quantity,
					},
				},
			},
			Insurance: domain.Insurance{
				UnitValue: fp.UnitValue,
				Currency:  fp.Currency,
			},
		}
	}
	return packages
}

func mapDocuments(flattenedDocs table.JSONReference) []domain.Document {
	docs := make([]domain.Document, len(flattenedDocs))
	for i, doc := range flattenedDocs {
		docs[i] = domain.Document{
			Value: doc.Value,
			Type:  doc.Type,
		}
	}
	return docs
}

func mapTransportRequirements(requirements []string) []domain.Reference {
	refs := make([]domain.Reference, len(requirements))
	for i, req := range requirements {
		refs[i] = domain.Reference{
			Value: req,
		}
	}
	return refs
}

func mapVisits(flattenedVisits []FlattenedVisit) []domain.Visit {
	visits := make([]domain.Visit, len(flattenedVisits))
	for i, fv := range flattenedVisits {
		visits[i] = domain.Visit{
			Date: fv.Date,
			TimeRange: domain.TimeRange{
				StartTime: fv.TimeRangeStart,
				EndTime:   fv.TimeRangeEnd,
			},
		}
	}
	return visits
}
