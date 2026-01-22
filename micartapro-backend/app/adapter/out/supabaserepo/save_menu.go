package supabaserepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/supabasecli"
	"micartapro/app/shared/sharedcontext"

	"github.com/google/uuid"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/supabase-community/supabase-go"
)

type SaveMenu func(ctx context.Context, menu events.MenuCreateRequest) error

func init() {
	ioc.Registry(NewSaveMenu, supabasecli.NewSupabaseClient)
}

func NewSaveMenu(supabase *supabase.Client) SaveMenu {
	return func(ctx context.Context, menu events.MenuCreateRequest) error {
		// Obtener user_id del contexto
		userIDStr, ok := sharedcontext.UserIDFromContext(ctx)
		if !ok || userIDStr == "" {
			return errors.New("user_id is required but not found in context")
		}

		// Validar y convertir el ID del menú a UUID
		menuID, err := uuid.Parse(menu.ID)
		if err != nil {
			return fmt.Errorf("invalid menu ID format: %w", err)
		}

		// Validar user_id como UUID
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return fmt.Errorf("invalid user_id format: %w", err)
		}

		// Obtener todas las versiones del menú para determinar el siguiente número de versión
		var existingVersions []struct {
			VersionNumber int `json:"version_number"`
		}

		data, _, err := supabase.From("menu_versions").
			Select("version_number", "", false).
			Eq("menu_id", menuID.String()).
			Execute()

		if err != nil && err.Error() != "PGRST116" && err.Error() != "no rows in result set" {
			return fmt.Errorf("error querying menu versions: %w", err)
		}

		// Verificar si el menú ya existe consultando si tiene versiones
		isNewMenu := len(data) == 0
		nextVersionNumber := 1
		
		if !isNewMenu {
			if err := json.Unmarshal(data, &existingVersions); err == nil {
				// Encontrar la versión máxima
				maxVersion := 0
				for _, v := range existingVersions {
					if v.VersionNumber > maxVersion {
						maxVersion = v.VersionNumber
					}
				}
				if maxVersion > 0 {
					nextVersionNumber = maxVersion + 1
				}
			}
		}

		// 1. Primero hacer upsert en la tabla menus (sin current_version_id todavía)
		// Esto asegura que el menú existe antes de insertar la versión
		menuRecord := map[string]interface{}{
			"id":      menuID.String(),
			"user_id": userID.String(),
		}

		// Upsert usando el ID como clave única
		_, _, err = supabase.From("menus").
			Upsert(menuRecord, "id", "", "").
			Execute()

		if err != nil {
			return fmt.Errorf("error upserting menu: %w", err)
		}

		// 2. Convertir el contenido del menú a JSONB
		menuContentBytes, err := json.Marshal(menu)
		if err != nil {
			return fmt.Errorf("error marshaling menu content: %w", err)
		}

		// 3. Crear nueva versión del menú
		// Usar version_id del contexto si existe, sino crear uno nuevo
		var versionID uuid.UUID
		if versionIDStr, ok := sharedcontext.VersionIDFromContext(ctx); ok && versionIDStr != "" {
			var err error
			versionID, err = uuid.Parse(versionIDStr)
			if err != nil {
				return fmt.Errorf("invalid version_id format from context: %w", err)
			}
		} else {
			versionID = uuid.New()
		}
		
		versionRecord := map[string]interface{}{
			"id":             versionID.String(),
			"menu_id":        menuID.String(),
			"version_number": nextVersionNumber,
			"content":        json.RawMessage(menuContentBytes),
		}

		// 4. Insertar la nueva versión usando Upsert con la clave única (menu_id, version_number)
		// Esto asegura que no se dupliquen versiones
		_, _, err = supabase.From("menu_versions").
			Upsert(versionRecord, "menu_id,version_number", "", "").
			Execute()

		if err != nil {
			return fmt.Errorf("error upserting menu version: %w", err)
		}

		// 5. Solo actualizar current_version_id si es un menú nuevo (primera versión)
		// Para menús existentes, la versión quedará como borrador hasta que se confirme con otro endpoint
		if isNewMenu {
			updateRecord := map[string]interface{}{
				"current_version_id": versionID.String(),
			}

			_, _, err = supabase.From("menus").
				Update(updateRecord, "", "").
				Eq("id", menuID.String()).
				Execute()

			if err != nil {
				return fmt.Errorf("error updating menu current_version_id: %w", err)
			}
		}

		return nil
	}
}
