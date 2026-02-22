# iLEAP Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/way-platform/ileap-go)](https://pkg.go.dev/github.com/way-platform/ileap-go)
[![GoReportCard](https://goreportcard.com/badge/github.com/way-platform/ileap-go)](https://goreportcard.com/report/github.com/way-platform/ileap-go)
[![CI](https://github.com/way-platform/ileap-go/actions/workflows/release.yaml/badge.svg)](https://github.com/way-platform/ileap-go/actions/workflows/release.yaml)

A Go SDK and server implementation for logistics emissions data compatible with the [iLEAP Technical Specifications](https://sine-fdn.github.io/ileap-extension/).

## SDK

### Features

* Go Client and Server for [PACT Product Footprints](https://wbcsd.github.io/tr/2024/data-exchange-protocol-20241024/#dt-pf) with [iLEAP extensions](https://sine-fdn.github.io/ileap-extension/#pcf-mapping).
* Go Client and Server for [iLEAP Transport Activity Data](https://sine-fdn.github.io/ileap-extension/#dt-tad).

### Installing

```bash
$ go get github.com/way-platform/ileap-go@latest
```

### Using the Client

```go
client := ileap.NewClient(
    ileap.WithBaseURL(os.Getenv("BASE_URL")),
    ileap.WithOAuth2(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")),
)

footprint, err := client.GetFootprint(context.Background(), &ileap.GetFootprintRequest{
    ID: "91715e5e-fd0b-4d1c-8fab-76290c46e6ed",
})
if err != nil {
    // Handle error.
}
fmt.Println(footprint)
```

### Using the Server

```go
server := ileap.NewServer(
    ileap.WithFootprintHandler(myFootprintHandler),
    ileap.WithTADHandler(myTADHandler),
    ileap.WithTokenValidator(myTokenValidator),
)

log.Fatal(http.ListenAndServe(":8080", server))
```

#### Pre-built Handlers

The `handlers/` directory provides pre-built implementations that can be plugged directly into the server:

* **`ileapdemo`**: Provides demo implementations of `FootprintHandler`, `TADHandler`, `TokenIssuer`, `TokenValidator`, and `OIDCProvider` loaded with sample data and static credentials. Ideal for testing and local development.
* **`ileapclerk`**: Provides `TokenIssuer`, `TokenValidator`, and `OIDCProvider` implementations that delegate authentication to [Clerk](https://clerk.com/) via the Clerk Frontend API.

### Developing

#### Build project

The project is built using [Mage](https://magefile.org). See [magefile.go](./magefile.go).

```bash
$ ./tools/mage build
```

For all available build tasks, see:

```bash
$ ./tools/mage
```

## CLI tool

<img src="docs/cli.gif" />

### Installing

```bash
$ go install github.com/way-platform/ileap-go/cmd/ileap@latest
```

Prebuilt binaries for Linux, Windows, and Mac are available from the [Releases](https://github.com/way-platform/ileap-go/releases).

### Using

Start a local demo server in the background:

```bash
$ ileap demo-server --port 8080 &
INFO iLEAP demo server listening address=:8080
```

Log in to the local demo server:

```bash
$ ileap auth login \
  --base-url http://localhost:8080 \
  --client-id hello \
  --client-secret pathfinder

Logged in to http://localhost:8080.
```

Fetch a product footprint:

```bash
$ ileap footprint 91715e5e-fd0b-4d1c-8fab-76290c46e6ed
{
  "companyName": "My Corp",
  "created": "2022-03-01T09:32:20Z",
  "id": "91715e5e-fd0b-4d1c-8fab-76290c46e6ed",
  "organizationName": "My Corp",
  "pcf": {
    "biogenicCarbonContent": "0.41",
    "fossilGhgEmissions": "1.5",
    "pCfExcludingBiogenic": "1.63",
    "pCfIncludingBiogenic": "1.85",
    "unitaryProductAmount": "1"
  },
  "productDescription": "Bio-Ethanol 98%, corn feedstock (bulk - no packaging)",
  "specVersion": "2.0.0",
  "status": "Active"
}
```
