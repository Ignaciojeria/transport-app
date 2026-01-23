package imageuploader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/infrastructure/supabasecli"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	supabase "github.com/supabase-community/supabase-go"
)

type UploadImage func(ctx context.Context, imageBytes []byte, fileName string) (string, error)

func init() {
	ioc.Registry(NewImageUploader, observability.NewObservability, configuration.NewConf, supabasecli.NewSupabaseClient)
}

func NewImageUploader(obs observability.Observability, conf configuration.Conf, supabaseClient *supabase.Client) (UploadImage, error) {
	return func(ctx context.Context, imageBytes []byte, fileName string) (string, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "upload_image")
		defer span.End()

		userID, ok := sharedcontext.UserIDFromContext(spanCtx)
		if !ok || userID == "" {
			return "", fmt.Errorf("userID is required but not found in context")
		}

		obs.Logger.InfoContext(spanCtx, "uploading_image", "fileName", fileName, "size_bytes", len(imageBytes), "userID", userID)

		// Construir la ruta del archivo: userID/filename (igual que en micartapro-console)
		// El bucket ya se llama "menu-photos", no necesitamos incluirlo en el path
		filePath := fmt.Sprintf("%s/%s", userID, fileName)

		// Subir la imagen al bucket "menu-photos" usando HTTP directo
		// Esto nos da control total sobre los headers de autenticación, igual que en save_menu_order.go
		bucketName := "menu-photos"
		uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", conf.SUPABASE_PROJECT_URL, bucketName, filePath)

		req, err := http.NewRequestWithContext(spanCtx, "POST", uploadURL, bytes.NewReader(imageBytes))
		if err != nil {
			return "", fmt.Errorf("error creating request: %w", err)
		}

		// Configurar headers de autenticación igual que en save_menu_order.go
		// Esto funciona para las tablas, así que debería funcionar para Storage también
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("apikey", conf.SUPABASE_BACKEND_API_KEY)
		req.Header.Set("Authorization", "Bearer "+conf.SUPABASE_BACKEND_API_KEY)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_uploading_image", "error", err, "fileName", fileName)
			return "", fmt.Errorf("error uploading image to Supabase Storage: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error reading response: %w", err)
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			obs.Logger.ErrorContext(spanCtx, "error_uploading_image", "status", resp.StatusCode, "body", string(body), "fileName", fileName)
			return "", fmt.Errorf("error uploading image to Supabase Storage: status %d, body: %s", resp.StatusCode, string(body))
		}

		// Construir la URL pública de la imagen
		// La URL pública de Supabase Storage sigue el formato:
		// https://<project-ref>.supabase.co/storage/v1/object/public/<bucket>/<path>
		publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s",
			conf.SUPABASE_PROJECT_URL,
			bucketName,
			filePath)

		obs.Logger.InfoContext(spanCtx, "image_uploaded_successfully", "publicURL", publicURL, "fileName", fileName, "status", resp.StatusCode, "body", string(body))

		return publicURL, nil
	}, nil
}
