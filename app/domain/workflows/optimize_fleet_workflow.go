package workflows

import (
	"context"
	"encoding/json"
	"errors"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type OptimizeFleetWorkflow struct {
	IdempotencyKey string
	NextInput      []byte
	storjBucket    *storjbucket.TransportAppBucket
	fsm            *fsm.FSM
}

func init() {
	ioc.Registry(
		NewOptimizeFleetWorkflow,
		storjbucket.NewTransportAppBucket,
	)
}

func NewOptimizeFleetWorkflow(
	storjBucket *storjbucket.TransportAppBucket,
) (OptimizeFleetWorkflow, error) {
	return OptimizeFleetWorkflow{
		storjBucket: storjBucket,
	}, nil
}

func (w OptimizeFleetWorkflow) Restore(ctx context.Context, idempotencyKey string) (OptimizeFleetWorkflow, error) {
	w.IdempotencyKey = idempotencyKey

	// Obtener token desde el contexto
	token, ok := sharedcontext.BucketTokenFromContext(ctx)
	if !ok {
		return w, errors.New("token del bucket no encontrado en el contexto")
	}

	// Intentar recuperar el estado desde StorJ bucket
	data, err := w.storjBucket.DownloadWithToken(ctx, token, idempotencyKey)

	var transition string
	var nextInput []byte

	if err != nil {
		// Si no existe el estado, usar el estado inicial
		transition = w.OptimizationStarted()
		nextInput = nil
	} else {
		// Deserializar el estado guardado
		var fsmState domain.FSMState
		if err := json.Unmarshal(data, &fsmState); err != nil {
			return w, err
		}
		transition = fsmState.State
		nextInput = fsmState.NextInput
	}

	// Si no hay estado guardado, usar el estado inicial
	if transition == "" {
		transition = w.OptimizationStarted()
	}

	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.OptimizationCompleted(), Src: []string{w.OptimizationStarted()}, Dst: w.OptimizationCompleted()},
		},
		fsm.Callbacks{},
	)
	w.NextInput = nextInput
	return w, nil
}

func (w OptimizeFleetWorkflow) WorkflowName() string {
	return "optimize_fleet_workflow"
}

func (w OptimizeFleetWorkflow) OptimizationCompleted() string {
	return "optimization_completed"
}

func (w OptimizeFleetWorkflow) OptimizationStarted() string {
	return "optimization_started"
}

func (w OptimizeFleetWorkflow) SaveState(ctx context.Context) error {
	// Obtener token desde el contexto
	token, ok := sharedcontext.BucketTokenFromContext(ctx)
	if !ok {
		return errors.New("token del bucket no encontrado en el contexto")
	}

	// Serializar el estado actual
	fsmState := w.Map(ctx)
	data, err := json.Marshal(fsmState)
	if err != nil {
		return err
	}

	// Guardar en StorJ bucket usando el mismo patrón que en el publisher
	return w.storjBucket.UploadWithToken(ctx, token, w.IdempotencyKey, data)
}

func (w OptimizeFleetWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		NextInput:      w.NextInput,
		State:          w.fsm.Current(),
	}
}

func (w OptimizeFleetWorkflow) TransitionToOptimizationCompleted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OptimizationCompleted())
}

func (w OptimizeFleetWorkflow) CanTransitionToOptimizationCompleted() bool {
	return w.fsm.Can(w.OptimizationCompleted())
}

func (w OptimizeFleetWorkflow) IsOptimizationCompleted() bool {
	return w.fsm.Current() == w.OptimizationCompleted()
}

func (w OptimizeFleetWorkflow) IsOptimizationStarted() bool {
	return w.fsm.Current() == w.OptimizationStarted()
}

// CompleteOptimization completa la optimización de la flota
func (w OptimizeFleetWorkflow) SetOptimizationCompletedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OptimizationCompleted())
}
