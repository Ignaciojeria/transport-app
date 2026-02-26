package subtitles

import (
	"regexp"
	"strings"
)

const (
	defaultMaxCharsPerLine = 42
	minCharsPerLine        = 10  // Evitar líneas muy cortas (ej: "lo crean.")
	minWordsPerLine        = 2   // O al menos 2 palabras
)

// SplitSegmentIntoLinesWithOverflow divide en líneas respetando maxLines según overflowStrategy.
// overflowStrategy: ALLOW (no hacer nada), REBALANCE (aumentar maxChars para fewer lines),
// SHRINK (fusionar líneas), REBALANCE_THEN_SHRINK (rebalance primero, luego shrink si no alcanza).
func SplitSegmentIntoLinesWithOverflow(
	text string, start, end float64, maxCharsPerLine, maxLines int, overflowStrategy string,
) []SubtitleLine {
	lines := SplitSegmentIntoLines(text, start, end, maxCharsPerLine)
	if maxLines <= 0 || overflowStrategy == "" || overflowStrategy == OverflowAllow || len(lines) <= maxLines {
		return lines
	}

	switch overflowStrategy {
	case OverflowRebalance, OverflowRebalanceShrink:
		// Rebalance: permitir más chars por línea para reducir número de líneas
		for tryChars := maxCharsPerLine + 8; tryChars <= 80 && len(lines) > maxLines; tryChars += 8 {
			lines = SplitSegmentIntoLines(text, start, end, tryChars)
		}
		if overflowStrategy == OverflowRebalanceShrink && len(lines) > maxLines {
			lines = shrinkToMaxLines(lines, maxLines)
		}
	case OverflowShrink:
		lines = shrinkToMaxLines(lines, maxLines)
	}
	return lines
}

// SplitSegmentIntoLinesWithWordTiming divide en líneas usando timestamps reales por palabra.
// Si words está vacío, delega a SplitSegmentIntoLinesWithOverflow (distribución proporcional).
// Prefiere cortar en pausas (gap > minPauseSec entre palabras) para mejor sincronía con audio.
func SplitSegmentIntoLinesWithWordTiming(
	text string, start, end float64, words []SubtitleWord,
	maxCharsPerLine, maxLines int, overflowStrategy string, minPauseSec float64,
) []SubtitleLine {
	if len(words) == 0 {
		return SplitSegmentIntoLinesWithOverflow(text, start, end, maxCharsPerLine, maxLines, overflowStrategy)
	}
	if minPauseSec <= 0 {
		minPauseSec = 0.35
	}

	lines := splitTextIntoLinesWithWordTiming(words, maxCharsPerLine, minPauseSec)
	if len(lines) == 0 {
		return []SubtitleLine{{Text: text, Start: start, End: end}}
	}
	if len(lines) > maxLines && overflowStrategy != "" && overflowStrategy != OverflowAllow {
		lines = applyOverflowToWordTimedLines(lines, words, maxLines, overflowStrategy, maxCharsPerLine)
	}
	return lines
}

// splitTextIntoLinesWithWordTiming agrupa palabras en líneas respetando maxChars y prefiriendo cortes en pausas.
func splitTextIntoLinesWithWordTiming(words []SubtitleWord, maxChars int, minPauseSec float64) []SubtitleLine {
	if len(words) == 0 {
		return nil
	}
	if maxChars <= 0 {
		maxChars = defaultMaxCharsPerLine
	}

	var lines []SubtitleLine
	var lineWords []SubtitleWord
	var lineStart float64
	lineStart = words[0].Start
	lastEnd := words[0].End

	for _, w := range words {
		wordLen := len(w.Text)
		if len(lineWords) > 0 {
			wordLen++ // espacio antes
		}
		pause := w.Start - lastEnd

		// Cortar si: (1) línea excede maxChars, o (2) hay pausa y línea tiene longitud mínima
		shouldBreak := false
		if len(lineWords) > 0 {
			lineLen := 0
			for _, lw := range lineWords {
				lineLen += len(lw.Text) + 1
			}
			lineLen--
			if lineLen+wordLen > maxChars {
				shouldBreak = true
			} else if pause >= minPauseSec && len(lineWords) >= minWordsPerLine && lineLen >= minCharsPerLine {
				// Preferir cortar en pausas naturales (respetar silencios del audio)
				shouldBreak = true
			}
		}

		if shouldBreak && len(lineWords) > 0 {
			lineEnd := lineWords[len(lineWords)-1].End
			lineText := wordsToText(lineWords)
			lines = append(lines, SubtitleLine{Text: lineText, Start: lineStart, End: lineEnd})
			lineWords = nil
			lineStart = w.Start
		}

		lineWords = append(lineWords, w)
		lastEnd = w.End
	}

	if len(lineWords) > 0 {
		lineEnd := lineWords[len(lineWords)-1].End
		lines = append(lines, SubtitleLine{Text: wordsToText(lineWords), Start: lineStart, End: lineEnd})
	}
	return lines
}

func wordsToText(words []SubtitleWord) string {
	var b strings.Builder
	for i, w := range words {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(w.Text)
	}
	return b.String()
}

func applyOverflowToWordTimedLines(lines []SubtitleLine, words []SubtitleWord, maxLines int, overflowStrategy string, maxCharsPerLine int) []SubtitleLine {
	if len(lines) <= maxLines {
		return lines
	}
	switch overflowStrategy {
	case OverflowRebalance, OverflowRebalanceShrink:
		for tryChars := maxCharsPerLine + 8; tryChars <= 80 && len(lines) > maxLines; tryChars += 8 {
			lines = splitTextIntoLinesWithWordTiming(words, tryChars, 0.35)
			if len(lines) <= maxLines {
				return lines
			}
		}
		if overflowStrategy == OverflowRebalanceShrink {
			return shrinkToMaxLines(lines, maxLines)
		}
	case OverflowShrink:
		return shrinkToMaxLines(lines, maxLines)
	}
	return lines
}

// shrinkToMaxLines fusiona líneas excedentes en la última permitida.
func shrinkToMaxLines(lines []SubtitleLine, maxLines int) []SubtitleLine {
	if len(lines) <= maxLines {
		return lines
	}
	result := make([]SubtitleLine, maxLines)
	copy(result, lines[:maxLines-1])
	last := &result[maxLines-1]
	last.Text = lines[maxLines-1].Text
	last.Start = lines[maxLines-1].Start
	last.End = lines[len(lines)-1].End
	for i := maxLines; i < len(lines); i++ {
		last.Text = last.Text + " " + lines[i].Text
		last.End = lines[i].End
	}
	return result
}

// SplitSegmentIntoLines divide el texto en líneas y distribuye el tiempo proporcionalmente.
// Aplica rebalance: evita líneas < minCharsPerLine o < minWordsPerLine; fusiona con la anterior.
func SplitSegmentIntoLines(text string, start, end float64, maxCharsPerLine int) []SubtitleLine {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	if maxCharsPerLine <= 0 {
		maxCharsPerLine = defaultMaxCharsPerLine
	}

	lines := splitTextIntoLines(text, maxCharsPerLine)
	if len(lines) == 0 {
		return []SubtitleLine{{Text: text, Start: start, End: end}}
	}
	if len(lines) == 1 {
		return []SubtitleLine{{Text: lines[0], Start: start, End: end}}
	}

	// Distribuir tiempo proporcionalmente
	totalLen := 0
	for _, l := range lines {
		totalLen += len(l)
	}
	duration := end - start

	var raw []SubtitleLine
	currStart := start
	for i, line := range lines {
		ratio := float64(len(line)) / float64(totalLen)
		lineDuration := duration * ratio
		currEnd := currStart + lineDuration
		if i == len(lines)-1 {
			currEnd = end
		}
		raw = append(raw, SubtitleLine{Text: line, Start: currStart, End: currEnd})
		currStart = currEnd
	}

	return rebalanceShortLines(raw)
}

// rebalanceShortLines fusiona líneas muy cortas con la anterior.
func rebalanceShortLines(lines []SubtitleLine) []SubtitleLine {
	if len(lines) <= 1 {
		return lines
	}
	var result []SubtitleLine
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		words := strings.Fields(line.Text)
		isShort := len(line.Text) < minCharsPerLine || len(words) < minWordsPerLine

		if isShort && len(result) > 0 {
			// Fusionar con la línea anterior
			prev := &result[len(result)-1]
			prev.Text = prev.Text + " " + line.Text
			prev.End = line.End
		} else {
			result = append(result, line)
		}
	}
	return result
}

// splitTextIntoLines divide texto en líneas respetando palabras y puntuación.
func splitTextIntoLines(text string, maxChars int) []string {
	words := regexp.MustCompile(`\s+`).Split(text, -1)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	var current string
	for _, w := range words {
		if w == "" {
			continue
		}
		sep := " "
		if current == "" {
			sep = ""
		}
		candidate := current + sep + w
		if len(candidate) <= maxChars {
			current = candidate
			continue
		}
		if current != "" {
			lines = append(lines, strings.TrimSpace(current))
		}
		// Palabra muy larga: partirla
		if len(w) > maxChars {
			for len(w) > maxChars {
				lines = append(lines, w[:maxChars])
				w = w[maxChars:]
			}
			current = w
		} else {
			current = w
		}
	}
	if current != "" {
		lines = append(lines, strings.TrimSpace(current))
	}
	return lines
}
