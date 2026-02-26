package speechtotext

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/gcs"
	"micartapro/app/shared/infrastructure/observability"

	"cloud.google.com/go/storage"

	speechv2 "cloud.google.com/go/speech/apiv2"
	speechpbv2 "cloud.google.com/go/speech/apiv2/speechpb"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/durationpb"
)

const transcriptionBucket = "micartapro-images"
const transcriptionPrefix = "audio/transcription-temp/"
const chirp3Model = "chirp_3"

// SegmentWord palabra con timing (Chirp).
type SegmentWord struct {
	Text  string
	Start float64
	End   float64
}

// SubtitleSegment representa un segmento de subtítulo con timing.
type SubtitleSegment struct {
	Text                string        `json:"text"`
	Start               float64       `json:"start"`
	End                 float64       `json:"end"`
	DurationSeconds     float64       `json:"-"` // Recalculado: end - start (nunca confiar en el modelo)
	Words               []SegmentWord `json:"-"`
	SilenceAfterSeconds float64       `json:"-"` // Segundos entre este segmento y el siguiente (para trailer gaps)
	Confidence          float64       `json:"-"` // Confianza del segmento (0-1)
	Emotion             string        `json:"-"` // happy, sad, angry, neutral
	EmotionConfidence   float64       `json:"-"` // Si < 0.7 → tratar como neutral
}

// TranscribeAudio transcribe audio desde una URL y retorna segmentos con timing.
type TranscribeAudio func(ctx context.Context, audioURL string, languageCode string) ([]SubtitleSegment, float64, error)

// ChirpTranscribe envuelve la implementación Chirp.
type ChirpTranscribe struct {
	TranscribeAudio TranscribeAudio
}

// GeminiTranscribe envuelve la implementación Gemini.
type GeminiTranscribe struct {
	TranscribeAudio TranscribeAudio
}

func init() {
	ioc.Registry(NewSpeechToText, observability.NewObservability, gcs.NewClient, configuration.NewConf)
}

func NewSpeechToText(obs observability.Observability, gcsClient *storage.Client, conf configuration.Conf) (ChirpTranscribe, error) {
	impl := func(ctx context.Context, audioURL string, languageCode string) ([]SubtitleSegment, float64, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "transcribe_audio")
		defer span.End()

		if strings.TrimSpace(audioURL) == "" {
			return nil, 0, fmt.Errorf("audioUrl is required")
		}
		if lang := strings.TrimSpace(languageCode); lang == "" {
			languageCode = "es-ES"
		}

		// Si la URL ya es de nuestro bucket GCS, usar gs:// directamente (evita descarga + re-upload)
		gcsURI := ""
		if gcsURI, _ = urlToGCSURI(audioURL, transcriptionBucket); gcsURI != "" {
			obs.Logger.InfoContext(spanCtx, "using_existing_gcs_uri", "uri", gcsURI)
		} else {
			// Descargar y subir a GCS
			if gcsClient == nil {
				return nil, 0, fmt.Errorf("GCS client required for Chirp 3 transcription")
			}
			req, err := http.NewRequestWithContext(spanCtx, http.MethodGet, audioURL, nil)
			if err != nil {
				return nil, 0, fmt.Errorf("creating request: %w", err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, 0, fmt.Errorf("downloading audio: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return nil, 0, fmt.Errorf("audio download failed: status %d", resp.StatusCode)
			}
			audioBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, 0, fmt.Errorf("reading audio: %w", err)
			}
			obs.Logger.InfoContext(spanCtx, "audio_downloaded", "bytes", len(audioBytes), "url", audioURL)

			urlPath := strings.Split(audioURL, "?")[0]
			contentType := resp.Header.Get("Content-Type")
			ext := ".mp3"
			switch {
			case strings.HasSuffix(strings.ToLower(urlPath), ".wav") || strings.Contains(contentType, "wav"):
				ext = ".wav"
			case strings.HasSuffix(strings.ToLower(urlPath), ".flac") || strings.Contains(contentType, "flac"):
				ext = ".flac"
			case strings.HasSuffix(strings.ToLower(urlPath), ".mp3") || strings.Contains(contentType, "mpeg") || strings.Contains(contentType, "mp3"):
				ext = ".mp3"
			}
			gcsObjectPath := transcriptionPrefix + uuid.New().String() + ext
			if err := uploadToGCS(spanCtx, gcsClient, transcriptionBucket, gcsObjectPath, audioBytes, contentTypeFromExt(ext), obs); err != nil {
				return nil, 0, fmt.Errorf("uploading audio to GCS: %w", err)
			}
			gcsURI = fmt.Sprintf("gs://%s/%s", transcriptionBucket, gcsObjectPath)
			obs.Logger.InfoContext(spanCtx, "audio_uploaded_to_gcs", "uri", gcsURI)
			defer func() {
				delCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = gcsClient.Bucket(transcriptionBucket).Object(gcsObjectPath).Delete(delCtx)
			}()
		}

		// Cliente Speech V2: Chirp 3 solo soporta us, eu o global (NO us-central1).
		// Ver https://cloud.google.com/speech-to-text/v2/docs/chirp_3-model
		speechRegion := mapToSpeechRegion(conf.GOOGLE_PROJECT_LOCATION)
		var clientOpts []option.ClientOption
		if speechRegion != "global" {
			endpoint := speechRegion + "-speech.googleapis.com"
			clientOpts = append(clientOpts, option.WithEndpoint(endpoint))
			obs.Logger.InfoContext(spanCtx, "speech_client_endpoint", "endpoint", endpoint, "region", speechRegion)
		}
		client, err := speechv2.NewClient(spanCtx, clientOpts...)
		if err != nil {
			return nil, 0, fmt.Errorf("creating speech client: %w", err)
		}
		defer client.Close()

		recognizer := fmt.Sprintf("projects/%s/locations/%s/recognizers/_", conf.GOOGLE_PROJECT_ID, speechRegion)
		config := &speechpbv2.RecognitionConfig{
			DecodingConfig: &speechpbv2.RecognitionConfig_AutoDecodingConfig{
				AutoDecodingConfig: &speechpbv2.AutoDetectDecodingConfig{},
			},
			LanguageCodes: []string{languageCode},
			Model:         chirp3Model,
			Features: &speechpbv2.RecognitionFeatures{
				EnableWordTimeOffsets: true,
			},
		}

		batchReq := &speechpbv2.BatchRecognizeRequest{
			Recognizer: recognizer,
			Config:     config,
			Files: []*speechpbv2.BatchRecognizeFileMetadata{
				{AudioSource: &speechpbv2.BatchRecognizeFileMetadata_Uri{Uri: gcsURI}},
			},
			RecognitionOutputConfig: &speechpbv2.RecognitionOutputConfig{
				Output: &speechpbv2.RecognitionOutputConfig_InlineResponseConfig{
					InlineResponseConfig: &speechpbv2.InlineOutputConfig{},
				},
			},
		}

		op, err := client.BatchRecognize(spanCtx, batchReq)
		if err != nil {
			return nil, 0, fmt.Errorf("starting recognition: %w", err)
		}
		obs.Logger.InfoContext(spanCtx, "batch_recognize_started", "gcsUri", gcsURI, "recognizer", recognizer)

		ctxWait, cancel := context.WithTimeout(spanCtx, 10*time.Minute)
		defer cancel()
		batchResp, err := op.Wait(ctxWait)
		if err != nil {
			return nil, 0, fmt.Errorf("recognition failed: %w", err)
		}

		fileResult, ok := batchResp.Results[gcsURI]
		if !ok || fileResult == nil {
			// Log keys disponibles por si el formato de clave difiere
			keys := make([]string, 0, len(batchResp.Results))
			for k := range batchResp.Results {
				keys = append(keys, k)
			}
			obs.Logger.WarnContext(spanCtx, "speech_api_empty_results", "url", audioURL, "gcsUri", gcsURI, "resultKeys", keys)
			return nil, 0, nil
		}
		if fileErr := fileResult.GetError(); fileErr != nil {
			obs.Logger.ErrorContext(spanCtx, "speech_api_file_error", "message", fileErr.GetMessage(), "gcsUri", gcsURI)
			return nil, 0, fmt.Errorf("recognition error for file: %s", fileErr.GetMessage())
		}
		inlineResult := fileResult.GetInlineResult()
		if inlineResult == nil {
			obs.Logger.WarnContext(spanCtx, "speech_api_no_inline_result", "url", audioURL)
			return nil, 0, nil
		}
		transcript := inlineResult.GetTranscript()
		if transcript == nil || len(transcript.Results) == 0 {
			obs.Logger.WarnContext(spanCtx, "speech_api_empty_transcript", "url", audioURL)
			return nil, 0, nil
		}

		segments, durationSec := buildSegmentsFromV2Results(transcript.Results, obs, spanCtx)
		obs.Logger.InfoContext(spanCtx, "transcription_complete", "segments", len(segments), "duration_sec", durationSec, "model", chirp3Model)
		return segments, durationSec, nil
	}
	return ChirpTranscribe{TranscribeAudio: impl}, nil
}

func buildSegmentsFromV2Results(results []*speechpbv2.SpeechRecognitionResult, obs observability.Observability, ctx context.Context) ([]SubtitleSegment, float64) {
	var segments []SubtitleSegment
	var durationSec float64

	for _, result := range results {
		if len(result.Alternatives) == 0 {
			continue
		}
		alt := result.Alternatives[0]
		words := alt.Words
		if len(words) == 0 {
			transcript := strings.TrimSpace(alt.Transcript)
			if transcript != "" {
				estDuration := float64(len(transcript)) / 12.0
				if estDuration > durationSec {
					durationSec = estDuration
				}
				segments = append(segments, SubtitleSegment{
					Text:  transcript,
					Start: 0,
					End:   estDuration,
				})
				obs.Logger.InfoContext(ctx, "speech_fallback_no_words", "transcript_len", len(transcript))
			}
			continue
		}

		const maxWordsPerSegment = 8
		const minWordsPerSegment = 3
		segmentWords := []*speechpbv2.WordInfo{}
		segmentStart := durationToSeconds(words[0].StartOffset)
		lastEnd := segmentStart

		for _, w := range words {
			startSec := durationToSeconds(w.StartOffset)
			endSec := durationToSeconds(w.EndOffset)
			if endSec > durationSec {
				durationSec = endSec
			}
			pause := startSec - lastEnd
			if len(segmentWords) >= maxWordsPerSegment || (len(segmentWords) >= minWordsPerSegment && pause > 0.5) {
				if len(segmentWords) > 0 {
					segEnd := durationToSeconds(segmentWords[len(segmentWords)-1].EndOffset)
					seg := SubtitleSegment{
						Text:  wordsToTextV2(segmentWords),
						Start: segmentStart,
						End:   segEnd,
					}
					seg.Words = wordsToSubtitleWords(segmentWords)
					segments = append(segments, seg)
				}
				segmentWords = nil
				segmentStart = startSec
			}
			segmentWords = append(segmentWords, w)
			lastEnd = endSec
		}
		if len(segmentWords) > 0 {
			segEnd := durationToSeconds(segmentWords[len(segmentWords)-1].EndOffset)
			seg := SubtitleSegment{
				Text:  wordsToTextV2(segmentWords),
				Start: segmentStart,
				End:   segEnd,
			}
			seg.Words = wordsToSubtitleWords(segmentWords)
			segments = append(segments, seg)
		}
	}
	return segments, durationSec
}

func durationToSeconds(d *durationpb.Duration) float64 {
	if d == nil {
		return 0
	}
	return d.AsDuration().Seconds()
}

func wordsToTextV2(words []*speechpbv2.WordInfo) string {
	var b strings.Builder
	for i, w := range words {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(w.Word)
	}
	return strings.TrimSpace(b.String())
}

func wordsToSubtitleWords(words []*speechpbv2.WordInfo) []SegmentWord {
	out := make([]SegmentWord, len(words))
	for i, w := range words {
		out[i] = SegmentWord{
			Text:  w.Word,
			Start: durationToSeconds(w.StartOffset),
			End:   durationToSeconds(w.EndOffset),
		}
	}
	return out
}

func uploadToGCS(ctx context.Context, client *storage.Client, bucket, objectPath string, data []byte, contentType string, obs observability.Observability) error {
	w := client.Bucket(bucket).Object(objectPath).NewWriter(ctx)
	w.ContentType = contentType
	if _, err := w.Write(data); err != nil {
		_ = w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	obs.Logger.InfoContext(ctx, "gcs_upload_complete", "bucket", bucket, "object", objectPath, "bytes", len(data))
	return nil
}

func contentTypeFromExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	default:
		return "audio/mpeg"
	}
}

// mapToSpeechRegion mapea GOOGLE_PROJECT_LOCATION a región válida de Chirp 3.
// Chirp 3 solo soporta: us, eu, global (y algunas en preview).
// us-central1, us-east1, etc. no son válidos; se mapean a "us".
func mapToSpeechRegion(loc string) string {
	loc = strings.TrimSpace(strings.ToLower(loc))
	switch {
	case loc == "" || strings.HasPrefix(loc, "us"):
		return "us"
	case strings.HasPrefix(loc, "eu"):
		return "eu"
	default:
		return "global"
	}
}

// urlToGCSURI convierte URL de storage.googleapis.com a gs:// si es de nuestro bucket.
// Retorna ("", false) si no aplica.
func urlToGCSURI(url, bucket string) (string, bool) {
	url = strings.TrimSpace(url)
	if url == "" {
		return "", false
	}
	// https://storage.googleapis.com/micartapro-images/audio/xxx/file.wav
	prefix := "https://storage.googleapis.com/" + bucket + "/"
	if !strings.HasPrefix(url, prefix) {
		// También soportar http
		prefixHTTP := "http://storage.googleapis.com/" + bucket + "/"
		if !strings.HasPrefix(url, prefixHTTP) {
			return "", false
		}
		prefix = prefixHTTP
	}
	path := strings.TrimPrefix(url, prefix)
	path = strings.Split(path, "?")[0] // quitar query params
	if path == "" {
		return "", false
	}
	return fmt.Sprintf("gs://%s/%s", bucket, path), false // false = no cleanup (archivo ya existía)
}
