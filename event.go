package ileap

import (
	"encoding/json"
	"time"
)

// Event is a PACT event.
//
// See: https://docs.carbon-transparency.org/tr/data-exchange-protocol/latest/#action-events
// See: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/bindings/http-protocol-binding.md#32-structured-content-mode
type Event struct {
	// Type is the type of the event.
	Type EventType `json:"type"`
	// Specversion is the version of the CloudEvents specification that the event uses.
	Specversion string `json:"specversion"`
	// ID is a unique identifier for the event.
	ID string `json:"id"`
	// Source is the source of the event.
	Source string `json:"source"`
	// Time is the time the event occurred.
	Time time.Time `json:"time"`
	// Data is the event data as raw JSON.
	Data json.RawMessage `json:"data"`
}

// EventType is the type of the event.
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
