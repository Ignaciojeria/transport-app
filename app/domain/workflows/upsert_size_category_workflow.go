package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertSizeCategoryWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertSizeCategoryWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertSizeCategoryWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertSizeCategoryWorkflow, error) {
	return UpsertSizeCategoryWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertSizeCategoryWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertSizeCategoryWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertSizeCategoryStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.SizeCategoryUpserted(), Src: []string{w.UpsertSizeCategoryStarted()}, Dst: w.SizeCategoryUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertSizeCategoryWorkflow) WorkflowName() string {
	return "upsert_size_category_workflow"
}

func (w UpsertSizeCategoryWorkflow) SizeCategoryUpserted() string {
	return "size_category_upserted"
}

func (w UpsertSizeCategoryWorkflow) UpsertSizeCategoryStarted() string {
	return "upsert_size_category_started"
}

func (w UpsertSizeCategoryWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertSizeCategoryWorkflow) CanTransitionToSizeCategoryUpserted() bool {
	return w.fsm.Can(w.SizeCategoryUpserted())
}

func (w UpsertSizeCategoryWorkflow) IsSizeCategoryUpserted() bool {
	return w.fsm.Current() == w.SizeCategoryUpserted()
}

func (w UpsertSizeCategoryWorkflow) IsUpsertSizeCategoryStarted() bool {
	return w.fsm.Current() == w.UpsertSizeCategoryStarted()
}

// SetSizeCategoryUpsertedTransition completa el upsert de size category
func (w UpsertSizeCategoryWorkflow) SetSizeCategoryUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.SizeCategoryUpserted())
}