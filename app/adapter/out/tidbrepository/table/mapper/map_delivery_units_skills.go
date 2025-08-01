package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapDeliveryUnitsSkills(ctx context.Context, order domain.Order) []table.DeliveryUnitsSkills {
	var skills []table.DeliveryUnitsSkills
	for _, du := range order.DeliveryUnits {
		if len(du.Skills) == 0 {
			skills = append(skills, table.DeliveryUnitsSkills{
				Skill:           "",
				DeliveryUnitDoc: du.DocID(ctx).String(),
			})
			continue
		}
		for _, skill := range du.Skills {
			skills = append(skills, table.DeliveryUnitsSkills{
				Skill:           string(skill),
				DeliveryUnitDoc: du.DocID(ctx).String(),
			})
		}
	}
	if len(order.DeliveryUnits) == 0 {
		emptyDu := domain.DeliveryUnit{}
		deliveryUnitDoc := emptyDu.DocID(ctx).String()
		skills = append(skills, table.DeliveryUnitsSkills{
			Skill:           "",
			DeliveryUnitDoc: deliveryUnitDoc,
		})
	}
	return skills
}
