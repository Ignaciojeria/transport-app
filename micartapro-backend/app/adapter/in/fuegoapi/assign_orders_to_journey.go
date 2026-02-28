package fuegoapi

import (
	"errors"

	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/usecase/order"
	"net/http"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

// AssignOrdersToJourneyRequest es el body del POST para asignar órdenes a la jornada activa.
type AssignOrdersToJourneyRequest struct {
	AggregateIDs []int64 `json:"aggregateIds"`
}

// AssignOrdersToJourneyResponse es la respuesta con el resultado por orden.
type AssignOrdersToJourneyResponse struct {
	Results []AssignOrderResult `json:"results"`
}

// AssignOrderResult resultado de asignar una orden.
type AssignOrderResult struct {
	AggregateID int64 `json:"aggregateId"`
	OrderNumber *int  `json:"orderNumber,omitempty"`
	Assigned    bool  `json:"assigned"`
}

func init() {
	ioc.Register(assignOrdersToJourneyHandler)
}

func assignOrdersToJourneyHandler(
	s httpserver.Server,
	obs observability.Observability,
	assignOrdersUseCase order.AssignOrdersToJourney,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Post(s.Manager, "/api/menus/{menuId}/journeys/assign-orders",
		func(c fuego.ContextWithBody[AssignOrdersToJourneyRequest]) (AssignOrdersToJourneyResponse, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "assignOrdersToJourney")
			defer span.End()

			menuID := c.PathParam("menuId")
			if menuID == "" {
				return AssignOrdersToJourneyResponse{}, fuego.HTTPError{
					Title:  "menuId is required",
					Detail: "menuId path parameter is required",
					Status: http.StatusBadRequest,
				}
			}

			body, err := c.Body()
			if err != nil {
				return AssignOrdersToJourneyResponse{}, fuego.HTTPError{
					Title:  "invalid body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			results, err := assignOrdersUseCase(spanCtx, menuID, body.AggregateIDs)
			if err != nil {
				if errors.Is(err, order.ErrUnauthorized) {
					return AssignOrdersToJourneyResponse{}, fuego.HTTPError{
						Title:  "unauthorized",
						Detail: err.Error(),
						Status: http.StatusUnauthorized,
					}
				}
				if errors.Is(err, order.ErrMenuNotFound) {
					return AssignOrdersToJourneyResponse{}, fuego.HTTPError{
						Title:  "menu not found",
						Detail: err.Error(),
						Status: http.StatusNotFound,
					}
				}
				if errors.Is(err, order.ErrNoActiveJourney) {
					return AssignOrdersToJourneyResponse{}, fuego.HTTPError{
						Title:  "no active journey",
						Detail: "no hay jornada abierta para asignar órdenes",
						Status: http.StatusConflict,
					}
				}
				obs.Logger.ErrorContext(spanCtx, "error assigning orders to journey", "error", err, "menuID", menuID)
				return AssignOrdersToJourneyResponse{}, fuego.HTTPError{
					Title:  "error assigning orders",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			respResults := make([]AssignOrderResult, len(results))
			for i, r := range results {
				respResults[i] = AssignOrderResult{
					AggregateID: r.AggregateID,
					OrderNumber: r.OrderNumber,
					Assigned:    r.Assigned,
				}
			}

			return AssignOrdersToJourneyResponse{Results: respResults}, nil
		},
		option.Summary("Assign pending orders to active journey"),
		option.Description("Assigns orders with journey_id IS NULL to the active journey. Returns result per order."),
		option.Tags("menu", "journey", "orders"),
		option.Path("menuId", "string", param.Required()),
		option.Middleware(jwtAuthMiddleware),
	)
}
