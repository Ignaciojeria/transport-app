package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertNodeInfoHeadersWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertNodeInfoHeadersWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertNodeInfoHeadersWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertNodeInfoHeadersWorkflow, error) {
	return UpsertNodeInfoHeadersWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertNodeInfoHeadersWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertNodeInfoHeadersWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertNodeInfoHeadersStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.NodeInfoHeadersUpserted(), Src: []string{w.UpsertNodeInfoHeadersStarted()}, Dst: w.NodeInfoHeadersUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertNodeInfoHeadersWorkflow) WorkflowName() string {
	return "upsert_node_info_headers_workflow"
}

func (w UpsertNodeInfoHeadersWorkflow) NodeInfoHeadersUpserted() string {
	return "node_info_headers_upserted"
}

func (w UpsertNodeInfoHeadersWorkflow) UpsertNodeInfoHeadersStarted() string {
	return "upsert_node_info_headers_started"
}

func (w UpsertNodeInfoHeadersWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertNodeInfoHeadersWorkflow) TransitionToNodeInfoHeadersUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.NodeInfoHeadersUpserted())
}

func (w UpsertNodeInfoHeadersWorkflow) CanTransitionToNodeInfoHeadersUpserted() bool {
	return w.fsm.Can(w.NodeInfoHeadersUpserted())
}

func (w UpsertNodeInfoHeadersWorkflow) IsNodeInfoHeadersUpserted() bool {
	return w.fsm.Current() == w.NodeInfoHeadersUpserted()
}

func (w UpsertNodeInfoHeadersWorkflow) IsUpsertNodeInfoHeadersStarted() bool {
	return w.fsm.Current() == w.UpsertNodeInfoHeadersStarted()
}

// SetNodeInfoHeadersUpsertedTransition completa el upsert de node info headers
func (w UpsertNodeInfoHeadersWorkflow) SetNodeInfoHeadersUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.NodeInfoHeadersUpserted())
}