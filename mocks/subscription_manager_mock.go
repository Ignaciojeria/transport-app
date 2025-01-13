package mocks

import (
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"

	"cloud.google.com/go/pubsub"
)

// MockSubscriptionManager is a mock implementation of the SubscriptionManager interface.
type MockSubscriptionManager struct {
	SubscriptionFunc         func(id string) *pubsub.Subscription
	WithMessageProcessorFunc func(mp subscriptionwrapper.MessageProcessor) subscriptionwrapper.SubscriptionManager
	WithPushHandlerFunc      func(path string) subscriptionwrapper.SubscriptionManager
	StartFunc                func() (subscriptionwrapper.SubscriptionManager, error)
}

// Ensure MockSubscriptionManager implements SubscriptionManager.
var _ subscriptionwrapper.SubscriptionManager = &MockSubscriptionManager{}

// New creates a new mock instance of a subscription, simulating the behavior of the interface's New method.
func (m *MockSubscriptionManager) Subscription(id string) *pubsub.Subscription {
	if m.SubscriptionFunc != nil {
		return m.SubscriptionFunc(id)
	}
	return nil // or return a dummy *pubsub.Subscription if needed
}

// From simulates the behavior of the From method of the interface, allowing chaining.
func (m *MockSubscriptionManager) WithMessageProcessor(mp subscriptionwrapper.MessageProcessor) subscriptionwrapper.SubscriptionManager {
	if m.WithMessageProcessorFunc != nil {
		return m.WithMessageProcessorFunc(mp)
	}
	return m // Return self to allow chaining
}

// WithPushHandler simulates the behavior of the interface's WithPushHandler method, allowing chaining.
func (m *MockSubscriptionManager) WithPushHandler(path string) subscriptionwrapper.SubscriptionManager {
	if m.WithPushHandlerFunc != nil {
		return m.WithPushHandlerFunc(path)
	}
	return m // Return self to allow chaining
}

// Start simulates the behavior of the Start method of the interface, initializing the subscription.
func (m *MockSubscriptionManager) Start(subRef *pubsub.Subscription) (subscriptionwrapper.SubscriptionManager, error) {
	if m.StartFunc != nil {
		return m.StartFunc()
	}
	return m, nil // or return an error if needed
}
