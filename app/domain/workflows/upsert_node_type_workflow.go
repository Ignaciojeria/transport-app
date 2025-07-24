package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertNodeTypeWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertNodeTypeWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertNodeTypeWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertNodeTypeWorkflow, error) {
	return UpsertNodeTypeWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertNodeTypeWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertNodeTypeWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertNodeTypeStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.NodeTypeUpserted(), Src: []string{w.UpsertNodeTypeStarted()}, Dst: w.NodeTypeUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertNodeTypeWorkflow) WorkflowName() string {
	return "upsert_node_type_workflow"
}

func (w UpsertNodeTypeWorkflow) NodeTypeUpserted() string {
	return "node_type_upserted"
}

func (w UpsertNodeTypeWorkflow) UpsertNodeTypeStarted() string {
	return "upsert_node_type_started"
}

func (w UpsertNodeTypeWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertNodeTypeWorkflow) CanTransitionToNodeTypeUpserted() bool {
	return w.fsm.Can(w.NodeTypeUpserted())
}

func (w UpsertNodeTypeWorkflow) IsNodeTypeUpserted() bool {
	return w.fsm.Current() == w.NodeTypeUpserted()
}

func (w UpsertNodeTypeWorkflow) IsUpsertNodeTypeStarted() bool {
	return w.fsm.Current() == w.UpsertNodeTypeStarted()
}

// SetNodeTypeUpsertedTransition completa el upsert de node type
func (w UpsertNodeTypeWorkflow) SetNodeTypeUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.NodeTypeUpserted())
}