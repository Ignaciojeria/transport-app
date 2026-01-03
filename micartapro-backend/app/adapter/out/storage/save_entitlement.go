package storage

import (
	"context"
	"encoding/json"

	"micartapro/app/domain"
	"micartapro/app/shared/infrastructure/gcs"
	"micartapro/app/shared/infrastructure/observability"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SaveEntitlement func(ctx context.Context, entitlement domain.Entitlement) error

func init() {
	ioc.Registry(NewSaveEntitlement,
		observability.NewObservability,
		gcs.NewClient)
}

func NewSaveEntitlement(obs observability.Observability, gcs *storage.Client) SaveEntitlement {
	return func(ctx context.Context, entitlement domain.Entitlement) error {
		obs.Logger.InfoContext(ctx, "save_entitlement", "entitlement", entitlement, "userID", entitlement.UserID)
		bucket := gcs.Bucket("micartapro-menus")
		objectPath := entitlement.UserID + "/entitlement.json"
		object := bucket.Object(objectPath)
		// Upsert: si ya existe, lo sobrescribe.
		writer := object.NewWriter(ctx)
		writer.ContentType = "application/json"
		writer.CacheControl = "no-cache, max-age=0"

		if err := json.NewEncoder(writer).Encode(entitlement); err != nil {
			writer.Close()
			obs.Logger.ErrorContext(ctx, "error_encoding_entitlement", "error", err, "objectPath", objectPath)
			return err
		}

		if err := writer.Close(); err != nil {
			obs.Logger.ErrorContext(ctx, "error_closing_writer", "error", err, "objectPath", objectPath)
			return err
		}

		obs.Logger.InfoContext(ctx, "entitlement_saved_successfully", "userID", entitlement.UserID)
		return nil
	}
}
