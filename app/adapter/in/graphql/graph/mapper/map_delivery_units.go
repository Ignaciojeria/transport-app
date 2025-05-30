package mapper

import (
	"context"
	"transport-app/app/adapter/in/graphql/graph/model"
	"transport-app/app/adapter/out/tidbrepository/projectionresult"
)

func MapDeliveryUnits(ctx context.Context, deliveryUnits []projectionresult.DeliveryUnitsProjectionResult) []*model.DeliveryUnitsReport {
	deliveryUnitsReport := make([]*model.DeliveryUnitsReport, len(deliveryUnits))

	for i, du := range deliveryUnits {
		report := &model.DeliveryUnitsReport{
			Commerce:             &du.Commerce,
			Consumer:             &du.Consumer,
			Channel:              &du.Channel,
			Status:               &du.Status,
			DeliveryInstructions: &du.OrderDeliveryInstructions,
			ReferenceID:          du.OrderReferenceID,
			GroupBy: &model.GroupBy{
				Type:  &du.OrderGroupByType,
				Value: &du.OrderGroupByValue,
			},
			CollectAvailabilityDate: &model.CollectAvailabilityDate{
				Date: &du.OrderCollectAvailabilityDate,
				TimeRange: &model.TimeRange{
					StartTime: &du.OrderCollectAvailabilityDateStartTime,
					EndTime:   &du.OrderCollectAvailabilityDateEndTime,
				},
			},
			Destination: &model.Location{
				AddressInfo: &model.AddressInfo{
					AddressLine1: &du.DestinationAddressLine1,
					AddressLine2: &du.DestinationAddressLine2,
					District:     &du.DestinationDistrict,
					Province:     &du.DestinationProvince,
					State:        &du.DestinationState,
					TimeZone:     &du.DestinationTimeZone,
					ZipCode:      &du.DestinationZipCode,
					Coordinates: &model.Coordinates{
						Latitude:  &du.DestinationCoordinatesLatitude,
						Longitude: &du.DestinationCoordinatesLongitude,
						Source:    &du.DestinationCoordinatesSource,
						Confidence: &model.Confidence{
							Level:   &du.DestinationCoordinatesConfidenceLevel,
							Message: &du.DestinationCoordinatesConfidenceMessage,
							Reason:  &du.DestinationCoordinatesConfidenceReason,
						},
					},
					Contact: &model.Contact{
						Email:      &du.DestinationContactEmail,
						FullName:   &du.DestinationContactFullName,
						NationalID: &du.DestinationContactNationalID,
						Phone:      &du.DestinationContactPhone,
					},
				},
			},
			Origin: &model.Location{
				AddressInfo: &model.AddressInfo{
					AddressLine1: &du.OriginAddressLine1,
					AddressLine2: &du.OriginAddressLine2,
					District:     &du.OriginDistrict,
					Province:     &du.OriginProvince,
					State:        &du.OriginState,
					TimeZone:     &du.OriginTimeZone,
					ZipCode:      &du.OriginZipCode,
					Coordinates: &model.Coordinates{
						Latitude:  &du.OriginCoordinatesLatitude,
						Longitude: &du.OriginCoordinatesLongitude,
						Source:    &du.OriginCoordinatesSource,
						Confidence: &model.Confidence{
							Level:   &du.OriginCoordinatesConfidenceLevel,
							Message: &du.OriginCoordinatesConfidenceMessage,
							Reason:  &du.OriginCoordinatesConfidenceReason,
						},
					},
					Contact: &model.Contact{
						Email:      &du.OriginContactEmail,
						FullName:   &du.OriginContactFullName,
						NationalID: &du.OriginContactNationalID,
						Phone:      &du.OriginContactPhone,
					},
				},
			},
			OrderType: &model.OrderType{
				Type:        &du.OrderType,
				Description: &du.OrderTypeDescription,
			},
			DeliveryUnit: &model.DeliveryUnit{
				Lpn: &du.LPN,
			},
			PromisedDate: &model.PromisedDate{
				ServiceCategory: &du.OrderPromisedDateServiceCategory,
				DateRange: &model.DateRange{
					StartDate: &du.OrderPromisedDateStartDate,
					EndDate:   &du.OrderPromisedDateEndDate,
				},
				TimeRange: &model.TimeRange{
					StartTime: &du.OrderPromisedDateStartTime,
					EndTime:   &du.OrderPromisedDateEndTime,
				},
			},
		}

		// Map references if they exist
		if du.OrderReferences != nil {
			references := make([]*model.Reference, len(du.OrderReferences))
			for j, ref := range du.OrderReferences {
				references[j] = &model.Reference{
					Type:  &ref.Type,
					Value: &ref.Value,
				}
			}
			report.References = references
		}

		// Map delivery unit labels if they exist
		if du.DeliveryUnitLabels != nil {
			labels := make([]*model.Label, len(du.DeliveryUnitLabels))
			for j, label := range du.DeliveryUnitLabels {
				labels[j] = &model.Label{
					Type:  &label.Type,
					Value: &label.Value,
				}
			}
			report.DeliveryUnit.Labels = labels
		}

		// Map extra fields if they exist
		if du.ExtraFields != nil {
			extraFields := make([]*model.KeyValuePair, 0, len(du.ExtraFields))
			for key, value := range du.ExtraFields {
				extraFields = append(extraFields, &model.KeyValuePair{
					Key:   key,
					Value: value,
				})
			}
			report.ExtraFields = extraFields
		}

		deliveryUnitsReport[i] = report
	}

	return deliveryUnitsReport
}
