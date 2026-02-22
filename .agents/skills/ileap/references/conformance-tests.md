# iLEAP Conformance Testing

The iLEAP Technical Specifications are designed to be incrementally adopted and realized in host systems for different stakeholder groups. Conformance testing ensures interoperability and syntactic correctness across these systems.

## Kinds of iLEAP Conformance

There are two kinds of iLEAP conformance defined:

1. **iLEAP Emissions Data Conformance**
2. **iLEAP Activity Data Conformance**

For each kind of iLEAP Conformance, a host system achieves conformity by:
1. Realizing the normative statements and relevant test cases defined below.
2. Successfully passing Bilateral Testing (an interoperability test where two different host system implementations verify their ability to work together).
3. Successfully passing the [Automated Conformance Testing](https://act.sine.dev).

To achieve iLEAP Emissions Data Conformance, host systems are further **REQUIRED** to support the following features:
- Fetching Transport Activity Data (TAD) from another host system.
- Using it to generate Transport Operation Category (TOC), Hub Operation Category (HOC), or ShipmentFootprint data.

It is **RECOMMENDED** for host systems to conform to both kinds of iLEAP conformance.

### Requirements for Formal iLEAP / SFC Certification

When applying for formal certification (e.g., the SFC Certification Scheme extension for iLEAP Tool Providers), the testing requirements strictly map to the data transactions (DT):

*   **MANDATORY**: To achieve certification, implementations *must* pass tests for **DT#1** and **DT#2** (ShipmentFootprint, TOC, HOC), which corresponds to the **iLEAP Emissions Data Conformance** track (Test Cases **TC001-TC003**), alongside the baseline **PACT Required Test Cases** (specifically Authentication, Action ListFootprints, and Action GetFootprint tests).
*   **OPTIONAL (but recommended)**: Passing tests for **DT#3** (TAD), which corresponds to the **iLEAP Activity Data Conformance** track (Test Cases **TC004-TC008**), is currently optional for certification but recommended.

*Note: During the pilot phase, if an implementation fails the mandatory criteria, providers are given a two-week timeframe for remediation, allowing for a maximum of three certification test runs.*

---

## 1. PACT Conformance Tests

Since the iLEAP Technical Specifications were conceived as an extension to the PACT Data Model and Data Exchange Protocol, any iLEAP conformant implementation must also implement the PACT Data Exchange Protocol ([PACT Technical Specifications V2.1](https://wbcsd.github.io/tr/2023/data-exchange-protocol-20231207/)) or above.

For that reason, the required tests in the [PACT Conformance Testing Checklist](https://wbcsd.github.io/pact-conformance-testing/checklist.html#required-tests) MUST also be performed for both iLEAP Emissions Data Conformance and iLEAP Activity Data Conformance.

---

## 2. iLEAP Emissions Data Conformance

This conformance level includes all [PACT Required Test Cases](https://wbcsd.github.io/pact-conformance-testing/checklist.html#required-tests) and the following specific iLEAP test cases:

### Test Case 001: Get ProductFootprint with ShipmentFootprint
Tests the target host system's ability to return `ProductFootprints` with `ShipmentFootprint`s as extensions.

*   **Request**: A `ListFootprints` `GET` request must be sent to the `/2/footprints` endpoint of the test target host system with a valid access token and the syntax specified in [PACT Technical Specifications V2.1 §api-action-list-request](https://wbcsd.github.io/tr/2023/data-exchange-protocol-20231207/#api-action-list-request).
*   **Expected Response**: The target host system must respond with `200 OK` with a JSON body containing a list of `ProductFootprints` (as per [PACT Technical Specifications V2.1 §api-action-list-response](https://wbcsd.github.io/tr/2023/data-exchange-protocol-20231207/#api-action-list-response)). Those which include `productIds` with the format specified for [ShipmentFootprints](https://sine-fdn.github.io/ileap-extension/#pcf-mapping-sf) must be conformant with the Data Model specified for ShipmentFootprint.

### Test Case 002: Get ProductFootprint with TOC
Tests the target host system's ability to return `ProductFootprints` with `TOC`s as extensions.

*   **Request**: A `ListFootprints` `GET` request must be sent to the `/2/footprints` endpoint of the test target host system with a valid access token and the syntax specified in [PACT Technical Specifications V2.1 §api-action-list-request](https://wbcsd.github.io/tr/2023/data-exchange-protocol-20231207/#api-action-list-request).
*   **Expected Response**: The target host system must respond with `200 OK` with a JSON body containing a list of `ProductFootprints`. This list must include all the `TOCs` referenced in the `ShipmentFootprint`s returned in Test Case 001 (identified through the `productIds` field) and may include others. `ProductFootprints` which include `productIds` with the format specified for [TOCs](https://sine-fdn.github.io/ileap-extension/#pcf-mapping-toc) must be conformant with the TOC Data Model.

### Test Case 003: Get ProductFootprint with HOC
Tests the target host system's ability to return `ProductFootprints` with `HOC`s as extensions.

*   **Request**: A `ListFootprints` `GET` request must be sent to the `/2/footprints` endpoint of the test target host system with a valid access token and the syntax specified in [PACT Technical Specifications V2.1 §api-action-list-request](https://wbcsd.github.io/tr/2023/data-exchange-protocol-20231207/#api-action-list-request).
*   **Expected Response**: The target host system must respond with `200 OK` with a JSON body containing a list of `ProductFootprints`. This list must include all the `HOCs` referenced in the `ShipmentFootprint`s returned in Test Case 001 (identified through the `productIds` field) and may include others. `ProductFootprints` which include `productIds` with the format specified for [HOCs](https://sine-fdn.github.io/ileap-extension/#pcf-mapping-hoc) must be conformant with the HOC Data Model.

---

## 3. iLEAP Activity Data Conformance

This conformance level includes all [PACT Required Test Cases](https://wbcsd.github.io/pact-conformance-testing/checklist.html#required-tests) and the following specific iLEAP test cases:

### Test Case 004: Get All TransportActivityData
Tests the target host system's ability to return all `TransportActivityData`.

*   **Request**: A `TransportActivityData` `GET` request must be sent to the `/2/ileap/tad` endpoint of the test target host system with a valid access token and the syntax specified in the [Action TAD Request](https://sine-fdn.github.io/ileap-extension/#action-tad-request) specification. The access token must be obtained through the [PACT's Authentication Flow](https://wbcsd.github.io/tr/2023/data-exchange-protocol-20231207/#api-auth). This can be tested through PACT's Test Cases [001](https://wbcsd.github.io/pact-conformance-testing/checklist.html#tc001) and [002](https://wbcsd.github.io/pact-conformance-testing/checklist.html#tc002).
*   **Expected Response**: The test target host system must respond with `200 OK` and a JSON body containing a list of `TransportActivityData`, in conformance with the [Action TAD Response](https://sine-fdn.github.io/ileap-extension/#action-tad-response) and following the data model specified for TAD.

### Test Case 005: Get Filtered List of TransportActivityData
Tests the target host system's ability to return a filtered list of `TransportActivityData`.

*   **Request**: A `TransportActivityData` `GET` request must be sent to the `/2/ileap/tad` endpoint of the test target host system with a valid access token and the syntax specified in the [Action TAD Request](https://sine-fdn.github.io/ileap-extension/#action-tad-request). The request must include a query parameter `filter`. Any property can be used as a filter, but it is recommended to iterate over all possible values of `TransportMode`:
    - `GET /2/ileap/tad?mode=Road HTTP/1.1`
    - `GET /2/ileap/tad?mode=Rail HTTP/1.1`
    - `GET /2/ileap/tad?mode=Air HTTP/1.1`
    - `GET /2/ileap/tad?mode=Sea HTTP/1.1`
    - `GET /2/ileap/tad?mode=InlandWaterway HTTP/1.1`
*   **Expected Response**: For at least one filter, the test target host system must respond with `200 OK` and a JSON body containing a list of `TransportActivityData` matching the filter, in conformance with the [Action TAD Response](https://sine-fdn.github.io/ileap-extension/#action-tad-response).

### Test Case 006: Get Limited List of TransportActivityData
Tests the target host system's ability to return a limited list of `TransportActivityData`.

*   **Request**: A `TransportActivityData` `GET` request must be sent to the `/2/ileap/tad` endpoint of the test target host system with a valid access token and the syntax specified in the [Action TAD Request](https://sine-fdn.github.io/ileap-extension/#action-tad-request). The request must include a query parameter `limit` with value `1`.
*   **Expected Response**: The test target host system must respond with `200 OK` and a JSON body containing one or zero `TransportActivityData`, in conformance with the [Action TAD Response](https://sine-fdn.github.io/ileap-extension/#action-tad-response).
    *   If a pagination link is returned, it must conform to the syntax specified in the [Action TAD Response](https://sine-fdn.github.io/ileap-extension/#action-tad-response). Upon calling that link, the target host system must respond with `200 OK` and a JSON body containing one or more `TransportActivityData`.
    *   If no pagination link is returned, a `GET` request sent to `/2/ileap/tad` (without the query parameter `limit`) must respond with `200 OK` and a JSON body containing exactly the same number of `TransportActivityData` as that returned in the first request (including the query parameter `limit` with value `1`).

### Test Case 007: Attempt TransportActivityData with Invalid Token
Tests the target host system's ability to handle a `TransportActivityData` request with an invalid access token.

*   **Request**: A `TransportActivityData` `GET` request must be sent to the `/2/ileap/tad` endpoint of the test target host system with an invalid access token and the syntax specified in the [Action TAD Request](https://sine-fdn.github.io/ileap-extension/#action-tad-request).
*   **Expected Response**: The test target host system must respond with a `403 Forbidden` status code and an `AccessDenied` error response code, as specified in the TadResponseBody format in the [Action TAD Response](https://sine-fdn.github.io/ileap-extension/#action-tad-response).

### Test Case 008: Attempt TransportActivityData with Expired Token
Tests the target host system's ability to handle a `TransportActivityData` request with an expired access token.

*   **Request**: A `TransportActivityData` `GET` request must be sent to the `/2/ileap/tad` endpoint of the test target host system with an expired access token and the syntax specified in the [Action TAD Request](https://sine-fdn.github.io/ileap-extension/#action-tad-request). The access token must have been obtained through the [PACT's Authentication Flow](https://wbcsd.github.io/tr/2023/data-exchange-protocol-20231207/#api-auth).
*   **Expected Response**: The test target host system must respond with a `401 Unauthorized` status code and an `ExpiredToken` error response code, as specified in the TadResponseBody format in the [Action TAD Response](https://sine-fdn.github.io/ileap-extension/#action-tad-response).
