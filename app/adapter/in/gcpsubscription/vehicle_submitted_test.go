package gcpsubscription

import (
	"transport-app/mocks"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"cloud.google.com/go/pubsub"
)

func TestVehicleSubmitted(t *testing.T) {

	// Prepare mock data model to be used in the test
	expectedDataModel := map[string]interface{}{"key": "value"}
	mockData, err := json.Marshal(expectedDataModel)
	if err != nil {
		t.Fatalf("Failed to marshal expectedDataModel: %v", err)
	}

	// Create a mock pubsub.Message using the mock data
	mockMessage := &pubsub.Message{
		Data: mockData,
		ID:   "123",
	}

	// Initialize the MessageProcessor with a mock subscription reference (not directly used in Pull)
	mockSubscription := &pubsub.Subscription{}

	mockMgr := &mocks.MockSubscriptionManager{
		SubscriptionFunc: func(id string) *pubsub.Subscription {
			return mockSubscription
		},
	}
	messageProcessor := newVehicleSubmitted(mockMgr)

	// Invoke the Pull method with a background context and the mock message
	statusCode, err := messageProcessor(context.Background(), mockMessage)

	// Assert that there's no error from MessageProcessor.Pull
	if err != nil {
		t.Errorf("Expected no errors from MessageProcessor.Pull, got: %v", err)
	}

	// Assert that the HTTP status code is 200 OK
	if statusCode != http.StatusOK {
		t.Errorf("Expected HTTP status code 200 OK, got: %d", statusCode)
	}
}
