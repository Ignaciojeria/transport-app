package storage

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/gcs"
	"micartapro/app/shared/infrastructure/observability"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/api/iterator"
)

var ErrMenuNotFound = errors.New("menu not found")

type GetLatestMenuById func(ctx context.Context, menuID string) (events.MenuCreateRequest, error)

func init() {
	ioc.Registry(NewGetLatestMenuById,
		observability.NewObservability,
		gcs.NewClient)
}

func NewGetLatestMenuById(obs observability.Observability, gcs *storage.Client) GetLatestMenuById {
	return func(ctx context.Context, menuID string) (events.MenuCreateRequest, error) {
		obs.Logger.InfoContext(ctx, "get_latest_menu_by_id", "menuID", menuID)
		bucket := gcs.Bucket("micartapro-menus")
		prefix := "menus/" + menuID + "/"
		latestPath := prefix + "latest.json"

		// Intentar leer latest.json para obtener el nombre del último archivo
		var latestObjectName string
		latestObject := bucket.Object(latestPath)
		latestReader, err := latestObject.NewReader(ctx)
		if err == nil {
			defer latestReader.Close()
			var latestData map[string]string
			if err := json.NewDecoder(latestReader).Decode(&latestData); err == nil {
				if filename, ok := latestData["filename"]; ok && filename != "" {
					latestObjectName = filename
					obs.Logger.InfoContext(ctx, "found_latest_pointer", "filename", filename)
				}
			}
		}

		// Si no se pudo leer latest.json, hacer fallback al listado
		if latestObjectName == "" {
			obs.Logger.InfoContext(ctx, "latest_json_not_found_fallback_to_listing", "prefix", prefix)
			query := &storage.Query{
				Prefix: prefix,
			}
			it := bucket.Objects(ctx, query)

			// Encontrar el último objeto (UUIDv7 más reciente lexicográficamente)
			for {
				attrs, err := it.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					obs.Logger.ErrorContext(ctx, "error_listing_objects", "error", err, "prefix", prefix)
					return events.MenuCreateRequest{}, err
				}

				// Extraer el nombre del archivo (sin la extensión .json)
				filename := strings.TrimPrefix(attrs.Name, prefix)
				if filename == "" || !strings.HasSuffix(filename, ".json") || filename == "latest.json" {
					continue
				}

				// Comparar lexicográficamente para encontrar el último UUIDv7
				if latestObjectName == "" || filename > latestObjectName {
					latestObjectName = filename
				}
			}
		}

		if latestObjectName == "" {
			obs.Logger.WarnContext(ctx, "menu_not_found", "menuID", menuID, "prefix", prefix)
			return events.MenuCreateRequest{}, ErrMenuNotFound
		}

		// Leer el último objeto encontrado
		objectPath := prefix + latestObjectName
		object := bucket.Object(objectPath)
		obs.Logger.InfoContext(ctx, "reading_latest_object", "objectPath", objectPath)

		reader, err := object.NewReader(ctx)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "error_reading_object", "error", err, "objectPath", objectPath)
			if errors.Is(err, storage.ErrObjectNotExist) {
				obs.Logger.WarnContext(ctx, "menu_not_found", "menuID", menuID, "objectPath", objectPath)
				return events.MenuCreateRequest{}, ErrMenuNotFound
			}
			return events.MenuCreateRequest{}, err
		}
		defer reader.Close()

		var menu events.MenuCreateRequest
		if err := json.NewDecoder(reader).Decode(&menu); err != nil {
			obs.Logger.ErrorContext(ctx, "error_decoding_menu", "error", err)
			return events.MenuCreateRequest{}, err
		}

		obs.Logger.Info("menu_found_successfully", "menuID", menuID, "objectPath", objectPath)
		return menu, nil
	}
}
