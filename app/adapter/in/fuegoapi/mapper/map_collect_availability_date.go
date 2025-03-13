package mapper

import "transport-app/app/domain"

func MapCollectAvailabilityDateToDomain(collectDate struct {
	Date      string `json:"date"`
	TimeRange struct {
		EndTime   string `json:"endTime"`
		StartTime string `json:"startTime"`
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

func MapCollectAvailabilityDateFromDomain(collectDate domain.CollectAvailabilityDate) struct {
	Date      string `json:"date"`
	TimeRange struct {
		EndTime   string `json:"endTime"`
		StartTime string `json:"startTime"`
	} `json:"timeRange"`
} {
	response := struct {
		Date      string `json:"date"`
		TimeRange struct {
			EndTime   string `json:"endTime"`
			StartTime string `json:"startTime"`
		} `json:"timeRange"`
	}{
		Date: collectDate.Date.Format("2006-01-02"),
	}
	response.TimeRange.StartTime = collectDate.TimeRange.StartTime
	response.TimeRange.EndTime = collectDate.TimeRange.EndTime
	return response
}
