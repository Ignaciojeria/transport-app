package domain

import (
	"context"
	"strings"
)

type State string

func (s State) DocID(ctx context.Context) DocumentID {
	return HashByCountry(ctx, "state", s.String())
}

func (s State) String() string {
	return string(s)
}

func (s State) IsEmpty() bool {
	return strings.TrimSpace(string(s)) == ""
}

func (s State) Equals(other State) bool {
	return s.String() == other.String()
}

func (s *State) UpdateIfChanged(new State) bool {
	if !s.Equals(new) && !new.IsEmpty() {
		*s = new
		return true
	}
	return false
}
