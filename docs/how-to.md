# How to serve iLEAP data with ileap-go

## Overview

This guide covers how to expose existing emission data — transport carbon footprints, activity data,
or both — via the iLEAP protocol using the ileap-go SDK.

## What the SDK provides

The SDK is an `http.Handler` you mount into your existing server or deploy standalone. It handles
the full PACT/iLEAP protocol: JSON envelopes, filter translation, `Link` header pagination, the OAuth2
token endpoint, OIDC discovery, and JWKS. You provide your data (via `ILeapServiceHandler`) and
your auth (via `AuthHandler`).

Filter translation is the key simplification:
- Incoming HTTP filters (including legacy OData on `/2/footprints`) are translated by the server.
- Handler methods receive simple request-local filters: `field_path` + `value`.
- The filter shape is aligned with iLEAP standalone field-path semantics.

```go
import ileap "github.com/way-platform/ileap-go"

srv := ileap.NewServer(
    ileap.WithServiceHandler(myHandler),
    ileap.WithAuthHandler(myAuth),
)
```

## Integration options

### Option A — inline integration

Implement `ILeapServiceHandler` against your existing data store and mount the iLEAP server into
your existing HTTP mux (`net/http`, Gin, Echo, Chi, etc.).

The interface has three methods:

```go
// ileapv1connect.ILeapServiceHandler
ListFootprints(context.Context, *ileapv1.ListFootprintsRequest) (*ileapv1.ListFootprintsResponse, error)
GetFootprint(context.Context, *ileapv1.GetFootprintRequest) (*ileapv1.GetFootprintResponse, error)
ListTransportActivityData(context.Context, *ileapv1.ListTransportActivityDataRequest) (*ileapv1.ListTransportActivityDataResponse, error)
```

A minimal database-backed implementation:

```go
import (
    "context"

    ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
)

type myHandler struct{ db *sql.DB }

func (h *myHandler) ListFootprints(
    ctx context.Context,
    req *ileapv1.ListFootprintsRequest,
) (*ileapv1.ListFootprintsResponse, error) {
    rows, err := h.db.QueryContext(ctx, "SELECT ... FROM footprints LIMIT ?", req.GetLimit())
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var footprints []*ileapv1.ProductFootprint
    for rows.Next() {
        fp := &ileapv1.ProductFootprint{}
        // scan into fp ...
        footprints = append(footprints, fp)
    }
    return &ileapv1.ListFootprintsResponse{Data: footprints}, nil
}

// GetFootprint and ListTransportActivityData follow the same pattern.
```

Mount into your existing mux:

```go
mux.Handle("/", srv)
http.ListenAndServe(":8080", mux)
```

### Option B — standalone gateway

Use this when your data lives in a separate service. Deploy the iLEAP server as a standalone
gateway: it speaks iLEAP HTTP externally and [Connect RPC](https://connectrpc.com) internally.

Implement a Connect RPC service (in Go or any language with a Connect implementation) that satisfies
`ILeapServiceHandler`. The `ileapconnect` package creates a client you can pass directly as the
handler:

```go
import (
    "log"
    "net/http"

    ileap "github.com/way-platform/ileap-go"
    "github.com/way-platform/ileap-go/handlers/ileapconnect"
)

srv := ileap.NewServer(
    ileap.WithServiceHandler(ileapconnect.NewClient("http://my-backend:9000")),
    ileap.WithAuthHandler(myAuth),
)
log.Fatal(http.ListenAndServe(":8080", srv))
```

The gateway forwards each authenticated iLEAP request to your internal Connect service, including
the bearer token.

## Authentication

`AuthHandler` covers token issuance, token validation, and OIDC discovery:

```go
// ileap.AuthHandler
type AuthHandler interface {
    IssueToken(ctx context.Context, clientID, clientSecret string) (*oauth2.Token, error)
    ValidateToken(ctx context.Context, token string) (*ileap.TokenInfo, error)
    OpenIDConfiguration(baseURL string) *ileap.OpenIDConfiguration
    JWKS() *ileap.JWKSet
}
```

### Clerk

The `ileapclerk` package implements `AuthHandler` using Clerk for token issuance, validation, and
JWKS:

```go
import "github.com/way-platform/ileap-go/handlers/ileapclerk"

clerkClient := ileapclerk.NewClient("your-instance.clerk.accounts.dev")
auth := ileapclerk.NewAuthHandler(clerkClient)
```

The SDK serves the OIDC discovery endpoint automatically.

### Custom implementation

Implement `AuthHandler` directly against your existing auth system. See
`handlers/ileapdemo/auth_provider.go` (`ileapdemo.AuthProvider`) for a self-contained reference
implementation (JWT signing with an embedded RSA keypair — useful as a starting point or in tests).

Contributions for other providers (Auth0, Keycloak, etc.) are welcome.

## Demo server

To see a working server before integrating your data:

```bash
go install github.com/way-platform/ileap-go/cmd/ileap@latest
ileap demo-server --port 8080
```

In another terminal:

```bash
ileap auth login --base-url http://localhost:8080 --client-id hello --client-secret pathfinder
ileap footprints
ileap tad
```

Or against the live demo:

```bash
ileap auth login --base-url https://demo.ileap.way.cloud --client-id hello --client-secret pathfinder
```

## Conformance testing

The SDK ships the PACT/iLEAP conformance suite as an importable Go test helper:

```go
import "github.com/way-platform/ileap-go/ileaptest"

func TestMyILeapServer(t *testing.T) {
    srv := startMyServer(t)
    ileaptest.RunConformanceTests(t, ileaptest.ConformanceTestConfig{
        ServerURL: srv.URL,
        Username:  "myuser",
        Password:  "mysecret",
    })
}
```

Or run against a deployed server:

```bash
ILEAP_SERVER_URL=https://your-server.example.com \
ILEAP_USERNAME=hello \
ILEAP_PASSWORD=pathfinder \
    go test -v ./ileaptest/...
```

## Further reading

- **Client usage** — consuming other iLEAP servers: see the README.
- **Data model** — ShipmentFootprint, TCE, TOC, HOC, TAD: see the
  [iLEAP Technical Specifications](https://sine-fdn.github.io/ileap-extension/) and the
  [PACT Data Exchange Protocol](https://wbcsd.github.io/tr/2023/data-exchange-protocol-20231207/).
- **API reference** — full Go API: [pkg.go.dev/github.com/way-platform/ileap-go](https://pkg.go.dev/github.com/way-platform/ileap-go).
