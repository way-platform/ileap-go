package ileap

// EventType is a PACT CloudEvents event type string.
type EventType string

// Event is a structured-mode CloudEvent payload accepted by POST /2/events.
type Event struct {
	Type        string
	Specversion string
	ID          string
	Source      string
	Data        []byte
}

// Known event types.
const (
	EventTypeRequestCreatedV3   EventType = "org.wbcsd.pact.ProductFootprint.RequestCreatedEvent.3"
	EventTypeRequestFulfilledV3 EventType = "org.wbcsd.pact.ProductFootprint.RequestFulfilledEvent.3"
	EventTypeRequestRejectedV3  EventType = "org.wbcsd.pact.ProductFootprint.RequestRejectedEvent.3"
	EventTypePublishedV3        EventType = "org.wbcsd.pact.ProductFootprint.PublishedEvent.3"

	EventTypeRequestCreatedV1   EventType = "org.wbcsd.pathfinder.ProductFootprintRequest.Created.v1"
	EventTypePublishedV1        EventType = "org.wbcsd.pathfinder.ProductFootprint.Published.v1"
	EventTypeRequestFulfilledV1 EventType = "org.wbcsd.pathfinder.ProductFootprintRequest.Fulfilled.v1"
	EventTypeRequestRejectedV1  EventType = "org.wbcsd.pathfinder.ProductFootprintRequest.Rejected.v1"
)

// IsKnownEventType reports whether the event type is supported by the server.
func IsKnownEventType(t EventType) bool {
	switch t {
	case EventTypeRequestCreatedV1,
		EventTypeRequestFulfilledV1,
		EventTypeRequestRejectedV1,
		EventTypePublishedV1,
		EventTypeRequestCreatedV3,
		EventTypeRequestFulfilledV3,
		EventTypeRequestRejectedV3,
		EventTypePublishedV3:
		return true
	default:
		return false
	}
}
