package ileap

// eventType is a PACT CloudEvents event type string.
type eventType string

// event is a structured-mode CloudEvent payload accepted by POST /2/events.
type event struct {
	Type        string
	Specversion string
	ID          string
	Source      string
	Data        []byte
}

// Known event types.
const (
	eventTypeRequestCreatedV1   eventType = "org.wbcsd.pathfinder.ProductFootprintRequest.Created.v1"
	eventTypePublishedV1        eventType = "org.wbcsd.pathfinder.ProductFootprint.Published.v1"
	eventTypeRequestFulfilledV1 eventType = "org.wbcsd.pathfinder.ProductFootprintRequest.Fulfilled.v1"
	eventTypeRequestRejectedV1  eventType = "org.wbcsd.pathfinder.ProductFootprintRequest.Rejected.v1"
)

// isKnownEventType reports whether the event type is supported by the server.
func isKnownEventType(t eventType) bool {
	switch t {
	case eventTypeRequestCreatedV1,
		eventTypeRequestFulfilledV1,
		eventTypeRequestRejectedV1,
		eventTypePublishedV1:
		return true
	default:
		return false
	}
}
