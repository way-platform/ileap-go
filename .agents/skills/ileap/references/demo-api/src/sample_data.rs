/*
* Copyright (c) Martin Pompéry
*
* This source code is licensed under the MIT license found in the
* LICENSE file in the crate's root directory of this source tree.
*/
use chrono::prelude::*;
use chrono::Duration;
use ileap_data_model::*;
use pact_data_model::*;
use rust_decimal_macros::dec;
use uuid::uuid;

lazy_static!(
    static ref EXAMPLE_1: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("91715e5e-fd0b-4d1c-8fab-76290c46e6ed")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(1),
        created: Utc.with_ymd_and_hms(2022, 3, 1, 9, 32, 20).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: Some(Utc.with_ymd_and_hms(2022, 3, 1, 9, 32, 20).unwrap()),
        validity_period_end: Some(Utc.with_ymd_and_hms(2024, 12, 31, 00, 00, 00).unwrap()),
        company_name: String::from("My Corp").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:uuid:69585GB6-56T9-6958-E526-6FDGZJHU1326".to_string()), Urn::from("urn:epc:id:sgln:562958.00000.4".to_string())]),
        product_description: "Bio-Ethanol 98%, corn feedstock (bulk - no packaging)".to_string(),
        product_ids: ProductIdSet(vec![Urn::from("urn:gtin:5695872369587".to_string())]),
        product_category_cpc: String::from("6398").into(),
        product_name_company: String::from("Green Ethanol").into(),
        comment: "".into(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::Liter,
            unitary_product_amount: dec!(1).into(),
            p_cf_excluding_biogenic: dec!(1.63).into(),
            p_cf_including_biogenic: Some(dec!(1.85).into()),
            fossil_ghg_emissions: dec!(1.5).into(),
            fossil_carbon_content: dec!(0).into(),
            biogenic_carbon_content: dec!(0.41).into(),
            d_luc_ghg_emissions: Some(dec!(0.8).into()),
            land_management_ghg_emissions: Some(dec!(0.6).into()),
            other_biogenic_ghg_emissions: Some(dec!(0.4).into()),
            i_luc_ghg_emissions: Some(dec!(0).into()),
            biogenic_carbon_withdrawal: Some(dec!(-1.5).into()),
            aircraft_ghg_emissions: Some(dec!(0.2).into()),
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp, DeprecatedCrossSectoralStandard::ISO14067]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![
                ProductOrSectorSpecificRule {
                    operator: ProductOrSectorSpecificRuleOperator::Other,
                    rule_names: vec![String::from("The Product Carbon Footprint Guideline for the Chemical Industry, v.2.0").into()].into(),
                    other_operator_name: Some(String::from("Tfs").into())
                }])),
            biogenic_accounting_methodology: Some(BiogenicAccountingMethodology::Ghgp),
            boundary_processes_description: String::from("1) Material acquisition and preprocessing, including growth of corn 2) Production: fuel consumption, electricity consumption, water consumption, process-generated direct emissions 3) Distribution and storage: transportation of the finished product from manufacturing site to storage site"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 00, 00, 00).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 00, 00, 00).unwrap(),
            geographic_scope: Some(GeographicScope::Regional(UNRegionOrSubregion::WesternEurope)),
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: Some("Using mass allocation following the product specific rule as per PACT Framework decision-making tree".to_string()),
            uncertainty_assessment_description: Some("A model of corn production is involved in predicting emissions from the production of the corn feedstock. Emissions of N2O due to application of nitrogen fertilizers are based on a linear modeling of interactions of the fertilizer with the soil and plant systems. As these interactions are more complicated than the model assumes, there is uncertainty regarding the emissions resulting from this model".to_string()),
            primary_data_share: Some(Percent::from(12.9)),
            dqi: Some(DataQualityIndicators {
                coverage_percent: Percent::from(78.0),
                technological_d_q_r: 1.6.into(),
                temporal_d_q_r: 2.6.into(),
                geographical_d_q_r: 1.8.into(),
                completeness_d_q_r: 1.7.into(),
                reliability_d_q_r: 2.1.into()
            }),
            assurance: Some(Assurance::default()),
        },
        extensions: None,
    };

    static ref EXAMPLE_2: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("61ff98c0-9e13-47d9-bb13-0b5381468165")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(1),
        created: Utc.with_ymd_and_hms(2022, 2, 22, 10, 47, 32).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: Some(Utc.with_ymd_and_hms(2022, 2, 22, 10, 47, 32).unwrap()),
        validity_period_end: Some(Utc.with_ymd_and_hms(2024, 12, 31, 00, 00, 00).unwrap()),
        company_name: String::from("My Corp").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:uuid:51131FB5-42A2-4267-A402-0ECFEFAD1619".to_string()), Urn::from("urn:epc:id:sgln:4063973.00000.8".to_string())]),
        product_description: "12L Bottle of bio-Ethanol 98%, corn feedstock (including 100% fossil plastic packaging)".to_string(),
        product_ids: ProductIdSet(vec![Urn::from("urn:gtin:4712345060507".to_string())]),
        product_category_cpc: String::from("3342").into(),
        product_name_company: String::from("Green Ethanol").into(),
        comment: "".into(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::Liter,
            unitary_product_amount: dec!(12).into(),
            p_cf_excluding_biogenic: dec!(1.75).into(),
            p_cf_including_biogenic: Some(dec!(1.97).into()),
            fossil_ghg_emissions: dec!(1.5).into(),
            fossil_carbon_content: dec!(0).into(),
            biogenic_carbon_content: dec!(0.41).into(),
            d_luc_ghg_emissions: Some(dec!(0.8).into()),
            land_management_ghg_emissions: Some(dec!(0.6).into()),
            other_biogenic_ghg_emissions: Some(dec!(0.4).into()),
            i_luc_ghg_emissions: Some(dec!(0).into()),
            biogenic_carbon_withdrawal: Some(dec!(-1.5).into()),
            aircraft_ghg_emissions: Some(dec!(0.2).into()),
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp, DeprecatedCrossSectoralStandard::ISO14067]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![
                ProductOrSectorSpecificRule {
                    operator: ProductOrSectorSpecificRuleOperator::Other,
                    rule_names: vec![String::from("The Product Carbon Footprint Guideline for the Chemical Industry, v.2.0").into()].into(),
                    other_operator_name: Some(String::from("Tfs").into())
                }])),
            biogenic_accounting_methodology: Some(BiogenicAccountingMethodology::Ghgp),
            boundary_processes_description: String::from("1) Material acquisition and preprocessing, including growth of corn 2) Production: fuel consumption, electricity consumption, water consumption, process-generated direct emissions 3) Distribution and storage: transportation of the finished product from manufacturing site to storage site"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 00, 00, 00).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 00, 00, 00).unwrap(),
            geographic_scope: Some(GeographicScope::Country(ISO3166CC(String::from("DE")))),
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.8),
            exempted_emissions_description: "Using the most conservative emission factor from a region with higher energy intensive grid for input A resulted in a contribution of 0.8% for this input. This is less than 1% and therefore considered under the cut off rule".to_string(),
            packaging_emissions_included: true,
            packaging_ghg_emissions: Some(dec!(0.12).into()),
            allocation_rules_description: Some("Using mass allocation following the product specific rule as per PACT Framework decision-making tree".to_string()),
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(16.8)),
            dqi: Some(DataQualityIndicators {
                coverage_percent: Percent::from(87.0),
                technological_d_q_r: 2.3.into(),
                temporal_d_q_r: 1.4.into(),
                geographical_d_q_r: 2.5.into(),
                completeness_d_q_r: 1.1.into(),
                reliability_d_q_r: 1.6.into()
            }),
            assurance: Some(Assurance::default()),
        },
        extensions: None,
    };

    // a footprint deprecated by EXAMPLE_4
    static ref EXAMPLE_3: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("fb77319f-2338-4338-868a-98b2206340ad")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(2),
        created: Utc.with_ymd_and_hms(2022, 3, 15, 11, 47, 32).unwrap(),
        updated: Some(Utc.with_ymd_and_hms(2023, 6, 27, 12, 12, 3).unwrap()),
        status: PfStatus::Deprecated,
        status_comment: Some("Replaced by a new version".to_string()),
        validity_period_start: Some(Utc.with_ymd_and_hms(2022, 3, 15, 11, 47, 32).unwrap()),
        validity_period_end: Some(Utc.with_ymd_and_hms(2023, 6, 27, 12, 12, 3).unwrap()),
        company_name: String::from("My Corp").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:uuid:25639HN5-58Q6-1238-S596-9STHZHZJ5623".to_string()), Urn::from("urn:epc:id:sgln:6957976.00000.1".to_string())]),
        product_description: "Cardboard box 50x40x40 cm".to_string(),
        product_ids: ProductIdSet(vec![Urn::from("urn:gtin:5268596541023".to_string())]),
        product_category_cpc: String::from("4365").into(),
        product_name_company: String::from("Cardboard504040").into(),
        comment: "".into(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::Kilogram,
            unitary_product_amount: dec!(0.8).into(),
            p_cf_excluding_biogenic: dec!(0.28).into(),
            p_cf_including_biogenic: Some(dec!(-0.28).into()),
            fossil_ghg_emissions: dec!(0.19).into(),
            fossil_carbon_content: dec!(0.08).into(),
            biogenic_carbon_content: dec!(0.44).into(),
            d_luc_ghg_emissions: Some(dec!(0.42).into()),
            land_management_ghg_emissions: Some(dec!(0.34).into()),
            other_biogenic_ghg_emissions: Some(dec!(0.2).into()),
            i_luc_ghg_emissions: Some(dec!(0.03).into()),
            biogenic_carbon_withdrawal: Some(dec!(-1.6).into()),
            aircraft_ghg_emissions: Some(dec!(0.08).into()),
            characterization_factors: CharacterizationFactors::Ar5,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR5").into(), String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![
                ProductOrSectorSpecificRule {
                    operator: ProductOrSectorSpecificRuleOperator::EPDInternational,
                    rule_names: vec![String::from("PCR cardboard").into()].into(),
                    other_operator_name: None
                }
            ])),
            biogenic_accounting_methodology: Some(BiogenicAccountingMethodology::Pef),
            boundary_processes_description: String::from("1) Material acquisition and preprocessing, including growth of trees 2) Production: fuel consumption, electricity consumption, water consumption, process-generated direct emissions 3) Distribution and storage: transportation of the finished product from manufacturing site to storage site"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 00, 00, 00).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 00, 00, 00).unwrap(),
            geographic_scope: Some(GeographicScope::Country(ISO3166CC("FR".to_string()))),
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Gabi").into(),
                version: String::from("2022").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: Some("No allocation used, process subdivision was possible".to_string()),
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(12.6)),
            dqi: Some(DataQualityIndicators {
                coverage_percent: Percent::from(83.0),
                technological_d_q_r: 1.8.into(),
                temporal_d_q_r: 1.2.into(),
                geographical_d_q_r: 1.9.into(),
                completeness_d_q_r: 1.7.into(),
                reliability_d_q_r: 1.4.into()
            }),
            assurance: Some(Assurance {
                assurance: true,
                coverage: Some(AssuranceCoverage::PcfSystem),
                level: Some(AssuranceLevel::Limited),
                boundary: Some(AssuranceBoundary::CradleToGate),
                provider_name: "My Auditor".to_string(),
                completed_at: Some(Utc.with_ymd_and_hms(2022, 12, 15, 00, 00, 00).unwrap()),
                standard_name: Some("ISO 14044".to_string()),
                comments: None
            }),
        },
        extensions: None,
    };

    // this is the PCF superseeding EXAMPLE_3
    static ref EXAMPLE_4: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("f369091a-aa5d-4248-9bd5-2812329e1ef1")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: Some(NonEmptyPfIdVec(vec![PfId(uuid!("fb77319f-2338-4338-868a-98b2206340ad"))])),
        version: VersionInteger(1),
        created: Utc.with_ymd_and_hms(2023, 6, 27, 12, 12, 3).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: Some(Utc.with_ymd_and_hms(2023, 2, 1, 00, 00, 00).unwrap()),
        validity_period_end: Some(Utc.with_ymd_and_hms(2025, 8, 31, 00, 00, 00).unwrap()),
        company_name: String::from("My Corp").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:uuid:25639HN5-58Q6-1238-S596-9STHZHZJ5623".to_string()), Urn::from("urn:epc:id:sgln:6957976.00000.1".to_string())]),
        product_description: "Cardboard box 50x40x40 cm".to_string(),
        product_ids: ProductIdSet(vec![Urn::from("urn:gtin:5268596541023".to_string())]),
        product_category_cpc: String::from("4365").into(),
        product_name_company: String::from("Cardboard504040").into(),
        comment: "".into(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::Kilogram,
            unitary_product_amount: dec!(0.8).into(),
            p_cf_excluding_biogenic: dec!(0.32).into(),
            p_cf_including_biogenic: Some(dec!(-0.28).into()),
            fossil_ghg_emissions: dec!(0.23).into(),
            fossil_carbon_content: dec!(0.08).into(),
            biogenic_carbon_content: dec!(0.44).into(),
            d_luc_ghg_emissions: Some(dec!(0.42).into()),
            land_management_ghg_emissions: Some(dec!(0.34).into()),
            other_biogenic_ghg_emissions: Some(dec!(0.2).into()),
            i_luc_ghg_emissions: Some(dec!(0.03).into()),
            biogenic_carbon_withdrawal: Some(dec!(-1.6).into()),
            aircraft_ghg_emissions: Some(dec!(0.08).into()),
            characterization_factors: CharacterizationFactors::Ar5,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR5").into(), String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![
                ProductOrSectorSpecificRule {
                    operator: ProductOrSectorSpecificRuleOperator::EPDInternational,
                    rule_names: vec![String::from("PCR cardboard").into()].into(),
                    other_operator_name: None
                }
            ])),
            biogenic_accounting_methodology: Some(BiogenicAccountingMethodology::Pef),
            boundary_processes_description: String::from("1) Material acquisition and preprocessing, including growth of trees 2) Production: fuel consumption, electricity consumption, water consumption, process-generated direct emissions 3) Distribution and storage: transportation of the finished product from manufacturing site to storage site"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 00, 00, 00).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 00, 00, 00).unwrap(),
            geographic_scope: Some(GeographicScope::Subdivision(String::from("FR-89").into())),
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Gabi").into(),
                version: String::from("2022").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: Some("No allocation used, process subdivision was possible".to_string()),
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(12.6)),
            dqi: Some(DataQualityIndicators {
                coverage_percent: Percent::from(83.0),
                technological_d_q_r: 1.8.into(),
                temporal_d_q_r: 1.2.into(),
                geographical_d_q_r: 1.9.into(),
                completeness_d_q_r: 1.7.into(),
                reliability_d_q_r: 1.4.into()
            }),
            assurance: Some(Assurance {
                assurance: true,
                coverage: Some(AssuranceCoverage::PcfSystem),
                level: Some(AssuranceLevel::Limited),
                boundary: Some(AssuranceBoundary::CradleToGate),
                provider_name: "My Auditor".to_string(),
                completed_at: Some(Utc.with_ymd_and_hms(2022, 12, 15, 00, 00, 00).unwrap()),
                standard_name: Some("ISO 14044".to_string()),
                comments: None
            }),
        },
        extensions: None,
    };

    // a ShipmentFooprint PCF (iLEAP Extension)
    static ref SHIPMENT_SIMPLE_EXAMPLE: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("d9be4477-e351-45b3-acd9-e1da05e6f633")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(0),
        created: Utc.with_ymd_and_hms(2022, 5, 22, 21, 47, 32).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: None,
        validity_period_end: None,
        company_name: String::from("Super Duper Transport Co.").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:epc:id:sgln:4063973.00000.8".to_string())]),
        product_description: String::from("Logistics emissions related to shipment with ID shipment-1"),
        product_ids: ProductIdSet(vec![Urn::from("urn:pathfinder:product:customcode:vendor-assigned:shipment:shipment-simple-1".to_string())]),
        product_category_cpc: String::from("83117").into(),
        product_name_company: String::from("Shipment with ID shipment-simple-1").into(),
        comment: "".to_string(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::TonKilometer,
            unitary_product_amount: dec!(16920).into(),
            p_cf_excluding_biogenic: dec!(1962.72).into(),
            p_cf_including_biogenic: None,
            fossil_ghg_emissions: dec!(1962.72).into(),
            fossil_carbon_content: dec!(1962.72).into(),
            biogenic_carbon_content: dec!(0).into(),
            d_luc_ghg_emissions: None,
            land_management_ghg_emissions: None,
            other_biogenic_ghg_emissions: None,
            i_luc_ghg_emissions: None,
            biogenic_carbon_withdrawal: None,
            aircraft_ghg_emissions: None,
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![])),
            biogenic_accounting_methodology: None,
            boundary_processes_description: String::from("SFC GLEC Framework-conforming (W2W CO2e emissions)"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 0, 0, 0).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 0, 0, 0).unwrap(),
            geographic_scope: None,
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.9.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: None,
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(56.12)),
            dqi: None,
            assurance: None
        },
        extensions: Some(vec![
            DataModelExtension {
                spec_version: SpecVersionString::from("2.0.0".to_string()),
                data_schema: "https://api.ileap.sine.dev/shipment-footprint.json".to_string(),
                documentation: Some("https://sine-fdn.github.io/ileap-extension/".to_string()),
                data: ILeapType::ShipmentFootprint(shipment_footprint(
                        ShipmentArgs {mass: "40000".to_string(),
                        shipment_id: "shipment-simple-1".to_string(),
                        tces: NonEmptyVec::<Tce>::from(vec![
                            tce(
                                TceArgs {
                                    tce_id: "abcdef".to_string(),
                                    prev_tce_ids: Some(vec![]),
                                    toc_id: Some("truck-40t-euro5-de".to_string()),
                                    hoc_id: None,
                                    shipment_id: "shipment-1".to_string(),
                                    mass: dec!(40000).into(),
                                    distance: GlecDistance::new_actual(dec!(423).into()),
                                    transport_activity: dec!(16920).into(),
                                    co2e_wtw: dec!(1962.72).into(),
                                    co2e_ttw: dec!(1505.88).into()
                                },
                            )
                        ])}
                    )
                )
            }
        ])
    };

    static ref SHIPMENT_MULTIMODAL_EXAMPLE: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("380baaac-6f47-471c-83d5-a9db87b8bede")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(0),
        created: Utc.with_ymd_and_hms(2022, 5, 22, 21, 47, 32).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: None,
        validity_period_end: None,
        company_name: String::from("Super Duper Transport Co.").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:epc:id:sgln:4063973.00000.8".to_string())]),
        product_description: String::from("Logistics emissions related to shipment with ID shipment-multi-modal-1"),
        product_ids: ProductIdSet(vec![Urn::from("urn:pathfinder:product:customcode:vendor-assigned:shipment:shipment-multi-modal-1".to_string())]),
        product_category_cpc: String::from("83117").into(),
        product_name_company: String::from("Shipment with ID shipment-multi-modal-1").into(),
        comment: "".to_string(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::TonKilometer,
            unitary_product_amount: dec!(33840).into(),
            p_cf_excluding_biogenic: dec!(3131.06).into(),
            p_cf_including_biogenic: None,
            fossil_ghg_emissions: dec!(3131.06).into(),
            fossil_carbon_content: dec!(3131.06).into(),
            biogenic_carbon_content: dec!(0).into(),
            d_luc_ghg_emissions: None,
            land_management_ghg_emissions: None,
            other_biogenic_ghg_emissions: None,
            i_luc_ghg_emissions: None,
            biogenic_carbon_withdrawal: None,
            aircraft_ghg_emissions: None,
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![])),
            biogenic_accounting_methodology: None,
            boundary_processes_description: String::from("SFC GLEC Framework-conforming (W2W CO2e emissions)"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 0, 0, 0).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 0, 0, 0).unwrap(),
            geographic_scope: None,
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.9.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: None,
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(56.12)),
            dqi: None,
            assurance: None
        },
        extensions: Some(vec![
            DataModelExtension {
                spec_version: SpecVersionString::from("2.0.0".to_string()),
                data_schema: "https://api.ileap.sine.dev/shipment-footprint.json".to_string(),
                documentation: Some("https://sine-fdn.github.io/ileap-extension/".to_string()),
                data: ILeapType::ShipmentFootprint(shipment_footprint(
                        ShipmentArgs {mass: "40000".to_string(),
                        shipment_id: "shipment-multi-modal-1".to_string(),
                        tces: NonEmptyVec::<Tce>::from(vec![
                            tce(
                                TceArgs {
                                    tce_id: "tce-1-toc-rail-1".to_string(),
                                    prev_tce_ids: Some(vec![]),
                                    hoc_id: None,
                                    toc_id: Some("toc-rail-1".to_string()),
                                    shipment_id: "shipment-multi-modal-1".to_string(),
                                    mass: dec!(40000).into(),
                                    distance: GlecDistance::new_actual(dec!(423).into()),
                                    transport_activity:  dec!(16920).into(),
                                    co2e_wtw:  dec!(118.44).into(),
                                    co2e_ttw: dec!(0).into() },
                            ),
                            tce(
                                TceArgs {
                                    tce_id: "tce-2-hoc-transshipment-1".to_string(),
                                    prev_tce_ids: Some(vec!["tce-1-toc-rail-1".to_string()]),
                                    toc_id: None,
                                    hoc_id: Some("hoc-transshipment-1".to_string()),
                                    shipment_id: "shipment-multi-modal-1".to_string(),
                                    mass: dec!(40000).into(),
                                    distance: GlecDistance::new_actual(dec!(0).into()),
                                    transport_activity: dec!(0).into(),
                                    co2e_wtw: dec!(1320).into(),
                                    co2e_ttw: dec!(400).into()
                                },
                            ),
                            tce(
                                TceArgs {
                                    tce_id: "tce-3-toc-road-1".to_string(),
                                    prev_tce_ids: Some(vec!["tce-2-hoc-transshipment-1".to_string()]),
                                    toc_id: Some("toc-road-1".to_string()),
                                    hoc_id: None,
                                    shipment_id: "shipment-multi-modal-1".to_string(),
                                    mass: dec!(40000).into(),
                                    distance: GlecDistance::new_actual(dec!(423).into()),
                                    transport_activity: dec!(16920).into(),
                                    co2e_wtw: dec!(1692.62).into(),
                                    co2e_ttw: dec!(1505.88).into() },
                            )
                        ])}
                    )
                )
            }
        ])
    };

    static ref SHIPMENT_PRE_AND_POST_LEGS_EXAMPLE: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("4861c2a5-68a5-4c40-bcc0-7e234ad24184")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(0),
        created: Utc.with_ymd_and_hms(2022, 5, 22, 21, 47, 32).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: None,
        validity_period_end: None,
        company_name: String::from("Super Duper Transport Co.").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:epc:id:sgln:4063973.00000.8".to_string())]),
        product_description: String::from("Logistics emissions related to shipment with ID pre-and-post-legs-1"),
        product_ids: ProductIdSet(vec![Urn::from("urn:pathfinder:product:customcode:vendor-assigned:shipment:pre-and-post-legs-1".to_string())]),
        product_category_cpc: String::from("83117").into(),
        product_name_company: String::from("Shipment with ID pre-and-post-legs-1").into(),
        comment: "".to_string(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::TonKilometer,
            unitary_product_amount: dec!(320).into(),
            p_cf_excluding_biogenic: dec!(385.76).into(),
            p_cf_including_biogenic: None,
            fossil_ghg_emissions: dec!(385.76).into(),
            fossil_carbon_content: dec!(385.76).into(),
            biogenic_carbon_content: dec!(0).into(),
            d_luc_ghg_emissions: None,
            land_management_ghg_emissions: None,
            other_biogenic_ghg_emissions: None,
            i_luc_ghg_emissions: None,
            biogenic_carbon_withdrawal: None,
            aircraft_ghg_emissions: None,
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![])),
            biogenic_accounting_methodology: None,
            boundary_processes_description: String::from("SFC GLEC Framework-conforming (W2W CO2e emissions)"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 0, 0, 0).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 0, 0, 0).unwrap(),
            geographic_scope: None,
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.9.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: None,
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(0.7)),
            dqi: None,
            assurance: None
        },
        extensions: Some(vec![
            DataModelExtension {
                spec_version: SpecVersionString::from("2.0.0".to_string()),
                data_schema: "https://api.ileap.sine.dev/shipment-footprint.json".to_string(),
                documentation: Some("https://sine-fdn.github.io/ileap-extension/".to_string()),
                data: ILeapType::ShipmentFootprint(shipment_footprint(
                        ShipmentArgs {
                            mass: "4000".to_string(),
                            shipment_id: "pre-and-post-legs-1".to_string(),
                            tces: NonEmptyVec::<Tce>::from(vec![
                                tce(
                                    TceArgs {
                                        tce_id: "tce-1-pre-leg-1".to_string(),
                                        prev_tce_ids: Some(vec![]),
                                        toc_id: Some("toc-small-truck-1".to_string()),
                                        hoc_id: None,
                                        shipment_id: "pre-and-post-legs-1".to_string(),
                                        mass: dec!(4000).into(),
                                        distance: GlecDistance::new_actual(dec!(30).into()),
                                        transport_activity: dec!(120).into(),
                                        co2e_wtw: dec!(95.16).into(),
                                        co2e_ttw: dec!(73.08).into(),
                                    },
                                ),
                                tce(
                                    TceArgs {
                                        tce_id: "tce-2-hub-1".to_string(),
                                        prev_tce_ids: Some(vec!["tce-1-pre-leg-1".to_string()]),
                                        toc_id: None,
                                        hoc_id: Some("hoc-warehouse-1".to_string()),
                                        shipment_id: "pre-and-post-legs-1".to_string(),
                                        mass: dec!(4000).into(),
                                        distance: GlecDistance::new_actual(dec!(0).into()),
                                        transport_activity: dec!(0).into(),
                                        co2e_wtw: dec!(132).into(),
                                        co2e_ttw: dec!(40).into(),
                                    }
                                ),
                                tce(
                                    TceArgs {
                                        tce_id: "tce-3-post-leg-1".to_string(),
                                        prev_tce_ids: Some(vec!["tce-2-hub-1".to_string()]),
                                        toc_id: Some("toc-small-truck-1".to_string()),
                                        hoc_id: None,
                                        shipment_id: "pre-and-post-legs-1".to_string(),
                                        mass: dec!(4000).into(),
                                        distance: GlecDistance::new_actual(dec!(50).into()),
                                        transport_activity: dec!(200).into(),
                                        co2e_wtw: dec!(158.6).into(),
                                        co2e_ttw: dec!(121.8).into(),
                                    },
                                )
                            ])
                        }
                    )
                )
            }
        ])
    };

    // a TOC PCF (iLEAP Extension)
    static ref TOC_ROAD_EXAMPLE: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("f3c04ec8-b33a-43b1-9fa7-d6a448fd60af")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(0),
        created: Utc.with_ymd_and_hms(2022, 5, 22, 21, 47, 32).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: None,
        validity_period_end: None,
        company_name: String::from("Super Duper Transport Co.").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:epc:id:sgln:4063973.00000.8".to_string())]),
        product_description: String::from("Logistics emissions related to TOC with ID toc-road-1"),
        product_ids: ProductIdSet(vec![Urn::from("urn:pathfinder:product:customcode:vendor-assigned:toc:toc-road-1".to_string())]),
        product_category_cpc: String::from("83117").into(),
        product_name_company: String::from("TOC with ID toc-road-1").into(),
        comment: "".to_string(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::TonKilometer,
            unitary_product_amount: dec!(1).into(),
            p_cf_excluding_biogenic: dec!(0.116).into(),
            p_cf_including_biogenic: None,
            fossil_ghg_emissions: dec!(0.116).into(),
            fossil_carbon_content: dec!(0.116).into(),
            biogenic_carbon_content: dec!(0).into(),
            d_luc_ghg_emissions: None,
            land_management_ghg_emissions: None,
            other_biogenic_ghg_emissions: None,
            i_luc_ghg_emissions: None,
            biogenic_carbon_withdrawal: None,
            aircraft_ghg_emissions: None,
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![])),
            biogenic_accounting_methodology: None,
            boundary_processes_description: String::from("SFC GLEC Framework-conforming (W2W CO2e emissions)"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 0, 0, 0).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 0, 0, 0).unwrap(),
            geographic_scope: None,
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.9.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: None,
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(56.12)),
            dqi: None,
            assurance: None
        },
        extensions: Some(vec![
            DataModelExtension {
                spec_version: SpecVersionString::from("2.0.0".to_string()),
                data_schema: "https://api.ileap.sine.dev/toc.json".to_string(),
                documentation: Some("https://sine-fdn.github.io/ileap-extension/".to_string()),
                data: ILeapType::Toc(toc(
                        TocArgs {
                            toc_id: "toc-road-1".to_string(),
                            mode: TransportMode::Road,
                            load_factor: Some(dec!(0.6).to_string()),
                            empty_distance_factor: Some(dec!(0.3).to_string()),
                            temperature_control: Some(TemperatureControl::Refrigerated),
                            truck_loading_sequence: Some(TruckLoadingSequence::Ftl),
                            energy_carriers: vec![EnergyCarrier {
                                energy_carrier: EnergyCarrierType::Diesel,
                                feedstocks: Some(vec![Feedstock {
                                    feedstock: FeedstockType::Fossil,
                                    feedstock_share: None,
                                    region_provenance: Some("Europe".to_string()),
                                }]),
                                energy_consumption: None,
                                energy_consumption_unit: Some(EnergyConsumptionUnit::Kg),
                                emission_factor_wtw: dec!(4.13).into(),
                                emission_factor_ttw: dec!(3.17).into(),
                                relative_share: dec!(1.0).into(),
                            }].into(),
                            co2e_intensity_wtw: dec!(0.116).into(),
                            co2e_intensity_ttw: dec!(0.089).into(),
                            transport_activity_unit: TransportActivityUnit::Tkm
                        },
                    )
                ),
            }
        ])
    };

    static ref TOC_ROAD_40T_EXAMPLE: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("f9526d98-de57-4d24-a131-95fcef75defb")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(0),
        created: Utc.with_ymd_and_hms(2022, 5, 22, 21, 47, 32).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: None,
        validity_period_end: None,
        company_name: String::from("Super Duper Transport Co.").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:epc:id:sgln:4063973.00000.8".to_string())]),
        product_description: String::from("Logistics emissions related to TOC with ID truck-40t-euro5-de"),
        product_ids: ProductIdSet(vec![Urn::from("urn:pathfinder:product:customcode:vendor-assigned:toc:truck-40t-euro5-de".to_string())]),
        product_category_cpc: String::from("83117").into(),
        product_name_company: String::from("TOC with ID truck-40t-euro5-de").into(),
        comment: "".to_string(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::TonKilometer,
            unitary_product_amount: dec!(1).into(),
            p_cf_excluding_biogenic: dec!(0.116).into(),
            p_cf_including_biogenic: None,
            fossil_ghg_emissions: dec!(0.116).into(),
            fossil_carbon_content: dec!(0.116).into(),
            biogenic_carbon_content: dec!(0).into(),
            d_luc_ghg_emissions: None,
            land_management_ghg_emissions: None,
            other_biogenic_ghg_emissions: None,
            i_luc_ghg_emissions: None,
            biogenic_carbon_withdrawal: None,
            aircraft_ghg_emissions: None,
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![])),
            biogenic_accounting_methodology: None,
            boundary_processes_description: String::from("SFC GLEC Framework-conforming (W2W CO2e emissions)"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 0, 0, 0).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 0, 0, 0).unwrap(),
            geographic_scope: None,
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.9.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: None,
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(56.12)),
            dqi: None,
            assurance: None
        },
        extensions: Some(vec![
            DataModelExtension {
                spec_version: SpecVersionString::from("2.0.0".to_string()),
                data_schema: "https://api.ileap.sine.dev/toc.json".to_string(),
                documentation: Some("https://sine-fdn.github.io/ileap-extension/".to_string()),
                data: ILeapType::Toc(toc(
                        TocArgs {
                            toc_id: "truck-40t-euro5-de".to_string(),
                            mode: TransportMode::Road,
                            load_factor: Some(dec!(0.6).to_string()),
                            empty_distance_factor: Some(dec!(0.3).to_string()),
                            temperature_control: Some(TemperatureControl::Refrigerated),
                            truck_loading_sequence: Some(TruckLoadingSequence::Ftl),
                            energy_carriers: vec![EnergyCarrier {
                                energy_carrier: EnergyCarrierType::Diesel,
                                feedstocks: Some(vec![Feedstock {
                                    feedstock: FeedstockType::Fossil,
                                    feedstock_share: None,
                                    region_provenance: Some("Europe".to_string()),
                                }]),
                                energy_consumption: None,
                                energy_consumption_unit: Some(EnergyConsumptionUnit::Kg),
                                emission_factor_wtw: dec!(4.13).into(),
                                emission_factor_ttw: dec!(3.17).into(),
                                relative_share: dec!(1.0).into(),
                            }].into(),
                            co2e_intensity_wtw: dec!(0.116).into(),
                            co2e_intensity_ttw: dec!(0.089).into(),
                            transport_activity_unit: TransportActivityUnit::Tkm
                        },
                    )
                ),
            }
        ])
    };

    static ref TOC_SMALL_TRUCK_EXAMPLE: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("eff7bded-948d-4ed6-adca-fc4a8f0602a5")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(0),
        created: Utc.with_ymd_and_hms(2022, 5, 22, 21, 47, 32).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: None,
        validity_period_end: None,
        company_name: String::from("Super Duper Transport Co.").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:epc:id:sgln:4063973.00000.8".to_string())]),
        product_description: String::from("Logistics emissions related to TOC with ID toc-small-truck-1"),
        product_ids: ProductIdSet(vec![Urn::from("urn:pathfinder:product:customcode:vendor-assigned:toc:toc-small-truck-1".to_string())]),
        product_category_cpc: String::from("83117").into(),
        product_name_company: String::from("TOC with ID toc-small-truck-1").into(),
        comment: "".to_string(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::TonKilometer,
            unitary_product_amount: dec!(1).into(),
            p_cf_excluding_biogenic: dec!(0.793).into(),
            p_cf_including_biogenic: None,
            fossil_ghg_emissions: dec!(0.793).into(),
            fossil_carbon_content: dec!(0.793).into(),
            biogenic_carbon_content: dec!(0).into(),
            d_luc_ghg_emissions: None,
            land_management_ghg_emissions: None,
            other_biogenic_ghg_emissions: None,
            i_luc_ghg_emissions: None,
            biogenic_carbon_withdrawal: None,
            aircraft_ghg_emissions: None,
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![])),
            biogenic_accounting_methodology: None,
            boundary_processes_description: String::from("SFC GLEC Framework-conforming (W2W CO2e emissions)"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 0, 0, 0).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 0, 0, 0).unwrap(),
            geographic_scope: None,
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.9.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: None,
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(56.12)),
            dqi: None,
            assurance: None
        },
        extensions: Some(vec![
            DataModelExtension {
                spec_version: SpecVersionString::from("2.0.0".to_string()),
                data_schema: "https://api.ileap.sine.dev/toc.json".to_string(),
                documentation: Some("https://sine-fdn.github.io/ileap-extension/".to_string()),
                data: ILeapType::Toc(toc(
                        TocArgs {
                            toc_id: "toc-small-truck-1".to_string(),
                            mode: TransportMode::Road,
                            load_factor: Some(dec!(0.2).to_string()),
                            empty_distance_factor: Some(dec!(0.1).to_string()),
                            temperature_control: Some(TemperatureControl::Ambient),
                            truck_loading_sequence: Some(TruckLoadingSequence::Ftl),
                            energy_carriers: vec![EnergyCarrier {
                                energy_carrier: EnergyCarrierType::Diesel,
                                feedstocks: Some(vec![Feedstock {
                                    feedstock: FeedstockType::Fossil,
                                    feedstock_share: None,
                                    region_provenance: Some("Europe".to_string()),
                                }]),
                                energy_consumption: None,
                                energy_consumption_unit: Some(EnergyConsumptionUnit::Kg),
                                emission_factor_wtw: dec!(4.13).into(),
                                emission_factor_ttw: dec!(3.17).into(),
                                relative_share: dec!(1.0).into(),
                            }].into(),
                            co2e_intensity_wtw: dec!(0.793).into(),
                            co2e_intensity_ttw: dec!(0.609).into(),
                            transport_activity_unit: TransportActivityUnit::Tkm
                        },
                    )
                ),
            }
        ])
    };

    static ref TOC_RAIL_EXAMPLE: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("02e10995-569a-4af9-b7a5-7c46dccb0fd3")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(0),
        created: Utc.with_ymd_and_hms(2022, 5, 22, 21, 47, 32).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: None,
        validity_period_end: None,
        company_name: String::from("Super Duper Transport Co.").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:epc:id:sgln:4063973.00000.8".to_string())]),
        product_description: String::from("Logistics emissions related to TOC with ID toc-rail-1"),
        product_ids: ProductIdSet(vec![Urn::from("urn:pathfinder:product:customcode:vendor-assigned:toc:toc-rail-1".to_string())]),
        product_category_cpc: String::from("83117").into(),
        product_name_company: String::from("TOC with ID toc-rail-1").into(),
        comment: "".to_string(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::TonKilometer,
            unitary_product_amount: dec!(1).into(),
            p_cf_excluding_biogenic: dec!(0.007).into(),
            p_cf_including_biogenic: None,
            fossil_ghg_emissions: dec!(0.007).into(),
            fossil_carbon_content: dec!(0.007).into(),
            biogenic_carbon_content: dec!(0).into(),
            d_luc_ghg_emissions: None,
            land_management_ghg_emissions: None,
            other_biogenic_ghg_emissions: None,
            i_luc_ghg_emissions: None,
            biogenic_carbon_withdrawal: None,
            aircraft_ghg_emissions: None,
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![])),
            biogenic_accounting_methodology: None,
            boundary_processes_description: String::from("SFC GLEC Framework-conforming (W2W CO2e emissions)"),
            reference_period_start: Utc.with_ymd_and_hms(2021, 1, 1, 0, 0, 0).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2022, 1, 1, 0, 0, 0).unwrap(),
            geographic_scope: None,
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.9.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: None,
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(100.0)),
            dqi: None,
            assurance: None
        },
        extensions: Some(vec![
            DataModelExtension {
                spec_version: SpecVersionString::from("2.0.0".to_string()),
                data_schema: "https://api.ileap.sine.dev/toc.json".to_string(),
                documentation: Some("https://sine-fdn.github.io/ileap-extension/".to_string()),
                data:
                    ILeapType::Toc(toc(
                        TocArgs {
                            toc_id: "toc-rail-1".to_string(),
                            mode: TransportMode::Rail,
                            load_factor: Some(dec!(0.6).to_string()),
                            empty_distance_factor: Some(dec!(0.33).to_string()),
                            temperature_control: Some(TemperatureControl::Ambient),
                            truck_loading_sequence: None,
                            energy_carriers: vec![EnergyCarrier {
                                energy_carrier: EnergyCarrierType::Electric,
                                feedstocks: Some(vec![Feedstock {
                                    feedstock: FeedstockType::Grid,
                                    feedstock_share: None,
                                    region_provenance: Some("Europe".to_string()),
                                }]),
                                energy_consumption: None,
                                energy_consumption_unit: Some(EnergyConsumptionUnit::MJ),
                                emission_factor_wtw: dec!(97).into(),
                                emission_factor_ttw: dec!(0).into(),
                                relative_share: dec!(1.0).into(),
                            }].into(),
                            co2e_intensity_wtw: dec!(0.007).into(),
                            co2e_intensity_ttw: dec!(0).into(),
                            transport_activity_unit: TransportActivityUnit::Tkm
                        }
                    )),
            }
        ])
    };

    static ref HOC_TRANSSHIPMENT_EXAMPLE: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("46e095f0-f73d-4ace-adba-2ec8bf305339")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(0),
        created: Utc.with_ymd_and_hms(2024, 5, 22, 21, 47, 32).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: None,
        validity_period_end: None,
        company_name: String::from("Hub Co.").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:epc:id:sgln:4063973.00000.9".to_string())]),
        product_description: String::from("Logistics emissions related to HOC with ID hoc-transshipment-1"),
        product_ids: ProductIdSet(vec![Urn::from("urn:pathfinder:product:customcode:vendor-assigned:hoc:hoc-transshipment-1".to_string())]),
        product_category_cpc: String::from("83117").into(),
        product_name_company:String::from("HOC with ID hoc-transshipment-1").into(),
        comment: "".to_string(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::Kilogram,
            unitary_product_amount: dec!(1000).into(),
            p_cf_excluding_biogenic: dec!(33).into(),
            p_cf_including_biogenic: None,
            fossil_ghg_emissions: dec!(33).into(),
            fossil_carbon_content: dec!(33).into(),
            biogenic_carbon_content: dec!(0).into(),
            d_luc_ghg_emissions: None,
            land_management_ghg_emissions: None,
            other_biogenic_ghg_emissions: None,
            i_luc_ghg_emissions: None,
            biogenic_carbon_withdrawal: None,
            aircraft_ghg_emissions: None,
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![])),
            biogenic_accounting_methodology: None,
            boundary_processes_description: String::from("SFC GLEC Framework-conforming (W2W CO2e emissions)"),
            reference_period_start: Utc.with_ymd_and_hms(2023, 1, 1, 0, 0, 0).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2024, 1, 1, 0, 0, 0).unwrap(),
            geographic_scope: None,
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.9.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: None,
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(100.0)),
            dqi: None,
            assurance: None
        },
        extensions: Some(vec![DataModelExtension {
            spec_version: SpecVersionString::from("2.0.0".to_string()),
            data_schema: "https://api.ileap.sine.dev/hoc.json".to_string(),
            documentation: Some("https://sine-fdn.github.io/ileap-extension/".to_string()),
            data: ILeapType::Hoc(hoc(HocArgs {
                hoc_id:  "hoc-transshipment-1".to_string(),
                hub_type: HubType::Transshipment,
                temperature_control: Some(TemperatureControl::Refrigerated),
                inbound_transport_mode: Some(TransportMode::Road),
                outbound_transport_mode: Some(TransportMode::Rail),
            }))
        }])
    };


    static ref HOC_WAREHOUSE_EXAMPLE: ProductFootprint<ILeapType> = ProductFootprint {
        id: PfId(uuid!("2cf0a291-ad7f-43df-b0d7-b9554b0a9a02")),
        spec_version: SpecVersionString::from("2.0.0".to_string()),
        preceding_pf_ids: None,
        version: VersionInteger(0),
        created: Utc.with_ymd_and_hms(2024, 5, 22, 21, 47, 32).unwrap(),
        updated: None,
        status: PfStatus::Active,
        status_comment: None,
        validity_period_start: None,
        validity_period_end: None,
        company_name: String::from("Hub Co.").into(),
        company_ids: CompanyIdSet(vec![Urn::from("urn:epc:id:sgln:4063973.00000.9".to_string())]),
        product_description: String::from("Logistics emissions related to HOC with ID hoc-warehouse-1"),
        product_ids: ProductIdSet(vec![Urn::from("urn:pathfinder:product:customcode:vendor-assigned:hoc:hoc-warehouse-1".to_string())]),
        product_category_cpc: String::from("83117").into(),
        product_name_company:String::from("HOC with ID hoc-warehouse-1").into(),
        comment: "".to_string(),
        pcf: CarbonFootprint {
            declared_unit: DeclaredUnit::Kilogram,
            unitary_product_amount: dec!(1000).into(),
            p_cf_excluding_biogenic: dec!(33).into(),
            p_cf_including_biogenic: None,
            fossil_ghg_emissions: dec!(33).into(),
            fossil_carbon_content: dec!(33).into(),
            biogenic_carbon_content: dec!(0).into(),
            d_luc_ghg_emissions: None,
            land_management_ghg_emissions: None,
            other_biogenic_ghg_emissions: None,
            i_luc_ghg_emissions: None,
            biogenic_carbon_withdrawal: None,
            aircraft_ghg_emissions: None,
            characterization_factors: CharacterizationFactors::Ar6,
            ipcc_characterization_factors_sources: IpccCharacterizationFactorsSources::from(vec![String::from("AR6").into()]),
            cross_sectoral_standards_used: CrossSectoralStandardSet(vec![DeprecatedCrossSectoralStandard::Ghgp]),
            product_or_sector_specific_rules: Some(ProductOrSectorSpecificRuleSet(vec![])),
            biogenic_accounting_methodology: None,
            boundary_processes_description: String::from("SFC GLEC Framework-conforming (W2W CO2e emissions)"),
            reference_period_start: Utc.with_ymd_and_hms(2023, 1, 1, 0, 0, 0).unwrap(),
            reference_period_end: Utc.with_ymd_and_hms(2024, 1, 1, 0, 0, 0).unwrap(),
            geographic_scope: None,
            secondary_emission_factor_sources: Some(EmissionFactorDSSet(vec![EmissionFactorDS {
                name: String::from("Ecoinvent").into(),
                version: String::from("3.9.1").into(),
            }])),
            exempted_emissions_percent: ExemptedEmissionsPercent(0.0),
            exempted_emissions_description: "".to_string(),
            packaging_emissions_included: false,
            packaging_ghg_emissions: None,
            allocation_rules_description: None,
            uncertainty_assessment_description: None,
            primary_data_share: Some(Percent::from(100.0)),
            dqi: None,
            assurance: None
        },
        extensions: Some(vec![DataModelExtension {
            spec_version: SpecVersionString::from("2.0.0".to_string()),
            data_schema: "https://api.ileap.sine.dev/hoc.json".to_string(),
            documentation: Some("https://sine-fdn.github.io/ileap-extension/".to_string()),
            data: ILeapType::Hoc(
                hoc(HocArgs {
                hoc_id:  "hoc-warehouse-1".to_string(),
                hub_type: HubType::Warehouse,
                temperature_control: Some(TemperatureControl::Ambient),
                inbound_transport_mode: Some(TransportMode::Road),
                outbound_transport_mode: Some(TransportMode::Road),
            }))
        }])
    };

);

struct ShipmentArgs {
    mass: String,
    shipment_id: String,
    tces: NonEmptyVec<Tce>,
}

fn shipment_footprint(
    ShipmentArgs {
        mass,
        shipment_id,
        tces,
    }: ShipmentArgs,
) -> ShipmentFootprint {
    ShipmentFootprint {
        mass,
        volume: None,
        shipment_id,
        tces,
    }
}

struct TocArgs {
    toc_id: String,
    mode: TransportMode,
    load_factor: Option<String>,
    empty_distance_factor: Option<String>,
    temperature_control: Option<TemperatureControl>,
    truck_loading_sequence: Option<TruckLoadingSequence>,
    energy_carriers: NonEmptyVec<EnergyCarrier>,
    co2e_intensity_wtw: WrappedDecimal,
    co2e_intensity_ttw: WrappedDecimal,
    transport_activity_unit: TransportActivityUnit,
}

fn toc(
    TocArgs {
        toc_id,
        mode,
        load_factor,
        empty_distance_factor,
        temperature_control,
        truck_loading_sequence,
        energy_carriers,
        co2e_intensity_wtw,
        co2e_intensity_ttw,
        transport_activity_unit,
    }: TocArgs,
) -> Toc {
    Toc {
        toc_id,
        certifications: None,
        description: None,
        mode,
        load_factor,
        empty_distance_factor,
        temperature_control,
        truck_loading_sequence,
        air_shipping_option: None,
        flight_length: None,
        energy_carriers,
        co2e_intensity_wtw,
        co2e_intensity_ttw,
        transport_activity_unit,
    }
}

struct HocArgs {
    hoc_id: String,
    hub_type: HubType,
    temperature_control: Option<TemperatureControl>,
    inbound_transport_mode: Option<TransportMode>,
    outbound_transport_mode: Option<TransportMode>,
}

fn hoc(
    HocArgs {
        hoc_id,
        hub_type,
        temperature_control,
        inbound_transport_mode,
        outbound_transport_mode,
    }: HocArgs,
) -> Hoc {
    Hoc {
        hoc_id,
        description: None,
        certifications: None,
        hub_type,
        temperature_control,
        hub_location: None,
        inbound_transport_mode,
        outbound_transport_mode,
        packaging_or_tr_eq_type: None,
        packaging_or_tr_eq_amount: None,
        energy_carriers: vec![
            EnergyCarrier {
                energy_carrier: EnergyCarrierType::Diesel,
                feedstocks: None,
                energy_consumption: None,
                energy_consumption_unit: Some(EnergyConsumptionUnit::Kg),
                emission_factor_wtw: dec!(4.13).into(),
                emission_factor_ttw: dec!(3.17).into(),
                relative_share: dec!(0.5).into(),
            },
            EnergyCarrier {
                energy_carrier: EnergyCarrierType::Electric,
                feedstocks: None,
                energy_consumption: None,
                energy_consumption_unit: Some(EnergyConsumptionUnit::MJ),
                emission_factor_wtw: dec!(97).into(),
                emission_factor_ttw: dec!(0).into(),
                relative_share: dec!(0.5).into(),
            },
        ]
        .into(),
        co2e_intensity_wtw: dec!(33).into(),
        co2e_intensity_ttw: dec!(10).into(),
        hub_activity_unit: HubActivityUnit::Tonnes,
    }
}

struct TceArgs {
    tce_id: String,
    prev_tce_ids: Option<Vec<String>>,
    toc_id: Option<String>,
    hoc_id: Option<String>,
    shipment_id: String,
    mass: WrappedDecimal,
    distance: GlecDistance,
    transport_activity: WrappedDecimal,
    co2e_wtw: WrappedDecimal,
    co2e_ttw: WrappedDecimal,
}

fn tce(
    TceArgs {
        tce_id,
        prev_tce_ids,
        toc_id,
        hoc_id,
        shipment_id,
        mass,
        distance,
        transport_activity,
        co2e_wtw,
        co2e_ttw,
    }: TceArgs,
) -> Tce {
    Tce {
        tce_id,
        prev_tce_ids,
        toc_id,
        hoc_id,
        shipment_id,
        consignment_id: None,
        mass,
        packaging_or_tr_eq_type: None,
        packaging_or_tr_eq_amount: None,
        distance,
        origin: None,
        destination: None,
        transport_activity,
        departure_at: None,
        arrival_at: None,
        flight_no: None,
        voyage_no: None,
        incoterms: None,
        co2e_wtw,
        co2e_ttw,
        nox_ttw: None,
        sox_ttw: None,
        ch4_ttw: None,
        pm_ttw: None,
    }
}

lazy_static! {
    pub(crate) static ref PCF_DEMO_DATA: Vec<ProductFootprint<ILeapType>> = vec![
        EXAMPLE_1.clone(),
        EXAMPLE_2.clone(),
        EXAMPLE_3.clone(),
        EXAMPLE_4.clone(),
        SHIPMENT_SIMPLE_EXAMPLE.clone(),
        SHIPMENT_MULTIMODAL_EXAMPLE.clone(),
        SHIPMENT_PRE_AND_POST_LEGS_EXAMPLE.clone(),
        TOC_RAIL_EXAMPLE.clone(),
        TOC_ROAD_EXAMPLE.clone(),
        TOC_ROAD_40T_EXAMPLE.clone(),
        TOC_SMALL_TRUCK_EXAMPLE.clone(),
        HOC_TRANSSHIPMENT_EXAMPLE.clone(),
        HOC_WAREHOUSE_EXAMPLE.clone(),
    ];
}

fn demo_tad_base(activity_id: String, consignment_id: String, feedstock: FeedstockType) -> Tad {
    Tad {
        activity_id,
        consignment_ids: vec![consignment_id],
        distance: GlecDistance::new_actual(dec!(656.0).into()),
        mass: Some(dec!(1000.0).into()),
        origin: Location {
            street: Some("Bredowstraße".into()),
            zip: Some("10551".into()),
            city: "Berlin".into(),
            country: ISO3166CC(String::from("DE")),
            iata: None,
            locode: Some(Locode("DEBER".into())),
            uic: None,
            lat: Some(dec!(52.52437).into()),
            lng: Some(dec!(13.41053).into()),
        },
        destination: Location {
            street: Some("Keizersgracht".into()),
            zip: Some("1017".into()),
            city: "Amsterdam".into(),
            country: ISO3166CC(String::from("NL")),
            iata: None,
            locode: Some(Locode("NLAMS".into())),
            uic: None,
            lat: Some(dec!(52.37403).into()),
            lng: Some(dec!(4.88969).into()),
        },
        departure_at: Utc::now(),
        arrival_at: Utc::now() + Duration::days(10),
        mode: TransportMode::Road,
        packaging_or_tr_eq_type: Some(PackagingOrTrEqType::Pallet),
        packaging_or_tr_eq_amount: Some(10),
        // energy_carrier: EnergyCarrier {
        //     energy_carrier: "Diesel".into(),
        //     feedstocks: Some(vec![Feedstock {
        //         feedstock,
        //         feedstock_share: Some(1.0),
        //     }]),
        //     energy_consumption: Some(dec!(10.0).into()),
        //     energy_consumption_unit: Some("kg".into()),
        //     co2e_intensity_wtw: dec!(10.0).into(),
        //     co2e_intensity_ttw: dec!(10.0).into(),
        // },
        load_factor: Some(dec!(0.8).into()),
        empty_distance_factor: Some(dec!(0.1).into()),
        energy_carriers: Some(NonEmptyVec::from(vec![EnergyCarrier {
            energy_carrier: EnergyCarrierType::Diesel,
            feedstocks: Some(vec![Feedstock {
                feedstock,
                feedstock_share: Some(WrappedDecimal(dec!(1.0))),
                region_provenance: Some("Europe".to_string()),
            }]),
            energy_consumption: Some(WrappedDecimal(dec!(10.496))),
            energy_consumption_unit: Some(EnergyConsumptionUnit::L),
            emission_factor_wtw: dec!(4.13).into(),
            emission_factor_ttw: dec!(3.17).into(),
            relative_share: dec!(1.0).into(),
        }])),
        temperature_control: None,
    }
}

lazy_static! {
    pub(crate) static ref ILEAP_TAD_DEMO_DATA: Vec<Tad> = {
        let mut demo_data = vec![];
        for i in 1..10 {
            demo_data.push(demo_tad_base(
                i.to_string(),
                i.to_string(),
                FeedstockType::Fossil,
            ));
        }
        demo_data.push(demo_tad_base(
            "10".to_string(),
            "10".to_string(),
            FeedstockType::CookingOil,
        ));
        demo_data
    };
}

#[test]
fn test_skip_serializing_if_none() {
    for json in PCF_DEMO_DATA.iter() {
        let serialized = serde_json::to_string(&json).unwrap();
        assert!(!serialized.contains("null"));
    }

    for json in ILEAP_TAD_DEMO_DATA.iter() {
        let serialized = serde_json::to_string(&json).unwrap();
        assert!(!serialized.contains("null"));
    }
}
