package subtitles

import "strings"

// AddWordsToLinesToSegments añade words[] a todas las líneas de los segmentos.
// Si segmentWords[i] tiene datos (timestamps reales de Chirp), los usa; si no, distribuye proporcionalmente.
func AddWordsToLinesToSegments(segments []SubtitleSegment, segmentWords [][]SubtitleWord) {
	for i := range segments {
		if len(segmentWords) > i && len(segmentWords[i]) > 0 {
			assignRealWordsToLines(&segments[i], segmentWords[i])
		} else {
			AddWordsToLines(segments[i].Lines)
		}
	}
}

// assignRealWordsToLines reparte los timestamps reales de Chirp entre las líneas del segmento.
func assignRealWordsToLines(seg *SubtitleSegment, words []SubtitleWord) {
	if len(seg.Lines) == 0 || len(words) == 0 {
		return
	}
	wordIdx := 0
	for j := range seg.Lines {
		line := &seg.Lines[j]
		lineWords := strings.Fields(line.Text)
		if len(lineWords) == 0 {
			continue
		}
		line.Words = make([]SubtitleWord, 0, len(lineWords))
		for k := 0; k < len(lineWords) && wordIdx < len(words); k++ {
			line.Words = append(line.Words, SubtitleWord{
				Text:  words[wordIdx].Text,
				Start: words[wordIdx].Start,
				End:   words[wordIdx].End,
			})
			wordIdx++
		}
	}
}

// AddWordsToLines divide cada línea en palabras con timing proporcional.
// Para estilo ALEX_HORMOZI: palabra activa en amarillo, scale, stroke negro.
func AddWordsToLines(lines []SubtitleLine) {
	for i := range lines {
		line := &lines[i]
		words := strings.Fields(line.Text)
		if len(words) <= 1 {
			continue
		}
		duration := line.End - line.Start
		totalLen := 0
		for _, w := range words {
			totalLen += len(w) + 1
		}
		if totalLen <= 0 {
			continue
		}
		line.Words = make([]SubtitleWord, len(words))
		currStart := line.Start
		for j, w := range words {
			ratio := float64(len(w)+1) / float64(totalLen)
			wordDur := duration * ratio
			currEnd := currStart + wordDur
			if j == len(words)-1 {
				currEnd = line.End
			}
			line.Words[j] = SubtitleWord{Text: w, Start: currStart, End: currEnd}
			currStart = currEnd
		}
	}
}

// ApplyAlexHormoziStyle aplica el estilo Alex Hormozi: MAYÚSCULAS, BOTTOM, FADE_IN.
// AddWordsToLines debe llamarse antes para tener words[].
func ApplyAlexHormoziStyle(segments []SubtitleSegment) {
	for i := range segments {
		seg := &segments[i]
		seg.Text = strings.ToUpper(seg.Text)
		seg.Placement = PlacementBottom
		seg.Animation = AnimationFadeIn
		for j := range seg.Lines {
			seg.Lines[j].Text = strings.ToUpper(seg.Lines[j].Text)
			seg.Lines[j].Placement = PlacementBottom
			seg.Lines[j].Animation = AnimationFadeIn
			for k := range seg.Lines[j].Words {
				seg.Lines[j].Words[k].Text = strings.ToUpper(seg.Lines[j].Words[k].Text)
			}
		}
	}
}
