# PACT v2.1.0 Implementation Guide

Actionable implementation steps distilled from the PACT Technical
Specifications for PCF Data Exchange (Version 2.1.0).

Source: `pact-spec-v2/index.bs`

## Version Alignment

- iLEAP requires **PACT v2.1.0**
- `ProductFootprint.specVersion` MUST be `"2.1.0"`
- `DataModelExtension.specVersion` MUST be `"2.0.0"` (extension spec
  version, distinct from the PACT spec version)
- The PACT spec is sometimes called "v2" informally -- this refers to
  the v2.x line (v2.0.0, v2.0.1, v2.1.0). There is **no v3** in
  production use.

## Authentication Flow

### Overview

1. **Discover** token endpoint via OpenID Connect
2. **Authenticate** with OAuth2 Client Credentials
3. **Use** Bearer token on all subsequent API calls

### Step 1: OpenID Discovery

```
GET /.well-known/openid-configuration HTTP/1.1
Host: {auth-hostname}
```

Response (minimum required fields):

```json
{
  "issuer": "https://example.com/",
  "token_endpoint": "https://example.com/auth/token",
  "jwks_uri": "https://example.com/jwks"
}
```

If OpenID discovery is unavailable, fall back to
`{auth-hostname}/{auth-subpath}/auth/token`.

### Step 2: Obtain Token

```
POST /auth/token HTTP/1.1
Host: {auth-hostname}
Content-Type: application/x-www-form-urlencoded
Authorization: Basic {base64(client_id:client_secret)}

grant_type=client_credentials
```

Success response:

```json
{
  "access_token": "<JWT>",
  "token_type": "bearer"
}
```

Error response (invalid credentials):

```json
{
  "error": "invalid_client",
  "error_description": "Authentication failed"
}
```

### Step 3: Use Bearer Token

All subsequent API calls include:

```
Authorization: Bearer {access_token}
```

Tokens SHOULD expire. Clients MUST re-authenticate when expired.

## Action ListFootprints

```
GET {subpath}/2/footprints?[$filter={expr}][&limit={n}] HTTP/1.1
Host: {hostname}
Authorization: Bearer {token}
```

### Filtering (optional)

- Parameter name: `$filter`
- Syntax: OData v4 subset
- Supported operators: `eq`, `lt`, `le`, `gt`, `ge`, `and`
- Lambda operator: `any` on `companyIds`, `productIds`
- Supported properties: `created`, `updated`, `productCategoryCpc`,
  `pcf/geographyCountry`, `pcf/referencePeriodStart`,
  `pcf/referencePeriodEnd`

Examples:

```
$filter=productCategoryCpc eq '83117'
$filter=pcf/geographyCountry eq 'DE'
$filter=productIds/any(pid:(pid eq 'urn:...'))
```

Host systems MAY ignore filters or return `NotImplemented`.

### Pagination (mandatory)

- Host systems MUST implement server-side pagination
- `limit` parameter controls page size
- If more results exist, response includes a `Link` header:

```
Link: <https://example.com/2/footprints?limit=10&offset=10>; rel="next"
```

- Pagination links MUST be absolute URLs
- Links MUST be valid for at least 180 seconds
- Clients SHOULD NOT assume previous pagination links remain valid

### Response

- Status 200: complete result
- Status 202: incomplete result (client may retry with backoff)
- Body: `{ "data": [ProductFootprint, ...] }`
- Returns latest version of each footprint (max `version` value)

## Action GetFootprint

```
GET {subpath}/2/footprints/{id} HTTP/1.1
Host: {hostname}
Authorization: Bearer {token}
```

### Response

- Status 200 + body: `{ "data": ProductFootprint }`
- Returns latest version if footprint was updated
- Error: `NoSuchFootprint` (404) if ID doesn't exist

## Action Events

```
POST {subpath}/2/events HTTP/1.1
Host: {hostname}
Authorization: Bearer {token}
Content-Type: application/cloudevents+json; charset=UTF-8
```

Support is OPTIONAL. Non-implementing hosts SHOULD return
`NotImplemented`.

### Event Types

| Event | Type String | Direction |
|---|---|---|
| PF Update | `org.wbcsd.pathfinder.ProductFootprint.Published.v1` | Owner -> Recipient |
| PF Request | `org.wbcsd.pathfinder.ProductFootprintRequest.Created.v1` | Recipient -> Owner |
| PF Response | `org.wbcsd.pathfinder.ProductFootprintRequest.Fulfilled.v1` | Owner -> Recipient |
| PF Response Error | `org.wbcsd.pathfinder.ProductFootprintRequest.Rejected.v1` | Owner -> Recipient |

### PF Update Event

```json
{
  "type": "org.wbcsd.pathfinder.ProductFootprint.Published.v1",
  "specversion": "1.0",
  "id": "{event-id}",
  "source": "//{hostname}/{subpath}",
  "time": "2024-01-15T10:00:00Z",
  "data": {
    "pfIds": ["uuid-1", "uuid-2"]
  }
}
```

### PF Request Event

```json
{
  "type": "org.wbcsd.pathfinder.ProductFootprintRequest.Created.v1",
  "specversion": "1.0",
  "id": "{event-id}",
  "source": "//{hostname}/{subpath}",
  "data": {
    "pf": { "productIds": ["urn:..."] },
    "comment": "optional request comment"
  }
}
```

## Error Codes

| Code | HTTP Status | Description |
|---|---|---|
| `AccessDenied` | 403 | Invalid or missing token |
| `BadRequest` | 400 | Malformed request |
| `NoSuchFootprint` | 404 | Footprint ID not found |
| `NotImplemented` | 400 | Unsupported filter or feature |
| `TokenExpired` | 401 | Expired access token |
| `InternalError` | 500 | Unexpected server error |

Error response body format:

```json
{
  "code": "AccessDenied",
  "message": "Access denied"
}
```

Note: `NotImplemented` uses HTTP 400, not 501. This is per the PACT
spec.

## ProductFootprint Data Model

Key properties (see spec for full list):

| Property | Type | Req | Notes |
|---|---|---|---|
| `id` | UUID v4 | M | Unique footprint identifier |
| `specVersion` | String | M | Must be `"2.1.0"` |
| `version` | Integer | M | Incremented on updates |
| `created` | DateTime | M | UTC timestamp |
| `updated` | DateTime | O | Set only after updates |
| `status` | String | M | `"Active"` or `"Deprecated"` |
| `companyName` | String | M | Data owner company name |
| `companyIds` | URN[] | M | Non-empty set of company URNs |
| `productDescription` | String | M | Free-form product description |
| `productIds` | URN[] | M | Non-empty set of product URNs |
| `productCategoryCpc` | String | M | UN CPC code (iLEAP: `"83117"`) |
| `productNameCompany` | String | M | Trade name |
| `comment` | String | M | Calculation notes, audit info |
| `pcf` | CarbonFootprint | M | The carbon footprint data |
| `extensions` | DataModelExtension[] | O | iLEAP data goes here |

## CarbonFootprint Data Model

Key properties:

| Property | Type | Req | Notes |
|---|---|---|---|
| `declaredUnit` | DeclaredUnit | M | iLEAP: `"ton kilometer"` or `"kilogram"` |
| `unitaryProductAmount` | Decimal | M | Strictly > 0 |
| `pCfExcludingBiogenic` | Decimal | M | kgCO2e per declared unit |
| `fossilGhgEmissions` | Decimal | M | kgCO2e per declared unit |
| `fossilCarbonContent` | Decimal | M | kgC per declared unit |
| `biogenicCarbonContent` | Decimal | M | kgC per declared unit |
| `characterizationFactors` | String | M | `"AR5"` or `"AR6"` |
| `crossSectoralStandardsUsed` | String[] | M | e.g. `["ISO Standard 14083"]` |
| `boundaryProcessesDescription` | String | M | Free-form |
| `referencePeriodStart` | DateTime | M | Start of validity |
| `referencePeriodEnd` | DateTime | M | End of validity |
| `exemptedEmissionsPercent` | Number | M | 0-5 |
| `exemptedEmissionsDescription` | String | M | Can be empty string |
| `packagingEmissionsIncluded` | Boolean | M | iLEAP: `false` |

**DeclaredUnit** values: `liter`, `kilogram`, `cubic meter`,
`kilowatt hour`, `megajoule`, `ton kilometer`, `square meter`

## DataModelExtension

```json
{
  "specVersion": "2.0.0",
  "dataSchema": "https://api.ileap.sine.dev/shipment-footprint.json",
  "data": { ... }
}
```

- `specVersion`: always `"2.0.0"` (extension spec version)
- `dataSchema`: URL pointing to the JSON schema
- `data`: the extension payload (iLEAP data type)

## Navigating the Raw Spec

Source: `pact-spec-v2/index.bs` (2520 lines, Bikeshed format)

### Section Index

| Section | Topic | Lines |
|---|---|---|
| `#intro` | Introduction, scope, audience | 16-70 |
| `#terminology` | Term definitions | 73-121 |
| `#conformance` | Conformance requirements | 125-134 |
| `#data-model` | Data model overview | 138-161 |
| `#dt-pf` | ProductFootprint properties | 163-334 |
| `#dt-carbonfootprint` | CarbonFootprint properties | 337-646 |
| `#dt-dataqualityindicators` | DataQualityIndicators | 649-741 |
| `#dt-assurance` | Assurance | 744-824 |
| `#dt-datamodelextension` | DataModelExtension | 827-852 |
| `#dt-declaredunit` | DeclaredUnit enum | 1089-1121 |
| `#dt-decimal` | Decimal type (JSON String) | 1349-1362 |
| `#lifecycle` | ProductFootprint lifecycle | 1383-1497 |
| `#api` | HTTP API overview | 1499-1504 |
| `#api-auth` | Authentication flow | 1539-1564 |
| `#api-action-auth` | Action Authenticate | 1591-1703 |
| `#api-action-list` | Action ListFootprints | 1706-1861 |
| `#api-action-get` | Action GetFootprint | 1863-1922 |
| `#api-action-events` | Action Events | 1925-2134 |
| `#api-error-responses` | Error codes table | 2135-2277 |
| `#changelog` | Version history | 2288-2463 |

### Grep Patterns

```sh
# Find property definitions
grep -n '<dfn>' pact-spec-v2/index.bs | head -40

# Find all mandatory properties
grep -n '<td>M$' pact-spec-v2/index.bs

# Find MUST requirements
grep -n 'MUST' pact-spec-v2/index.bs | head -30

# Find error codes
grep -n 'dfn>.*Denied\|dfn>.*Request\|dfn>.*Expired' pact-spec-v2/index.bs

# Find enum values
grep -n 'dfn>' pact-spec-v2/index.bs | grep -i 'liter\|kilogram\|ton\|meter'
```
