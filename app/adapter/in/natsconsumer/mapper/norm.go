package mapper

import (
	"strings"
)

// ------------ helpers ------------
func norm(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = deaccent(s)
	repl := strings.NewReplacer("-", " ", "_", " ")
	s = repl.Replace(s)
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}
