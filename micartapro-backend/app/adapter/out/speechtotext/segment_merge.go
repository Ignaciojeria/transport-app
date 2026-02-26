package speechtotext

import (
	"strings"
	"unicode"
)

// MergeFragmentedSegments fusiona segmentos que fueron cortados a mitad de frase,
// mejorando la coherencia para subtítulos. Ej: "contenido," + "lo crean" → un solo segmento.
func MergeFragmentedSegments(segments []SubtitleSegment) []SubtitleSegment {
	if len(segments) <= 1 {
		computeSilenceAfter(segments)
		return segments
	}
	merged := make([]SubtitleSegment, 0, len(segments))
	merged = append(merged, segments[0])
	for i := 1; i < len(segments); i++ {
		prev := &merged[len(merged)-1]
		curr := segments[i]
		if shouldMerge(prev, &curr) {
			prev.Text = strings.TrimSpace(prev.Text) + " " + strings.TrimSpace(curr.Text)
			prev.End = curr.End
			prev.Words = append(prev.Words, curr.Words...)
		} else {
			merged = append(merged, curr)
		}
	}
	computeSilenceAfter(merged)
	return merged
}

// computeSilenceAfter rellena SilenceAfterSeconds para cada segmento (next.start - this.end).
func computeSilenceAfter(segments []SubtitleSegment) {
	for i := 0; i < len(segments); i++ {
		if i+1 < len(segments) {
			gap := segments[i+1].Start - segments[i].End
			if gap > 0 {
				segments[i].SilenceAfterSeconds = gap
			}
		} else {
			segments[i].SilenceAfterSeconds = 0
		}
		// Recalcular DurationSeconds siempre (nunca confiar en el modelo)
		segments[i].DurationSeconds = segments[i].End - segments[i].Start
	}
}

func shouldMerge(prev, next *SubtitleSegment) bool {
	prevTrim := strings.TrimSpace(prev.Text)
	nextTrim := strings.TrimSpace(next.Text)
	if prevTrim == "" || nextTrim == "" {
		return false
	}
	// Gap grande (>1.2s): probablemente pausa intencional, no fusionar
	gap := next.Start - prev.End
	if gap > 1.2 {
		return false
	}
	// Prev termina en puntuación de fin de frase → no fusionar
	if endsSentence(prevTrim) && !startsWithLowercase(nextTrim) {
		return false
	}
	// Prev termina en coma/sin punto y next empieza con minúscula → misma frase
	if (endsWithCommaOrConjunction(prevTrim) || !endsSentence(prevTrim)) && startsWithLowercase(nextTrim) {
		return true
	}
	// Prev termina en coma, gap pequeño, next corto → continuación (ej: listas)
	if endsWithComma(prevTrim) && gap < 0.9 && len(nextTrim) < 55 {
		return true
	}
	return false
}

func endsSentence(s string) bool {
	s = strings.TrimRight(s, " ")
	if s == "" {
		return false
	}
	last := s[len(s)-1]
	return last == '.' || last == '!' || last == '?'
}

func endsWithComma(s string) bool {
	s = strings.TrimRight(s, " ")
	if s == "" {
		return false
	}
	return s[len(s)-1] == ','
}

func endsWithCommaOrConjunction(s string) bool {
	s = strings.TrimRight(s, " ")
	if s == "" {
		return false
	}
	last := s[len(s)-1]
	return last == ',' || last == ';' || last == ':'
}

func startsWithLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return unicode.IsLower(r)
		}
	}
	return false
}
