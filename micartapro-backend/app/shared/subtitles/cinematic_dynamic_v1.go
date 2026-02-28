package subtitles

import (
	"strings"
)

func pseudoRandV1(seed int, text string) float64 {
	h := seed
	for _, c := range text {
		h = h*31 + int(c)
	}
	if h < 0 {
		h = -h
	}
	return float64(h%1000) / 1000.0
}

// Fases narrativas (porcentaje de duración total).
const (
	phaseHookPctEnd     = 0.12 // 0% – 12%
	phaseBodyPctEnd     = 0.70 // 12% – 70%
	phaseEmphasisPctEnd = 0.90 // 70% – 90%
	// 90% – 100% = CTA
)

const (
	maxCenterIntervalSec = 6.0
	shortPhraseChars     = 28
)

// getPhase devuelve la fase narrativa según el porcentaje de avance.
func getPhase(startPct float64) string {
	if startPct < phaseHookPctEnd {
		return "HOOK"
	}
	if startPct < phaseBodyPctEnd {
		return "BODY"
	}
	if startPct < phaseEmphasisPctEnd {
		return "EMPHASIS"
	}
	return "CTA"
}

// hasEmphasisPhrase indica si el segmento contiene alguna frase de emphasis.
func hasEmphasisPhrase(seg *SubtitleSegment, emphasisPhrases []string) bool {
	if seg.Emphasis != nil && (len(seg.Emphasis.Phrases) > 0 || seg.Emphasis.Level != "") {
		return true
	}
	for _, p := range emphasisPhrases {
		if strings.Contains(strings.ToLower(seg.Text), strings.ToLower(strings.TrimSpace(p))) {
			return true
		}
	}
	return false
}

// pickContrastCandidate selecciona el índice del bloque con mayor carga semántica (longitud) en BODY.
// Solo se considera 1 bloque como CONTRAST.
func pickContrastCandidate(segments []SubtitleSegment, durationSec float64) int {
	bodyStart := durationSec * phaseHookPctEnd
	bodyEnd := durationSec * phaseBodyPctEnd

	maxLen := 0
	candidate := -1
	for i := range segments {
		seg := &segments[i]
		mid := (seg.Start + seg.End) / 2
		if mid < bodyStart || mid > bodyEnd {
			continue
		}
		ln := len(seg.Text)
		if ln > maxLen {
			maxLen = ln
			candidate = i
		}
	}
	return candidate
}

// longSentenceChars umbral para considerar una frase "larga" cuando avoidCenterLongSentences.
const longSentenceChars = 60

// ApplyCinematicDynamicV1 aplica el preset CINEMATIC_DYNAMIC_V1: ritmo narrativo, énfasis claro, dinamismo.
// Heurísticas determinísticas sin ML.
// constraints: si avoidCenterLongSentences=true, evita CENTER para frases largas; preferCenterForEmphasis prioriza CENTER para énfasis.
func ApplyCinematicDynamicV1(
	segments []SubtitleSegment,
	emphasisPhrases []string,
	durationSec float64,
	creativity float64, // 0-1: para excepción BODY (frase corta → CENTER con prob baja)
	constraints *LayoutConstraints,
) {
	if durationSec <= 0 || len(segments) == 0 {
		return
	}

	contrastIdx := pickContrastCandidate(segments, durationSec)
	contrastUsed := false
	lastCenterAt := -999.0

	for i := range segments {
		seg := &segments[i]
		startPct := seg.Start / durationSec
		phase := getPhase(startPct)
		hasEmphasis := hasEmphasisPhrase(seg, emphasisPhrases)

		// Animación según placement (CENTER→SCALE_IN, BOTTOM/TOP→FADE_IN)
		placement := PlacementBottom
		animation := AnimationFadeIn
		size := ""

		avoidLong := constraints != nil && constraints.AvoidCenterLongSentences
		preferCenterEmphasis := constraints == nil || constraints.PreferCenterForEmphasis
		isLong := len(seg.Text) > longSentenceChars

		if hasEmphasis {
			// avoidCenterLongSentences: no poner frases largas en CENTER aunque tengan énfasis
			if avoidLong && isLong {
				placement = PlacementBottom
				animation = AnimationFadeIn
				size = SizeL
			} else if preferCenterEmphasis {
				placement = PlacementCenter
				animation = AnimationScaleIn
				size = SizeXL
			} else {
				placement = PlacementBottom
				animation = AnimationFadeIn
				size = SizeL
			}
			if seg.Emphasis != nil {
				seg.Emphasis.Level = SizeXL
			} else {
				// ApplyEmphasis ya corrió; si llegamos aquí, emphasisPhrases matcheó
				matched := make([]string, 0, len(emphasisPhrases))
				for _, p := range emphasisPhrases {
					if strings.Contains(strings.ToLower(seg.Text), strings.ToLower(strings.TrimSpace(p))) {
						matched = append(matched, strings.TrimSpace(p))
					}
				}
				seg.Emphasis = &EmphasisObject{Phrases: matched, Level: SizeXL}
			}
		} else if phase == "HOOK" {
			if avoidLong && isLong {
				placement = PlacementBottom
				animation = AnimationFadeIn
				size = SizeL
			} else {
				placement = PlacementCenter
				animation = AnimationScaleIn
				size = SizeL
			}
		} else if phase == "CTA" {
			if avoidLong && isLong {
				placement = PlacementBottom
				animation = AnimationFadeIn
				size = SizeXL
			} else {
				placement = PlacementCenter
				animation = AnimationScaleIn
				size = SizeXL
			}
		} else if phase == "BODY" {
			if contrastIdx == i && !contrastUsed {
				placement = PlacementTop
				animation = AnimationFadeIn
				size = SizeL
				contrastUsed = true
			} else {
				// BODY: BOTTOM por defecto
				// Excepción: frase < 28 chars → CENTER con probabilidad baja (creativity-driven)
				// máximo 1 CENTER cada 6 segundos
				elapsedSinceCenter := seg.Start - lastCenterAt
				allowCenter := elapsedSinceCenter >= maxCenterIntervalSec
				shortPhrase := len(seg.Text) < shortPhraseChars
				randVal := pseudoRandV1(i, seg.Text)
				useCenter := allowCenter && shortPhrase && creativity >= 0.5 && randVal < (0.15*creativity) // baja prob

				if useCenter {
					placement = PlacementCenter
					animation = AnimationScaleIn
					size = SizeL
					lastCenterAt = seg.Start
				} else {
					placement = PlacementBottom
					animation = AnimationFadeIn
					size = SizeM
				}
			}
		} else {
			// EMPHASIS sin emphasisPhrases: tratar como BODY
			placement = PlacementBottom
			animation = AnimationFadeIn
			size = SizeM
		}

		seg.Placement = placement
		seg.Animation = animation
		if size != "" {
			seg.Size = size
		}

		// Lines heredan; solo overrides de size para líneas con emphasis
		for j := range seg.Lines {
			if seg.Size != "" {
				seg.Lines[j].Size = seg.Size
			}
			if seg.Lines[j].Emphasis != nil && seg.Lines[j].Emphasis.Level == SizeBig {
				seg.Lines[j].Size = SizeBig
			}
		}
	}

	// CTA: último bloque o último con emphasis → CENTER, SCALE_IN, XL
	// Ya aplicado en el loop para phase==CTA. Asegurar que el último bloque sea CTA si está en 90-100%
	if len(segments) > 0 {
		last := &segments[len(segments)-1]
		lastPct := last.Start / durationSec
		if lastPct >= 0.9 {
			avoidLong := constraints != nil && constraints.AvoidCenterLongSentences
			if avoidLong && len(last.Text) > longSentenceChars {
				last.Placement = PlacementBottom
				last.Animation = AnimationFadeIn
			} else {
				last.Placement = PlacementCenter
				last.Animation = AnimationScaleIn
			}
			last.Size = SizeXL
			if hasEmphasisPhrase(last, emphasisPhrases) && last.Emphasis != nil {
				last.Emphasis.Level = SizeXL
			}
		}
	}
}
