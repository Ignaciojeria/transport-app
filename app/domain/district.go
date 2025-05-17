package domain

import (
	"context"
	"strings"
)

type District string

func (s District) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, s.String())
}

func (s District) String() string {
	return string(s)
}

func (d District) IsEmpty() bool {
	return strings.TrimSpace(string(d)) == ""
}

func (s District) Equals(other District) bool {
	return s.String() == other.String()
}

func (s District) UpdateIfChanged(new District) (District, bool) {
	if !s.Equals(new) && !new.IsEmpty() {
		return new, true
	}
	return s, false
}
