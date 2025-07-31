package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/adapter/out/vroom"
	"transport-app/app/domain/optimization"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OptimizeFleetWorkflow func(ctx context.Context, input optimization.FleetOptimization) error

func init() {
	ioc.Registry(
		NewOptimizeFleetWorkflow,
		workflows.NewOptimizeFleetWorkflow,
		vroom.NewOptimize,
		storjbucket.NewTransportAppBucket,
		tidbrepository.NewSaveFSMTransition,
	)
}

func NewOptimizeFleetWorkflow(
	domainWorkflow workflows.OptimizeFleetWorkflow,
	optimize vroom.Optimize,
	storjBucket *storjbucket.TransportAppBucket,
	saveFSMTransition tidbrepository.SaveFSMTransition,
) OptimizeFleetWorkflow {
	return func(ctx context.Context, input optimization.FleetOptimization) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetOptimizationCompletedTransition(ctx); err != nil {
			return fmt.Errorf("failed to set optimization completed transition: %w", err)
		}
		routeRequests, err := optimize(ctx, input)
		if err != nil {
			return err
		}
		token, ok := sharedcontext.BucketTokenFromContext(ctx)
		if !ok {
			return fmt.Errorf("bucket token not found in context")
		}

		for _, v := range routeRequests {
			routeBytes, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("error marshaling route request: %w", err)
			}
			err = storjBucket.UploadWithToken(ctx, token, v.ReferenceID, routeBytes)
			if err != nil {
				return fmt.Errorf("error uploading route request: %w", err)
			}
		}
		fsmState := workflow.Map(ctx)
		err = saveFSMTransition(ctx, fsmState)
		if err != nil {
			return fmt.Errorf("failed to save FSM transition: %w", err)
		}
		fmt.Printf("Optimización completada. Se generaron %d rutas.\n", len(routeRequests))
		return nil
	}
}

// saveRouteRequestsToFile guarda las rutas y el input de optimización en un archivo JSON
func saveRouteRequestsToFile(routeRequests []request.UpsertRouteRequest, input optimization.FleetOptimization) error {
	// Crear estructura con rutas y input para análisis
	optimizationResult := struct {
		Timestamp     string                         `json:"timestamp"`
		Input         optimization.FleetOptimization `json:"input"`
		RouteRequests []request.UpsertRouteRequest   `json:"routeRequests"`
		RouteCount    int                            `json:"routeCount"`
	}{
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Input:         input,
		RouteRequests: routeRequests,
		RouteCount:    len(routeRequests),
	}

	// Convertir a JSON formateado
	jsonData, err := json.MarshalIndent(optimizationResult, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Crear nombre de archivo con timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("fleet_optimization_routes_%s.json", timestamp)

	// Escribir el JSON en el archivo
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	fmt.Printf("Rutas de optimización guardadas en: %s\n", filename)
	return nil
}
