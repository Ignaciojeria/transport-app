package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertNodeReferencesWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertNodeReferencesWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertNodeReferencesWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertNodeReferencesWorkflow, error) {
	return UpsertNodeReferencesWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertNodeReferencesWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertNodeReferencesWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertNodeReferencesStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.NodeReferencesUpserted(), Src: []string{w.UpsertNodeReferencesStarted()}, Dst: w.NodeReferencesUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertNodeReferencesWorkflow) WorkflowName() string {
	return "upsert_node_references_workflow"
}

func (w UpsertNodeReferencesWorkflow) NodeReferencesUpserted() string {
	return "node_references_upserted"
}

func (w UpsertNodeReferencesWorkflow) UpsertNodeReferencesStarted() string {
	return "upsert_node_references_started"
}

func (w UpsertNodeReferencesWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertNodeReferencesWorkflow) CanTransitionToNodeReferencesUpserted() bool {
	return w.fsm.Can(w.NodeReferencesUpserted())
}

func (w UpsertNodeReferencesWorkflow) IsNodeReferencesUpserted() bool {
	return w.fsm.Current() == w.NodeReferencesUpserted()
}

func (w UpsertNodeReferencesWorkflow) IsUpsertNodeReferencesStarted() bool {
	return w.fsm.Current() == w.UpsertNodeReferencesStarted()
}

// SetNodeReferencesUpsertedTransition completa el upsert de node references
func (w UpsertNodeReferencesWorkflow) SetNodeReferencesUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.NodeReferencesUpserted())
}