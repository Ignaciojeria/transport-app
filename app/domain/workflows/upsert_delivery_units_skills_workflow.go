package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertDeliveryUnitsSkillsWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertDeliveryUnitsSkillsWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertDeliveryUnitsSkillsWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertDeliveryUnitsSkillsWorkflow, error) {
	return UpsertDeliveryUnitsSkillsWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertDeliveryUnitsSkillsWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertDeliveryUnitsSkillsWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertDeliveryUnitsSkillsStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.DeliveryUnitsSkillsUpserted(), Src: []string{w.UpsertDeliveryUnitsSkillsStarted()}, Dst: w.DeliveryUnitsSkillsUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertDeliveryUnitsSkillsWorkflow) WorkflowName() string {
	return "upsert_delivery_units_skills_workflow"
}

func (w UpsertDeliveryUnitsSkillsWorkflow) DeliveryUnitsSkillsUpserted() string {
	return "delivery_units_skills_upserted"
}

func (w UpsertDeliveryUnitsSkillsWorkflow) UpsertDeliveryUnitsSkillsStarted() string {
	return "upsert_delivery_units_skills_started"
}

func (w UpsertDeliveryUnitsSkillsWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertDeliveryUnitsSkillsWorkflow) TransitionToDeliveryUnitsSkillsUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.DeliveryUnitsSkillsUpserted())
}

func (w UpsertDeliveryUnitsSkillsWorkflow) CanTransitionToDeliveryUnitsSkillsUpserted() bool {
	return w.fsm.Can(w.DeliveryUnitsSkillsUpserted())
}

func (w UpsertDeliveryUnitsSkillsWorkflow) IsDeliveryUnitsSkillsUpserted() bool {
	return w.fsm.Current() == w.DeliveryUnitsSkillsUpserted()
}

func (w UpsertDeliveryUnitsSkillsWorkflow) IsUpsertDeliveryUnitsSkillsStarted() bool {
	return w.fsm.Current() == w.UpsertDeliveryUnitsSkillsStarted()
}

// SetDeliveryUnitsSkillsUpsertedTransition completa el upsert de delivery units skills
func (w UpsertDeliveryUnitsSkillsWorkflow) SetDeliveryUnitsSkillsUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.DeliveryUnitsSkillsUpserted())
}