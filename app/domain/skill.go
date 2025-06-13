package domain

import "context"

type Skill string

func (s Skill) DocumentID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, string(s))
}
