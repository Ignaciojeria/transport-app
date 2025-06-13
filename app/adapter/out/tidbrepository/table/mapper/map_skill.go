package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapSkills(ctx context.Context, skills []domain.Skill) []table.Skill {
	var skillsRecords []table.Skill
	if len(skills) == 0 {
		emptySkill := domain.Skill("")
		skillsRecords = append(skillsRecords, table.Skill{
			Name:       "",
			TenantID:   sharedcontext.TenantIDFromContext(ctx),
			DocumentID: emptySkill.DocumentID(ctx).String(),
		})
		return skillsRecords
	}

	for _, skill := range skills {
		skillsRecords = append(skillsRecords, table.Skill{
			Name:       string(skill),
			TenantID:   sharedcontext.TenantIDFromContext(ctx),
			DocumentID: skill.DocumentID(ctx).String(),
		})
	}
	return skillsRecords
}
