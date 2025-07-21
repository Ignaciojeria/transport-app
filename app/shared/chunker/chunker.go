package chunker

import (
	"bytes"
	"errors"

	"github.com/google/uuid"
)

const MaxChunkSize = 734003 // 0.7 MB in bytes

type Chunk struct {
	ID   uuid.UUID
	Data []byte
	Idx  int
}

// Divide los datos en chunks de hasta 0.7 MB
func SplitBytes(data []byte) ([]Chunk, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	var chunks []Chunk
	for i := 0; i < len(data); i += MaxChunkSize {
		end := i + MaxChunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, Chunk{
			ID:   uuid.New(),
			Data: data[i:end],
			Idx:  len(chunks),
		})
	}
	return chunks, nil
}

// Reconstruye el byte array original
func ReconstructBytes(chunks []Chunk) ([]byte, error) {
	var buffer bytes.Buffer
	for _, chunk := range chunks {
		buffer.Write(chunk.Data)
	}
	return buffer.Bytes(), nil
}
