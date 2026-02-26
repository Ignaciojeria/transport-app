package subtitles

import (
	"strings"
)

// pseudoRand genera valor 0-1 determinístico desde índice y texto (mismo input = mismo output).
func pseudoRand(seed int, text string) float64 {
	h := seed
	for _, c := range text {
		h = h*31 + int(c)
	}
	if h < 0 {
		h = -h
	}
	return float64(h%1000) / 1000.0
}

// ApplyPlacementStrategy asigna placement, animation y size por segmento según creativity.
// creativity: 0.0 = todo BOTTOM (seguro), 0.5 = algo dinámico, 0.8 = creativo, 1.0 = muy cinematográfico.
// dynamicRules: dampening para no cambiar placement demasiado seguido.
func ApplyPlacementStrategy(
	segments []SubtitleSegment,
	creativity float64,
	emphasisPhrases []string,
	constraints *LayoutConstraints,
	dynamicRules *DynamicRules,
) {
	if creativity <= 0 {
		for i := range segments {
			segments[i].Placement = PlacementBottom
			segments[i].Animation = AnimationNone
		}
		return
	}

	preferCenter := constraints != nil && constraints.PreferCenterForEmphasis
	avoidLong := constraints != nil && constraints.AvoidCenterLongSentences
	minSec := 0.0
	maxCenterInRow := 0
	if dynamicRules != nil {
		minSec = dynamicRules.MinSecondsBetweenPlacementChanges
		maxCenterInRow = dynamicRules.MaxCenterBlocksInRow
	}
	if maxCenterInRow <= 0 {
		maxCenterInRow = 3
	}

	lastPlacement := PlacementBottom
	lastPlacementAt := -999.0
	centerInRow := 0

	for i := range segments {
		seg := &segments[i]
		hasEmphasis := seg.Emphasis != nil && (len(seg.Emphasis.Phrases) > 0 || seg.Emphasis.Level != "")
		if !hasEmphasis && len(emphasisPhrases) > 0 {
			for _, p := range emphasisPhrases {
				if strings.Contains(strings.ToLower(seg.Text), strings.ToLower(strings.TrimSpace(p))) {
					hasEmphasis = true
					break
				}
			}
		}

		isQuestion := strings.Contains(seg.Text, "?") || strings.Contains(seg.Text, "¿")
		duration := seg.End - seg.Start
		isShort := duration < 3.0 || len(seg.Text) < 50
		isLong := len(seg.Text) > 80

		wanted := PlacementBottom
		wantedAnim := AnimationFadeIn
		wantedSize := ""

		if hasEmphasis && preferCenter {
			wanted = PlacementCenter
			wantedAnim = AnimationScaleIn
			wantedSize = SizeXL
		} else if hasEmphasis {
			wanted = PlacementCenter
			wantedAnim = AnimationScaleIn
			wantedSize = SizeBig
		} else if isQuestion && creativity >= 0.5 {
			wanted = PlacementCenter
			wantedAnim = AnimationScaleIn
		} else if avoidLong && isLong {
			wanted = PlacementBottom
			wantedAnim = AnimationFadeIn
		} else if creativity >= 0.8 {
			choice := pseudoRand(i, seg.Text)
			if isShort && choice < 0.4 {
				wanted = PlacementCenter
				wantedAnim = AnimationFadeIn
			} else if choice < 0.2 {
				wanted = PlacementTop
				wantedAnim = AnimationFadeIn
			} else {
				wanted = PlacementBottom
				wantedAnim = AnimationFadeIn
			}
		} else if creativity >= 0.5 {
			if isShort && pseudoRand(i, seg.Text) < 0.3 {
				wanted = PlacementCenter
				wantedAnim = AnimationFadeIn
			} else {
				wanted = PlacementBottom
				wantedAnim = AnimationFadeIn
			}
		} else {
			wanted = PlacementBottom
			wantedAnim = AnimationFadeIn
		}

		// Dampening: respetar minSecondsBetweenPlacementChanges y maxCenterBlocksInRow
		elapsed := seg.Start - lastPlacementAt
		if minSec > 0 && elapsed < minSec && wanted != lastPlacement {
			wanted = lastPlacement
			wantedAnim = AnimationFadeIn
		}
		if wanted == PlacementCenter {
			if centerInRow >= maxCenterInRow {
				wanted = PlacementBottom
				wantedAnim = AnimationFadeIn
			} else {
				centerInRow++
			}
		} else {
			centerInRow = 0
		}
		if wanted != lastPlacement {
			lastPlacementAt = seg.Start
			lastPlacement = wanted
		}

		seg.Placement = wanted
		seg.Animation = wantedAnim
		if wantedSize != "" {
			seg.Size = wantedSize
		}

		// Lines heredan placement/animation del segmento; solo overrides (Size cuando línea tiene emphasis)
		for j := range seg.Lines {
			if seg.Size != "" {
				seg.Lines[j].Size = seg.Size
			}
			if seg.Lines[j].Emphasis != nil && seg.Lines[j].Emphasis.Level == SizeBig {
				seg.Lines[j].Size = SizeBig
			}
		}
	}
}
