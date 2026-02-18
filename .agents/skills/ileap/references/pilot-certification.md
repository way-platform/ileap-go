# iLEAP Certification Pilot Plan

## 1. Project description

The project consists in extending the current SFC Certification Scheme to cover approval of Tool Providers (TPs) for verification services against the iLEAP protocol (iLEAP Technical Specifications). Currently, the Scheme only covers methodology control of Tool Providers for certification services against ISO 14083 and GLEC Framework.

iLEAP is co-sponsored by Smart Freight Centre and the SINE Foundation. Together, both organizations from the secretariat of iLEAP and are responsible for the management and growth of the project.

The iLEAP community includes all types of relevant stakeholders, including Tool Providers. These are iLEAP implementers, in that they commit to implementing iLEAP conformant software.

The GLEC Framework was used as a methodological foundation to publish the ISO standard on GHG accounting for transport operations. The iLEAP protocol was developed to offer a digital counterpart of the methodology standards to enable interoperable data exchange.

SFC, together with the SINE Foundation and the community, have now released the version 1 of the iLEAP Technical Specifications and want to expand the SFC Certification Scheme to cover it. The iLEAP Technical Specifications (v1.0.0) encompass the protocol specifications, the iLEAP data model, and conformance tests. Yet, the iLEAP team also published supporting materials, such as:

1. The iLEAP Community guide
2. A demo API
3. SINE’s Automated Conformance Testing tool (ACT)
4. The iLEAP Whitepaper (v1)

During the pilot phase, the iLEAP Certification process will test implementations of the iLEAP Technical Specifications (v1.0.0) by running them against ACT and performing the tests documented in the Technical Specifications’ Appendix C. iLEAP Certification will be construed as an add-on to SFC Certification. Only Data Transactions (DT) 1 and 2 will be subjected to mandatory testing. DT#3 will be subjected to optional testing. DT#4 remains out of scope for the pilot.

The pilot provides an opportunity for Tool Providers to certify their implementation of iLEAP, while contributing to shape the certification process and its further development.

In the future, iLEAP Certification shall be further integrated into the SFC Certification process, aiming at a thorough end-to-end testing, the details of which are yet to be defined.

## 2. Project details

| Applicable standard(s): | iLEAP Technical Specifications v1.0.0 |
| :--- | :--- |
| MOU in place? | N/A |
| SFC team contacts: | Stichting Smart Freight Centre<br>SFC Team |
| SINE team contact: | SINE Foundation e.V.<br>SINE Team |
| TP team contact: | [Participant Organization]<br>[Participant Name] <[Participant Email]> |

## 3. Useful documents

*   iLEAP Community Guide
*   iLEAP Tech Specs
*   ACT (Automated Conformance Testing) Web UI
*   ACT (Automated Conformance Testing) CLI tool
*   iLEAP demo API (and source code)
*   iLEAP Data Model Rust implementation
*   PACT Tests documentation

## 4. Roadmap (to be adapted in Q2)

| Timeline | Milestone |
| :--- | :--- |
| Q1 2026 | Pilot: two pilots are run with SFC Certified tool providers selected from a shortlist (free of charge) |
| Q2 to end of Q3 2026 | Rollout: the remaining tool providers in the shortlist are invited to go through the certification process (with a significant discount) |
| 2027-2028 | Progressive integration of iLEAP Certification as a mandatory part of SFC certification, depending on the maturity and development of the ecosystem |

## 5. Appendix: iLEAP Certification Pilot Details

### 5.1. Scope

**Timeline:** The pilot phase will run during Q1, as early as possible.

**Participants:** The pilot will encompass two participants, running pilot certifications independently.

**Coverage:** The pilot will test iLEAP implementations, covering all data transactions and data types in v1.0.0 (see Success Criteria below).

**Limitations:** Only SFC Certified tool providers are eligible.

**Cost:** Certifications during the pilot phase will be free.

### 5.2. Participant Selection

**Eligibility**

To be eligible to participate in the iLEAP Certification pilot, tool providers must satisfy the following two pre-requisites:

1. Being SFC Certified (i.e., having a valid, less than 12 months old, SFC Certification);
2. Having implemented at least iLEAP Data Transactions #1 (TCE Data Exchange) and #2 (TOC/HOC Data Exchange).

**Application Process**

Tool providers receiving this document are on the shortlist and have expressed their interest in participating in the iLEAP Certification pilot. The iLEAP team will prioritize providers who can commit to the earliest date. Candidates will be notified of their selection status by the mid-February 2026 the latest.

### 5.3. Pilot Certification Process

**Pre-requisites**

*   Being SFC Certified (i.e., having a valid, less than 12 months old, SFC Certification);
*   Having implemented at least iLEAP Data Transactions #1 (TCE Data Exchange) and #2 (TOC/HOC Data Exchange).

**Process**

iLEAP implementations will be tested using the Automated Conformance Testing (ACT):

*   PACT-related tests, performed using the PACT Conformance Service.
*   iLEAP-related tests, documented in the Appendix C of the iLEAP Technical Specifications.
*   The test runs will be performed asynchronously or during a call.

The test results will be evaluated by the iLEAP team, leading to two possible outcomes:

*   Passed: if the success criteria are met (see below), the tool will be considered iLEAP Certified.
*   Failed: If the success criteria are not met (see below), the tool provider will be given a two-week timeframe (subject to negotiation) to implement the necessary API changes and rerun the certification tests. This remediation cycle may be repeated up to two times, allowing for a maximum of three certification test runs.
    *   After three failed attempts: If concerns arise about ACT, the tests, or the Technical Specifications, the iLEAP team will review potential adjustments. Tool providers may also be offered the opportunity to retest at a later date at no additional cost.
*   Two calls (30-45 minutes each) will be scheduled: one before testing begins and one after results are available.

**Success Criteria**

A tool provider is considered iLEAP Certified if and only if its iLEAP implementation:

*   Passes essential PACT tests. These are the PACT-related tests that are necessary to enable iLEAP Data Exchanges:
    *   Authentication tests;
    *   Action ListFootprints tests;
    *   Action GetFootprint tests.
*   Passes iLEAP Tests related to DT#1 and DT#2

In addition, tests related to DT#3 will be run. Passing them is optional, but recommended. DT#4-related tests are out of scope for the pilot.

**iLEAP Conformance Matrix for Certification Pilot**

*(Note: The matrix image was removed. Refer to official documentation for the matrix.)*

### 5.4. Pilot Participant Expectations

Tool providers participating in the pilot are expected to:

*   Obtain iLEAP Conformance:
    *   As detailed above, this does not mean that the iLEAP implementation passes all the relevant tests at first try. Tool providers are expected to commit resources to pass all the tests during Q1.
    *   If iLEAP Conformance cannot be achieved due to substantive issues (detailed below), providers may participate in a complimentary certification round after the pilot concludes and any necessary changes are implemented.
        *   Substantive issues include, but are not limited to:
            *   Disagreements about the Technical Specifications that warrant community discussion;
            *   Conflicts between the testing methodology and the provider's implementation approach.
*   Provide feedback about the certification process, including:
    *   The quality of the materials;
    *   The team support;
    *   The tests being performed.
*   Provide suggestions regarding how iLEAP Certification might be improved.

### 5.5. Feedback

We'll collect feedback continuously through calls and ongoing communication. After the pilot concludes, we'll distribute a survey to gather structured input on the process and future directions.

### 5.6. Certification Validity

iLEAP Certifications obtained during the pilot will remain valid for the duration of the corresponding SFC Certification validity period.

If an SFC Certification has a short validity period (e.g., 3 months), we will work with you to either extend the iLEAP Certification or provide complimentary retesting as part of the pilot program.

### 5.7. Next Steps

**Post-pilot improvements (Q2 2026):** Following the pilot phase, we will implement changes and improvements based on participant feedback.

**Rollout phase (Q2-Q3 2026):** The rollout will begin immediately after (and may briefly overlap with) the pilot phase, with completion targeted for end of Q3 2026.

**Full operations (Q4 2026 onward):** The iLEAP Certification process will be fully operational starting Q4 2026.

**Mandatory integration (2027-2028):** During 2027 and 2028, we aim to turn iLEAP Certification into a mandatory component of SFC certification.
