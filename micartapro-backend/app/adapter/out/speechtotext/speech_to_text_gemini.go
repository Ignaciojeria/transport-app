package speechtotext

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"micartapro/app/shared/infrastructure/ai"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/genai"
)

const geminiAudioModel = "gemini-2.5-pro"
const maxInlineAudioBytes = 20 * 1024 * 1024 // 20 MB
const maxDownloadBytes = 50 * 1024 * 1024    // 50 MB — límite de descarga para evitar OOM

// geminiTranscriptionResponse estructura de la respuesta JSON de Gemini.
type geminiTranscriptionResponse struct {
	Summary  string `json:"summary"`
	Segments []struct {
		Speaker            string  `json:"speaker"`
		StartSec           float64 `json:"start_seconds"`
		EndSec             float64 `json:"end_seconds"`
		DurationSec        float64 `json:"duration_seconds,omitempty"`
		SilenceAfterSec    float64 `json:"silence_after_seconds,omitempty"`
		Content            string  `json:"content"`
		Language           string  `json:"language"`
		LanguageCode       string  `json:"language_code"`
		LanguageConfidence float64 `json:"language_confidence,omitempty"`
		Emotion            string  `json:"emotion"`
		EmotionConfidence  float64 `json:"emotion_confidence,omitempty"`
		Confidence         float64 `json:"confidence,omitempty"`
		Words              []struct {
			Text       string  `json:"text"`
			Start      float64 `json:"start"`
			End        float64 `json:"end"`
			Confidence float64 `json:"confidence,omitempty"`
		} `json:"words,omitempty"`
	} `json:"segments"`
}

func init() {
	ioc.Registry(NewSpeechToTextGemini, ai.NewClient, observability.NewObservability)
}

// NewSpeechToTextGemini transcribe audio usando Gemini Audio Understanding.
// Alternativa a Chirp cuando este falla o para probar. Usa el mismo tipo TranscribeAudio.
func NewSpeechToTextGemini(genaiClient *genai.Client, obs observability.Observability) (GeminiTranscribe, error) {
	impl := func(ctx context.Context, audioURL string, languageCode string) ([]SubtitleSegment, float64, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "transcribe_audio_gemini")
		defer span.End()

		if strings.TrimSpace(audioURL) == "" {
			return nil, 0, fmt.Errorf("audioUrl is required")
		}
		if strings.TrimSpace(languageCode) == "" {
			languageCode = "es-ES"
		}

		if genaiClient == nil {
			return nil, 0, fmt.Errorf("genai client is not initialized (check GOOGLE_PROJECT_ID)")
		}

		// 1. Descargar audio
		audioBytes, mimeType, err := downloadAudio(spanCtx, audioURL, obs)
		if err != nil {
			return nil, 0, err
		}
		obs.Logger.InfoContext(spanCtx, "audio_downloaded_for_gemini", "bytes", len(audioBytes), "mimeType", mimeType)

		// 2. Construir parts: prompt + audio (inline o subido)
		parts, err := buildAudioParts(spanCtx, genaiClient, audioBytes, mimeType, obs)
		if err != nil {
			return nil, 0, err
		}

		prompt := buildTranscriptionPrompt(languageCode)
		contents := []*genai.Content{
			{
				Role:  "user",
				Parts: append([]*genai.Part{genai.NewPartFromText(prompt)}, parts...),
			},
		}

		// 3. Schema para respuesta estructurada (start_seconds, end_seconds por segmento)
		responseSchema := transcriptionResponseSchema()

		config := &genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema:   responseSchema,
		}

		obs.Logger.InfoContext(spanCtx, "gemini_transcription_start", "model", geminiAudioModel)
		resp, err := genaiClient.Models.GenerateContent(spanCtx, geminiAudioModel, contents, config)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "gemini_transcription_error", "error", err)
			return nil, 0, fmt.Errorf("gemini transcription: %w", err)
		}

		text := resp.Text()
		if text == "" {
			obs.Logger.WarnContext(spanCtx, "gemini_empty_response")
			return nil, 0, nil
		}

		segments, transDuration, err := parseGeminiResponse(text, obs, spanCtx)
		if err != nil {
			return nil, 0, err
		}
		// Usar duración real del archivo si es mayor (evita 48s cuando el audio dura 50s con silencio al final)
		actualDuration := getAudioDurationSeconds(audioBytes, mimeType)
		durationSec := transDuration
		if actualDuration > transDuration {
			durationSec = actualDuration
			obs.Logger.InfoContext(spanCtx, "using_actual_audio_duration", "actual_sec", actualDuration, "trans_sec", transDuration)
		}
		// Confiamos en el modelo: no clampear duration
		obs.Logger.InfoContext(spanCtx, "transcription_complete", "segments", len(segments), "duration_sec", durationSec, "model", geminiAudioModel)
		return segments, durationSec, nil
	}
	return GeminiTranscribe{TranscribeAudio: impl}, nil
}

func downloadAudio(ctx context.Context, audioURL string, obs observability.Observability) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, audioURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("creating request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("downloading audio: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("audio download failed: status %d", resp.StatusCode)
	}
	limitReader := io.LimitReader(resp.Body, maxDownloadBytes+1)
	audioBytes, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, "", fmt.Errorf("reading audio: %w", err)
	}
	if int64(len(audioBytes)) > maxDownloadBytes {
		return nil, "", fmt.Errorf("audio file exceeds max size (%d MB)", maxDownloadBytes/(1024*1024))
	}

	urlPath := strings.Split(audioURL, "?")[0]
	contentType := resp.Header.Get("Content-Type")
	mimeType := audioMimeType(urlPath, contentType)
	return audioBytes, mimeType, nil
}

func audioMimeType(urlPath, contentType string) string {
	lower := strings.ToLower(urlPath)
	ct := strings.ToLower(contentType)
	switch {
	case strings.HasSuffix(lower, ".wav") || strings.Contains(ct, "wav"):
		return "audio/wav"
	case strings.HasSuffix(lower, ".flac") || strings.Contains(ct, "flac"):
		return "audio/flac"
	case strings.HasSuffix(lower, ".mp3") || strings.Contains(ct, "mpeg") || strings.Contains(ct, "mp3"):
		return "audio/mp3"
	case strings.HasSuffix(lower, ".ogg") || strings.Contains(ct, "ogg"):
		return "audio/ogg"
	default:
		return "audio/wav"
	}
}

func buildAudioParts(ctx context.Context, client *genai.Client, audioBytes []byte, mimeType string, obs observability.Observability) ([]*genai.Part, error) {
	if len(audioBytes) <= maxInlineAudioBytes {
		obs.Logger.InfoContext(ctx, "gemini_using_inline_audio", "bytes", len(audioBytes))
		return []*genai.Part{genai.NewPartFromBytes(audioBytes, mimeType)}, nil
	}
	// Archivo grande: subir vía Files API
	obs.Logger.InfoContext(ctx, "gemini_uploading_audio", "bytes", len(audioBytes))
	uploaded, err := client.Files.Upload(ctx, bytes.NewReader(audioBytes), &genai.UploadFileConfig{
		MIMEType: mimeType,
	})
	if err != nil {
		return nil, fmt.Errorf("uploading audio to Gemini: %w", err)
	}
	return []*genai.Part{genai.NewPartFromURI(uploaded.URI, uploaded.MIMEType)}, nil
}

func buildTranscriptionPrompt(languageCode string) string {
	langHint := "Spanish (es-ES)"
	if strings.HasPrefix(strings.ToLower(languageCode), "en") {
		langHint = "English"
	} else if strings.HasPrefix(strings.ToLower(languageCode), "pt") {
		langHint = "Portuguese"
	}
	if strings.TrimSpace(languageCode) == "" {
		languageCode = "es-ES"
	}
	return fmt.Sprintf(`Act as an expert audio transcriber. Generate a high-fidelity transcription with precise timing for subtitles.

AUDIO CONTEXT:
- Language: %s (language_code: %s). Include language_confidence (0-1) if multi-language.
- Task: Speech-to-Text with speaker diarization and word-level timestamps.

CRITICAL — Timing accuracy:
1) start_seconds = exact moment the first phoneme is heard. Place it slightly earlier (0.05–0.1 s before audible onset) so subtitles appear just before the speaker speaks.
2) end_seconds = when the last word ends. Use decimals (e.g., 2.34, 5.67). Round to max 3 decimal places.
3) Preserve pauses: if there is 1 s of silence, there MUST be a 1 s gap between segments. Do NOT compress time.
4) Words array: each word has {text, start, end} in seconds. Words are sequential (no overlaps), within segment bounds.
5) Do NOT output text for silence. Only output segments where speech is clearly audible.
6) Segments must be strictly sequential: segment[i].start_seconds >= segment[i-1].end_seconds. No overlaps allowed.

Speaker diarization:
- Use SPEAKER_01, SPEAKER_02, SPEAKER_03, etc. consistently. Do not invent names.

Cleanliness:
- Do NOT transcribe filler words (uh, um, ah, eh) unless they carry emotional weight or are part of a phrase.

Segmentation (TRANSCRIPTION layer — do not optimize for visual styling):
- Each segment = complete sentence OR complete clause that can stand alone on screen.
- Split ONLY at: (1) sentence end (. ! ?), (2) speaker change, (3) long pause (>= 2.0 seconds).
- Do NOT split on commas unless there is a long pause. Keep punctuation in content.

Emotion label:
- One per segment: happy, sad, angry, neutral. If unsure, use neutral. Include emotion_confidence (0-1). If < 0.7, treat as neutral.

Output JSON — each segment must include:
- speaker, start_seconds, end_seconds, duration_seconds (= end - start)
- silence_after_seconds (= seconds between this end and next start; 0 for last segment)
- content, language, language_code, language_confidence (optional)
- emotion, emotion_confidence, confidence (segment 0-1)
- words: [{text, start, end, confidence?}]`, langHint, languageCode)
}

func transcriptionResponseSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"summary": {
				Type:        genai.TypeString,
				Description: "A concise summary of the audio content.",
			},
			"segments": {
				Type:        genai.TypeArray,
				Description: "List of transcribed segments (transcription layer).",
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"speaker":               {Type: genai.TypeString, Description: "Speaker identifier: SPEAKER_01, SPEAKER_02, etc."},
						"start_seconds":         {Type: genai.TypeNumber, Description: "Start time in seconds (max 3 decimals)"},
						"end_seconds":           {Type: genai.TypeNumber, Description: "End time in seconds (max 3 decimals)"},
						"duration_seconds":      {Type: genai.TypeNumber, Description: "end_seconds - start_seconds"},
						"silence_after_seconds": {Type: genai.TypeNumber, Description: "Seconds between this end and next start; 0 for last segment"},
						"content":               {Type: genai.TypeString, Description: "Transcribed text with punctuation"},
						"language":              {Type: genai.TypeString},
						"language_code":         {Type: genai.TypeString},
						"language_confidence":   {Type: genai.TypeNumber, Description: "Language detection confidence (0-1), optional"},
						"emotion": {
							Type:        genai.TypeString,
							Description: "happy, sad, angry, or neutral. Use neutral if unsure.",
							Enum:        []string{"happy", "sad", "angry", "neutral"},
						},
						"emotion_confidence": {
							Type:        genai.TypeNumber,
							Description: "Confidence in emotion label (0-1). If < 0.7, use neutral.",
						},
						"confidence": {
							Type:        genai.TypeNumber,
							Description: "Segment transcription confidence (0-1).",
						},
						"words": {
							Type:        genai.TypeArray,
							Description: "Word-level timing. Each word: text, start, end (seconds).",
							Items: &genai.Schema{
								Type: genai.TypeObject,
								Properties: map[string]*genai.Schema{
									"text":       {Type: genai.TypeString, Description: "Word text"},
									"start":      {Type: genai.TypeNumber, Description: "Start time in seconds"},
									"end":        {Type: genai.TypeNumber, Description: "End time in seconds"},
									"confidence": {Type: genai.TypeNumber, Description: "Word confidence (0-1), optional"},
								},
								Required: []string{"text", "start", "end"},
							},
						},
					},
					Required: []string{"speaker", "start_seconds", "end_seconds", "content", "language", "language_code", "emotion"},
				},
			},
		},
		Required: []string{"summary", "segments"},
	}
}

const timestampPrecision = 1000.0 // 3 decimales: round(x*1000)/1000

// roundTo3Decimals redondea a máximo 3 decimales (evita 2.3456789123).
func roundTo3Decimals(x float64) float64 {
	return math.Round(x*timestampPrecision) / timestampPrecision
}

// normalizeWordTimestamps pasa los timestamps del modelo sin ajustes. Confiamos en el modelo.
func normalizeWordTimestamps(words []SegmentWord, _, _ float64) []SegmentWord {
	if len(words) == 0 {
		return nil
	}
	result := make([]SegmentWord, 0, len(words))
	for _, w := range words {
		if strings.TrimSpace(w.Text) == "" {
			continue
		}
		result = append(result, SegmentWord{
			Text:  w.Text,
			Start: roundTo3Decimals(w.Start),
			End:   roundTo3Decimals(w.End),
		})
	}
	return result
}

func parseGeminiResponse(text string, obs observability.Observability, ctx context.Context) ([]SubtitleSegment, float64, error) {
	var parsed geminiTranscriptionResponse
	if err := json.Unmarshal([]byte(text), &parsed); err != nil {
		preview := text
		if len(text) > 200 {
			preview = text[:200]
		}
		obs.Logger.ErrorContext(ctx, "gemini_response_parse_error", "error", err, "text_preview", preview)
		return nil, 0, fmt.Errorf("parsing gemini response: %w", err)
	}

	var segments []SubtitleSegment
	var durationSec float64
	for _, s := range parsed.Segments {
		if s.Content == "" {
			continue
		}
		start := roundTo3Decimals(s.StartSec)
		end := roundTo3Decimals(s.EndSec)
		if end > durationSec {
			durationSec = end
		}

		duration := roundTo3Decimals(end - start)
		silenceAfter := roundTo3Decimals(s.SilenceAfterSec)

		emotion := s.Emotion
		if s.EmotionConfidence > 0 && s.EmotionConfidence < 0.7 {
			emotion = "neutral"
		}
		seg := SubtitleSegment{
			Text:                strings.TrimSpace(s.Content),
			Start:               start,
			End:                 end,
			DurationSeconds:     duration,
			SilenceAfterSeconds: silenceAfter,
			Confidence:          s.Confidence,
			Emotion:             emotion,
			EmotionConfidence:   s.EmotionConfidence,
		}
		if len(s.Words) > 0 {
			rawWords := make([]SegmentWord, len(s.Words))
			for j, w := range s.Words {
				rawWords[j] = SegmentWord{
					Text:  w.Text,
					Start: roundTo3Decimals(w.Start),
					End:   roundTo3Decimals(w.End),
				}
			}
			seg.Words = normalizeWordTimestamps(rawWords, start, end)
		}
		segments = append(segments, seg)
	}
	return segments, durationSec, nil
}
