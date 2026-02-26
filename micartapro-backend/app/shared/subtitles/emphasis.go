package subtitles

import "strings"

// ApplyEmphasis marca segmentos y líneas que contienen las frases de emphasis.
// emphasisPhrases: ej. ["Human Content AI"]. Output: EmphasisObject { phrases, level }.
func ApplyEmphasis(segments []SubtitleSegment, emphasisPhrases []string) {
	if len(emphasisPhrases) == 0 {
		return
	}
	for i := range segments {
		var segPhrases []string
		for _, phrase := range emphasisPhrases {
			phrase = strings.TrimSpace(phrase)
			if phrase == "" {
				continue
			}
			if strings.Contains(strings.ToLower(segments[i].Text), strings.ToLower(phrase)) {
				segPhrases = append(segPhrases, phrase)
			}
			// Marcar líneas que contienen la frase (solo level, hereda phrases del segmento)
			for j := range segments[i].Lines {
				if strings.Contains(strings.ToLower(segments[i].Lines[j].Text), strings.ToLower(phrase)) {
					segments[i].Lines[j].Emphasis = &EmphasisObject{Level: SizeBig}
				}
			}
		}
		if len(segPhrases) > 0 {
			segments[i].Emphasis = &EmphasisObject{Phrases: segPhrases, Level: SizeXL}
		}
	}
}
