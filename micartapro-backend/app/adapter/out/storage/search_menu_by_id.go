package storage

import (
	"context"
	"encoding/json"
	"errors"

	"micartapro/app/domain"
	"micartapro/app/shared/infrastructure/gcs"
	"micartapro/app/shared/infrastructure/observability"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

var ErrMenuNotFound = errors.New("menu not found")

type SearchMenuById func(ctx context.Context, menuID string) (domain.MenuCreateRequest, error)

func init() {
	ioc.Registry(NewSearchMenuById,
		observability.NewObservability,
		gcs.NewClient)
}

func NewSearchMenuById(obs observability.Observability, gcs *storage.Client) SearchMenuById {
	return func(ctx context.Context, menuID string) (domain.MenuCreateRequest, error) {
		obs.Logger.Info("search_menu_by_id", "menuID", menuID)
		bucket := gcs.Bucket("micartapro-menus")
		objectPath := "menus/" + menuID + ".json"
		object := bucket.Object(objectPath)
		obs.Logger.Info("searching_object", "objectPath", objectPath)

		reader, err := object.NewReader(ctx)
		if err != nil {
			obs.Logger.Error("error_reading_object", "error", err, "objectPath", objectPath)
			if errors.Is(err, storage.ErrObjectNotExist) {
				obs.Logger.Warn("menu_not_found", "menuID", menuID, "objectPath", objectPath)
				return domain.MenuCreateRequest{}, ErrMenuNotFound
			}
			return domain.MenuCreateRequest{}, err
		}
		defer reader.Close()

		var menu domain.MenuCreateRequest
		if err := json.NewDecoder(reader).Decode(&menu); err != nil {
			obs.Logger.Error("error_decoding_menu", "error", err)
			return domain.MenuCreateRequest{}, err
		}

		obs.Logger.Info("menu_found_successfully", "menuID", menuID)
		return menu, nil
	}
}
