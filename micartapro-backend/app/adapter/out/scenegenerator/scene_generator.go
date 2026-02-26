package scenegenerator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"micartapro/app/adapter/out/imagegenerator"
	"micartapro/app/shared/infrastructure/ai"
	"micartapro/app/shared/infrastructure/gcs"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/genai"
)

const sceneAspectRatio = "9:16" // Vertical para reels
const veoModel = "veo-3.1-generate-001"
const videoPollInterval = 10 * time.Second
const videoPollTimeout = 6 * time.Minute
const videoMaxRetries = 3

// ErrVideoGenerationSkipped se retorna cuando la generación falla tras max reintentos.
// El caller debe omitir ese segmento y continuar con el siguiente.
var ErrVideoGenerationSkipped = errors.New("video generation skipped after retries")

// MediaType indica si generar imagen o video.
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

// GenerateSceneImage genera una imagen para una escena de reel.
// fullTranscript: contexto global del video (coherencia visual).
// fragmentText: fragmento específico de esta escena.
// aspectRatio: "9:16" por defecto si vacío.
// style: "illustration" = ilustraciones/dibujos; "" o "photorealistic" = fotorrealista (default).
type GenerateSceneImage func(ctx context.Context, fullTranscript string, fragmentText string, sceneIndex int, aspectRatio string, style string) (imageURL string, err error)

// GenerateSceneMedia genera imagen o video según mediaType.
// Retorna la URL pública del recurso generado (imagen o video).
type GenerateSceneMedia func(ctx context.Context, fullTranscript string, fragmentText string, sceneIndex int, aspectRatio string, style string, mediaType MediaType) (mediaURL string, err error)

// SceneGenerator agrupa generación de imagen y video para escenas.
type SceneGenerator struct {
	// GenerateImage genera solo imagen (compatibilidad).
	GenerateImage GenerateSceneImage
	// GenerateMedia genera imagen o video según mediaType.
	GenerateMedia GenerateSceneMedia
}

func init() {
	ioc.Registry(NewSceneGenerator, ai.NewClient, observability.NewObservability, gcs.NewClient)
}

func NewSceneGenerator(client *genai.Client, obs observability.Observability, gcsClient *storage.Client) (*SceneGenerator, error) {
	generateMedia := func(ctx context.Context, fullTranscript string, fragmentText string, sceneIndex int, aspectRatio string, style string, mediaType MediaType) (string, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "generate_scene_media")
		defer span.End()

		userID, ok := sharedcontext.UserIDFromContext(spanCtx)
		if !ok || userID == "" {
			return "", fmt.Errorf("userID is required but not found in context")
		}

		fragmentText = strings.TrimSpace(fragmentText)
		if fragmentText == "" {
			return "", fmt.Errorf("fragment text cannot be empty")
		}

		ar := sceneAspectRatio
		if aspectRatio != "" {
			ar = aspectRatio
		}

		prompt := buildScenePrompt(fullTranscript, fragmentText, style)

		if mediaType == MediaTypeVideo {
			return generateSceneVideo(spanCtx, client, gcsClient, obs, userID, sceneIndex, prompt, ar)
		}
		return generateSceneImage(spanCtx, client, gcsClient, obs, userID, sceneIndex, prompt, ar)
	}

	generateImage := func(ctx context.Context, fullTranscript string, fragmentText string, sceneIndex int, aspectRatio string, style string) (string, error) {
		return generateMedia(ctx, fullTranscript, fragmentText, sceneIndex, aspectRatio, style, MediaTypeImage)
	}

	return &SceneGenerator{
		GenerateImage: generateImage,
		GenerateMedia: generateMedia,
	}, nil
}

func generateSceneImage(ctx context.Context, client *genai.Client, gcsClient *storage.Client, obs observability.Observability, userID string, sceneIndex int, prompt string, aspectRatio string) (string, error) {
	obs.Logger.InfoContext(ctx, "generating_scene_image", "sceneIndex", sceneIndex, "fragment", prompt[:min(60, len(prompt))]+"...")

	config := &genai.GenerateImagesConfig{
		AspectRatio:    aspectRatio,
		NumberOfImages: 1,
	}

	resp, err := client.Models.GenerateImages(ctx, "imagen-4.0-ultra-generate-001", prompt, config)
	if err != nil {
		obs.Logger.ErrorContext(ctx, "error_generating_scene_image", "error", err, "sceneIndex", sceneIndex)
		return "", fmt.Errorf("error generating scene image: %w", err)
	}

	if len(resp.GeneratedImages) == 0 {
		return "", fmt.Errorf("no scene image generated")
	}

	img := resp.GeneratedImages[0]
	if img.Image == nil || len(img.Image.ImageBytes) == 0 {
		return "", fmt.Errorf("generated scene image is empty")
	}

	imgBytes := img.Image.ImageBytes
	mimeType := "image/png"

	uploadURL, publicURL, _, err := imagegenerator.GenerateSignedWriteURLForReelScene(ctx, gcsClient, obs, userID, sceneIndex, mimeType)
	if err != nil {
		return "", fmt.Errorf("error generating upload URL: %w", err)
	}

	if err := uploadBytes(ctx, uploadURL, imgBytes, mimeType, obs, "scene_image"); err != nil {
		return "", err
	}

	obs.Logger.InfoContext(ctx, "scene_image_uploaded", "publicURL", publicURL, "sceneIndex", sceneIndex)
	return publicURL, nil
}

func generateSceneVideo(ctx context.Context, client *genai.Client, gcsClient *storage.Client, obs observability.Observability, userID string, sceneIndex int, prompt string, aspectRatio string) (string, error) {
	videoConfig := &genai.GenerateVideosConfig{
		AspectRatio:      aspectRatio,
		NumberOfVideos:  1,
		Resolution:       "720p",
		DurationSeconds: ptr(int32(6)),
		GenerateAudio:   ptr(false), // Sin audio: los clips se combinan con el TTS/subtítulos existentes
	}

	var lastErr error
	for attempt := 1; attempt <= videoMaxRetries; attempt++ {
		obs.Logger.InfoContext(ctx, "generating_scene_video", "sceneIndex", sceneIndex, "attempt", attempt, "maxRetries", videoMaxRetries, "fragment", prompt[:min(60, len(prompt))]+"...")

		publicURL, err := tryGenerateSceneVideo(ctx, client, gcsClient, obs, userID, sceneIndex, prompt, videoConfig)
		if err == nil {
			return publicURL, nil
		}
		lastErr = err
		obs.Logger.WarnContext(ctx, "video_generation_attempt_failed", "sceneIndex", sceneIndex, "attempt", attempt, "error", err)

		if attempt < videoMaxRetries {
			backoff := time.Duration(attempt*5) * time.Second
			obs.Logger.InfoContext(ctx, "video_generation_retry", "sceneIndex", sceneIndex, "nextAttempt", attempt+1, "backoffSec", backoff.Seconds())
			time.Sleep(backoff)
		}
	}

	obs.Logger.ErrorContext(ctx, "video_generation_skipped_after_retries", "sceneIndex", sceneIndex, "error", lastErr)
	return "", ErrVideoGenerationSkipped
}

func tryGenerateSceneVideo(ctx context.Context, client *genai.Client, gcsClient *storage.Client, obs observability.Observability, userID string, sceneIndex int, prompt string, videoConfig *genai.GenerateVideosConfig) (string, error) {
	operation, err := client.Models.GenerateVideos(ctx, veoModel, prompt, nil, videoConfig)
	if err != nil {
		return "", fmt.Errorf("error generating scene video: %w", err)
	}

	// Poll hasta que el video esté listo
	deadline := time.Now().Add(videoPollTimeout)
	for !operation.Done {
		if time.Now().After(deadline) {
			return "", fmt.Errorf("video generation timed out after %v", videoPollTimeout)
		}
		obs.Logger.InfoContext(ctx, "waiting_for_video_generation", "sceneIndex", sceneIndex)
		time.Sleep(videoPollInterval)
		operation, err = client.Operations.GetVideosOperation(ctx, operation, nil)
		if err != nil {
			return "", fmt.Errorf("error polling video operation: %w", err)
		}
	}

	if operation.Error != nil {
		return "", fmt.Errorf("video generation failed: %v", operation.Error)
	}

	if operation.Response == nil || len(operation.Response.GeneratedVideos) == 0 {
		return "", fmt.Errorf("no scene video generated")
	}

	generatedVideo := operation.Response.GeneratedVideos[0]
	if generatedVideo.Video == nil {
		return "", fmt.Errorf("generated scene video is empty")
	}

	// Descargar bytes si vienen por URI (Vertex AI)
	if len(generatedVideo.Video.VideoBytes) == 0 && generatedVideo.Video.URI != "" {
		_, err := client.Files.Download(ctx, genai.NewDownloadURIFromVideo(generatedVideo.Video), nil)
		if err != nil {
			return "", fmt.Errorf("error downloading generated video: %w", err)
		}
	}

	videoBytes := generatedVideo.Video.VideoBytes
	if len(videoBytes) == 0 {
		return "", fmt.Errorf("generated scene video has no bytes")
	}

	mimeType := "video/mp4"
	uploadURL, publicURL, _, err := imagegenerator.GenerateSignedWriteURLForReelSceneVideo(ctx, gcsClient, obs, userID, sceneIndex, mimeType)
	if err != nil {
		return "", fmt.Errorf("error generating upload URL: %w", err)
	}

	if err := uploadBytes(ctx, uploadURL, videoBytes, mimeType, obs, "scene_video"); err != nil {
		return "", err
	}

	obs.Logger.InfoContext(ctx, "scene_video_uploaded", "publicURL", publicURL, "sceneIndex", sceneIndex)
	return publicURL, nil
}

func uploadBytes(ctx context.Context, uploadURL string, data []byte, mimeType string, obs observability.Observability, logKind string) error {
	req, err := http.NewRequestWithContext(ctx, "PUT", uploadURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error creating upload request: %w", err)
	}
	req.Header.Set("Content-Type", mimeType)

	httpClient := &http.Client{Timeout: 60 * time.Second}
	uploadResp, err := httpClient.Do(req)
	if err != nil {
		obs.Logger.ErrorContext(ctx, "error_uploading_"+logKind, "error", err)
		return fmt.Errorf("error uploading %s: %w", logKind, err)
	}
	body, _ := io.ReadAll(uploadResp.Body)
	uploadResp.Body.Close()

	if uploadResp.StatusCode < 200 || uploadResp.StatusCode >= 300 {
		obs.Logger.ErrorContext(ctx, "error_upload_status", "status", uploadResp.StatusCode, "body", string(body))
		return fmt.Errorf("error uploading %s: status %d: %s", logKind, uploadResp.StatusCode, string(body))
	}
	return nil
}

func ptr[T any](v T) *T { return &v }

// globalContext define el estilo visual fotorrealista (SaaS, MiCartaPro).
const globalContextPhotorealistic = "CRITICAL: Do NOT include any text, subtitles, captions, words, letters, numbers, labels, or written content in the image. The image must be purely visual—no readable text whatsoever. Text and subtitles will be added separately in post-production. " +
	"Aspect ratio 9:16 (vertical, portrait format for reels/shorts). " +
	"Vertical infographic scene, SaaS product explainer style. " +
	"Digital restaurant ordering platform MiCartaPro. " +
	"Dark minimalist background with smooth gradient from black to deep blue. " +
	"Clean modern UI elements: dashboard panels, graphs, mobile phone mockups, glowing interface elements. " +
	"Maximize horizontal space usage—fill the width effectively. " +
	"Minimalist composition, cinematic lighting, futuristic SaaS aesthetic. " +
	"Inspired by Stripe, Linear, Apple product presentations. " +
	"No clutter, no busy background, focus on clarity and readability. " +
	"Photorealistic rendering, high contrast subject, dark non-distracting background. " +
	"Remember: image must be text-free—visual elements only. "

// globalContextIllustration define el estilo visual inspirado en Malika Favre.
const globalContextIllustration = "CRITICAL: Do NOT include any text, subtitles, captions, words, letters, numbers, labels, or written content in the image. The image must be purely visual—no readable text whatsoever. " +
	"Aspect ratio 9:16 (vertical, portrait format for reels/shorts). " +
	"In the style of artist Malika Favre: bold minimalist vector illustration. " +
	"Limited color palette: 2-3 solid colors maximum (e.g. black, white, one accent like coral, teal or mustard). " +
	"Fill the entire frame with color—no empty spaces, no white or blank areas. Full composition, rich and saturated. " +
	"Flat design, no gradients, no textures. Clean geometric shapes, elegant curves, suggestive and sophisticated. " +
	"Retro 60s/70s pop art influence, high contrast, graphic poster aesthetic. " +
	"Minimalist composition, one strong visual concept. Elegant and refined. " +
	"Remember: image must be text-free—visual elements only. "

// buildScenePrompt construye el prompt con 2 capas: global context + scene context.
// fullTranscript: narrativa completa (coherencia visual).
// fragmentText: concepto específico de esta escena.
// style: "illustration" usa estilo ilustración; otro valor usa fotorrealista.
func buildScenePrompt(fullTranscript string, fragmentText string, style string) string {
	fullTranscript = strings.TrimSpace(fullTranscript)
	globalCtx := globalContextPhotorealistic
	if strings.TrimSpace(strings.ToLower(style)) == "illustration" {
		globalCtx = globalContextIllustration
	}
	sceneContext := "Scene showing a visual representation of: \"" + fragmentText + "\". "
	if fullTranscript != "" && fullTranscript != fragmentText {
		sceneContext = "Full video narrative: " + truncate(fullTranscript, 200) + ". " +
			"This scene specifically shows: \"" + fragmentText + "\". "
	}
	return globalCtx + sceneContext + "No text in the image."
}

func truncate(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
