package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapSkills(ctx context.Context, skills []domain.Skill) []table.Skill {
	skillsRecords := make([]table.Skill, len(skills))
	for i, skill := range skills {
		skillsRecords[i] = table.Skill{
			Name:       string(skill),
			TenantID:   sharedcontext.TenantIDFromContext(ctx),
			DocumentID: string(skill.DocumentID(ctx)),
		}
	}
	return skillsRecords
}
