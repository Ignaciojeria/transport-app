package views

import (
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	"github.com/biter777/countries"
	"github.com/paulmach/orb"
)

type FlattenedOrderView struct {
	OrderID                      int64               `gorm:"column:order_id"`
	ReferenceID                  string              `gorm:"column:reference_id"`
	OrganizationID               int64               `gorm:"column:organization_id"`
	OrganizationCountry          string              `gorm:"column:organization_country"`
	HeaderID                     int64               `gorm:"column:header_id"`
	CommerceName                 string              `gorm:"column:commerce_name"`
	ConsumerName                 string              `gorm:"column:consumer_name"`
	OrderStatusID                int64               `gorm:"column:order_status_id"`
	OrderStatus                  string              `gorm:"column:order_status"`
	OrderTypeID                  int64               `gorm:"column:order_type_id"`
	OrderType                    string              `gorm:"column:order_type"`
	OrderSequenceNumber          *int                `gorm:"column:order_sequence_number"`
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
	OriginLocality               string              `gorm:"column:origin_locality"`
	OriginDistrict               string              `gorm:"column:origin_district"`
	OriginZipCode                string              `gorm:"column:origin_zipcode"`
	OriginLatitude               float64             `gorm:"column:origin_latitude"`
	OriginLongitude              float64             `gorm:"column:origin_longitude"`
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
	DestinationLocality          string              `gorm:"column:destination_locality"`
	DestinationDistrict          string              `gorm:"column:destination_district"`
	DestinationZipCode           string              `gorm:"column:destination_zipcode"`
	DestinationLatitude          float64             `gorm:"column:destination_latitude"`
	DestinationLongitude         float64             `gorm:"column:destination_longitude"`
	DestinationTimeZone          string              `gorm:"column:destination_timezone"`
	DestinationNodeInfoID        int64               `gorm:"column:destination_node_info_id"`
	DestinationNodeReferenceID   string              `gorm:"column:destination_node_reference_id"`
	DestinationNodeName          string              `gorm:"column:destination_node_name"`
	DestinationNodeType          string              `gorm:"column:destination_node_type"`
	DestinationNodeOperatorID    int64               `gorm:"column:destination_node_operator_id"`
	DestinationNodeOperatorName  string              `gorm:"column:destination_node_operator_name"`
	Items                        table.JSONItems     `gorm:"column:items"`
	CollectAvailabilityDate      time.Time           `gorm:"column:collect_availability_date"`
	CollectStartTime             string              `gorm:"column:collect_start_time"`
	CollectEndTime               string              `gorm:"column:collect_end_time"`
	PromisedStartDate            time.Time           `gorm:"column:promised_start_date"`
	PromisedEndDate              time.Time           `gorm:"column:promised_end_date"`
	PromisedStartTime            string              `gorm:"column:promised_start_time"`
	PromisedEndTime              string              `gorm:"column:promised_end_time"`
	TransportRequirements        table.JSONReference `gorm:"column:transport_requirements"`
	//Campos nuevos
	PlanID                  int64                  `gorm:"column:plan_id"`
	PlanReferenceID         string                 `gorm:"column:plan_reference_id"`
	PlannedDate             time.Time              `gorm:"column:planned_date"`
	PlanStartLocation       table.JSONPlanLocation `gorm:"column:plan_start_location"`
	RouteID                 int64                  `gorm:"column:route_id"`
	RouteReferenceID        string                 `gorm:"column:route_reference_id"`
	RouteEndLocation        table.JSONPlanLocation `gorm:"column:route_end_location"`
	RouteEndNodeReferenceID string                 `gorm:"column:route_end_node_reference_id"`
	RouteAccountID          int64                  `gorm:"column:route_account_id"`
	RouteAccountReferenceID string                 `gorm:"column:route_account_reference_id"`
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

func (o FlattenedOrderView) ToOrder(packages []FlattenedPackageView, refs []FlattenedOrderReferenceView) domain.Order {
	// Mapear referencias
	references := make([]domain.Reference, len(refs))
	for i, ref := range refs {
		references[i] = domain.Reference{
			ID:    ref.ReferenceID,
			Type:  ref.Type,
			Value: ref.Value,
		}
	}

	// Mapear requisitos de transporte
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
		Headers: domain.Headers{
			ID: o.HeaderID,
			Organization: domain.Organization{
				ID:      o.OrganizationID,
				Country: countries.ByName(o.OrganizationCountry),
			},
			Consumer: o.ConsumerName,
			Commerce: o.CommerceName,
		},
		SequenceNumber:       o.OrderSequenceNumber,
		DeliveryInstructions: o.DeliveryInstructions,
		OrderStatus: domain.OrderStatus{
			ID:     o.OrderStatusID,
			Status: o.OrderStatus,
		},
		OrderType: domain.OrderType{
			ID:          o.OrderTypeID,
			Type:        o.OrderType,
			Description: o.OrderTypeDescription,
		},

		// ✅ Mapeo COMPLETO del contacto y dirección de ORIGEN
		Origin: domain.NodeInfo{
			ReferenceID: domain.ReferenceID(o.OriginNodeReferenceID),
			Name:        o.OriginNodeName,
			NodeType: domain.NodeType{
				Value: o.OriginNodeType,
			},
			Contact: domain.Contact{
				ID:           o.OriginContactID,
				FullName:     o.OriginContactName,
				PrimaryPhone: o.OriginContactPhone,
				PrimaryEmail: o.OriginContactEmail,
				NationalID:   o.OriginContactNationalID,
				Documents:    mapDocuments(o.OriginContactDocuments),
			},
			AddressInfo: domain.AddressInfo{
				AddressLine1: o.OriginAddressLine1,
				//	AddressLine2: o.OriginAddressLine2,
				//	AddressLine3: o.OriginAddressLine3,
				State:    o.OriginState,
				Province: o.OriginProvince,
				//	Locality:     o.OriginLocality,
				District: o.OriginDistrict,
				ZipCode:  o.OriginZipCode,
				Location: orb.Point{
					o.OriginLongitude,
					o.OriginLatitude,
				},
				TimeZone: o.OriginTimeZone,
			},
		},

		// ✅ Mapeo COMPLETO del contacto y dirección de DESTINO
		Destination: domain.NodeInfo{
			ReferenceID: domain.ReferenceID(o.DestinationNodeReferenceID),
			Name:        o.DestinationNodeName,
			NodeType: domain.NodeType{
				Value: o.DestinationNodeType,
			},
			Contact: domain.Contact{
				ID:           o.DestinationContactID,
				FullName:     o.DestinationContactName,
				PrimaryPhone: o.DestinationContactPhone,
				PrimaryEmail: o.DestinationContactEmail,
				NationalID:   o.DestinationContactNationalID,
				Documents:    mapDocuments(o.DestinationContactDocuments),
			},
			AddressInfo: domain.AddressInfo{
				AddressLine1: o.DestinationAddressLine1,
				//	AddressLine2: o.DestinationAddressLine2,
				//	AddressLine3: o.DestinationAddressLine3,
				State:    o.DestinationState,
				Province: o.DestinationProvince,
				//	Locality:     o.DestinationLocality,
				District: o.DestinationDistrict,
				ZipCode:  o.DestinationZipCode,
				Location: orb.Point{
					o.DestinationLongitude,
					o.DestinationLatitude,
				},
				TimeZone: o.DestinationTimeZone,
			},
		},

		// ✅ Mapeo de Plan y Ruta
		Plan: domain.Plan{
			ID:          o.PlanID,
			PlannedDate: o.PlannedDate,
			ReferenceID: o.PlanReferenceID,
			Origin: domain.NodeInfo{
				AddressInfo: domain.AddressInfo{
					Location: orb.Point{o.PlanStartLocation.Longitude, o.PlanStartLocation.Latitude},
				},
			},
			Routes: []domain.Route{
				{
					ID:          o.RouteID,
					ReferenceID: o.RouteReferenceID,
					Destination: domain.NodeInfo{
						ReferenceID: domain.ReferenceID(o.RouteEndNodeReferenceID),
						AddressInfo: domain.AddressInfo{
							Location: orb.Point{o.RouteEndLocation.Longitude, o.RouteEndLocation.Latitude},
						},
					},
					Organization: domain.Organization{
						ID: o.RouteAccountID,
					},
				},
			},
		},

		// ✅ Mapeo de paquetes, referencias e ítems
		Packages:              mapPackages(packages),
		Items:                 mapJSONItems(o.Items),
		References:            references,
		TransportRequirements: transportReqs,

		// ✅ Mapeo de disponibilidad de recolección
		CollectAvailabilityDate: domain.CollectAvailabilityDate{
			Date: o.CollectAvailabilityDate,
			TimeRange: domain.TimeRange{
				StartTime: o.CollectStartTime,
				EndTime:   o.CollectEndTime,
			},
		},

		// ✅ Mapeo de fecha prometida
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
