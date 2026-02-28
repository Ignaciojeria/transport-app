package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/journey"
	"net/http"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

// CreateJourneyRequest es el body del POST para abrir una jornada.
type CreateJourneyRequest struct {
	OpenedBy string  `json:"openedBy"` // "USER" | "SYSTEM"
	Reason   *string `json:"reason"`   // opcional, ej. "Apertura manual"
}

// CreateJourneyResponse es la respuesta 201 con la jornada creada (campos en camelCase para la API).
type CreateJourneyResponse struct {
	ID       string  `json:"id"`
	MenuID   string  `json:"menuId"`
	Status   string  `json:"status"`
	OpenedAt string  `json:"openedAt"` // RFC3339
	OpenedBy string  `json:"openedBy"`
	Reason   *string `json:"reason,omitempty"`
}

func init() {
	ioc.Register(createJourneyHandler)
}

func createJourneyHandler(
	s httpserver.Server,
	obs observability.Observability,
	getActiveJourney supabaserepo.GetActiveJourneyByMenuID,
	insertJourney supabaserepo.InsertJourney,
	userHasMenu supabaserepo.UserHasMenu,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Post(s.Manager, "/api/menus/{menuId}/journeys",
		func(c fuego.ContextWithBody[CreateJourneyRequest]) (CreateJourneyResponse, error) {
			ctx := c.Context()
			spanCtx, span := obs.Tracer.Start(ctx, "createJourney")
			defer span.End()

			menuID := c.PathParam("menuId")
			if menuID == "" {
				return CreateJourneyResponse{}, fuego.HTTPError{
					Title:  "menuId is required",
					Detail: "menuId path parameter is required",
					Status: http.StatusBadRequest,
				}
			}

			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return CreateJourneyResponse{}, fuego.HTTPError{
					Title:  "unauthorized",
					Detail: "user id not found in context",
					Status: http.StatusUnauthorized,
				}
			}

			hasMenu, err := userHasMenu(spanCtx, userID, menuID)
			if err != nil || !hasMenu {
				return CreateJourneyResponse{}, fuego.HTTPError{
					Title:  "menu not found",
					Detail: "menu not found or you do not own it",
					Status: http.StatusNotFound,
				}
			}

			body, err := c.Body()
			if err != nil {
				return CreateJourneyResponse{}, fuego.HTTPError{
					Title:  "invalid body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			openedBy := journey.OpenedBy(body.OpenedBy)
			if openedBy != journey.OpenedByUser && openedBy != journey.OpenedBySystem {
				return CreateJourneyResponse{}, fuego.HTTPError{
					Title:  "invalid openedBy",
					Detail: "openedBy must be USER or SYSTEM",
					Status: http.StatusBadRequest,
				}
			}

			// Si ya hay una jornada abierta para este menú, fallar
			active, err := getActiveJourney(spanCtx, menuID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error checking active journey", "error", err, "menuID", menuID)
				return CreateJourneyResponse{}, fuego.HTTPError{
					Title:  "error checking journey",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			if active != nil {
				return CreateJourneyResponse{}, fuego.HTTPError{
					Title:  "journey already open",
					Detail: "ya existe una jornada abierta para ese menú",
					Status: http.StatusConflict, // 409
				}
			}

			j, err := insertJourney(spanCtx, menuID, openedBy, body.Reason)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error creating journey", "error", err, "menuID", menuID)
				return CreateJourneyResponse{}, fuego.HTTPError{
					Title:  "error creating journey",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			resp := CreateJourneyResponse{
				ID:       j.ID,
				MenuID:   j.MenuID,
				Status:   string(j.Status),
				OpenedAt: j.OpenedAt.Format("2006-01-02T15:04:05Z07:00"),
				OpenedBy: string(j.OpenedBy),
				Reason:   j.OpenedReason,
			}
			return resp, nil
		},
		option.Summary("Open a new journey for the menu"),
		option.Description("Creates a new OPEN journey for the given menu. Fails with 409 if there is already an open journey for this menu."),
		option.Tags("menu", "journey"),
		option.Middleware(jwtAuthMiddleware),
	)
}
