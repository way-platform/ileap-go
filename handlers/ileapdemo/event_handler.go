package ileapdemo

import (
	"context"
	"fmt"

	"github.com/way-platform/ileap-go"
)

// EventHandler implements ileap.EventHandler for the demo server.
type EventHandler struct{}

// HandleEvent processes an incoming PACT event.
func (h *EventHandler) HandleEvent(_ context.Context, event ileap.Event) error {
	switch event.Type {
	case ileap.EventTypeRequestCreatedV1:
		// TODO: Handle RequestCreated.
	case ileap.EventTypeRequestFulfilledV1:
		// TODO: Handle RequestFulfilled.
	case ileap.EventTypeRequestRejectedV1:
		// TODO: Handle RequestRejected.
	case ileap.EventTypePublishedV1:
		// TODO: Handle Published.
	default:
		return fmt.Errorf("invalid event type: %s: %w", event.Type, ileap.ErrBadRequest)
	}
	return nil
}
