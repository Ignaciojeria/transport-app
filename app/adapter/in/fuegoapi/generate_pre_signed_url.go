package fuegoapi

import (
	"fmt"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/shared/infrastructure/storj"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
)

func init() {
	ioc.Registry(generatePreSignedUrl, httpserver.New, storjbucket.NewTransportAppBucket)
}

func generatePreSignedUrl(s httpserver.Server, storjManager storj.UplinkManager) {
	fuego.Post(s.Manager, "/images/pre-signed-url",
		func(c fuego.ContextWithBody[request.GeneratePreSignedUrlRequest]) (response.GeneratePreSignedUrlResponse, error) {
			ctx := c.Context()
			req, err := c.Body()
			if err != nil {
				return response.GeneratePreSignedUrlResponse{}, fmt.Errorf("failed to parse request body: %w", err)
			}

			var resp response.GeneratePreSignedUrlResponse

			// TTL para upload: 48 horas (realista para usuarios)
			uploadTTL := 48 * time.Hour
			
			// TTL para download: 10 años (prácticamente sin expiración para reportería)
			downloadTTL := 10 * 365 * 24 * time.Hour
			
			uploadExpiresAt := time.Now().Add(uploadTTL)
			downloadExpiresAt := time.Now().Add(downloadTTL)

			// Preparar todas las claves de objeto para generación batch
			var objectKeys []string
			var urlInfos []struct {
				ContentType string
				ObjectKey   string
			}

			for _, uploadUrl := range req.UploadUrls {
				for i := 0; i < uploadUrl.Count; i++ {
					// Generar un nombre único para cada archivo
					fileID := uuid.New()
					objectKey := fmt.Sprintf("orders/%s/lpn/%s/%s_%d_%s.%s",
						req.OrderReferenceID,
						req.Lpn,
						uploadUrl.Type,
						i+1,
						fileID.String()[:8],
						getFileExtension(uploadUrl.ContentType))

					objectKeys = append(objectKeys, objectKey)
					urlInfos = append(urlInfos, struct {
						ContentType string
						ObjectKey   string
					}{
						ContentType: uploadUrl.ContentType,
						ObjectKey:   objectKey,
					})
				}
			}

			// Generar URLs de upload (1 hora)
			uploadURLs, err := storjManager.GeneratePreSignedURLsBatch(ctx, objectKeys, uploadTTL)
			if err != nil {
				return response.GeneratePreSignedUrlResponse{}, fmt.Errorf("failed to generate upload URLs: %w", err)
			}

			// Generar URLs de descarga (7 días) - en batch para mejor rendimiento
			var downloadObjectKeys []string
			for _, urlInfo := range urlInfos {
				downloadObjectKeys = append(downloadObjectKeys, urlInfo.ObjectKey)
			}
			
			downloadURLs := make([]string, len(downloadObjectKeys))
			for i, objectKey := range downloadObjectKeys {
				downloadURL, err := storjManager.GeneratePublicDownloadURL(ctx, objectKey, downloadTTL)
				if err != nil {
					return response.GeneratePreSignedUrlResponse{}, fmt.Errorf("failed to generate download URL: %w", err)
				}
				downloadURLs[i] = downloadURL
			}

			// Construir respuesta con las URLs generadas
			for i, urlInfo := range urlInfos {

				urlEntry := struct {
					UploadURL        string `json:"uploadUrl" example:"https://gateway.storjshare.io/bucket/file.jpg?signature=upload"`
					DownloadURL      string `json:"downloadUrl" example:"https://gateway.storjshare.io/bucket/file.jpg?signature=download"`
					Method           string `json:"method" example:"PUT"`
					ContentType      string `json:"contentType" example:"image/webp"`
					UploadExpiresAt  string `json:"uploadExpiresAt" example:"2025-08-30T12:00:00Z"`
					DownloadExpiresAt string `json:"downloadExpiresAt" example:"2035-08-28T12:00:00Z"`
					FileKey          string `json:"fileKey" example:"orders/123/lpn/456/delivery_1_abc123.webp"`
				}{
					UploadURL:        uploadURLs[i],
					DownloadURL:      downloadURLs[i],
					Method:           "PUT", 
					ContentType:      urlInfo.ContentType,
					UploadExpiresAt:  uploadExpiresAt.Format(time.RFC3339),
					DownloadExpiresAt: downloadExpiresAt.Format(time.RFC3339),
					FileKey:          urlInfo.ObjectKey,
				}
				
				resp.UploadUrls = append(resp.UploadUrls, urlEntry)
			}

			return resp, nil
		}, option.Summary("generate pre signed url"), option.Tags("orders"))
}

func getFileExtension(contentType string) string {
	switch contentType {
	case "image/webp":
		return "webp"
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	default:
		return "bin"
	}
}
