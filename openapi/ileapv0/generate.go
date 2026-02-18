// Package ileapv0 generates the iLEAP v0 API client.
package ileapv0

//go:generate sh -c "go tool openapi-overlay apply overlay.yaml openapi.original.json > api.json"
//go:generate go tool oapi-codegen -config cfg.yaml api.json
