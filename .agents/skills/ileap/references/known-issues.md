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

## Known Permanent ACT Failures (Not Implementation Bugs)

Treat these as known failures when evaluating ACT results, as they often fail even on the SINE demo API at `api.ileap.sine.dev`:

| TC | Reason |
|---|---|
| PACT TC8 | Requires short-lived / expired tokens; not practical to provision |
| PACT TC18 | OIDC discovery flow; fails on the SINE reference API too |
| PACT TC19 | Structurally depends on TC18 |

## `biogenicAccountingMethodology` Valid Values

Valid values: `GHGP`, `ISO`, `PEF`, `Quantis`.

Note: the correct value is `GHGP` (not `GHPG`) — common typo.

## `hubType` Enum Is PascalCase

The `hubType` field on HOC must use PascalCase values (e.g.,
`"Transshipment"`, `"Warehouse"`, `"StorageAndTransshipment"`).
Using lowercase values (like `"transshipment"`) will cause ACT TC 003
to fail with an "unknown variant" error.

## OData `$filter` with Timezone Offsets (TC 20)

ACT sends timestamps with `+00:00` timezone offsets (not just `Z`) in OData filter values (e.g., `updated+lt+'2023-06-27T13:00:00.000+00:00'`). Standard URL query decoders (like Go's `url.ParseQuery`) will decode the `+` character as a space, resulting in `updated lt '2023-06-27T13:00:00.000 00:00'`. 

If your filter parser splits by whitespace (e.g., using `strings.Fields`), it will incorrectly split the timestamp into multiple tokens and fail to parse the filter.
**Fix**: Use a limited split (like `strings.SplitN(data, " ", 3)` in Go) to keep the RHS intact regardless of internal spaces.

## CloudEvents Validation (TC 21)

The `POST /2/events` handler must strictly validate the incoming event:
1. **`data` field**: Must be present **and not null**. ACT explicitly tests that an event with `"data": null` or an absent `data` field returns `400 BadRequest`.
2. **`specversion` field**: Must be exactly `"1.0"`. ACT sends an event with `specversion: "0.3"` to test rejection. Checking only for an empty string will fail this test; you must enforce `specversion == "1.0"`.
