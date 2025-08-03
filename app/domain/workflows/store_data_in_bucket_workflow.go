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

type StoreDataInBucketWorkflow struct {
	IdempotencyKey string
	NextInput      []byte
	storjBucket    *storjbucket.TransportAppBucket
	fsm            *fsm.FSM
}

func init() {
	ioc.Registry(NewStoreDataInBucketWorkflow,
		storjbucket.NewTransportAppBucket)
}

func NewStoreDataInBucketWorkflow(
	storjBucket *storjbucket.TransportAppBucket,
) (StoreDataInBucketWorkflow, error) {
	return StoreDataInBucketWorkflow{
		storjBucket: storjBucket,
	}, nil
}

func (w StoreDataInBucketWorkflow) Restore(ctx context.Context, idempotencyKey string) (StoreDataInBucketWorkflow, error) {
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
		transition = w.StoreDataInBucketStarted()
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
		transition = w.StoreDataInBucketStarted()
	}

	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.DataStoredInBucket(), Src: []string{w.StoreDataInBucketStarted()}, Dst: w.DataStoredInBucket()},
		},
		fsm.Callbacks{},
	)
	w.NextInput = nextInput
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

func (w StoreDataInBucketWorkflow) SaveState(ctx context.Context) error {
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

func (w StoreDataInBucketWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		NextInput:      w.NextInput,
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
