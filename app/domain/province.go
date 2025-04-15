package domain

import (
	"context"
	"strings"
)

type Province string

func (s Province) DocID(ctx context.Context) DocumentID {
	return HashByCountry(ctx, "province", s.String())
}

func (s Province) String() string {
	return string(s)
}

func (p Province) IsEmpty() bool {
	return strings.TrimSpace(string(p)) == ""
}

func (s Province) Equals(other Province) bool {
	return s.String() == other.String()
}

func (s *Province) UpdateIfChanged(new Province) bool {
	if !s.Equals(new) && !new.IsEmpty() {
		*s = new
		return true
	}
	return false
}
