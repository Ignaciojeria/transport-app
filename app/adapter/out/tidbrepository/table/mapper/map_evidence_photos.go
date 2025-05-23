package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapEvidencePhotosTable(ctx context.Context, photos domain.EvidencePhotos) table.JSONEvidencePhotos {
	result := make(table.JSONEvidencePhotos, len(photos))
	for i, photo := range photos {
		result[i] = table.EvidencePhoto{
			URL:     photo.URL,
			Type:    photo.Type,
			TakenAt: &photo.TakenAt,
		}
	}
	return result
}
