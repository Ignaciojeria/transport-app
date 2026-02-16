package order

import (
	"context"
	"errors"

	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrMenuNotFound    = errors.New("menu not found")
	ErrNoActiveJourney = errors.New("no active journey")
)

// AssignOrdersToJourney asigna órdenes pendientes a la jornada activa del menú.
type AssignOrdersToJourney func(ctx context.Context, menuID string, aggregateIDs []int64) ([]supabaserepo.AssignOrdersResult, error)

func init() {
	ioc.Registry(
		NewAssignOrdersToJourney,
		supabaserepo.NewAssignOrdersToJourney,
		supabaserepo.NewGetActiveJourneyByMenuID,
		supabaserepo.NewUserHasMenu,
	)
}

func NewAssignOrdersToJourney(
	assignRepo supabaserepo.AssignOrdersToJourney,
	getActiveJourney supabaserepo.GetActiveJourneyByMenuID,
	userHasMenu supabaserepo.UserHasMenu,
) AssignOrdersToJourney {
	return func(ctx context.Context, menuID string, aggregateIDs []int64) ([]supabaserepo.AssignOrdersResult, error) {
		userID, ok := sharedcontext.UserIDFromContext(ctx)
		if !ok || userID == "" {
			return nil, ErrUnauthorized
		}

		hasMenu, err := userHasMenu(ctx, userID, menuID)
		if err != nil || !hasMenu {
			return nil, ErrMenuNotFound
		}

		active, err := getActiveJourney(ctx, menuID)
		if err != nil {
			return nil, err
		}
		if active == nil {
			return nil, ErrNoActiveJourney
		}

		if len(aggregateIDs) == 0 {
			return []supabaserepo.AssignOrdersResult{}, nil
		}

		return assignRepo(ctx, menuID, active.ID, aggregateIDs)
	}
}
