// Package pactv3 generates the PACT v3 API client.
package pactv3

//go:generate echo flattening...
//go:generate npx openapi-flattener@v1.0.3 -s 00-original.yaml -o 01-flattened.yaml

//go:generate echo applying overlay...
//go:generate sh -c "go tool openapi-overlay apply overlay.yaml 00-original.yaml > 02-overlayed.yaml"

//go:generate echo downconverting...
//go:generate npx @apiture/openapi-down-convert@v0.14.1 --input 02-overlayed.yaml --output 03-downconverted.yaml

//go:generate echo generating code...
//go:generate go tool oapi-codegen -config config.yaml 03-downconverted.yaml

//go:generate echo postprocessing...
//go:generate sed -i -f postprocess.sed pactv3.gen.go
