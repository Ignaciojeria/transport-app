package subtitletiming

import (
	"context"
	"regexp"
	"strings"

	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

// GenerateSubtitleTiming genera segmentos de subtítulos con timing estimado.
// Input: transcript (texto original) + duración del audio en segundos.
// Output: segmentos { text, start, end } para sincronizar con el audio.
// No usa STT: el transcript es la fuente de verdad. Distribuye el tiempo proporcionalmente por frases.
type GenerateSubtitleTiming func(ctx context.Context, transcript string, durationSeconds float64) ([]SubtitleSegment, error)

// SubtitleSegment representa un segmento de subtítulo con timing.
type SubtitleSegment struct {
	Text  string  `json:"text"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

func init() {
	ioc.Registry(NewSubtitleTimingAgent, observability.NewObservability)
}

func NewSubtitleTimingAgent(obs observability.Observability) (GenerateSubtitleTiming, error) {
	return func(ctx context.Context, transcript string, durationSeconds float64) ([]SubtitleSegment, error) {
		spanCtx, span := obs.Tracer.Start(ctx, "generate_subtitle_timing")
		defer span.End()

		transcript = strings.TrimSpace(transcript)
		if transcript == "" {
			return nil, nil
		}
		if durationSeconds <= 0 {
			return []SubtitleSegment{{Text: transcript, Start: 0, End: 0}}, nil
		}

		// Dividir en frases (por . ! ? o saltos de línea)
		phrases := splitIntoPhrases(transcript)
		if len(phrases) == 0 {
			return []SubtitleSegment{{Text: transcript, Start: 0, End: durationSeconds}}, nil
		}

		// Calcular caracteres totales para distribución proporcional
		totalChars := 0
		for _, p := range phrases {
			totalChars += len(strings.TrimSpace(p))
		}
		if totalChars == 0 {
			return []SubtitleSegment{{Text: transcript, Start: 0, End: durationSeconds}}, nil
		}

		// Asignar tiempo proporcional a cada frase
		segments := make([]SubtitleSegment, 0, len(phrases))
		elapsed := 0.0
		for _, p := range phrases {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			ratio := float64(len(p)) / float64(totalChars)
			segDuration := durationSeconds * ratio
			segments = append(segments, SubtitleSegment{
				Text:  p,
				Start: elapsed,
				End:   elapsed + segDuration,
			})
			elapsed += segDuration
		}

		// Ajustar el último segmento para que end = durationSeconds (evitar redondeo)
		if len(segments) > 0 {
			segments[len(segments)-1].End = durationSeconds
		}

		obs.Logger.InfoContext(spanCtx, "subtitle_timing_generated", "segments", len(segments), "duration_sec", durationSeconds)
		return segments, nil
	}, nil
}

// splitIntoPhrases divide el texto en frases por puntuación (. ! ?) o saltos de línea.
func splitIntoPhrases(text string) []string {
	// Normalizar: reemplazar múltiples espacios/saltos por uno
	text = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(text), " ")
	// Dividir por . ! ? manteniendo el delimitador en la frase
	re := regexp.MustCompile(`[^.!?]+[.!?]*`)
	matches := re.FindAllString(text, -1)
	result := make([]string, 0, len(matches))
	for _, m := range matches {
		m = strings.TrimSpace(m)
		if m != "" {
			result = append(result, m)
		}
	}
	if len(result) == 0 && text != "" {
		return []string{text}
	}
	return result
}
