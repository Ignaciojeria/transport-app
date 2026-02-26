package subtitles

import (
	"regexp"
	"strings"
)

const (
	trailerMaxWordsPerFragment = 5
	trailerMinWordsPerFragment = 1
	TrailerMicroGapMin         = 0.2
	TrailerMicroGapMax         = 0.4
	trailerCenterPct           = 0.70
	trailerBottomPct           = 0.20
	trailerTopPct              = 0.10
)

// trailerRand determina placement de forma pseudo-aleatoria por índice.
func trailerRand(i int, text string) float64 {
	h := i * 31
	for _, c := range text {
		h = h*31 + int(c)
	}
	if h < 0 {
		h = -h
	}
	return float64(h%1000) / 1000.0
}

// fragmentByPunctuation divide por comas, puntos y coma, signos de interrogación.
var fragmentByPunctuation = regexp.MustCompile(`[,;]+\s*|[\?\.]+\s*`)

// FragmentSegmentsForTrailer fragmenta segmentos de forma agresiva para estilo TRAILER.
// Cada fragmento se convierte en un segmento con 1 línea. maxLines efectivo = 1.
func FragmentSegmentsForTrailer(segments []SubtitleSegment) []SubtitleSegment {
	var out []SubtitleSegment
	for _, seg := range segments {
		text := strings.TrimSpace(seg.Text)
		if text == "" {
			continue
		}
		duration := seg.End - seg.Start
		if duration <= 0 {
			duration = 1.0
		}

		// Primero dividir por puntuación
		parts := fragmentByPunctuation.Split(text, -1)
		var fragments []string
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			words := strings.Fields(p)
			if len(words) <= trailerMaxWordsPerFragment {
				fragments = append(fragments, p)
				continue
			}
			// Fragmentar por grupos de palabras
			for i := 0; i < len(words); i += trailerMaxWordsPerFragment {
				end := i + trailerMaxWordsPerFragment
				if end > len(words) {
					end = len(words)
				}
				chunk := strings.Join(words[i:end], " ")
				if len(chunk) >= trailerMinWordsPerFragment {
					fragments = append(fragments, chunk)
				}
			}
		}

		if len(fragments) == 0 {
			fragments = []string{text}
		}

		// Distribuir tiempo proporcionalmente
		totalLen := 0
		for _, f := range fragments {
			totalLen += len(f) + 1
		}
		if totalLen <= 0 {
			totalLen = 1
		}

		currStart := seg.Start
		for i, f := range fragments {
			ratio := float64(len(f)+1) / float64(totalLen)
			segDur := duration * ratio
			currEnd := currStart + segDur
			if i == len(fragments)-1 {
				currEnd = seg.End
			}

			out = append(out, SubtitleSegment{
				Text:  f,
				Start: currStart,
				End:   currEnd,
				Lines: []SubtitleLine{{Text: f, Start: currStart, End: currEnd}},
			})
			currStart = currEnd
		}
	}
	return out
}

// AddMicroGaps crea micro gaps (0.2–0.4s) entre segmentos acortando el END del anterior.
// No desplaza el START: los subtítulos siguen sincronizados con el audio.
func AddMicroGaps(segments []SubtitleSegment, gapMin, gapMax float64) {
	if len(segments) < 2 || gapMin <= 0 {
		return
	}
	if gapMax < gapMin {
		gapMax = gapMin
	}

	for i := 1; i < len(segments); i++ {
		prev := &segments[i-1]
		curr := &segments[i]
		actualGap := curr.Start - prev.End
		targetGap := gapMin
		if gapMax > gapMin {
			r := float64((i*17)%100) / 100.0
			targetGap = gapMin + (gapMax-gapMin)*r
		}
		if actualGap < targetGap {
			// Acortar prev.End para crear gap visual; NO desplazar curr (mantiene sync con audio)
			trim := targetGap - actualGap
			minDuration := 0.25
			maxTrim := prev.End - prev.Start - minDuration
			if maxTrim > 0 && trim > 0 {
				actualTrim := trim
				if actualTrim > maxTrim {
					actualTrim = maxTrim
				}
				newEnd := prev.End - actualTrim
				if newEnd > prev.Start {
					prev.End = newEnd
					for j := range prev.Lines {
						prev.Lines[j].End = newEnd
					}
				}
			}
		}
	}
}

// ApplyTrailerV1 aplica el preset TRAILER_V1: fragmentación, CENTER dominante, énfasis XL+.
// Reglas: maxLines=1, 70% CENTER, 20% BOTTOM, 10% TOP, emphasis solo en CENTER con XL o XXL.
func ApplyTrailerV1(segments []SubtitleSegment, emphasisPhrases []string, durationSec float64) {
	if len(segments) == 0 {
		return
	}

	for i := range segments {
		seg := &segments[i]
		hasEmphasis := hasEmphasisPhrase(seg, emphasisPhrases)
		r := trailerRand(i, seg.Text)

		// Placement: 70% CENTER, 20% BOTTOM, 10% TOP
		var placement string
		if r < trailerCenterPct {
			placement = PlacementCenter
		} else if r < trailerCenterPct+trailerBottomPct {
			placement = PlacementBottom
		} else {
			placement = PlacementTop
		}

		// Emphasis: siempre CENTER, SCALE_IN, XL o XXL
		if hasEmphasis {
			placement = PlacementCenter
			seg.Animation = AnimationScaleIn
			seg.Size = SizeXXL
			if seg.Emphasis != nil {
				seg.Emphasis.Level = SizeXXL
			} else {
				matched := make([]string, 0, len(emphasisPhrases))
				for _, p := range emphasisPhrases {
					if strings.Contains(strings.ToLower(seg.Text), strings.ToLower(strings.TrimSpace(p))) {
						matched = append(matched, strings.TrimSpace(p))
					}
				}
				seg.Emphasis = &EmphasisObject{Phrases: matched, Level: SizeXXL}
			}
		} else {
			seg.Animation = AnimationFadeIn
			if placement == PlacementCenter && (r < 0.3 || len(seg.Text) < 20) {
				seg.Animation = AnimationScaleIn
				seg.Size = SizeXL
			} else if seg.Size == "" {
				seg.Size = SizeL
			}
		}

		seg.Placement = placement
		for j := range seg.Lines {
			seg.Lines[j].Placement = placement
			seg.Lines[j].Animation = seg.Animation
			if seg.Size != "" {
				seg.Lines[j].Size = seg.Size
			}
			if seg.Lines[j].Emphasis != nil && seg.Emphasis != nil {
				seg.Lines[j].Emphasis.Level = seg.Emphasis.Level
			}
		}
	}
}
