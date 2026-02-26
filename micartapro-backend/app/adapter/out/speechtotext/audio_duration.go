package speechtotext

import (
	"bytes"
	"encoding/binary"
)

// getWavDurationSeconds devuelve la duración real del audio WAV en segundos.
// Parsea el header RIFF/WAV. Retorna 0 si no es WAV válido o hay error.
func getWavDurationSeconds(data []byte) float64 {
	if len(data) < 44 {
		return 0
	}
	if string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		return 0
	}
	var byteRate uint32
	var dataSize uint32
	pos := 12
	for pos+8 <= len(data) {
		chunkID := string(data[pos : pos+4])
		chunkSize := binary.LittleEndian.Uint32(data[pos+4 : pos+8])
		pos += 8
		if pos+int(chunkSize) > len(data) {
			break
		}
		switch chunkID {
		case "fmt ":
			if chunkSize >= 16 && pos+12 <= len(data) {
				byteRate = binary.LittleEndian.Uint32(data[pos+8 : pos+12])
			}
		case "data":
			dataSize = chunkSize
		}
		pos += int(chunkSize)
		if pos&1 != 0 {
			pos++
		}
	}
	if byteRate > 0 && dataSize > 0 {
		return float64(dataSize) / float64(byteRate)
	}
	return 0
}

// getAudioDurationSeconds devuelve la duración real del audio si se puede detectar.
// Soporta WAV. Retorna 0 para formatos no soportados.
func getAudioDurationSeconds(data []byte, mimeType string) float64 {
	switch {
	case mimeType == "audio/wav" || bytes.HasPrefix(data, []byte("RIFF")):
		return getWavDurationSeconds(data)
	}
	return 0
}
