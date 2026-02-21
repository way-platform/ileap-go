// Package ileapv1 generates the iLEAP v1 API client.
package ileapv1

//go:generate sh -c "go tool -modfile=../../tools/go.mod openapi-overlay apply overlay.yaml openapi.original.json > api.json"
//go:generate go tool -modfile=../../tools/go.mod oapi-codegen -config cfg.yaml api.json
//go:generate sed -i -f postprocess.sed api.gen.go
