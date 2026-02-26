package subtitles

// DirectionInput es el contrato de entrada para el SubtitleDirectionEngine.
// Transforma []SubtitleSegment (transcription) en segmentos dirigidos (placement, animation, size).
type DirectionInput struct {
	Style             string   // ONE_LINE | LINE_BY_LINE | CINEMATIC_DYNAMIC | TRAILER | ALEX_HORMOZI
	DirectionPreset   string   // CINEMATIC_DYNAMIC_V1 | TRAILER_V1 | "" (genérico)
	PlacementStrategy string   // FIXED | DYNAMIC
	DefaultPlacement  string   // BOTTOM | CENTER | TOP (cuando FIXED)
	EmphasisPhrases   []string // Frases a destacar
	DurationSec       float64  // Duración total del audio
	Creativity        float64  // 0-1: 0=seguro, 1=cinematográfico
	Constraints       *LayoutConstraints
	DynamicRules      *DynamicRules
	// SegmentWords: para ALEX_HORMOZI, timestamps por palabra por segmento
	SegmentWords [][]SubtitleWord `json:"-"`
}

// BuildSubtitleDirection transforma segmentos de transcripción en segmentos dirigidos.
// Input: []SubtitleSegment (con Text, Start, End, Lines ya poblados)
// Output: mismos segmentos con Placement, Animation, Size, Emphasis aplicados
//
// Pipeline determinístico:
// 1. Preproceso por estilo (TRAILER fragmenta, ALEX_HORMOZI añade words)
// 2. ApplyEmphasis
// 3. Preset de dirección (CINEMATIC_DYNAMIC_V1, TRAILER_V1) o PlacementStrategy genérico
// 4. Fallback FIXED (todo defaultPlacement)
func BuildSubtitleDirection(segments []SubtitleSegment, input DirectionInput) []SubtitleSegment {
	if len(segments) == 0 {
		return segments
	}
	subs := segments

	switch input.Style {
	case StyleTrailer:
		subs = FragmentSegmentsForTrailer(subs)
		AddMicroGaps(subs, TrailerMicroGapMin, TrailerMicroGapMax)
		ApplyEmphasis(subs, input.EmphasisPhrases)
		ApplyTrailerV1(subs, input.EmphasisPhrases, input.DurationSec)
	case StyleAlexHormozi:
		ApplyEmphasis(subs, input.EmphasisPhrases)
		if len(input.SegmentWords) > 0 {
			AddWordsToLinesToSegments(subs, input.SegmentWords)
		}
		ApplyAlexHormoziStyle(subs)
		return subs
	default:
		ApplyEmphasis(subs, input.EmphasisPhrases)
	}

	if input.PlacementStrategy != PlacementStrategyDynamic {
		applyFixedPlacement(subs, input.DefaultPlacement)
		return subs
	}

	creativity := input.Creativity
	if creativity < 0 {
		creativity = 0.8
	}
	if creativity > 1 {
		creativity = 1
	}

	switch input.DirectionPreset {
	case DirectionPresetDocumentary:
		applyFixedPlacement(subs, PlacementBottom)
		return subs
	case DirectionPresetCinematicDynamicV1:
		constraints := input.Constraints
		if constraints == nil {
			constraints = &LayoutConstraints{}
		}
		ApplyCinematicDynamicV1(subs, input.EmphasisPhrases, input.DurationSec, creativity, constraints)
	case DirectionPresetTrailerV1:
		// TRAILER ya aplicado arriba
	default:
		constraints := input.Constraints
		if constraints == nil {
			constraints = &LayoutConstraints{}
		}
		dynRules := input.DynamicRules
		if dynRules == nil {
			dynRules = &DynamicRules{}
		}
		if dynRules.MinSecondsBetweenPlacementChanges <= 0 && input.Style == StyleCinematicDynamic {
			dynRules.MinSecondsBetweenPlacementChanges = 2.0
		}
		if dynRules.MaxCenterBlocksInRow <= 0 {
			dynRules.MaxCenterBlocksInRow = 2
		}
		ApplyPlacementStrategy(subs, creativity, input.EmphasisPhrases, constraints, dynRules)
	}

	return subs
}

// applyFixedPlacement asigna el mismo placement y animation a todos los segmentos.
func applyFixedPlacement(segments []SubtitleSegment, placement string) {
	if placement == "" {
		placement = PlacementBottom
	}
	for i := range segments {
		segments[i].Placement = placement
		segments[i].Animation = AnimationFadeIn
		for j := range segments[i].Lines {
			segments[i].Lines[j].Placement = placement
			segments[i].Lines[j].Animation = AnimationFadeIn
		}
	}
}
