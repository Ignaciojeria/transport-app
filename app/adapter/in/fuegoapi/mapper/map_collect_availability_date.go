package mapper

import "transport-app/app/domain"

func MapCollectAvailabilityDateToDomain(collectDate struct {
	Date      string `json:"date" example:"2025-03-30"`
	TimeRange struct {
		EndTime   string `json:"endTime" example:"09:00"`
		StartTime string `json:"startTime" example:"19:00"`
	} `json:"timeRange"`
}) domain.CollectAvailabilityDate {
	return domain.CollectAvailabilityDate{
		Date: MapDateStringToTime(collectDate.Date),
		TimeRange: domain.TimeRange{
			StartTime: collectDate.TimeRange.StartTime,
			EndTime:   collectDate.TimeRange.EndTime,
		},
	}
}
