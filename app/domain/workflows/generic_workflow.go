package workflows

import (
	"context"
	"fmt"
	"strings"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type WorkflowConfig struct {
	Name         string
	EntityName   string
	StartedState string
	CompletedState string
	UseStorjBucket bool
}

type GenericWorkflow struct {
	config               WorkflowConfig
	IdempotencyKey       string
	NextInput           []byte
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                 *fsm.FSM
}

func init() {
	ioc.Registry(NewGenericWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewGenericWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (GenericWorkflow, error) {
	return GenericWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w GenericWorkflow) WithConfig(config WorkflowConfig) GenericWorkflow {
	w.config = config
	return w
}

// Método principal que recibe la configuración directamente
func (w GenericWorkflow) Initialize(ctx context.Context, idempotencyKey string, config WorkflowConfig) (GenericWorkflow, error) {
	w.config = config
	return w.Restore(ctx, idempotencyKey)
}

func (w GenericWorkflow) Restore(ctx context.Context, idempotencyKey string) (GenericWorkflow, error) {
	w.IdempotencyKey = idempotencyKey

	var transition string
	var nextInput []byte
	var err error

	if w.config.UseStorjBucket {
		transition, nextInput, err = w.restoreFromStorj(ctx, idempotencyKey)
	} else {
		transition, err = w.restoreFromTiDB(ctx, idempotencyKey)
	}

	if err != nil {
		return w, err
	}

	if transition == "" {
		transition = w.config.StartedState
	}

	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.config.CompletedState, Src: []string{w.config.StartedState}, Dst: w.config.CompletedState},
		},
		fsm.Callbacks{},
	)
	w.NextInput = nextInput
	return w, nil
}

func (w GenericWorkflow) restoreFromTiDB(ctx context.Context, idempotencyKey string) (string, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.config.Name)
	if err != nil {
		return "", err
	}
	return lastTransition.State, nil
}

func (w GenericWorkflow) restoreFromStorj(ctx context.Context, idempotencyKey string) (string, []byte, error) {
	// Storj workflow storage deprecated - return default state
	return w.config.StartedState, nil, nil
}

func (w GenericWorkflow) WorkflowName() string {
	return w.config.Name
}

func (w GenericWorkflow) StartedState() string {
	return w.config.StartedState
}

func (w GenericWorkflow) CompletedState() string {
	return w.config.CompletedState
}

func (w GenericWorkflow) SaveState(ctx context.Context) error {
	// Storj workflow storage deprecated - no-op
	return nil
}

func (w GenericWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		NextInput:      w.NextInput,
		State:          w.fsm.Current(),
	}
}

func (w GenericWorkflow) TransitionToCompleted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.config.CompletedState)
}

func (w GenericWorkflow) CanTransitionToCompleted() bool {
	return w.fsm.Can(w.config.CompletedState)
}

func (w GenericWorkflow) IsCompleted() bool {
	return w.fsm.Current() == w.config.CompletedState
}

func (w GenericWorkflow) IsStarted() bool {
	return w.fsm.Current() == w.config.StartedState
}

func (w GenericWorkflow) SetCompletedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.config.CompletedState)
}

// Métodos de conveniencia para diferentes tipos de workflows
func CreateUpsertWorkflow(entityName string) WorkflowConfig {
	return WorkflowConfig{
		Name:           fmt.Sprintf("upsert_%s_workflow", strings.ToLower(entityName)),
		EntityName:     entityName,
		StartedState:   fmt.Sprintf("upsert_%s_started", strings.ToLower(entityName)),
		CompletedState: fmt.Sprintf("%s_upserted", strings.ToLower(entityName)),
		UseStorjBucket: false,
	}
}

func CreateWorkflow(entityName, action string) WorkflowConfig {
	return WorkflowConfig{
		Name:           fmt.Sprintf("%s_%s_workflow", strings.ToLower(action), strings.ToLower(entityName)),
		EntityName:     entityName,
		StartedState:   fmt.Sprintf("%s_%s_started", strings.ToLower(action), strings.ToLower(entityName)),
		CompletedState: fmt.Sprintf("%s_%s", strings.ToLower(entityName), strings.ToLower(action)),
		UseStorjBucket: false,
	}
}

func CreateStorjWorkflow(name, startedState, completedState string) WorkflowConfig {
	return WorkflowConfig{
		Name:           name,
		StartedState:   startedState,
		CompletedState: completedState,
		UseStorjBucket: true,
	}
}