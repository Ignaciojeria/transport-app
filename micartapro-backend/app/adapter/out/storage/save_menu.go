package storage

import (
	"context"
	"encoding/json"

	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/ioc"
)

type SaveMenu func(ctx context.Context, menu events.MenuCreateRequest) error

func init() {
	ioc.Register(NewSaveMenu)
}

func NewSaveMenu(obs observability.Observability, gcs *storage.Client) SaveMenu {
	return func(ctx context.Context, menu events.MenuCreateRequest) error {
		idempotencyKey, _ := sharedcontext.IdempotencyKeyFromContext(ctx)
		userID, ok := sharedcontext.UserIDFromContext(ctx)
		if !ok || userID == "" {
			obs.Logger.Error("user_id_not_found_in_context", "error", "userID is required but not found in context")
			return nil // o podrías retornar un error si prefieres
		}
		obs.Logger.InfoContext(ctx, "save_menu", "menu", menu, "idempotencyKey", idempotencyKey, "userID", userID)
		bucket := gcs.Bucket("micartapro-menus")
		objectPath := userID + "/menus/" + menu.ID + "/" + idempotencyKey + ".json"
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

		// Actualizar latest.json para apuntar al último archivo guardado
		latestPath := userID + "/menus/" + menu.ID + "/latest.json"
		latestObject := bucket.Object(latestPath)
		latestWriter := latestObject.NewWriter(ctx)
		latestWriter.ContentType = "application/json"
		latestWriter.CacheControl = "no-cache, max-age=0"

		latestData := map[string]string{"filename": idempotencyKey + ".json"}
		if err := json.NewEncoder(latestWriter).Encode(latestData); err != nil {
			latestWriter.Close()
			obs.Logger.Error("error_encoding_latest", "error", err, "latestPath", latestPath)
			return err
		}

		if err := latestWriter.Close(); err != nil {
			obs.Logger.Error("error_closing_latest_writer", "error", err, "latestPath", latestPath)
			return err
		}

		obs.Logger.Info("menu_saved_successfully", "menuID", menu.ID)
		return nil
	}
}
