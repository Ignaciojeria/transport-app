package texttospeech

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"micartapro/app/shared/infrastructure/ai"
	"micartapro/app/shared/infrastructure/gcs"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"google.golang.org/genai"
)

const ttsModel = "gemini-2.5-flash-preview-tts"

// Voces predefinidas de Gemini TTS (ej: Sadachbia, Puck, Kore, Zephyr, Aoede, etc.)
const defaultVoice = "Sadachbia"

// delayBetweenTTSLines evita 429 (quota exceeded) al hacer muchas llamadas por minuto.
const delayBetweenTTSLines = 6 * time.Second

// GenerateSpeech convierte texto a audio usando Gemini TTS, lo sube a GCS y retorna la URL pública y duración.
type GenerateSpeech func(ctx context.Context, text string, opts *GenerateSpeechOptions) (*GenerateSpeechResult, error)

// GenerateSpeechResult resultado de la síntesis de voz.
type GenerateSpeechResult struct {
	AudioURL        string  `json:"audioUrl"`
	DurationSeconds float64 `json:"durationSeconds"`
}

// LineTiming representa el timing exacto de una línea (desde TTS por línea).
type LineTiming struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// GenerateSpeechForLinesResult resultado cuando se genera TTS por línea (sync exacto).
type GenerateSpeechForLinesResult struct {
	AudioURL        string       `json:"audioUrl"`
	DurationSeconds float64      `json:"durationSeconds"`
	LineTimings     []LineTiming `json:"lineTimings"` // start/end exacto por línea
}

// GenerateSpeechForLines genera TTS por línea y concatena. Retorna timing exacto por segmento.
type GenerateSpeechForLines func(ctx context.Context, lines []string, opts *GenerateSpeechOptions) (*GenerateSpeechForLinesResult, error)

// GenerateSpeechOptions opciones para la síntesis de voz.
type GenerateSpeechOptions struct {
	VoiceName   string // nombre de la voz (ej. Kore, Puck, Aoede). Default: Kore
	LanguageCode string // BCP-47 (ej. "es", "en"). Default: "es" (español)
	StylePrompt string // instrucciones de estilo (ej. "Say cheerfully:", "Say in a spooky whisper:")
}

func init() {
	ioc.Registry(NewTextToSpeech, ai.NewClient, observability.NewObservability, gcs.NewClient)
	ioc.Registry(NewTextToSpeechForLines, ai.NewClient, observability.NewObservability, gcs.NewClient)
}

func NewTextToSpeech(
	genaiClient *genai.Client,
	obs observability.Observability,
	gcsClient *storage.Client,
) (GenerateSpeech, error) {
	return func(ctx context.Context, text string, opts *GenerateSpeechOptions) (*GenerateSpeechResult, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "generate_speech")
		defer span.End()

		userID, ok := sharedcontext.UserIDFromContext(spanCtx)
		if !ok || userID == "" {
			return nil, fmt.Errorf("userID is required but not found in context")
		}

		if strings.TrimSpace(text) == "" {
			return nil, fmt.Errorf("text is required and cannot be empty")
		}

		if genaiClient == nil {
			return nil, fmt.Errorf("genai client is not initialized (check GOOGLE_PROJECT_ID)")
		}

		voiceName := defaultVoice
		if opts != nil && opts.VoiceName != "" {
			voiceName = opts.VoiceName
		}

		lang := "es"
		if opts != nil && opts.LanguageCode != "" {
			lang = opts.LanguageCode
		}

		// Construir el prompt: idioma (default español) + estilo opcional + texto
		// Ej: "Say in Spanish: Hola mundo" o "Say in Spanish, cheerfully: ¡Buenos días!"
		promptText := buildTTSPrompt(lang, opts != nil && opts.StylePrompt != "", opts.StylePrompt, text)

		obs.Logger.InfoContext(spanCtx, "generating_speech", "text_length", len(text), "voice", voiceName, "language", lang)

		pcmBytes, err := synthesizeToPCM(spanCtx, genaiClient, promptText, voiceName, obs)
		if err != nil {
			return nil, err
		}

		// Duración: PCM s16le 24kHz mono = 48000 bytes/seg
		const bytesPerSecond = 24000 * 2
		durationSeconds := float64(len(pcmBytes)) / float64(bytesPerSecond)

		obs.Logger.InfoContext(spanCtx, "speech_synthesized", "pcm_bytes", len(pcmBytes), "duration_sec", durationSeconds)

		// 3. Convertir PCM (s16le, 24000Hz, mono) a WAV
		wavBytes := pcmToWAV(pcmBytes, 24000, 1, 16)

		// 4. Subir a GCS
		contentType := "audio/wav"
		fileName := fmt.Sprintf("%d.wav", time.Now().UnixNano())
		uploadURL, publicURL, _, err := GenerateSignedWriteURLAudio(spanCtx, gcsClient, obs, userID, fileName, contentType)
		if err != nil {
			return nil, fmt.Errorf("error generating upload URL: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(spanCtx, "PUT", uploadURL, bytes.NewReader(wavBytes))
		if err != nil {
			return nil, fmt.Errorf("error creating upload request: %w", err)
		}
		httpReq.Header.Set("Content-Type", contentType)

		httpClient := &http.Client{Timeout: 30 * time.Second}
		uploadResp, err := httpClient.Do(httpReq)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_uploading_audio", "error", err)
			return nil, fmt.Errorf("error uploading audio: %w", err)
		}
		body, _ := io.ReadAll(uploadResp.Body)
		uploadResp.Body.Close()

		if uploadResp.StatusCode < 200 || uploadResp.StatusCode >= 300 {
			obs.Logger.ErrorContext(spanCtx, "error_upload_status", "status", uploadResp.StatusCode, "response_body", string(body))
			return nil, fmt.Errorf("error uploading audio: status %d: %s", uploadResp.StatusCode, string(body))
		}

		obs.Logger.InfoContext(spanCtx, "audio_uploaded_successfully", "publicURL", publicURL, "size_bytes", len(wavBytes))
		return &GenerateSpeechResult{AudioURL: publicURL, DurationSeconds: durationSeconds}, nil
	}, nil
}

// NewTextToSpeechForLines genera TTS por línea, concatena PCM y retorna timing exacto por línea.
func NewTextToSpeechForLines(
	genaiClient *genai.Client,
	obs observability.Observability,
	gcsClient *storage.Client,
) (GenerateSpeechForLines, error) {
	const bytesPerSecond = 24000 * 2 // PCM s16le 24kHz mono
	return func(ctx context.Context, lines []string, opts *GenerateSpeechOptions) (*GenerateSpeechForLinesResult, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "generate_speech_for_lines")
		defer span.End()

		userID, ok := sharedcontext.UserIDFromContext(spanCtx)
		if !ok || userID == "" {
			return nil, fmt.Errorf("userID is required but not found in context")
		}
		if len(lines) == 0 {
			return nil, fmt.Errorf("lines cannot be empty")
		}
		if genaiClient == nil {
			return nil, fmt.Errorf("genai client is not initialized (check GOOGLE_PROJECT_ID)")
		}

		voiceName := defaultVoice
		if opts != nil && opts.VoiceName != "" {
			voiceName = opts.VoiceName
		}
		lang := "es"
		if opts != nil && opts.LanguageCode != "" {
			lang = opts.LanguageCode
		}

		var allPCM []byte
		lineTimings := make([]LineTiming, 0, len(lines))
		start := 0.0

		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				lineTimings = append(lineTimings, LineTiming{Start: start, End: start})
				continue
			}
			// Rate limit: Gemini TTS tiene ~10 req/min. Delay entre líneas evita 429.
			if i > 0 {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delayBetweenTTSLines):
				}
			}
			stylePrompt := "fluently and briskly, with minimal pauses between phrases"
			if opts != nil && opts.StylePrompt != "" {
				stylePrompt = opts.StylePrompt
			}
			promptText := buildTTSPrompt(lang, true, stylePrompt, line)
			pcm, err := synthesizeToPCM(spanCtx, genaiClient, promptText, voiceName, obs)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", i+1, err)
			}
			dur := float64(len(pcm)) / float64(bytesPerSecond)
			lineTimings = append(lineTimings, LineTiming{Start: start, End: start + dur})
			start += dur
			allPCM = append(allPCM, pcm...)
		}

		durationSeconds := float64(len(allPCM)) / float64(bytesPerSecond)
		obs.Logger.InfoContext(spanCtx, "speech_synthesized_for_lines", "lines", len(lines), "pcm_bytes", len(allPCM), "duration_sec", durationSeconds)

		wavBytes := pcmToWAV(allPCM, 24000, 1, 16)
		contentType := "audio/wav"
		fileName := fmt.Sprintf("%d.wav", time.Now().UnixNano())
		uploadURL, publicURL, _, err := GenerateSignedWriteURLAudio(spanCtx, gcsClient, obs, userID, fileName, contentType)
		if err != nil {
			return nil, fmt.Errorf("error generating upload URL: %w", err)
		}
		httpReq, err := http.NewRequestWithContext(spanCtx, "PUT", uploadURL, bytes.NewReader(wavBytes))
		if err != nil {
			return nil, fmt.Errorf("error creating upload request: %w", err)
		}
		httpReq.Header.Set("Content-Type", contentType)
		httpClient := &http.Client{Timeout: 30 * time.Second}
		uploadResp, err := httpClient.Do(httpReq)
		if err != nil {
			obs.Logger.ErrorContext(spanCtx, "error_uploading_audio", "error", err)
			return nil, fmt.Errorf("error uploading audio: %w", err)
		}
		body, _ := io.ReadAll(uploadResp.Body)
		uploadResp.Body.Close()
		if uploadResp.StatusCode < 200 || uploadResp.StatusCode >= 300 {
			obs.Logger.ErrorContext(spanCtx, "error_upload_status", "status", uploadResp.StatusCode, "response_body", string(body))
			return nil, fmt.Errorf("error uploading audio: status %d: %s", uploadResp.StatusCode, string(body))
		}
		obs.Logger.InfoContext(spanCtx, "audio_uploaded_successfully", "publicURL", publicURL, "size_bytes", len(wavBytes))
		return &GenerateSpeechForLinesResult{AudioURL: publicURL, DurationSeconds: durationSeconds, LineTimings: lineTimings}, nil
	}, nil
}

// synthesizeToPCM llama a Gemini TTS y retorna PCM s16le 24kHz mono.
func synthesizeToPCM(ctx context.Context, genaiClient *genai.Client, promptText, voiceName string, obs observability.Observability) ([]byte, error) {
	contents := []*genai.Content{
		{Role: "user", Parts: []*genai.Part{{Text: promptText}}},
	}
	config := &genai.GenerateContentConfig{
		ResponseModalities: []string{string(genai.ModalityAudio)},
		SpeechConfig: &genai.SpeechConfig{
			VoiceConfig: &genai.VoiceConfig{
				PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{
					VoiceName: voiceName,
				},
			},
		},
	}
	resp, err := genaiClient.Models.GenerateContent(ctx, ttsModel, contents, config)
	if err != nil {
		obs.Logger.ErrorContext(ctx, "error_synthesizing_speech", "error", err, "text_preview", promptText[:min(50, len(promptText))]+"...")
		return nil, fmt.Errorf("error synthesizing speech: %w", err)
	}
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		obs.Logger.ErrorContext(ctx, "no_candidates_or_content", "text_preview", promptText[:min(50, len(promptText))]+"...")
		return nil, fmt.Errorf("no audio generated (empty response)")
	}
	var pcmBytes []byte
	for _, p := range resp.Candidates[0].Content.Parts {
		if p != nil && p.InlineData != nil && len(p.InlineData.Data) > 0 {
			pcmBytes = p.InlineData.Data
			break
		}
	}
	if len(pcmBytes) == 0 {
		obs.Logger.ErrorContext(ctx, "audio_content_empty", "text_preview", promptText[:min(50, len(promptText))]+"...")
		return nil, fmt.Errorf("synthesized audio is empty")
	}
	return pcmBytes, nil
}

// buildTTSPrompt construye el prompt con instrucción de idioma y estilo opcional.
func buildTTSPrompt(lang string, hasStyle bool, stylePrompt, text string) string {
	langBase := languageBase(lang)
	if !hasStyle || stylePrompt == "" {
		return langBase + " " + text
	}
	// Combinar idioma y estilo: "Say in Spanish, cheerfully: texto"
	style := strings.TrimPrefix(strings.TrimSpace(stylePrompt), "Say ")
	style = strings.TrimSuffix(style, ":")
	if style != "" {
		return langBase + ", " + style + ": " + text
	}
	return langBase + " " + text
}

func languageBase(lang string) string {
	switch strings.ToLower(lang) {
	case "es", "es-cl", "es-es", "es-mx":
		return "Say in Spanish"
	case "en", "en-us", "en-gb":
		return "Say in English"
	case "pt", "pt-br", "pt-pt":
		return "Say in Portuguese"
	case "fr":
		return "Say in French"
	case "de":
		return "Say in German"
	case "it":
		return "Say in Italian"
	case "ja":
		return "Say in Japanese"
	default:
		return "Say in Spanish" // fallback a español
	}
}

// pcmToWAV convierte PCM s16le a formato WAV con header.
func pcmToWAV(pcm []byte, sampleRate, numChannels, bitsPerSample int) []byte {
	dataSize := len(pcm)
	byteRate := sampleRate * numChannels * (bitsPerSample / 8)
	blockAlign := numChannels * (bitsPerSample / 8)
	// Chunk fmt: 16 bytes
	// Chunk data: 8 bytes header + dataSize
	// RIFF header: 4 + 4 + 4 = 12, fmt chunk: 4 + 4 + 16 = 24, data chunk: 4 + 4 + dataSize
	fileSize := 4 + 24 + 8 + dataSize // "WAVE" + fmt + data header + data

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, []byte("RIFF"))
	binary.Write(buf, binary.LittleEndian, uint32(fileSize))
	binary.Write(buf, binary.LittleEndian, []byte("WAVE"))
	binary.Write(buf, binary.LittleEndian, []byte("fmt "))
	binary.Write(buf, binary.LittleEndian, uint32(16)) // chunk size
	binary.Write(buf, binary.LittleEndian, uint16(1))  // PCM format
	binary.Write(buf, binary.LittleEndian, uint16(numChannels))
	binary.Write(buf, binary.LittleEndian, uint32(sampleRate))
	binary.Write(buf, binary.LittleEndian, uint32(byteRate))
	binary.Write(buf, binary.LittleEndian, uint16(blockAlign))
	binary.Write(buf, binary.LittleEndian, uint16(bitsPerSample))
	binary.Write(buf, binary.LittleEndian, []byte("data"))
	binary.Write(buf, binary.LittleEndian, uint32(dataSize))
	buf.Write(pcm)
	return buf.Bytes()
}

// GenerateSignedWriteURLAudio genera signed URL para subir audio a GCS.
func GenerateSignedWriteURLAudio(ctx context.Context, client *storage.Client, obs observability.Observability, userID string, fileName string, contentType string) (uploadURL string, publicURL string, objectPath string, err error) {
	if client == nil {
		return "", "", "", fmt.Errorf("GCS client is nil")
	}
	timestamp := time.Now().Unix()
	randomSuffix := fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
	objectPath = fmt.Sprintf("audio/%s/%d-%s-%s", userID, timestamp, randomSuffix, fileName)
	bucketName := "micartapro-images"

	opts := &storage.SignedURLOptions{
		Method:      "PUT",
		Expires:     time.Now().Add(15 * time.Minute),
		ContentType: contentType,
	}
	uploadURL, err = client.Bucket(bucketName).SignedURL(objectPath, opts)
	if err != nil {
		return "", "", "", fmt.Errorf("generating signed URL: %w", err)
	}
	publicURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)
	obs.Logger.InfoContext(ctx, "signed_audio_upload_url_generated", "objectPath", objectPath, "userID", userID, "contentType", contentType)
	return uploadURL, publicURL, objectPath, nil
}
