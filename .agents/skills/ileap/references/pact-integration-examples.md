# PACT Integration Examples

Annotated JSON examples showing how iLEAP data types are embedded as
PACT `DataModelExtension`s in `ProductFootprint` objects.

Distilled from [Appendix B](https://sine-fdn.github.io/ileap-extension/#appendix-b)
of the iLEAP Technical Specifications v1.0.0.

## Legend

- **PACT field (manual)** -- set by the data owner, not derivable from iLEAP data
- **PACT field (derived)** -- MUST be derived from the iLEAP extension per spec rules
- **iLEAP extension** -- the `extensions[].data` payload defined by the iLEAP data model

---

## 1. ShipmentFootprint as ProductFootprint Extension

From spec lines 2813-2889. Represents DT#1 (TCE Data Exchange).

```jsonc
{
  // --- PACT fields (manual) ---
  "id": "d9be4477-e351-45b3-acd9-e1da05e6f633",
  "specVersion": "2.0.0",
  "version": 0,
  "created": "2022-05-22T21:47:32Z",
  "status": "Active",
  "companyName": "Super Duper Transport Co.",
  "companyIds": [
    "urn:epc:id:sgln:4063973.00000.8"
  ],
  "productDescription": "Logistics emissions related to shipment with ID 1237890",
  "productNameCompany": "Shipment with ID 1237890",
  "comment": "",

  // --- PACT fields (derived from ShipmentFootprint) ---
  // productIds: URN encodes the shipmentId
  "productIds": [
    "urn:pathfinder:product:customcode:vendor-assigned:shipment:1237890"
  ],
  // productCategoryCpc: MUST always be "83117" for logistics
  "productCategoryCpc": "83117",

  "pcf": {
    // declaredUnit: MUST be "ton kilometer" for ShipmentFootprint
    "declaredUnit": "ton kilometer",
    // unitaryProductAmount: sum of TCE transportActivity values
    "unitaryProductAmount": "36.801",
    // pCfExcludingBiogenic: sum of TCE co2eWTW values
    "pCfExcludingBiogenic": "3.6801",
    // fossilGhgEmissions: same as pCfExcludingBiogenic
    "fossilGhgEmissions": "3.6801",
    // fossilCarbonContent: same as pCfExcludingBiogenic
    "fossilCarbonContent": "0",
    // biogenicCarbonContent: SHOULD be "0"
    "biogenicCarbonContent": "0",

    // --- PACT PCF fields (manual) ---
    "characterizationFactors": "AR6",
    "ipccCharacterizationFactorsSources": ["AR6"],
    "crossSectoralStandardsUsed": [
      "GHG Protocol Product standard"
    ],
    "productOrSectorSpecificRules": [],
    "boundaryProcessesDescription": "SFC GLEC Framework-conforming (W2W CO2e emissions)",
    "referencePeriodStart": "2021-01-01T00:00:00Z",
    "referencePeriodEnd": "2022-01-01T00:00:00Z",
    "secondaryEmissionFactorSources": [
      { "name": "Ecoinvent", "version": "3.9.1" }
    ],
    "exemptedEmissionsPercent": 0,
    "exemptedEmissionsDescription": "",
    // packagingEmissionsIncluded: MUST be false for ShipmentFootprint
    "packagingEmissionsIncluded": true,
    "primaryDataShare": 56.12
  },

  // --- iLEAP extension ---
  "extensions": [
    {
      "specVersion": "2.0.0",
      // dataSchema: MUST be this URL for ShipmentFootprint
      "dataSchema": "https://api.ileap.sine.dev/shipment-footprint.json",
      "documentation": "https://sine-fdn.github.io/ileap-extension/",
      "data": {
        "mass": "87",
        "shipmentId": "1237890",
        "tces": [
          {
            "tceId": "abcdef",
            "prevTceIds": [],
            "tocId": "truck-40t-euro5-de",
            "shipmentId": "1237890",
            "mass": "87",
            "distance": {
              "actual": "423"
            },
            "transportActivity": "3.6801",
            "co2eWTW": "3.6801",
            "co2eTTW": "3.2801"
          }
        ]
      }
    }
  ]
}
```

### Derivation rules for ShipmentFootprint

| PACT property | Derivation |
|---|---|
| `productIds` | `urn:pathfinder:product:customcode:vendor-assigned:shipment:{shipmentId}` |
| `productCategoryCpc` | Always `"83117"` |
| `pcf.declaredUnit` | Always `"ton kilometer"` |
| `pcf.unitaryProductAmount` | `sum(tces[].transportActivity)` |
| `pcf.pCfExcludingBiogenic` | `sum(tces[].co2eWTW)` |
| `pcf.fossilGhgEmissions` | Same as `pCfExcludingBiogenic` |
| `pcf.packagingEmissionsIncluded` | Always `false` |
| `extensions[].dataSchema` | `"https://api.ileap.sine.dev/shipment-footprint.json"` |

---

## 2. TOC as ProductFootprint Extension

From spec lines 2895-2969. Represents DT#2 (TOC Data Exchange).

```jsonc
{
  // --- PACT fields (manual) ---
  "id": "f3c04ec8-b33a-43b1-9fa7-d6a448fd60af",
  "specVersion": "2.0.0",
  "version": 0,
  "created": "2022-05-22T21:47:32Z",
  "status": "Active",
  "companyName": "Super Duper Transport Co.",
  "companyIds": [
    "urn:epc:id:sgln:4063973.00000.8"
  ],
  "productDescription": "Logistics emissions related to TOC with ID 4561230",
  "productNameCompany": "TOC with ID 4561230",
  "comment": "",

  // --- PACT fields (derived from TOC) ---
  // productIds: URN encodes the tocId
  "productIds": [
    "urn:pathfinder:product:customcode:vendor-assigned:toc:4561230"
  ],
  // productCategoryCpc: MUST always be "83117"
  "productCategoryCpc": "83117",

  "pcf": {
    // declaredUnit: MUST be "ton kilometer" for TOC
    "declaredUnit": "ton kilometer",
    // unitaryProductAmount: SHOULD be "1" (emissions per tkm)
    "unitaryProductAmount": "1",
    // pCfExcludingBiogenic: equals TOC co2eIntensityWTW
    "pCfExcludingBiogenic": "3.6801",
    // fossilGhgEmissions: same as pCfExcludingBiogenic
    "fossilGhgEmissions": "3.6801",
    "fossilCarbonContent": "0",
    "biogenicCarbonContent": "0",

    // --- PACT PCF fields (manual) ---
    "characterizationFactors": "AR6",
    "ipccCharacterizationFactorsSources": ["AR6"],
    "crossSectoralStandardsUsed": [
      "GHG Protocol Product standard"
    ],
    "productOrSectorSpecificRules": [],
    "boundaryProcessesDescription": "SFC GLEC Framework-conforming (W2W CO2e emissions)",
    "referencePeriodStart": "2021-01-01T00:00:00Z",
    "referencePeriodEnd": "2022-01-01T00:00:00Z",
    "secondaryEmissionFactorSources": [
      { "name": "Ecoinvent", "version": "3.9.1" }
    ],
    "exemptedEmissionsPercent": 0,
    "exemptedEmissionsDescription": "",
    // packagingEmissionsIncluded: MUST be false
    "packagingEmissionsIncluded": false,
    "primaryDataShare": 56.12
  },

  // --- iLEAP extension ---
  "extensions": [
    {
      "specVersion": "2.0.0",
      // dataSchema: MUST be this URL for TOC
      "dataSchema": "https://api.ileap.sine.dev/toc.json",
      "documentation": "https://sine-fdn.github.io/ileap-extension/",
      "data": {
        "tocId": "4561230",
        "certifications": ["ISO14083:2023"],
        "mode": "Road",
        "temperatureControl": "refrigerated",
        "truckLoadingSequence": "FTL",
        "energyCarriers": [
          {
            "energyCarrier": "Diesel",
            "relativeShare": "1",
            "emissionFactorWTW": "3.6801",
            "emissionFactorTTW": "3.2801"
          }
        ],
        "co2eIntensityWTW": "3.6801",
        "co2eIntensityTTW": "3.2801",
        "transportActivityUnit": "tkm"
      }
    }
  ]
}
```

### Derivation rules for TOC

| PACT property | Derivation |
|---|---|
| `productIds` | `urn:pathfinder:product:customcode:vendor-assigned:toc:{tocId}` |
| `productCategoryCpc` | Always `"83117"` |
| `pcf.declaredUnit` | Always `"ton kilometer"` |
| `pcf.unitaryProductAmount` | `"1"` (per tkm) |
| `pcf.pCfExcludingBiogenic` | `TOC.co2eIntensityWTW` |
| `pcf.fossilGhgEmissions` | Same as `pCfExcludingBiogenic` |
| `pcf.packagingEmissionsIncluded` | Always `false` |
| `extensions[].dataSchema` | `"https://api.ileap.sine.dev/toc.json"` |

---

## 3. HOC as ProductFootprint Extension

From spec lines 2975-3049. Represents DT#2 (HOC Data Exchange).

```jsonc
{
  // --- PACT fields (manual) ---
  "id": "f3c04ec8-b33a-43b1-9fa7-d6a448fd60af",
  "specVersion": "2.0.0",
  "version": 0,
  "created": "2022-05-22T21:47:32Z",
  "status": "Active",
  "companyName": "Super Duper Transport Co.",
  "companyIds": [
    "urn:epc:id:sgln:4063973.00000.8"
  ],
  "productDescription": "Logistics emissions related to HOC with ID 7890123",
  "productNameCompany": "HOC with ID 7890123",
  "comment": "",

  // --- PACT fields (derived from HOC) ---
  // productIds: URN encodes the hocId
  "productIds": [
    "urn:pathfinder:product:customcode:vendor-assigned:hoc:7890123"
  ],
  // productCategoryCpc: MUST always be "83117"
  "productCategoryCpc": "83117",

  "pcf": {
    // declaredUnit: MUST be "kilogram" for HOC
    "declaredUnit": "kilogram",
    // unitaryProductAmount: MUST be "1000" (per tonne leaving the hub)
    "unitaryProductAmount": "1000",
    // pCfExcludingBiogenic: equals HOC co2eIntensityWTW
    "pCfExcludingBiogenic": "3.6801",
    // fossilGhgEmissions: same as pCfExcludingBiogenic
    "fossilGhgEmissions": "3.6801",
    "fossilCarbonContent": "0",
    "biogenicCarbonContent": "0",

    // --- PACT PCF fields (manual) ---
    "characterizationFactors": "AR6",
    "ipccCharacterizationFactorsSources": ["AR6"],
    "crossSectoralStandardsUsed": [
      "GHG Protocol Product standard"
    ],
    "productOrSectorSpecificRules": [],
    "boundaryProcessesDescription": "SFC GLEC Framework-conforming (W2W CO2e emissions)",
    "referencePeriodStart": "2021-01-01T00:00:00Z",
    "referencePeriodEnd": "2022-01-01T00:00:00Z",
    "secondaryEmissionFactorSources": [
      { "name": "Ecoinvent", "version": "3.9.1" }
    ],
    "exemptedEmissionsPercent": 0,
    "exemptedEmissionsDescription": "",
    // packagingEmissionsIncluded: MUST be false
    "packagingEmissionsIncluded": false,
    "primaryDataShare": 56.12
  },

  // --- iLEAP extension ---
  "extensions": [
    {
      "specVersion": "2.0.0",
      // dataSchema: MUST be this URL for HOC
      "dataSchema": "https://api.ileap.sine.dev/hoc.json",
      "documentation": "https://sine-fdn.github.io/ileap-extension/",
      "data": {
        "hocId": "7890123",
        "certifications": ["ISO14083:2023"],
        "hubType": "Warehouse",
        "temperatureControl": "refrigerated",
        "energyCarriers": [
          {
            "energyCarrier": "Diesel",
            "relativeShare": "1",
            "emissionFactorWTW": "3.6801",
            "emissionFactorTTW": "3.2801"
          }
        ],
        "co2eIntensityWTW": "3.6801",
        "co2eIntensityTTW": "3.2801",
        "hubActivityUnit": "tonnes"
      }
    }
  ]
}
```

### Derivation rules for HOC

| PACT property | Derivation |
|---|---|
| `productIds` | `urn:pathfinder:product:customcode:vendor-assigned:hoc:{hocId}` |
| `productCategoryCpc` | Always `"83117"` |
| `pcf.declaredUnit` | Always `"kilogram"` |
| `pcf.unitaryProductAmount` | `"1000"` (per tonne) |
| `pcf.pCfExcludingBiogenic` | `HOC.co2eIntensityWTW` |
| `pcf.fossilGhgEmissions` | Same as `pCfExcludingBiogenic` |
| `pcf.packagingEmissionsIncluded` | Always `false` |
| `extensions[].dataSchema` | `"https://api.ileap.sine.dev/hoc.json"` |

---

## 4. TAD Response Body

From the `/2/ileap/tad` endpoint. Represents DT#3 (Transport Activity Data Exchange).
TAD is **not** embedded as a PACT extension -- it uses a dedicated endpoint.

```jsonc
{
  "data": [
    {
      "activityId": "1",
      "consignmentIds": ["consignment-1"],
      "distance": {
        "actual": "423"
      },
      "mass": "2300",
      "origin": {
        "city": "Hamburg",
        "country": "DE"
      },
      "destination": {
        "city": "Berlin",
        "country": "DE"
      },
      "departureAt": "2024-01-15T08:00:00Z",
      "arrivalAt": "2024-01-15T14:00:00Z",
      "mode": "Road",
      "energyCarriers": [
        {
          "energyCarrier": "Diesel",
          "relativeShare": "1",
          "feedstocks": [
            {
              "feedstock": "Fossil",
              "feedstockShare": "1"
            }
          ]
        }
      ]
    }
  ]
}
```

### Key differences from DT#1/DT#2

- TAD uses endpoint `GET /2/ileap/tad` (not `/2/footprints`)
- Response wraps data in `{ "data": [...] }`
- Supports query-parameter filtering (e.g., `?mode=Road`)
- Supports `limit` and pagination via `Link` header
- Contains no emissions data -- only activity data for downstream calculation
