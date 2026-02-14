package storage

import (
	"context"
	"fmt"

	"micartapro/app/shared/infrastructure/gcs"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

const journeyReportsBucket = "micartapro-journey-reports"

// UploadJourneyReport sube el XLSX a GCS y devuelve la URL pública.
// Requiere bucket micartapro-journey-reports en GCP (público o con IAM para lectura).
type UploadJourneyReport func(ctx context.Context, journeyID string, xlsxBytes []byte) (publicURL string, err error)

func init() {
	ioc.Registry(NewUploadJourneyReport, gcs.NewClient)
}

func NewUploadJourneyReport(gcsClient *storage.Client) UploadJourneyReport {
	return func(ctx context.Context, journeyID string, xlsxBytes []byte) (string, error) {
		objectPath := journeyID + "/report.xlsx"
		bucket := gcsClient.Bucket(journeyReportsBucket)
		object := bucket.Object(objectPath)

		writer := object.NewWriter(ctx)
		writer.ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		writer.CacheControl = "public, max-age=86400"

		if _, err := writer.Write(xlsxBytes); err != nil {
			writer.Close()
			return "", fmt.Errorf("writing report to GCS: %w", err)
		}
		if err := writer.Close(); err != nil {
			return "", fmt.Errorf("closing GCS writer: %w", err)
		}

		publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", journeyReportsBucket, objectPath)
		return publicURL, nil
	}
}
