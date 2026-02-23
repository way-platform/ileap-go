package ileap

// EventType is a PACT CloudEvents event type string.
type EventType string

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
