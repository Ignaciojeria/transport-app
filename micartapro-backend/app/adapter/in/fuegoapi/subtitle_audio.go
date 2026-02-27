package fuegoapi

import (
	"errors"
	"fmt"
	"strings"

	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/scenegenerator"
	"micartapro/app/adapter/out/speechtotext"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/adapter/out/texttospeech"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/shared/subtitles"
	"micartapro/app/usecase/billing"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
)

// SubtitleAudioRequest es el body del POST para subtitular un audio existente o generar desde guion.
// Modo 1: audioUrl → transcribe (Speech-to-Text) y genera subtítulos + imágenes.
// Modo 2: scriptText (sin audioUrl) → genera audio con TTS, subtítulos con timing exacto, e imágenes.
type SubtitleAudioRequest struct {
	AudioURL         string   `json:"audioUrl"`
	// ScriptText: guion en texto. Si audioUrl está vacío, se usa TTS para generar el audio.
	// Cada línea (separada por \n) se convierte en un segmento con timing exacto.
	ScriptText       string   `json:"scriptText,omitempty"`
	LanguageCode     string   `json:"languageCode,omitempty"`
	SubtitleStyle    string   `json:"subtitleStyle,omitempty"`
	MaxCharsPerLine  int      `json:"maxCharsPerLine,omitempty"`
	MaxLines         int      `json:"maxLines,omitempty"`
	// Placement fijo (cuando placementStrategy=FIXED)
	Placement string `json:"placement,omitempty"`
	// Estrategia creativa: FIXED = todo igual; DYNAMIC = placement por segmento
	PlacementStrategy string `json:"placementStrategy,omitempty"`
	// creativity: 0 = seguro (todo abajo), 0.5 = algo dinámico, 0.8 = creativo, 1.0 = muy cinematográfico
	Creativity float64 `json:"creativity,omitempty"`
	// Frases a destacar
	EmphasisPhrases []string `json:"emphasisPhrases,omitempty"`
	// Layout constraints para DYNAMIC
	AvoidCenterLongSentences bool `json:"avoidCenterLongSentences,omitempty"`
	PreferCenterForEmphasis  bool `json:"preferCenterForEmphasis,omitempty"`
	// overflowStrategy cuando lines excede maxLines: REBALANCE | SHRINK | REBALANCE_THEN_SHRINK
	OverflowStrategy string `json:"overflowStrategy,omitempty"`
	// Stickiness
	ExtendEndSeconds float64 `json:"extendEndSeconds,omitempty"`
	// timingOffsetSeconds: + = subtítulos aparecen antes, - = después (corrige desfase audio/subtítulos)
	TimingOffsetSeconds float64 `json:"timingOffsetSeconds,omitempty"`
	// Safe area (fracciones 0-1) para TOP/CENTER/BOTTOM
	SafeAreaTop    float64 `json:"safeAreaTop,omitempty"`
	SafeAreaBottom float64 `json:"safeAreaBottom,omitempty"`
	SafeAreaLeft   float64 `json:"safeAreaLeft,omitempty"`
	SafeAreaRight  float64 `json:"safeAreaRight,omitempty"`
	// Dampening de placement
	MinSecondsBetweenPlacementChanges float64 `json:"minSecondsBetweenPlacementChanges,omitempty"`
	MaxCenterBlocksInRow              int     `json:"maxCenterBlocksInRow,omitempty"`
	// Preset de dirección: CINEMATIC_DYNAMIC_V1 = heurísticas narrativas determinísticas
	DirectionPreset string `json:"directionPreset,omitempty"`
	// Incluir imagen por subtítulo. Si frameImages está vacío: genera imágenes con IA automáticamente.
	// Si frameImages tiene datos: asocia cada segmento al frame más cercano por timestamp.
	IncludeImagesPerSegment bool `json:"includeImagesPerSegment,omitempty"`
	// Frames con timestamp (segundos) e imagen. Opcional. Si vacío y includeImagesPerSegment=true, se generan con IA.
	FrameImages []FrameImage `json:"frameImages,omitempty"`
	// Estilo de imagen generada: "illustration" (default) o "photorealistic". Solo aplica cuando se generan con IA.
	ImageStyle string `json:"imageStyle,omitempty"`
	// OutputMediaType: "image" (default) o "video". Si video, genera clips con Veo 3.1 en lugar de imágenes.
	OutputMediaType string `json:"outputMediaType,omitempty"`
}

// FrameImage representa un frame de video con su timestamp para asociar a subtítulos
type FrameImage struct {
	TimestampSec float64 `json:"timestampSec"`
	ImageURL     string  `json:"imageUrl"`
}

func init() {
	ioc.Registry(
		subtitleAudioHandler,
		httpserver.New,
		observability.NewObservability,
		apimiddleware.NewJWTAuthMiddleware,
		speechtotext.NewTranscribeAudioProvider,
		texttospeech.NewTextToSpeechForLines,
		scenegenerator.NewSceneGenerator,
		supabaserepo.NewGetUserCredits,
		supabaserepo.NewConsumeCredits,
	)
}

func subtitleAudioHandler(
	s httpserver.Server,
	obs observability.Observability,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
	transcribeAudio speechtotext.TranscribeAudio,
	generateSpeechForLines texttospeech.GenerateSpeechForLines,
	sceneGenerator *scenegenerator.SceneGenerator,
	getUserCredits supabaserepo.GetUserCredits,
	consumeCredits supabaserepo.ConsumeCredits,
) {
	fuego.Post(s.Manager, "/api/speech/subtitle",
		func(c fuego.ContextWithBody[SubtitleAudioRequest]) (subtitles.TimelineResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "subtitleAudio")
			defer span.End()

			req, err := c.Body()
			if err != nil {
				return subtitles.TimelineResponse{}, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			audioURL := strings.TrimSpace(req.AudioURL)
			scriptText := strings.TrimSpace(req.ScriptText)
			useScriptMode := audioURL == "" && scriptText != ""

			if audioURL == "" && scriptText == "" {
				return subtitles.TimelineResponse{}, fuego.HTTPError{
					Title:  "audioUrl or scriptText required",
					Detail: "provide 'audioUrl' with a valid URL to an audio file, or 'scriptText' with your script to generate audio with TTS",
					Status: http.StatusBadRequest,
				}
			}
			if audioURL != "" && scriptText != "" {
				return subtitles.TimelineResponse{}, fuego.HTTPError{
					Title:  "provide one input only",
					Detail: "provide either 'audioUrl' (to transcribe) or 'scriptText' (to generate with TTS), not both",
					Status: http.StatusBadRequest,
				}
			}

			languageCode := req.LanguageCode
			if strings.TrimSpace(languageCode) == "" {
				languageCode = "es-ES"
			}

			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return subtitles.TimelineResponse{}, fuego.HTTPError{
					Title:  "user_id not found",
					Detail: "authentication required",
					Status: http.StatusUnauthorized,
				}
			}
			parsedUserID, err := uuid.Parse(userID)
			if err != nil {
				return subtitles.TimelineResponse{}, fuego.HTTPError{
					Title:  "invalid user_id",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			creditsNeeded := billing.CreditsPerSpeechSubtitle
			if useScriptMode {
				creditsNeeded = billing.CreditsPerSpeechTTS
			}
			userCredits, err := getUserCredits(spanCtx, parsedUserID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_getting_user_credits", "error", err)
				return subtitles.TimelineResponse{}, fuego.HTTPError{
					Title:  "error checking credits",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			if userCredits.Balance < creditsNeeded {
				return subtitles.TimelineResponse{}, fuego.HTTPError{
					Title:  "insufficient credits",
					Detail: "You don't have enough credits. Please purchase credits to continue.",
					Status: http.StatusPaymentRequired,
				}
			}

			var segments []speechtotext.SubtitleSegment
			var durationSec float64

			if useScriptMode {
				desc := "Generación de audio TTS desde guion"
				sourceID := "speech:tts:" + uuid.New().String()
				_, err = consumeCredits(spanCtx, billing.ConsumeCreditsRequest{
					UserID:      parsedUserID,
					Amount:      billing.CreditsPerSpeechTTS,
					Source:      "speech.tts",
					SourceID:    &sourceID,
					Description: &desc,
				})
				if err != nil {
					if err == supabaserepo.ErrInsufficientCredits {
						return subtitles.TimelineResponse{}, fuego.HTTPError{
							Title:  "insufficient credits",
							Detail: "You don't have enough credits for TTS.",
							Status: http.StatusPaymentRequired,
						}
					}
					obs.Logger.ErrorContext(spanCtx, "error_consuming_credits", "error", err)
					return subtitles.TimelineResponse{}, fuego.HTTPError{
						Title:  "error consuming credits",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}

				lines := splitScriptIntoLines(scriptText)
				if len(lines) == 0 {
					return subtitles.TimelineResponse{}, fuego.HTTPError{
						Title:  "scriptText is empty",
						Detail: "scriptText must contain at least one non-empty line",
						Status: http.StatusBadRequest,
					}
				}

				ttsOpts := &texttospeech.GenerateSpeechOptions{
					LanguageCode: languageCode,
				}
				result, err := generateSpeechForLines(spanCtx, lines, ttsOpts)
				if err != nil {
					obs.Logger.ErrorContext(spanCtx, "error_generating_tts", "error", err)
					return subtitles.TimelineResponse{}, fuego.HTTPError{
						Title:  "error generating speech",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}

				audioURL = result.AudioURL
				durationSec = result.DurationSeconds
				segments = make([]speechtotext.SubtitleSegment, len(lines))
				for i := range lines {
					segments[i] = speechtotext.SubtitleSegment{
						Text:  lines[i],
						Start: result.LineTimings[i].Start,
						End:   result.LineTimings[i].End,
					}
				}
			} else {
				desc := "Transcripción de audio a subtítulos"
				sourceID := "speech:subtitle:" + uuid.New().String()
				_, err = consumeCredits(spanCtx, billing.ConsumeCreditsRequest{
					UserID:      parsedUserID,
					Amount:      billing.CreditsPerSpeechSubtitle,
					Source:      "speech.subtitle",
					SourceID:    &sourceID,
					Description: &desc,
				})
				if err != nil {
					if err == supabaserepo.ErrInsufficientCredits {
						return subtitles.TimelineResponse{}, fuego.HTTPError{
							Title:  "insufficient credits",
							Detail: "You don't have enough credits for audio transcription.",
							Status: http.StatusPaymentRequired,
						}
					}
					obs.Logger.ErrorContext(spanCtx, "error_consuming_credits", "error", err)
					return subtitles.TimelineResponse{}, fuego.HTTPError{
						Title:  "error consuming credits",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}

				segments, durationSec, err = transcribeAudio(spanCtx, audioURL, languageCode)
				if err != nil {
					obs.Logger.ErrorContext(spanCtx, "error_transcribing_audio", "error", err, "audioUrl", audioURL)
					return subtitles.TimelineResponse{}, fuego.HTTPError{
						Title:  "error transcribing audio",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				segments = speechtotext.MergeFragmentedSegments(segments)
			}

			// Resolver estilo, layout y hints
			subtitleStyle := resolveSubtitleStyle(req.SubtitleStyle)
			maxChars := req.MaxCharsPerLine
			if maxChars <= 0 {
				maxChars = 32
			}
			maxLines := req.MaxLines
			if maxLines <= 0 {
				maxLines = 3
			}

			placementStrategy := resolvePlacementStrategy(req.PlacementStrategy, subtitleStyle)
			defaultPlacement := resolvePlacement(req.Placement)

			// Mapear a SubtitleSegment con lines[] cuando aplica
			overflow := resolveOverflowStrategy(req.OverflowStrategy)
			subs := make([]subtitles.SubtitleSegment, len(segments))
			for i, seg := range segments {
				subs[i] = subtitles.SubtitleSegment{
					Text:  seg.Text,
					Start: seg.Start,
					End:   seg.End,
				}
				if req.IncludeImagesPerSegment && len(req.FrameImages) > 0 {
					subs[i].ImageURL = findNearestFrameImage(seg.Start, req.FrameImages)
				}
				if subtitleStyle == subtitles.StyleTrailer {
					// TRAILER: fragmentación se hace después
				} else if subtitleStyle == subtitles.StyleLineByLine || subtitleStyle == subtitles.StyleLineByLineBig || subtitleStyle == subtitles.StyleAlexHormozi {
					subs[i].Lines = subtitles.SplitSegmentIntoLinesWithOverflow(seg.Text, seg.Start, seg.End, maxChars, maxLines, overflow)
				} else if subtitleStyle == subtitles.StyleCinematicDynamic {
					// CINEMATIC_DYNAMIC: respetar maxCharsPerLine y maxLines
					if len(seg.Words) > 0 {
						words := make([]subtitles.SubtitleWord, len(seg.Words))
						for j, w := range seg.Words {
							words[j] = subtitles.SubtitleWord{Text: w.Text, Start: w.Start, End: w.End}
						}
						subs[i].Lines = subtitles.SplitSegmentIntoLinesWithWordTiming(seg.Text, seg.Start, seg.End, words, maxChars, maxLines, overflow, 0.35)
					} else {
						subs[i].Lines = subtitles.SplitSegmentIntoLinesWithOverflow(seg.Text, seg.Start, seg.End, maxChars, maxLines, overflow)
					}
				}
			}

			directionPreset := resolveDirectionPreset(req.DirectionPreset, subtitleStyle)
			dynRules := &subtitles.DynamicRules{
				MinSecondsBetweenPlacementChanges: req.MinSecondsBetweenPlacementChanges,
				MaxCenterBlocksInRow:              req.MaxCenterBlocksInRow,
			}
			if dynRules.MinSecondsBetweenPlacementChanges <= 0 && subtitleStyle == subtitles.StyleCinematicDynamic {
				dynRules.MinSecondsBetweenPlacementChanges = 2.0
			}
			if dynRules.MaxCenterBlocksInRow <= 0 {
				dynRules.MaxCenterBlocksInRow = 2
			}

			segmentWords := make([][]subtitles.SubtitleWord, len(segments))
			for i, seg := range segments {
				if len(seg.Words) > 0 {
					segmentWords[i] = make([]subtitles.SubtitleWord, len(seg.Words))
					for j, w := range seg.Words {
						segmentWords[i][j] = subtitles.SubtitleWord{Text: w.Text, Start: w.Start, End: w.End}
					}
				}
			}

			subs = subtitles.BuildSubtitleDirection(subs, subtitles.DirectionInput{
				Style:             subtitleStyle,
				DirectionPreset:   directionPreset,
				PlacementStrategy: placementStrategy,
				DefaultPlacement:  defaultPlacement,
				EmphasisPhrases:   req.EmphasisPhrases,
				DurationSec:       durationSec,
				Creativity:        req.Creativity,
				Constraints: &subtitles.LayoutConstraints{
					AvoidCenterLongSentences: req.AvoidCenterLongSentences,
					PreferCenterForEmphasis:  req.PreferCenterForEmphasis,
				},
				DynamicRules:  dynRules,
				SegmentWords:  segmentWords,
			})

			// Generar imágenes o videos automáticamente por segmento cuando includeImagesPerSegment y no hay frameImages
			if req.IncludeImagesPerSegment && len(req.FrameImages) == 0 {
				outputMediaType := scenegenerator.MediaTypeImage
				if strings.TrimSpace(strings.ToLower(req.OutputMediaType)) == "video" {
					outputMediaType = scenegenerator.MediaTypeVideo
				}
				creditsPerSegment := billing.CreditsPerSceneImage
				sourceLabel := "speech.subtitle_segment_image"
				desc := "Imagen de segmento para subtítulos"
				if outputMediaType == scenegenerator.MediaTypeVideo {
					creditsPerSegment = billing.CreditsPerSceneVideo
					sourceLabel = "speech.subtitle_segment_video"
					desc = "Video de segmento para subtítulos"
				}

				fullTranscript := buildFullTranscript(subs)
				for i := range subs {
					credits, err := getUserCredits(spanCtx, parsedUserID)
					if err != nil || credits.Balance < creditsPerSegment {
						return subtitles.TimelineResponse{}, fuego.HTTPError{
							Title:  "insufficient credits",
							Detail: "You don't have enough credits for scene generation.",
							Status: http.StatusPaymentRequired,
						}
					}
					sourceID := "speech:subtitle_segment:" + uuid.New().String()
					_, err = consumeCredits(spanCtx, billing.ConsumeCreditsRequest{
						UserID:      parsedUserID,
						Amount:      creditsPerSegment,
						Source:      sourceLabel,
						SourceID:    &sourceID,
						Description: &desc,
					})
					if err != nil {
						if err == supabaserepo.ErrInsufficientCredits {
							return subtitles.TimelineResponse{}, fuego.HTTPError{
								Title:  "insufficient credits",
								Detail: "You don't have enough credits for scene generation.",
								Status: http.StatusPaymentRequired,
							}
						}
						obs.Logger.ErrorContext(spanCtx, "error_consuming_credits_segment_media", "error", err)
						return subtitles.TimelineResponse{}, fuego.HTTPError{
							Title:  "error consuming credits",
							Detail: err.Error(),
							Status: http.StatusInternalServerError,
						}
					}
					imageStyle := "illustration"
					if strings.TrimSpace(strings.ToLower(req.ImageStyle)) == "photorealistic" {
						imageStyle = "photorealistic"
					}
					mediaURL, err := sceneGenerator.GenerateMedia(spanCtx, fullTranscript, subs[i].Text, i, "9:16", imageStyle, outputMediaType)
					if err != nil {
						if errors.Is(err, scenegenerator.ErrVideoGenerationSkipped) {
							obs.Logger.WarnContext(spanCtx, "segment_video_skipped_after_retries", "segmentIndex", i)
							// Omitir este segmento, continuar con el siguiente (URL queda vacía)
						} else {
							obs.Logger.ErrorContext(spanCtx, "error_generating_segment_media", "error", err, "segmentIndex", i)
							return subtitles.TimelineResponse{}, fuego.HTTPError{
								Title:  "error generating segment media",
								Detail: err.Error(),
								Status: http.StatusInternalServerError,
							}
						}
					} else {
						if outputMediaType == scenegenerator.MediaTypeVideo {
							subs[i].VideoURL = mediaURL
						} else {
							subs[i].ImageURL = mediaURL
						}
					}
				}
			}

			trailerMaxLines := maxLines
			if subtitleStyle == subtitles.StyleTrailer {
				trailerMaxLines = 1
			}
			layout := subtitles.SubtitleLayout{
				PlacementStrategy: placementStrategy,
				DefaultPlacement:  defaultPlacement,
				MaxLines:          trailerMaxLines,
				MaxCharsPerLine:   maxChars,
				OverflowStrategy:  resolveOverflowStrategy(req.OverflowStrategy),
				Constraints: &subtitles.LayoutConstraints{
					AvoidCenterLongSentences: req.AvoidCenterLongSentences,
					PreferCenterForEmphasis:  req.PreferCenterForEmphasis,
				},
			}
			if req.SafeAreaTop > 0 || req.SafeAreaBottom > 0 || req.SafeAreaLeft > 0 || req.SafeAreaRight > 0 {
				layout.SafeArea = &subtitles.SafeArea{
					Top:    req.SafeAreaTop,
					Bottom: req.SafeAreaBottom,
					Left:   req.SafeAreaLeft,
					Right:  req.SafeAreaRight,
				}
			} else if placementStrategy == subtitles.PlacementStrategyDynamic {
				layout.SafeArea = &subtitles.SafeArea{Top: 0.12, Bottom: 0.14, Left: 0.06, Right: 0.06}
			}
			if placementStrategy == subtitles.PlacementStrategyDynamic {
				layout.DynamicRules = &subtitles.DynamicRules{
					MinSecondsBetweenPlacementChanges: req.MinSecondsBetweenPlacementChanges,
					MaxCenterBlocksInRow:              req.MaxCenterBlocksInRow,
				}
				if layout.DynamicRules.MinSecondsBetweenPlacementChanges <= 0 && subtitleStyle == subtitles.StyleCinematicDynamic {
					layout.DynamicRules.MinSecondsBetweenPlacementChanges = 2.0
				}
				if layout.DynamicRules.MaxCenterBlocksInRow <= 0 {
					layout.DynamicRules.MaxCenterBlocksInRow = 2
				}
			}

			var renderHints *subtitles.RenderHints
			if req.ExtendEndSeconds > 0 || req.TimingOffsetSeconds != 0 {
				renderHints = &subtitles.RenderHints{
					ExtendEndSeconds:    req.ExtendEndSeconds,
					TimingOffsetSeconds: req.TimingOffsetSeconds,
				}
			}

			// Aplicar theme accent a subs antes de mapear a scenes
			if subtitleStyle == subtitles.StyleTrailer {
				theme := subtitles.GetThemePreset(subtitles.ThemeTrailerDefault)
				if theme != nil && theme.AccentColor != "" {
					for i := range subs {
						if subs[i].Emphasis != nil && len(subs[i].Emphasis.Phrases) > 0 {
							subs[i].Emphasis.Color = theme.AccentColor
						}
					}
					for i := range subs {
						for j := range subs[i].Lines {
							if subs[i].Lines[j].Emphasis != nil && len(subs[i].Lines[j].Emphasis.Phrases) > 0 {
								subs[i].Lines[j].Emphasis.Color = theme.AccentColor
							}
						}
					}
				}
			}

			safeArea := subtitles.SafeArea{Top: 0.12, Bottom: 0.14, Left: 0.06, Right: 0.06}
			if layout.SafeArea != nil {
				safeArea = *layout.SafeArea
			}

			scenes := make([]subtitles.Scene, len(subs))
			for i, seg := range subs {
				placement := seg.Placement
				if placement == "" {
					placement = defaultPlacement
				}
				animation := seg.Animation
				if animation == "" {
					animation = subtitles.AnimationFadeIn
				}
				size := seg.Size
				if size == "" {
					size = subtitles.SizeL
				}

				lines := make([]subtitles.SceneSubtitleLine, len(seg.Lines))
				for j, ln := range seg.Lines {
					lineSize := ln.Size
					if lineSize == "" {
						lineSize = size
					}
					lines[j] = subtitles.SceneSubtitleLine{
						Text:     ln.Text,
						Start:    ln.Start,
						End:      ln.End,
						Size:     lineSize,
						Emphasis: ln.Emphasis,
					}
				}

				imageURL := seg.ImageURL
				if imageURL == "" && seg.VideoURL != "" {
					imageURL = seg.VideoURL
				}

				scenes[i] = subtitles.Scene{
					SceneID:  "scene_" + fmt.Sprintf("%d", i),
					Index:    i,
					Start:    seg.Start,
					End:      seg.End,
					Duration: seg.End - seg.Start,
					Voice: subtitles.SceneVoice{
						Text:         seg.Text,
						LanguageCode: languageCode,
					},
					Subtitle: subtitles.SceneSubtitle{
						Placement:  placement,
						Animation:  animation,
						Size:       size,
						Lines:      lines,
						Emphasis:   seg.Emphasis,
					},
					Visual: subtitles.SceneVisual{
						ImageURL:  imageURL,
						Animation: animation,
						VideoURL:  seg.VideoURL,
					},
					Layout: subtitles.SceneLayout{
						Placement: placement,
						SafeArea:  safeArea,
					},
				}
			}

			resp := subtitles.TimelineResponse{
				SubtitleSchemaVersion: subtitles.SchemaVersion,
				TimelineID:            "tl_" + uuid.New().String(),
				Version:               1,
				SubtitleStyle:         subtitleStyle,
				SubtitleLayout:        layout,
				RenderHints:           renderHints,
				Audio: subtitles.AudioInfo{
					URL:             audioURL,
					DurationSeconds: durationSec,
				},
				Scenes:          scenes,
				DurationSeconds: durationSec,
			}
			if directionPreset != "" {
				resp.DirectionPreset = directionPreset
			}
			if subtitleStyle == subtitles.StyleTrailer {
				resp.Theme = subtitles.GetThemePreset(subtitles.ThemeTrailerDefault)
			}
			return resp, nil
		},
		option.Summary("Subtitle existing audio"),
		option.Description("Transcribes an audio file from URL using Google Cloud Speech-to-Text and returns subtitles with timing"),
		option.Tags("speech"),
		option.Middleware(jwtAuthMiddleware),
	)
}

func resolveSubtitleStyle(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	if subtitles.ValidStyles[s] {
		return s
	}
	return subtitles.DefaultStyle
}

func resolvePlacement(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	if subtitles.ValidPlacements[s] {
		return s
	}
	return subtitles.PlacementBottom
}

func resolvePlacementStrategy(s string, style string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == subtitles.PlacementStrategyDynamic {
		return subtitles.PlacementStrategyDynamic
	}
	if s == subtitles.PlacementStrategyFixed {
		return subtitles.PlacementStrategyFixed
	}
	if style == subtitles.StyleCinematicDynamic || style == subtitles.StyleTrailer {
		return subtitles.PlacementStrategyDynamic
	}
	return subtitles.PlacementStrategyFixed
}

func resolveOverflowStrategy(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case subtitles.OverflowRebalance, subtitles.OverflowShrink, subtitles.OverflowRebalanceShrink:
		return s
	}
	return subtitles.OverflowRebalanceShrink
}

func resolveDirectionPreset(s string, style string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == subtitles.DirectionPresetCinematicDynamicV1 {
		return subtitles.DirectionPresetCinematicDynamicV1
	}
	if s == subtitles.DirectionPresetTrailerV1 {
		return subtitles.DirectionPresetTrailerV1
	}
	if s == subtitles.DirectionPresetDocumentary {
		return subtitles.DirectionPresetDocumentary
	}
	// Por defecto según estilo
	if style == subtitles.StyleCinematicDynamic {
		return subtitles.DirectionPresetCinematicDynamicV1
	}
	if style == subtitles.StyleTrailer {
		return subtitles.DirectionPresetTrailerV1
	}
	return ""
}

// splitScriptIntoLines divide el guion en líneas (por \n), filtra vacías.
func splitScriptIntoLines(script string) []string {
	raw := strings.Split(script, "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		t := strings.TrimSpace(line)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

// findNearestFrameImage devuelve la URL del frame cuyo timestamp está más cerca de segmentStart
func findNearestFrameImage(segmentStart float64, frames []FrameImage) string {
	if len(frames) == 0 {
		return ""
	}
	best := frames[0]
	bestDiff := abs(segmentStart - best.TimestampSec)
	for _, f := range frames[1:] {
		if f.ImageURL == "" {
			continue
		}
		diff := abs(segmentStart - f.TimestampSec)
		if diff < bestDiff {
			bestDiff = diff
			best = f
		}
	}
	return best.ImageURL
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func buildFullTranscript(subs []subtitles.SubtitleSegment) string {
	parts := make([]string, 0, len(subs))
	for _, s := range subs {
		if t := strings.TrimSpace(s.Text); t != "" {
			parts = append(parts, t)
		}
	}
	return strings.Join(parts, " ")
}
