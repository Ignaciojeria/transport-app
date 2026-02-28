package fuegoapi

import (
	"errors"
	"strings"

	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/scenegenerator"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/adapter/out/texttospeech"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/billing"
	"net/http"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
)

// SpeechScene representa una escena. Cada escena = 1 imagen + N segmentos de subtítulo (uno por línea).
type SpeechScene struct {
	Lines []string `json:"lines"`
}

// OutputConfig configura el formato de salida.
type OutputConfig struct {
	Type        string `json:"type"`                  // "image" = imágenes por escena, "video" = videos por escena (Veo 3.1)
	AspectRatio string `json:"aspectRatio,omitempty"` // "9:16" para reels
}

type GenerateSpeechRequest struct {
	Scenes       []SpeechScene `json:"scenes"`
	Output       *OutputConfig `json:"output,omitempty"`
	VoiceName    string        `json:"voiceName,omitempty"`
	LanguageCode string        `json:"languageCode,omitempty"`
	StylePrompt  string        `json:"stylePrompt,omitempty"`
}

type SceneSegment struct {
	Text     string  `json:"text"`
	Start    float64 `json:"start"`
	End      float64 `json:"end"`
	ImageURL string  `json:"imageUrl,omitempty"` // URL de la imagen (cuando output.type es "image")
	VideoURL string  `json:"videoUrl,omitempty"` // URL del video (cuando output.type es "video")
}

type GenerateSpeechResponse struct {
	AudioURL        string         `json:"audioUrl"`        // URL pública del audio generado (WAV)
	DurationSeconds float64        `json:"durationSeconds"` // Duración del audio en segundos
	Subtitles       []SceneSegment `json:"subtitles"`       // Segmentos con timing; imageUrl si output.type es "image"
}

func init() {
	ioc.Register(generateSpeechHandler)
}

func generateSpeechHandler(
	s httpserver.Server,
	obs observability.Observability,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
	generateSpeech texttospeech.GenerateSpeech,
	generateSpeechForLines texttospeech.GenerateSpeechForLines,
	sceneGenerator *scenegenerator.SceneGenerator,
	getUserCredits supabaserepo.GetUserCredits,
	consumeCredits supabaserepo.ConsumeCredits,
) {
	fuego.Post(s.Manager, "/api/speech/generate",
		func(c fuego.ContextWithBody[GenerateSpeechRequest]) (GenerateSpeechResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "generateSpeech")
			defer span.End()

			req, err := c.Body()
			if err != nil {
				return GenerateSpeechResponse{}, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			if len(req.Scenes) == 0 {
				return GenerateSpeechResponse{}, fuego.HTTPError{
					Title:  "scenes required",
					Detail: "provide 'scenes' array with at least one scene containing 'lines'",
					Status: http.StatusBadRequest,
				}
			}
			allLines := flattenLines(req.Scenes)
			if len(allLines) == 0 {
				return GenerateSpeechResponse{}, fuego.HTTPError{
					Title:  "scenes must have lines",
					Detail: "each scene must have at least one line in 'lines' array",
					Status: http.StatusBadRequest,
				}
			}
			fullText := strings.Join(allLines, " ")
			outputType := "image"
			if req.Output != nil && req.Output.Type != "" {
				outputType = req.Output.Type
			}
			generateMedia := outputType == "image" || outputType == "video"
			mediaType := scenegenerator.MediaTypeImage
			if outputType == "video" {
				mediaType = scenegenerator.MediaTypeVideo
			}
			aspectRatio := "9:16"
			if req.Output != nil && req.Output.AspectRatio != "" {
				aspectRatio = req.Output.AspectRatio
			}

			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return GenerateSpeechResponse{}, fuego.HTTPError{
					Title:  "user_id not found",
					Detail: "authentication required",
					Status: http.StatusUnauthorized,
				}
			}
			parsedUserID, err := uuid.Parse(userID)
			if err != nil {
				return GenerateSpeechResponse{}, fuego.HTTPError{
					Title:  "invalid user_id",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			// Calcular créditos necesarios: TTS + subtitle timing + (N escenas si generateMedia)
			// No conocemos N hasta después de TTS+subtitle, así que consumimos por etapas
			userCredits, err := getUserCredits(spanCtx, parsedUserID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_getting_user_credits", "error", err)
				return GenerateSpeechResponse{}, fuego.HTTPError{
					Title:  "error checking credits",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			minCredits := billing.CreditsPerSpeechTTS
			if generateMedia {
				creditsPerScene := billing.CreditsPerSceneImage
				if mediaType == scenegenerator.MediaTypeVideo {
					creditsPerScene = billing.CreditsPerSceneVideo
				}
				minCredits += len(req.Scenes) * creditsPerScene
			}
			if userCredits.Balance < minCredits {
				return GenerateSpeechResponse{}, fuego.HTTPError{
					Title:  "insufficient credits",
					Detail: "You don't have enough credits for speech generation. Please purchase credits to continue.",
					Status: http.StatusPaymentRequired,
				}
			}

			// Consumir créditos por TTS
			descTTS := "Generación de audio TTS"
			sourceIDTTS := "speech:tts:" + uuid.New().String()
			_, err = consumeCredits(spanCtx, billing.ConsumeCreditsRequest{
				UserID:      parsedUserID,
				Amount:      billing.CreditsPerSpeechTTS,
				Source:      "speech.tts",
				SourceID:    &sourceIDTTS,
				Description: &descTTS,
			})
			if err != nil {
				if err == supabaserepo.ErrInsufficientCredits {
					return GenerateSpeechResponse{}, fuego.HTTPError{
						Title:  "insufficient credits",
						Detail: "You don't have enough credits for speech generation.",
						Status: http.StatusPaymentRequired,
					}
				}
				obs.Logger.ErrorContext(spanCtx, "error_consuming_credits_tts", "error", err)
				return GenerateSpeechResponse{}, fuego.HTTPError{
					Title:  "error consuming credits",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			opts := &texttospeech.GenerateSpeechOptions{
				VoiceName:    req.VoiceName,
				LanguageCode: req.LanguageCode,
				StylePrompt:  req.StylePrompt,
			}

			// TTS por línea para timing exacto de subtítulos
			result, err := generateSpeechForLines(spanCtx, allLines, opts)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_generating_speech", "error", err, "text", fullText[:min(50, len(fullText))]+"...")
				return GenerateSpeechResponse{}, fuego.HTTPError{
					Title:  "error generating speech",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			segments := segmentsFromLineTimings(allLines, result.LineTimings)

			if generateMedia {
				creditsPerScene := billing.CreditsPerSceneImage
				sourceLabel := "speech.scene_image"
				descScene := "Imagen de escena para reel"
				if mediaType == scenegenerator.MediaTypeVideo {
					creditsPerScene = billing.CreditsPerSceneVideo
					sourceLabel = "speech.scene_video"
					descScene = "Video de escena para reel"
				}

				segIdx := 0
				for sceneIdx, scene := range req.Scenes {
					validLines := trimLines(scene.Lines)
					if len(validLines) == 0 {
						continue
					}
					promptText := strings.Join(validLines, " ")

					// Consumir créditos por cada escena
					credits, err := getUserCredits(spanCtx, parsedUserID)
					if err != nil || credits.Balance < creditsPerScene {
						return GenerateSpeechResponse{}, fuego.HTTPError{
							Title:  "insufficient credits",
							Detail: "You don't have enough credits for scene generation.",
							Status: http.StatusPaymentRequired,
						}
					}
					sourceIDScene := "speech:scene:" + uuid.New().String()
					_, err = consumeCredits(spanCtx, billing.ConsumeCreditsRequest{
						UserID:      parsedUserID,
						Amount:      creditsPerScene,
						Source:      sourceLabel,
						SourceID:    &sourceIDScene,
						Description: &descScene,
					})
					if err != nil {
						if err == supabaserepo.ErrInsufficientCredits {
							return GenerateSpeechResponse{}, fuego.HTTPError{
								Title:  "insufficient credits",
								Detail: "You don't have enough credits for scene generation.",
								Status: http.StatusPaymentRequired,
							}
						}
						obs.Logger.ErrorContext(spanCtx, "error_consuming_credits_scene", "error", err)
						return GenerateSpeechResponse{}, fuego.HTTPError{
							Title:  "error consuming credits",
							Detail: err.Error(),
							Status: http.StatusInternalServerError,
						}
					}

					mediaURL, err := sceneGenerator.GenerateMedia(spanCtx, fullText, promptText, sceneIdx, aspectRatio, "", mediaType)
					if err != nil {
						if errors.Is(err, scenegenerator.ErrVideoGenerationSkipped) {
							obs.Logger.WarnContext(spanCtx, "scene_video_skipped_after_retries", "sceneIndex", sceneIdx)
							// Omitir este segmento, continuar con el siguiente (URL queda vacía)
						} else {
							obs.Logger.ErrorContext(spanCtx, "error_generating_scene_media", "error", err, "sceneIndex", sceneIdx)
							return GenerateSpeechResponse{}, fuego.HTTPError{
								Title:  "error generating scene media",
								Detail: err.Error(),
								Status: http.StatusInternalServerError,
							}
						}
					} else {
						for j := 0; j < len(validLines); j++ {
							if mediaType == scenegenerator.MediaTypeVideo {
								segments[segIdx+j].VideoURL = mediaURL
							} else {
								segments[segIdx+j].ImageURL = mediaURL
							}
						}
					}
					segIdx += len(validLines)
				}
			}

			return GenerateSpeechResponse{
				AudioURL:        result.AudioURL,
				DurationSeconds: result.DurationSeconds,
				Subtitles:       segments,
			}, nil
		},
		option.Summary("Generate speech from text"),
		option.Description("Converts text to audio using Google Cloud Text-to-Speech and returns the public URL of the generated MP3"),
		option.Tags("speech"),
		option.Middleware(jwtAuthMiddleware),
	)
}

func flattenLines(scenes []SpeechScene) []string {
	var out []string
	for _, s := range scenes {
		for _, line := range s.Lines {
			t := strings.TrimSpace(line)
			if t != "" {
				out = append(out, t)
			}
		}
	}
	return out
}

func trimLines(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

// segmentsFromLineTimings construye segmentos con timing exacto desde TTS por línea.
func segmentsFromLineTimings(lines []string, timings []texttospeech.LineTiming) []SceneSegment {
	if len(lines) == 0 || len(timings) == 0 || len(lines) != len(timings) {
		return nil
	}
	segments := make([]SceneSegment, len(lines))
	for i := range lines {
		segments[i] = SceneSegment{
			Text:  strings.TrimSpace(lines[i]),
			Start: timings[i].Start,
			End:   timings[i].End,
		}
	}
	return segments
}
