package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type StoreDataInBucketWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewStoreDataInBucketWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewStoreDataInBucketWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (StoreDataInBucketWorkflow, error) {
	return StoreDataInBucketWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w StoreDataInBucketWorkflow) Restore(ctx context.Context, idempotencyKey string) (StoreDataInBucketWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.StoreDataInBucketStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.DataStoredInBucket(), Src: []string{w.StoreDataInBucketStarted()}, Dst: w.DataStoredInBucket()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w StoreDataInBucketWorkflow) WorkflowName() string {
	return "store_data_in_bucket_workflow"
}

func (w StoreDataInBucketWorkflow) DataStoredInBucket() string {
	return "data_stored_in_bucket"
}

func (w StoreDataInBucketWorkflow) StoreDataInBucketStarted() string {
	return "store_data_in_bucket_started"
}

func (w StoreDataInBucketWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w StoreDataInBucketWorkflow) CanTransitionToDataStoredInBucket() bool {
	return w.fsm.Can(w.DataStoredInBucket())
}

func (w StoreDataInBucketWorkflow) IsDataStoredInBucket() bool {
	return w.fsm.Current() == w.DataStoredInBucket()
}

func (w StoreDataInBucketWorkflow) IsStoreDataInBucketStarted() bool {
	return w.fsm.Current() == w.StoreDataInBucketStarted()
}

// SetDataStoredInBucketTransition completa el almacenamiento de datos en el bucket
func (w StoreDataInBucketWorkflow) SetDataStoredInBucketTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.DataStoredInBucket())
}
