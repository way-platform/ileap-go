# iLEAP Server Implementation Patterns

Common patterns for implementing an iLEAP-conformant server,
independent of language or framework.

## Authentication Pattern

### Key Generation

Generate an RSA key pair for JWT signing:

```sh
openssl genpkey -algorithm RSA -out keypair.pem -pkeyopt rsa_keygen_bits:2048
```

### Server-Side Auth Flow

1. **OpenID Discovery**: Serve `GET /.well-known/openid-configuration`
   returning `issuer`, `token_endpoint`, `jwks_uri`
2. **Token Endpoint**: Accept `POST /auth/token` with Basic Auth
   (`base64(client_id:client_secret)`) and
   `grant_type=client_credentials`. Issue a JWT signed with RS256.
3. **JWKS Endpoint**: Serve `GET /jwks` exposing the public key for
   token verification by peers.

### Token Validation

On each API request:
- Extract Bearer token from `Authorization` header
- Verify JWT signature against the signing key
- Check expiry (if set)
- Extract client identity from claims for access control

## Role-Based Data Filtering

iLEAP defines three roles with different data access needs:

| Role | Sees |
|---|---|
| Transport Service User (shipper) | ShipmentFootprints + TAD |
| Transport Service Organizer (LSP) | TOCs + HOCs + TAD |
| Global / admin | All footprints + TAD |

### Filtering Logic

- Determine role from authenticated client identity
- Filter `ProductFootprint` list by `extensions[].dataSchema` URL:
  - ShipmentFootprint schema URL -> visible to shippers
  - TOC/HOC schema URLs -> visible to organizers
- TAD endpoint: all authenticated users see TAD records

## TAD Endpoint Pattern

```
GET /2/ileap/tad?[mode=Road][&limit=10]
```

### Query Parameter Filtering

- Accept key-value query parameters for field-level filtering
- Match against TAD fields (e.g., `mode=Road` matches
  `tad.mode == "Road"`)
- Unknown parameters: ignore or return `NotImplemented`

### Pagination

- `limit` parameter controls max results per page
- If more results exist, return `Link` header:

```
Link: <https://example.com/2/ileap/tad?limit=10&offset=10>; rel="next"
```

- Use absolute URLs matching the request host

### Response Format

```json
{
  "data": [
    {
      "activityId": "tad-001",
      "consignmentIds": ["C-123"],
      "distance": { "actual": "450" },
      "origin": { "city": "Berlin" },
      "destination": { "city": "Munich" },
      "departureAt": "2024-01-15T08:00:00Z",
      "arrivalAt": "2024-01-15T14:00:00Z",
      "mode": "Road"
    }
  ]
}
```

## Error Response Pattern

All error responses use the same JSON shape:

```json
{
  "code": "AccessDenied",
  "message": "Human-readable description"
}
```

Map error codes to HTTP status:

| Code | HTTP Status |
|---|---|
| `AccessDenied` | 403 |
| `TokenExpired` | 401 |
| `BadRequest` | 400 |
| `NoSuchFootprint` | 404 |
| `NotImplemented` | 400 |
| `InternalError` | 500 |

## JSON Serialization: Omit Null Optional Fields

ACT and the SINE reference API expect optional fields to be omitted
from JSON, not serialized as `null`. When using code generators,
post-process generated types to ensure optional/nullable fields are
omitted when empty. Verify by comparing output field-for-field against
the SINE demo API (`api.ileap.sine.dev`).
