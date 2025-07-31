package domain

import (
	"time"

	"github.com/cockroachdb/errors"
)

// FSMState representa el estado actual de una máquina de estados finitos
type FSMState struct {
	TraceID        string
	IdempotencyKey string
	Workflow       string
	State          string
	NextInput      []byte
	CreatedAt      time.Time
}

// IsEmpty verifica si el estado está vacío (no hay transiciones previas)
func (s FSMState) IsEmpty() bool {
	return s.TraceID == "" && s.Workflow == "" && s.State == ""
}

// ErrNoTransitionsFound es retornado cuando no hay transiciones previas
var ErrNoTransitionsFound = errors.New("no transitions found for this trace")
