package eventprocessing

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	"micartapro/app/shared/infrastructure/httpserver"
)

// ------------------------------------------------------------
// Logger
// ------------------------------------------------------------

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// ------------------------------------------------------------
// Broker-agnostic Subscriber interface
// ------------------------------------------------------------

// ------------------------------------------------------------
// Pub/Sub implementation
// ------------------------------------------------------------

type PubSubSubscriber struct {
	client     *pubsub.Client
	httpServer httpserver.Server
}

func NewPubSubSubscriber(c *pubsub.Client, s httpserver.Server) Subscriber {
	return &PubSubSubscriber{client: c, httpServer: s}
}

// ------------------------------------------------------------
// Start: sets up both PULL and PUSH handling
// ------------------------------------------------------------

func (ps *PubSubSubscriber) Start(subscriptionName string, processor MessageProcessor) error {

	sub := ps.client.Subscriber(subscriptionName)
	sub.ReceiveSettings.MaxOutstandingMessages = 10

	// -------------------------------
	// PULL consumer
	// -------------------------------
	go func() {
		ctx := context.Background()
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {

			ce := ps.convertPullMessage(subscriptionName, msg)

			status, err := processor(ctx, ce)

			if err != nil || status >= 500 {
				msg.Nack()
				return
			}

			msg.Ack()
		})

		if err != nil {
			logger.Error("pubsub_receive_failed",
				"subscription", subscriptionName,
				"error", err.Error(),
			)
			time.Sleep(5 * time.Second)
			ps.Start(subscriptionName, processor) // retry
		}
	}()

	// -------------------------------
	// PUSH consumer via HTTP
	// -------------------------------
	path := "/subscription/" + subscriptionName
	httpserver.WrapPostStd(ps.httpServer, path, ps.makePushHandler(subscriptionName, processor))

	return nil
}

// ------------------------------------------------------------
// Pull message → CloudEvent
// ------------------------------------------------------------

func (ps *PubSubSubscriber) convertPullMessage(subName string, msg *pubsub.Message) cloudevents.Event {
	ce := cloudevents.NewEvent()

	ce.SetID(msg.ID)
	ce.SetType("google.pubsub.pull")
	ce.SetSource("gcp.pubsub/" + subName)
	ce.SetData(cloudevents.ApplicationJSON, msg.Data)

	for k, v := range msg.Attributes {
		ce.SetExtension(strings.ToLower(k), v)
	}

	return ce
}

// ------------------------------------------------------------
// Push Handler Factory
// ------------------------------------------------------------

func (ps *PubSubSubscriber) makePushHandler(subName string, processor MessageProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Native GCP push
		if r.Header.Get("X-Goog-Channel-ID") != "" {
			ps.handleNativePush(subName, processor, w, r)
			return
		}

		// Manual custom POST for testing
		ps.handleManualPush(subName, processor, w, r)
	}
}

// ------------------------------------------------------------
// Native Push
// ------------------------------------------------------------

func (ps *PubSubSubscriber) handleNativePush(subName string, processor MessageProcessor, w http.ResponseWriter, r *http.Request) {

	var envelope struct {
		Message struct {
			MessageID  string            `json:"messageId"`
			Data       []byte            `json:"data"`
			Attributes map[string]string `json:"attributes"`
		} `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ce := cloudevents.NewEvent()
	ce.SetID(envelope.Message.MessageID)
	ce.SetType("google.pubsub.push")
	ce.SetSource("gcp.pubsub/" + subName)
	ce.SetData(cloudevents.ApplicationJSON, envelope.Message.Data)

	for k, v := range envelope.Message.Attributes {
		ce.SetExtension(strings.ToLower(k), v)
	}

	status, err := processor(r.Context(), ce)

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(status)
}

// ------------------------------------------------------------
// Manual Push (webhook, testing, dev)
// ------------------------------------------------------------

func (ps *PubSubSubscriber) handleManualPush(subName string, processor MessageProcessor, w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ce := cloudevents.NewEvent()
	ce.SetID("")
	ce.SetType("manual.message")
	ce.SetSource("manual/" + subName)
	ce.SetData(cloudevents.ApplicationJSON, body)

	// Use request headers as CloudEvent extensions
	for key, values := range r.Header {
		if len(values) > 0 {
			ce.SetExtension(strings.ToLower(key), strings.Join(values, ","))
		}
	}

	status, err := processor(r.Context(), ce)

	if err != nil {
		w.WriteHeader(status)
		return
	}

	w.WriteHeader(status)
}
