package subscriptionwrapper

import (
	"transport-app/app/shared/infrastructure/gcppubsub"
	"transport-app/app/shared/infrastructure/httpserver"
	"encoding/json"
	"log/slog"
	"os"

	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SubscriptionManager interface {
	Subscription(id string) *pubsub.Subscription
	WithMessageProcessor(mp MessageProcessor) SubscriptionManager
	WithPushHandler(path string) SubscriptionManager
	Start(subscriptionRef *pubsub.Subscription) (SubscriptionManager, error)
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type SubscriptionWrapper struct {
	client           *pubsub.Client
	httpServer       httpserver.Server
	messageProcessor MessageProcessor
}

func init() {
	ioc.Registry(
		NewSubscriptionManager,
		gcppubsub.NewClient,
		httpserver.New,
	)
}
func NewSubscriptionManager(
	c *pubsub.Client,
	s httpserver.Server) SubscriptionManager {
	return &SubscriptionWrapper{client: c, httpServer: s}
}

func newSubscriptionManagerWithMessageProcessor(
	c *pubsub.Client,
	s httpserver.Server,
	mp MessageProcessor) SubscriptionManager {
	return &SubscriptionWrapper{client: c, httpServer: s, messageProcessor: mp}
}

func (sw *SubscriptionWrapper) Subscription(id string) *pubsub.Subscription {
	return sw.client.Subscription(id)
}

func (sw *SubscriptionWrapper) WithMessageProcessor(mp MessageProcessor) SubscriptionManager {
	return newSubscriptionManagerWithMessageProcessor(sw.client, sw.httpServer, mp)
}

func (s *SubscriptionWrapper) Start(subscriptionRef *pubsub.Subscription) (SubscriptionManager, error) {
	ctx := context.Background()
	if err := subscriptionRef.Receive(ctx, s.receive); err != nil {
		logger.Error(
			"subscription_signal_broken",
			"subscription_name", subscriptionRef.String(),
			"error", err.Error(),
		)
		time.Sleep(10 * time.Second)
		go s.Start(subscriptionRef)
		return s, err
	}
	return s, nil
}

func (s *SubscriptionWrapper) receive(ctx context.Context, m *pubsub.Message) {
	s.messageProcessor(ctx, m)
}

func (s *SubscriptionWrapper) WithPushHandler(path string) SubscriptionManager {
	httpserver.WrapPostStd(s.httpServer, path, s.pushHandler)
	return s
}

func (s *SubscriptionWrapper) pushHandler(w http.ResponseWriter, r *http.Request) {
	googleChannel := r.Header.Get("X-Goog-Channel-ID")

	var msg pubsub.Message

	if googleChannel != "" {
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, "error binding Pub/Sub message", http.StatusNoContent)
			return
		}

		statusCode, err := s.messageProcessor(r.Context(), &msg)
		if statusCode >= 500 && statusCode <= 599 {
			w.WriteHeader(statusCode)
			return
		}
		if err != nil {
			http.Error(w, "", http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	msg.Attributes = make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			msg.Attributes[strings.ToLower(key)] = strings.Join(values, ",")
		}
	}

	msg.Data = body
	statusCode, err := s.messageProcessor(r.Context(), &msg)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	w.WriteHeader(http.StatusOK)
}
