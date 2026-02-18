-- **DISCLAIMER**
--
-- This file was crated on a best-effort basis to provide an example of a possible database schema for the ILEAP project.
-- The SQL definitions are expected to be incomplete, errorneous.
-- The schema should be reviewed and validated by the community first.
--
-- Table for the iLEAP ShipmentFootprint Data Type
CREATE TABLE shipmentfootprint(
    mass varchar(255) NOT NULL,
    volume varchar(255),
    shipment_id varchar(255) PRIMARY KEY
);

-- Table for the iLEAP TOC Data Type
CREATE TABLE toc(
    toc_id varchar(255) PRIMARY KEY,
    certifications TEXT[], -- 'ISO14083:2023', 'GLECv2', 'GLECv3', or 'GLECv3.1'
    description text,
    mode varchar(20) NOT NULL CHECK (mode IN ('Road', 'Rail', 'Air', 'Sea', 'InlandWaterway')),
    load_factor DECIMAL(18, 6),
    empty_distance_factor varchar(255),
    temperature_control varchar(20) CHECK (temperature_control IN ('ambient', 'refrigerated', 'mixed')),
    truck_loading_sequence varchar(3) CHECK (truck_loading_sequence IN ('LTL', 'FTL')),
    air_shipping_option varchar(20) CHECK (air_shipping_option IN ('belly freight', 'freighter')),
    flight_length varchar(20) CHECK (flight_length IN ('short-haul', 'long-haul')),
    co2e_intensity_wtw DECIMAL(18, 6) NOT NULL,
    co2e_intensity_ttw DECIMAL(18, 6) NOT NULL,
    transport_activity_unit varchar(255) NOT NULL
);

-- Table for the iLEAP HOC Data Type
CREATE TABLE hoc(
    hoc_id varchar(255) PRIMARY KEY,
    description text,
    certifications TEXT[], -- 'ISO14083:2023', 'GLECv2', 'GLECv3', or 'GLECv3.1'
    hub_type varchar(50) NOT NULL CHECK (hub_type IN ('Transshipment', 'StorageAndTransshipment', 'Warehouse', 'LiquidBulkTerminal', 'MaritimeContainerTerminal')),
    temperature_control varchar(50) CHECK (temperature_control IN ('ambient', 'refrigerated', 'mixed')),
    hub_location_street varchar(255),
    hub_location_zip varchar(255),
    hub_location_city varchar(255) NOT NULL,
    hub_location_country varchar(255),
    hub_location_iata varchar(3),
    hub_location_locode varchar(5),
    hub_location_uic varchar(2),
    hub_location_lat DECIMAL(10, 8),
    hub_location_lng DECIMAL(11, 8),
    inbound_transport_mode varchar(50) CHECK (inbound_transport_mode IN ('Road', 'Rail', 'Air', 'Sea', 'InlandWaterway')),
    outbound_transport_mode varchar(50) CHECK (outbound_transport_mode IN ('Road', 'Rail', 'Air', 'Sea', 'InlandWaterway')),
    packaging_or_tr_eq_type varchar(255) CHECK (packaging_or_tr_eq_type IN ('Box', 'Pallet', 'Container-TEU', 'Container-FEU', 'Container')),
    packaging_or_tr_eq_amount int,
    co2e_intensity_wtw DECIMAL(18, 6) NOT NULL,
    co2e_intensity_ttw DECIMAL(18, 6) NOT NULL,
    hub_activity_unit varchar(255) NOT NULL
);

-- Table for the iLEAP TCE Data Type
CREATE TABLE tce(
    tce_id varchar(255) PRIMARY KEY,
    toc_id varchar(255),
    hoc_id varchar(255),
    shipment_id varchar(255) NOT NULL,
    consignment_id varchar(255),
    mass DECIMAL(18, 6) NOT NULL,
    packaging_or_tr_eq_type varchar(255) CHECK (packaging_or_tr_eq_type IN ('Box', 'Pallet', 'Container')),
    packaging_or_tr_eq_amount int,
    distance_type varchar(50) NOT NULL CHECK (distance_type IN ('actual', 'gcd', 'sfd')),
    distance_value DECIMAL(18, 6) NOT NULL,
    origin_street varchar(255),
    origin_zip varchar(255),
    origin_city varchar(255) NOT NULL,
    origin_country varchar(255) NOT NULL,
    origin_iata varchar(3),
    origin_locode varchar(5),
    origin_uic varchar(2),
    origin_lat DECIMAL(10, 8),
    origin_lng DECIMAL(11, 8),
    destination_street varchar(255),
    destination_zip varchar(255),
    destination_city varchar(255) NOT NULL,
    destination_country varchar(255) NOT NULL,
    destination_iata varchar(3),
    destination_locode varchar(5),
    destination_uic varchar(2),
    destination_lat DECIMAL(10, 8),
    destination_lng DECIMAL(11, 8),
    transport_activity DECIMAL(18, 6) NOT NULL,
    departure_at timestamp,
    arrival_at timestamp,
    flight_no varchar(255),
    voyage_no varchar(255),
    incoterms varchar(3) CHECK (incoterms IN ('EXW', 'FCA', 'CPT', 'CIP', 'DAP', 'DPU', 'DDP', 'FAS', 'FOB', 'CFR', 'CIF')),
    co2e_wtw DECIMAL(18, 6) NOT NULL,
    co2e_ttw DECIMAL(18, 6) NOT NULL,
    nox_ttw DECIMAL(18, 6),
    sox_ttw DECIMAL(18, 6),
    ch4_ttw DECIMAL(18, 6),
    pm_ttw DECIMAL(18, 6),
    FOREIGN KEY (shipment_id) REFERENCES ShipmentFootprint(shipment_id),
    FOREIGN KEY (toc_id) REFERENCES toc(toc_id),
    FOREIGN KEY (hoc_id) REFERENCES hoc(hoc_id)
);

-- Table for TAD
CREATE TABLE tad(
    activity_id varchar(255) PRIMARY KEY,
    consignment_ids text, -- Array of consignment IDs stored as text to workaround sqlite limitation(s), e.g. `consignment_id_1,consignment_id_2`
    distance_type varchar(50) NOT NULL CHECK (distance_type IN ('actual', 'gcd', 'sfd')),
    distance_value DECIMAL(18, 6) NOT NULL,
    mass DECIMAL(18, 6),
    load_factor DECIMAL(18, 6),
    empty_distance_factor DECIMAL(18, 6),
    origin_street varchar(255),
    origin_zip varchar(255),
    origin_city varchar(255) NOT NULL,
    origin_country varchar(255) NOT NULL,
    origin_iata varchar(3),
    origin_locode varchar(5),
    origin_uic varchar(2),
    origin_lat DECIMAL(10, 8),
    origin_lng DECIMAL(11, 8),
    destination_street varchar(255),
    destination_zip varchar(255),
    destination_city varchar(255) NOT NULL,
    destination_country varchar(255) NOT NULL,
    destination_iata varchar(3),
    destination_locode varchar(5),
    destination_uic varchar(2),
    destination_lat DECIMAL(10, 8),
    destination_lng DECIMAL(11, 8),
    departure_at timestamp,
    arrival_at timestamp,
    mode varchar(50) CHECK (mode IN ('Road', 'Rail', 'Air', 'Sea', 'InlandWaterway')),
    packaging_or_tr_eq_type varchar(255) CHECK (packaging_or_tr_eq_type IN ('Box', 'Pallet', 'Container')),
    packaging_or_tr_eq_amount int,
    energy_carrier varchar(50) CHECK (energy_carrier IN ('Diesel', 'HVO', 'Petrol', 'CNG', 'LNG', 'LPG', 'HFO', 'MGO', 'Aviation fuel', 'Hydrogen', 'Methanol', 'Electric')),
    temperature_control varchar(50) CHECK (temperature_control IN ('ambient', 'refrigerated'))
);

-- Table for EnergyCarrier (assuming the existing structure for EnergyCarrier)
CREATE TABLE energycarrier(
    id serial PRIMARY KEY,
    energy_carrier_type varchar(50) NOT NULL CHECK (energy_carrier_type IN ('Diesel', 'HVO', 'Petrol', 'CNG', 'LNG', 'LPG', 'HFO', 'MGO', 'Aviation fuel', 'Hydrogen', 'Methanol', 'Electric')),
    energy_consumption DECIMAL(18, 6),
    energy_consumption_unit varchar(255),
    emission_factor_wtw DECIMAL(18, 6) NOT NULL,
    emission_factor_ttw DECIMAL(18, 6) NOT NULL,
    relative_share DECIMAL(5, 2) CHECK (relative_share > 0 AND relative_share <= 1.00) NOT NULL
);

-- Table for Feedstock (assuming the existing structure for Feedstock)
CREATE TABLE feedstock(
    id serial PRIMARY KEY,
    feedstock_type varchar(50) NOT NULL CHECK (feedstock_type IN ('Fossil', 'Natural gas', 'Grid', 'Renewable electricity', 'Cooking oil')),
    feedstock_share DECIMAL(5, 2),
    region_provenance varchar(255)
);

-- Additional join tables (assuming the existing structure for join tables)
CREATE TABLE tocenergycarrier(
    toc_id varchar(255) NOT NULL,
    energy_carrier_id int NOT NULL,
    PRIMARY KEY (toc_id, energy_carrier_id),
    FOREIGN KEY (toc_id) REFERENCES toc(toc_id),
    FOREIGN KEY (energy_carrier_id) REFERENCES energycarrier(id)
);

CREATE TABLE hocenergycarrier(
    hoc_id varchar(255) NOT NULL,
    energy_carrier_id int NOT NULL,
    PRIMARY KEY (hoc_id, energy_carrier_id),
    FOREIGN KEY (hoc_id) REFERENCES hoc(hoc_id),
    FOREIGN KEY (energy_carrier_id) REFERENCES energycarrier(id)
);

CREATE TABLE energycarrierfeedstock(
    energy_carrier_id int NOT NULL,
    feedstock_id int NOT NULL,
    PRIMARY KEY (energy_carrier_id, feedstock_id),
    FOREIGN KEY (energy_carrier_id) REFERENCES energycarrier(id),
    FOREIGN KEY (feedstock_id) REFERENCES feedstock(id)
);

CREATE INDEX 'energycarrierfeedstock_feedstock_id' ON 'energycarrierfeedstock'('feedstock_id');

CREATE INDEX 'hocenergycarrier_energy_carrier_id' ON 'hocenergycarrier'('energy_carrier_id');

CREATE INDEX 'tce_hoc_id' ON 'tce'('hoc_id');

CREATE INDEX 'tce_toc_id' ON 'tce'('toc_id');

CREATE INDEX 'tce_shipment_id' ON 'tce'('shipment_id');

CREATE INDEX 'tocenergycarrier_energy_carrier_id' ON 'tocenergycarrier'('energy_carrier_id');
