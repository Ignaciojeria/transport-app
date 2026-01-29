package imagegenerator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/compute/metadata"
	iamcredentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"cloud.google.com/go/storage"
	"micartapro/app/shared/infrastructure/observability"
	"golang.org/x/oauth2/google"
)

// signingCreds holds either a private key (local) or a SignBytes function (Cloud Run / metadata).
type signingCreds struct {
	email      string
	privateKey []byte
	signBytes  func([]byte) ([]byte, error)
}

// getSigningCreds obtiene credenciales para firmar URLs: primero intenta GOOGLE_APPLICATION_CREDENTIALS,
// luego Application Default Credentials (en Cloud Run usa IAM SignBlob).
func getSigningCreds(ctx context.Context, obs observability.Observability) (*signingCreds, error) {
	// 1. Intentar archivo de credenciales (desarrollo local)
	credsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsPath != "" {
		credsData, err := os.ReadFile(credsPath)
		if err == nil {
			var credsJSON map[string]interface{}
			if err := json.Unmarshal(credsData, &credsJSON); err == nil {
				email, okE := credsJSON["client_email"].(string)
				key, okK := credsJSON["private_key"].(string)
				if okE && email != "" && okK && key != "" {
					obs.Logger.InfoContext(ctx, "using_credentials_file", "path", credsPath)
					return &signingCreds{email: email, privateKey: []byte(key)}, nil
				}
			}
		}
	}

	// 2. Application Default Credentials (Cloud Run, GCE, etc.)
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, fmt.Errorf("no credentials: %w", err)
	}

	// Si ADC incluye JSON (p. ej. desde otro origen), usar client_email y private_key
	if len(creds.JSON) > 0 {
		var credsJSON map[string]interface{}
		if err := json.Unmarshal(creds.JSON, &credsJSON); err == nil {
			email, okE := credsJSON["client_email"].(string)
			key, okK := credsJSON["private_key"].(string)
			if okE && email != "" && okK && key != "" {
				obs.Logger.InfoContext(ctx, "using_adc_json_credentials")
				return &signingCreds{email: email, privateKey: []byte(key)}, nil
			}
		}
	}

	// 3. Metadata server (Cloud Run / GCE): obtener email y usar IAM SignBlob.
	// En Cloud Run la cuenta de servicio debe tener el rol roles/iam.serviceAccountTokenCreator sobre sí misma.
	// Ejemplo (SA de este proyecto):
	//   gcloud iam service-accounts add-iam-policy-binding transport-app-sa@einar-404623.iam.gserviceaccount.com --member="serviceAccount:transport-app-sa@einar-404623.iam.gserviceaccount.com" --role="roles/iam.serviceAccountTokenCreator"
	if !metadata.OnGCE() {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS not set and not running on GCP (need credentials to sign URLs)")
	}
	email, err := metadata.Email("default")
	if err != nil || email == "" {
		return nil, fmt.Errorf("could not get service account email from metadata: %w", err)
	}
	obs.Logger.InfoContext(ctx, "using_metadata_and_iam_signblob", "email", email)

	iamClient, err := iamcredentials.NewIamCredentialsClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating IAM credentials client: %w", err)
	}
	defer iamClient.Close()

	// Nombre del recurso: projects/-/serviceAccounts/EMAIL (projects/- usa el proyecto por defecto)
	resourceName := fmt.Sprintf("projects/-/serviceAccounts/%s", email)
	signBytes := func(payload []byte) ([]byte, error) {
		resp, err := iamClient.SignBlob(ctx, &credentialspb.SignBlobRequest{
			Name:    resourceName,
			Payload: payload,
		})
		if err != nil {
			return nil, err
		}
		return resp.SignedBlob, nil
	}
	return &signingCreds{email: email, signBytes: signBytes}, nil
}

// applySigningCreds aplica creds a opts (GoogleAccessID y PrivateKey o SignBytes).
func applySigningCreds(opts *storage.SignedURLOptions, creds *signingCreds) {
	opts.GoogleAccessID = creds.email
	if len(creds.privateKey) > 0 {
		opts.PrivateKey = creds.privateKey
	} else {
		opts.SignBytes = creds.signBytes
	}
}

// FillSignedURLOptions obtiene credenciales (archivo, ADC o IAM SignBlob en Cloud Run) y las aplica a opts.
// Útil para que otros paquetes (p. ej. fuegoapi) generen signed URLs sin duplicar lógica.
func FillSignedURLOptions(ctx context.Context, obs observability.Observability, opts *storage.SignedURLOptions) error {
	creds, err := getSigningCreds(ctx, obs)
	if err != nil {
		return err
	}
	applySigningCreds(opts, creds)
	return nil
}

// GenerateSignedReadURL genera una signed URL de lectura para una imagen en GCS
// Si la URL ya es pública o no es de GCS, la retorna sin modificar
func GenerateSignedReadURL(ctx context.Context, obs observability.Observability, imageURL string) (string, error) {
	// Si la URL no es de GCS (storage.googleapis.com), retornarla sin modificar
	if !strings.Contains(imageURL, "storage.googleapis.com") {
		obs.Logger.InfoContext(ctx, "url_not_gcs", "url", imageURL)
		return imageURL, nil
	}

	// Extraer bucket y object path de la URL
	// Formato: https://storage.googleapis.com/<bucket>/<object-path>
	parts := strings.Split(imageURL, "storage.googleapis.com/")
	if len(parts) != 2 {
		obs.Logger.WarnContext(ctx, "invalid_gcs_url_format", "url", imageURL)
		return imageURL, nil // Retornar URL original si no podemos parsearla
	}

	pathParts := strings.SplitN(parts[1], "/", 2)
	if len(pathParts) != 2 {
		obs.Logger.WarnContext(ctx, "invalid_gcs_path_format", "url", imageURL)
		return imageURL, nil
	}

	bucketName := pathParts[0]
	objectPath := pathParts[1]

	// Generar signed URL para lectura (válida por 1 hora)
	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(1 * time.Hour),
	}
	creds, err := getSigningCreds(ctx, obs)
	if err != nil {
		obs.Logger.WarnContext(ctx, "no_gcs_credentials", "error", err, "message", "using original URL")
		return imageURL, nil
	}
	applySigningCreds(opts, creds)

	signedURL, err := storage.SignedURL(bucketName, objectPath, opts)
	if err != nil {
		obs.Logger.ErrorContext(ctx, "error_generating_signed_url", "error", err, "bucket", bucketName, "object", objectPath)
		return imageURL, nil // Retornar URL original si falla
	}

	obs.Logger.InfoContext(ctx, "signed_url_generated", "original_url", imageURL, "signed_url", signedURL[:50]+"...")
	return signedURL, nil
}

// GenerateSignedWriteURL genera una signed URL de escritura (PUT) para subir una imagen a GCS
// Retorna la signed URL, la URL pública y el objectPath
func GenerateSignedWriteURL(ctx context.Context, obs observability.Observability, userID string, fileName string, contentType string) (uploadURL string, publicURL string, objectPath string, err error) {
	// Construir la ruta del objeto: userID/timestamp-random-filename
	timestamp := time.Now().Unix()
	randomSuffix := fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
	objectPath = fmt.Sprintf("%s/%d-%s-%s", userID, timestamp, randomSuffix, fileName)

	bucketName := "micartapro-images"

	// Generar signed URL para PUT (subir)
	opts := &storage.SignedURLOptions{
		Method:      "PUT",
		Expires:     time.Now().Add(15 * time.Minute),
		ContentType: contentType,
	}
	creds, err := getSigningCreds(ctx, obs)
	if err != nil {
		return "", "", "", fmt.Errorf("obtaining signing credentials: %w", err)
	}
	applySigningCreds(opts, creds)

	uploadURL, err = storage.SignedURL(bucketName, objectPath, opts)
	if err != nil {
		return "", "", "", fmt.Errorf("error generating signed URL: %w", err)
	}

	// Construir la URL pública
	publicURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)

	obs.Logger.InfoContext(ctx, "signed_write_url_generated", "objectPath", objectPath, "userID", userID, "contentType", contentType)
	return uploadURL, publicURL, objectPath, nil
}
