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

type SaveMenu func(ctx context.Context, menu domain.MenuCreateRequest) error

func init() {
	ioc.Registry(NewSaveMenu,
		observability.NewObservability,
		gcs.NewClient)
}

func NewSaveMenu(obs observability.Observability, gcs *storage.Client) SaveMenu {
	return func(ctx context.Context, menu domain.MenuCreateRequest) error {
		obs.Logger.Info("save_menu", "menu", menu)
		bucket := gcs.Bucket("micartapro-menus")
		objectPath := "menus/" + menu.ID + ".json"
		object := bucket.Object(objectPath)
		// Upsert: si ya existe, lo sobrescribe.
		writer := object.NewWriter(ctx)
		// Con "Uniform bucket-level access" (UBLA) activado, las ACLs legacy están deshabilitadas.
		// Si necesitas lectura pública (sin autenticación), se debe configurar IAM a nivel bucket
		// (por ejemplo: allUsers -> roles/storage.objectViewer).
		writer.ContentType = "application/json"
		// Evita que caches (navegador/CDN) sirvan una versión vieja tras un overwrite.
		// Ajusta si quieres cachear (por ejemplo, public,max-age=60).
		writer.CacheControl = "no-cache, max-age=0"

		if err := json.NewEncoder(writer).Encode(menu); err != nil {
			writer.Close()
			obs.Logger.Error("error_encoding_menu", "error", err, "objectPath", objectPath)
			return err
		}

		if err := writer.Close(); err != nil {
			obs.Logger.Error("error_closing_writer", "error", err, "objectPath", objectPath)
			return err
		}

		obs.Logger.Info("menu_saved_successfully", "menuID", menu.ID)
		return nil
	}
}
