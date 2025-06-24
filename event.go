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
	// Id is a unique identifier for the event.
	Id string `json:"id"`
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
	EventTypeRequestCreated   EventType = "org.wbcsd.pact.ProductFootprint.RequestCreatedEvent.3"
	EventTypeRequestFulfilled EventType = "org.wbcsd.pact.ProductFootprint.RequestFulfilledEvent.3"
	EventTypeRequestRejected  EventType = "org.wbcsd.pact.ProductFootprint.RequestRejectedEvent.3"
	EventTypePublished        EventType = "org.wbcsd.pact.ProductFootprint.PublishedEvent.3"
)
