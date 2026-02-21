# Known Issues & Implementation Gotchas

Practical issues discovered while implementing iLEAP servers and
running ACT conformance tests.

## `organizationName` Field

PACT v2.1.0 defines `organizationName` as an optional alias for
`companyName` on ProductFootprint. Include `organizationName` in
ProductFootprint responses with the same value as `companyName`.

## Spec Version Confusion

- The current PACT spec is **v2.1.0** (released Dec 2023)
- `ProductFootprint.specVersion` must be `"2.1.0"`
- `DataModelExtension.specVersion` must be `"2.0.0"` (extension spec)
- There is **no PACT v3** in production. References to "v3" are
  typically future drafts or the iLEAP spec itself.
- The iLEAP spec references "PACT v2.1+" as its baseline.

## Error Code Split: PACT vs iLEAP

PACT and iLEAP endpoints use the same error code set but differ in
expected HTTP status for auth errors:

| Scenario | PACT spec says | Observed in ACT |
|---|---|---|
| Missing/invalid token | 403 `AccessDenied` | 403 `AccessDenied` |
| Expired token | 401 `TokenExpired` | 401 `TokenExpired` |
| Auth endpoint failure | 401 (RFC 6749) | 401 |

The iLEAP TAD endpoint (`/2/ileap/tad`) uses the same error codes as
the PACT endpoints. ACT tests both.

Note: PACT's `NotImplemented` maps to HTTP **400**, not 501. This is
intentional per spec -- it means the host doesn't support the requested
filter, not that the endpoint itself is unimplemented.

## TAD Field Evolution

The TAD data model has evolved across iLEAP spec revisions:

| Field | Status | Notes |
|---|---|---|
| `arrivalAt` | Now **required** | Was optional in earlier drafts |
| `departureAt` | Now **required** | Was optional in earlier drafts |
| `mode` | Now **required** | Was optional in earlier drafts |
| `feedstocks` | **Renamed** | Now `energyCarriers` (on TOC/HOC) |

If porting from an earlier implementation, ensure:
- TAD records include `arrivalAt`, `departureAt`, and `mode`
- TOC/HOC use `energyCarriers` not `feedstocks`

## OpenAPI Overlay Needs

The PACT OpenAPI spec (`pact-spec-v2/openapi.yaml`) and the generated
iLEAP OpenAPI spec need overlays for code generation:

- `CarbonFootprint` uses `oneOf` for fields that are conditionally
  required (`O*` fields). Most Go/TypeScript generators struggle with
  this -- flatten to optional fields.
- `ProductFootprint.id` is typed as `string` but should have
  `format: uuid` for proper type generation.
- iLEAP `Decimal` fields are JSON strings, not numbers. Generators may
  need custom type mappings.

## Pagination Link Header

The `Link` header for pagination must use **absolute URLs** with the
same host as the original request. Relative URLs will fail ACT tests.

Format: `Link: <https://example.com/2/footprints?limit=10&offset=10>; rel="next"`

## Optional Fields Must Be Omitted, Not Null

ACT and the SINE reference API expect optional fields to be absent from
the JSON response, not serialized as `null`. Ensure code generators and
serialization libraries omit nil/null optional fields.

## `secondaryEmissionFactorSources` Must Be Undefined

Per PACT spec, if no secondary data is used, `secondaryEmissionFactorSources`
MUST be absent from the JSON — not `null`, not an empty array.

## ACT Requires Publicly Reachable Server

ACT delegates PACT tests to an external service that must reach the API
over the internet. Localhost URLs cause PACT test failures. iLEAP-specific
tests (TC001-TC007) work locally.

## PACT TC8 and TC18: Known Reference Failures

PACT TC8 and TC18 fail even on the SINE demo API at
`api.ileap.sine.dev`. Treat these as known failures when evaluating
ACT results.

## `biogenicAccountingMethodology` Valid Values

Valid values: `GHGP`, `ISO`, `PEF`, `Quantis`.

Note: the correct value is `GHGP` (not `GHPG`) — common typo.
