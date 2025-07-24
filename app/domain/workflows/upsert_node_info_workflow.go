package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertNodeInfoWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertNodeInfoWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertNodeInfoWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertNodeInfoWorkflow, error) {
	return UpsertNodeInfoWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertNodeInfoWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertNodeInfoWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertNodeInfoStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.NodeInfoUpserted(), Src: []string{w.UpsertNodeInfoStarted()}, Dst: w.NodeInfoUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertNodeInfoWorkflow) WorkflowName() string {
	return "upsert_node_info_workflow"
}

func (w UpsertNodeInfoWorkflow) NodeInfoUpserted() string {
	return "node_info_upserted"
}

func (w UpsertNodeInfoWorkflow) UpsertNodeInfoStarted() string {
	return "upsert_node_info_started"
}

func (w UpsertNodeInfoWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertNodeInfoWorkflow) CanTransitionToNodeInfoUpserted() bool {
	return w.fsm.Can(w.NodeInfoUpserted())
}

func (w UpsertNodeInfoWorkflow) IsNodeInfoUpserted() bool {
	return w.fsm.Current() == w.NodeInfoUpserted()
}

func (w UpsertNodeInfoWorkflow) IsUpsertNodeInfoStarted() bool {
	return w.fsm.Current() == w.UpsertNodeInfoStarted()
}

// SetNodeInfoUpsertedTransition completa el upsert de node info
func (w UpsertNodeInfoWorkflow) SetNodeInfoUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.NodeInfoUpserted())
}