# iLEAP Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/way-platform/ileap-go)](https://pkg.go.dev/github.com/way-platform/ileap-go)
[![GoReportCard](https://goreportcard.com/badge/github.com/way-platform/ileap-go)](https://goreportcard.com/report/github.com/way-platform/ileap-go)
[![CI](https://github.com/way-platform/ileap-go/actions/workflows/ci.yaml/badge.svg)](https://github.com/way-platform/ileap-go/actions/workflows/ci.yaml)

A Go SDK for logistics emissions data compatible with the [iLEAP Technical
Specifications](https://sine-fdn.github.io/ileap-extension/).

## SDK

### Features

* Support for [PACT Product Footprints](https://wbcsd.github.io/tr/2024/data-exchange-protocol-20241024/#dt-pf) with [iLEAP extensions](https://sine-fdn.github.io/ileap-extension/#pcf-mapping).
* Support for [iLEAP Transport Activity Data](https://sine-fdn.github.io/ileap-extension/#dt-tad)

### Installing

```bash
$ go get github.com/way-platform/ileap-go@latest
```

### Using

```go
client := ileap.NewClient(
    ileap.WithBaseURL(os.Getenv("BASE_URL")),
    ileap.WithOAuth2(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")),
)
footprint, err := client.GetFootprint(context.Background(), &ileap.GetFootprintRequest{
    ID: "91715e5e-fd0b-4d1c-8fab-76290c46e6ed",
})
if err != nil {
    // TODO: Handle error.
}
fmt.Println(footprint)
```

### Developing

#### Build project

The project is built using [Mage](https://magefile.org), see
[magefile.go](./magefile.go).

```bash
$ go tool mage build
```

For all available build tasks, see:

```bash
$ go tool mage
```

## CLI tool

<img src="docs/cli.gif" />

### Installing

```bash
go install github.com/way-platform/ileap-go/cmd/ileap@latest
```

### Using

The following example logs in to the [SINE Foundation]()'s demo API.

```bash
$ ileap auth login \
  --base-url https://api.ileap.sine.dev \
  --client-id hello \
  --client-secret pathfinder

Logged in to https://api.ileap.sine.dev.
```

```bash
$ ileap footprint 91715e5e-fd0b-4d1c-8fab-76290c46e6ed
{
  "data": {
    "id": "91715e5e-fd0b-4d1c-8fab-76290c46e6ed",
    "specVersion": "2.0.0",
    "version": 1,
    "created": "2022-03-01T09:32:20Z",
    "status": "Active",
    "validityPeriodStart": "2022-03-01T09:32:20Z",
    "validityPeriodEnd": "2024-12-31T00:00:00Z",
    "companyName": "My Corp",
    "companyIds": [
      "urn:uuid:69585GB6-56T9-6958-E526-6FDGZJHU1326",
      "urn:epc:id:sgln:562958.00000.4"
    ],
    "productDescription": "Bio-Ethanol 98%, corn feedstock (bulk - no packaging)",
    "..."
}
```
