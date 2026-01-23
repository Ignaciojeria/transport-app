package supabaserepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

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

		// 3. Verificar si ya existe una versión con este menu_id y version_number
		var existingVersion []struct {
			ID string `json:"id"`
		}

		data, _, err = supabase.From("menu_versions").
			Select("id", "", false).
			Eq("menu_id", menuID.String()).
			Eq("version_number", strconv.Itoa(nextVersionNumber)).
			Execute()

		var versionID uuid.UUID
		if err == nil && len(data) > 0 {
			// Ya existe una versión con este menu_id y version_number
			if err := json.Unmarshal(data, &existingVersion); err == nil && len(existingVersion) > 0 {
				// Usar el ID existente
				versionID, err = uuid.Parse(existingVersion[0].ID)
				if err != nil {
					return fmt.Errorf("invalid existing version_id format: %w", err)
				}
			}
		}

		// 4. Si no existe, crear un nuevo version_id
		// Usar version_id del contexto si existe y no hay versión existente, sino crear uno nuevo
		if versionID == (uuid.UUID{}) {
			if versionIDStr, ok := sharedcontext.VersionIDFromContext(ctx); ok && versionIDStr != "" {
				var err error
				versionID, err = uuid.Parse(versionIDStr)
				if err != nil {
					return fmt.Errorf("invalid version_id format from context: %w", err)
				}
			} else {
				versionID = uuid.New()
			}
		}
		
		versionRecord := map[string]interface{}{
			"id":             versionID.String(),
			"menu_id":        menuID.String(),
			"version_number": nextVersionNumber,
			"content":        json.RawMessage(menuContentBytes),
		}

		// 5. Hacer upsert usando id como clave primaria
		// Si ya existe una versión con ese id, se actualizará
		// Si no existe, se insertará nueva
		_, _, err = supabase.From("menu_versions").
			Upsert(versionRecord, "id", "", "").
			Execute()

		if err != nil {
			return fmt.Errorf("error upserting menu version: %w", err)
		}

		// 5. Verificar si el menú tiene current_version_id
		// Si no tiene (es null o vacío), actualizarlo con la nueva versión
		var currentMenu []struct {
			CurrentVersionID *string `json:"current_version_id"`
		}

		data, _, err = supabase.From("menus").
			Select("current_version_id", "", false).
			Eq("id", menuID.String()).
			Execute()

		if err != nil {
			return fmt.Errorf("error querying menu for current_version_id: %w", err)
		}

		shouldUpdateCurrentVersion := isNewMenu
		if !isNewMenu && len(data) > 0 {
			if err := json.Unmarshal(data, &currentMenu); err == nil && len(currentMenu) > 0 {
				// Si current_version_id es null o vacío, actualizarlo
				if currentMenu[0].CurrentVersionID == nil || *currentMenu[0].CurrentVersionID == "" {
					shouldUpdateCurrentVersion = true
				}
			}
		}

		// 6. Actualizar current_version_id si es necesario
		// Para menús nuevos o menús existentes sin current_version_id
		if shouldUpdateCurrentVersion {
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
