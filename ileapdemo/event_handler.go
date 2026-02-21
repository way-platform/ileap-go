package ileapdemo

import (
	"context"
	"fmt"

	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/ileapserver"
)

// EventHandler implements ileapserver.EventHandler for the demo server.
type EventHandler struct{}

// HandleEvent processes an incoming PACT event.
func (h *EventHandler) HandleEvent(_ context.Context, event ileapserver.Event) error {
	switch ileap.EventType(event.Type) {
	case ileap.EventTypeRequestCreatedV1:
		// TODO: Handle RequestCreated.
	case ileap.EventTypeRequestFulfilledV1:
		// TODO: Handle RequestFulfilled.
	case ileap.EventTypeRequestRejectedV1:
		// TODO: Handle RequestRejected.
	case ileap.EventTypePublishedV1:
		// TODO: Handle Published.
	default:
		return fmt.Errorf("invalid event type: %s: %w", event.Type, ileapserver.ErrBadRequest)
	}
	return nil
}
