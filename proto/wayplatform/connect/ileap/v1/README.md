# iLEAP Protobuf Schemas

Protobuf schemas for the iLEAP Technical Specifications, extending the PACT
Data Exchange Protocol with logistics emissions data types conforming to
ISO 14083 and the GLEC Framework v3.1.

## Overview

> Transport Service Users, Transport Service Organizers, and Transport
> Operators have a growing need to calculate and report logistics emissions
> with greater accuracy and transparency. They want to transition from default
> and modeled data to primary data that accurately reflects the emissions of
> their operations.
>
> However, despite the methodological interoperability that has been achieved
> with ISO 14083:2023 and the GLEC Framework, data access processes remain
> predominantly manual and differ between companies.
>
> Because of this, accessing and exchanging logistics emissions data based on
> primary data is challenging and costly.
>
> The iLEAP Technical Specifications address these challenges by
>
> 1. "working backwards" from business cases
> 2. translating business cases into related data transactions enabling their
>    realization through independent software implementations, known as host
>    systems
> 3. enabling the interoperable flow of data between host systems by defining
>    a data exchange protocol and a data model that is rooted in the ISO 14083
>    and the GLEC Framework V3.1.
>
> -- iLEAP Technical Specifications, Introduction

### PACT Interoperability

Following the methodological alignment of the GLEC Framework with the PACT
Framework, the iLEAP Technical Specifications are designed to be interoperable
with the PACT Data Exchange Protocol and its PCF data model.

This integration reduces the friction between host systems used by Transport
Service Organizers or Transport Operators and those used by Transport Service
Users, enabling a holistic flow and approach to the management of carbon
emissions.

Implementing PACT is a prerequisite for implementing iLEAP, specifically the
authentication flows.

## Roles

| Role | Description |
|------|-------------|
| **Transport Service User** | Refers to the party that purchases and/or utilizes a transport service. It could be a shipper or a Transport Service Organizer. See ISO 14083, Section 3.1.33. |
| **Transport Service Organizer** | Refers to the party providing transport services, where some of the operations are subcontracted to a third party, usually a Transport Operator. See ISO 14083, Section 3.1.32. |
| **Transport Operator** | Refers to the party that carries out the transport service. See ISO 14083, Section 3.1.30. |

## Data Transactions

| DT | Data Type | Provider | Consumer | Endpoint |
|----|-----------|----------|----------|----------|
| DT#1 | ShipmentFootprint (TCEs) | Operator/Organizer | Service User | `GET /2/footprints` |
| DT#2 | TOC or HOC | Operator/Organizer | Organizer/User | `GET /2/footprints` |
| DT#3 | TAD | Operator | Organizer/User | `GET /2/ileap/tad` |

- **DT#1**: TCE-level emissions for a single shipment, wrapped in a ShipmentFootprint
  and embedded as a PACT DataModelExtension in a ProductFootprint.
- **DT#2**: Emission intensity data at TOC (transport) or HOC (hub) cluster level,
  embedded as a PACT DataModelExtension in a ProductFootprint.
- **DT#3**: Raw activity data (distance, mass, energy) for parties that cannot yet
  calculate emissions. Exchanged via a dedicated endpoint, NOT as PACT extensions.

## Data Model

The iLEAP Data Model is composed of five main data types:

| Type | File | Description |
|------|------|-------------|
| `ShipmentFootprint` | `shipment_footprint.proto` | Collection of TCEs for a single shipment |
| `TCE` | `tce.proto` | Transport Chain Element -- one leg in a shipment |
| `TOC` | `toc.proto` | Transport Operation Category -- emission intensities for a transport operation class |
| `HOC` | `hoc.proto` | Hub Operation Category -- emission intensities for a hub operation class |
| `TAD` | `tad.proto` | Transport Activity Data -- raw activity data without emissions |

Supporting types:

| Type | File | Description |
|------|------|-------------|
| `EnergyCarrier` | `energy_carrier.proto` | Energy carrier with emission factors and feedstocks |
| `Feedstock` | `feedstock.proto` | Feedstock of an energy carrier |
| `GLECDistance` | `glec_distance.proto` | Distance per GLEC Framework (actual, GCD, SFD) |
| `Location` | `location.proto` | Geographic location |

PACT base types:

| Type | File | Description |
|------|------|-------------|
| `ProductFootprint` | `product_footprint.proto` | PACT ProductFootprint container |
| `CarbonFootprint` | `carbon_footprint.proto` | PACT CarbonFootprint (PCF data) |
| `DataModelExtension` | `data_model_extension.proto` | PACT extension envelope for iLEAP data |

## Key Invariant: Decimal Type

All numeric values in the iLEAP data model are of type `Decimal` -- JSON
strings matching `^-?\d+(\.\d+)?$`. Never use JSON numbers for iLEAP numeric
fields.

> Note: The JSON String encoding is necessary to avoid floating point rounding
> errors.
>
> -- iLEAP Technical Specifications, section "Data Type Decimal"

In these protobuf schemas, all Decimal fields use `string` as the proto type.

## PACT Integration Rules

ShipmentFootprint, TOC, and HOC are stored as PACT DataModelExtensions; one
ProductFootprint MUST NOT contain more than one iLEAP data type.

### Common Rules (all iLEAP types)

- `productCategoryCpc`: MUST always be `"83117"`
- `packagingEmissionsIncluded`: MUST be `false`
- `biogenicCarbonContent`: SHOULD be `"0"`
- `extensions[].specVersion`: MUST be `"2.0.0"`

### Per-type Mapping

| PACT Field | ShipmentFootprint | TOC | HOC |
|------------|-------------------|-----|-----|
| `productIds` | `urn:...:shipment:{shipmentId}` | `urn:...:toc:{tocId}` | `urn:...:hoc:{hocId}` |
| `pcf.declaredUnit` | `"ton kilometer"` | `"ton kilometer"` | `"kilogram"` |
| `pcf.unitaryProductAmount` | `sum(tces[].transportActivity)` | `"1"` | `"1000"` |
| `pcf.pCfExcludingBiogenic` | `sum(tces[].co2eWTW)` | `co2eIntensityWTW` | `co2eIntensityWTW` |
| `extensions[].dataSchema` | `https://api.ileap.sine.dev/shipment-footprint.json` | `https://api.ileap.sine.dev/toc.json` | `https://api.ileap.sine.dev/hoc.json` |

### productIds URN Format

```
urn:pathfinder:product:customcode:vendor-assigned:{type}:{id}
urn:pathfinder:product:customcode:buyer-assigned:{type}:{id}
```

Where `{type}` is `shipment`, `toc`, or `hoc`.

## HTTP API

### Authentication

1. Discover token endpoint: `GET /.well-known/openid-configuration`
2. Obtain token: `POST /auth/token` with Basic Auth + `grant_type=client_credentials`
3. Use token: `Authorization: Bearer {token}` on all API calls

### PACT Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/2/footprints` | List ProductFootprints (with iLEAP extensions) |
| `GET` | `/2/footprints/{id}` | Get single ProductFootprint |
| `POST` | `/2/events` | Async event notifications |

Supports `$filter` query parameter (OData v4 subset) and `limit` pagination
with `Link` header.

### iLEAP Endpoint

```
GET /2/ileap/tad?[filter params]&limit={n}
```

- **Filtering**: query parameters as key-value pairs (e.g., `?mode=Road`)
- **Pagination**: `limit` param + `Link: <url>; rel="next"` header
- **Response**: `{ "data": [TAD, ...] }`

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `AccessDenied` | 403 | Invalid or missing token |
| `TokenExpired` | 401 | Expired access token |
| `BadRequest` | 400 | Malformed request |
| `NoSuchFootprint` | 404 | Footprint ID not found |
| `NotImplemented` | 400 | Unsupported filter or feature (note: HTTP 400, not 501) |

## Normative References

- **iLEAP Technical Specifications**: https://sine-fdn.github.io/ileap-extension/
- **PACT Data Exchange Protocol v2.1.0**: https://wbcsd.github.io/tr/2023/data-exchange-protocol-20231207/
- **PACT Data Model Extensions**: https://wbcsd.github.io/data-model-extensions/spec/
- **GLEC Framework v3.1**: https://www.smartfreightcentre.org/en/our-programs/global-logistics-emissions-council/glec-framework/
- **ISO 14083:2023**: Greenhouse gases -- Quantification and reporting of greenhouse gas emissions arising from transport chain operations
- **SINE Foundation Demo API**: https://api.ileap.sine.dev
