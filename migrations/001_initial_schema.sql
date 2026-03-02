-- ============================================================================
-- GENAPP Migration: CICS/COBOL to PostgreSQL
-- Migration 001: Initial Schema
--
-- This schema maps the COBOL copybook structures (lgcmarea.cpy, lgpolicy.cpy)
-- to PostgreSQL tables. Field sizes are derived from COBOL PIC definitions.
-- ============================================================================

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- SEQUENCES
-- These replace CICS Named Counter Server functionality
-- Starting at 1000000001 to match COBOL PIC 9(10) with leading zeros
-- ============================================================================

CREATE SEQUENCE customer_num_seq START 1000000001 MAXVALUE 9999999999;
CREATE SEQUENCE policy_num_seq START 1000000001 MAXVALUE 9999999999;
CREATE SEQUENCE claim_num_seq START 1000000001 MAXVALUE 9999999999;

-- ============================================================================
-- CUSTOMERS TABLE
-- Maps from DB2-CUSTOMER structure in lgpolicy.cpy
--
-- COBOL field mappings:
--   DB2-FIRSTNAME     PIC X(10)   -> first_name VARCHAR(10)
--   DB2-LASTNAME      PIC X(20)   -> last_name VARCHAR(20)
--   DB2-DATEOFBIRTH   PIC X(10)   -> date_of_birth DATE
--   DB2-HOUSENAME     PIC X(20)   -> house_name VARCHAR(20)
--   DB2-HOUSENUMBER   PIC X(4)    -> house_number VARCHAR(4)
--   DB2-POSTCODE      PIC X(8)    -> postcode VARCHAR(8)
--   DB2-PHONE-MOBILE  PIC X(20)   -> phone_mobile VARCHAR(20)
--   DB2-PHONE-HOME    PIC X(20)   -> phone_home VARCHAR(20)
--   DB2-EMAIL-ADDRESS PIC X(100)  -> email_address VARCHAR(100)
-- ============================================================================

CREATE TABLE customers (
    id              SERIAL PRIMARY KEY,
    customer_num    VARCHAR(10) UNIQUE NOT NULL,  -- PIC 9(10)
    first_name      VARCHAR(10),                  -- PIC X(10)
    last_name       VARCHAR(20),                  -- PIC X(20)
    date_of_birth   DATE,                         -- PIC X(10) as date
    house_name      VARCHAR(20),                  -- PIC X(20)
    house_number    VARCHAR(4),                   -- PIC X(4)
    postcode        VARCHAR(8),                   -- PIC X(8)
    phone_home      VARCHAR(20),                  -- PIC X(20)
    phone_mobile    VARCHAR(20),                  -- PIC X(20)
    email_address   VARCHAR(100),                 -- PIC X(100)
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_customers_customer_num ON customers(customer_num);
CREATE INDEX idx_customers_last_name ON customers(last_name);

-- ============================================================================
-- POLICIES TABLE (Master)
-- Maps from DB2-POLICY structure in lgpolicy.cpy
--
-- COBOL field mappings:
--   DB2-POLICYTYPE    PIC X       -> policy_type CHAR(1)
--   DB2-POLICYNUMBER  PIC 9(10)   -> policy_num VARCHAR(10)
--   DB2-ISSUEDATE     PIC X(10)   -> issue_date DATE
--   DB2-EXPIRYDATE    PIC X(10)   -> expiry_date DATE
--   DB2-LASTCHANGED   PIC X(26)   -> last_changed TIMESTAMP
--   DB2-BROKERID      PIC 9(10)   -> broker_id VARCHAR(10)
--   DB2-BROKERSREF    PIC X(10)   -> brokers_ref VARCHAR(10)
--   DB2-PAYMENT       PIC 9(6)    -> payment DECIMAL(8,2)
-- ============================================================================

CREATE TABLE policies (
    id              SERIAL PRIMARY KEY,
    policy_num      VARCHAR(10) UNIQUE NOT NULL,  -- PIC 9(10)
    customer_num    VARCHAR(10) NOT NULL,         -- PIC 9(10)
    policy_type     CHAR(1) NOT NULL,             -- PIC X: E=Endowment, H=House, M=Motor, C=Commercial
    issue_date      DATE,                         -- PIC X(10)
    expiry_date     DATE,                         -- PIC X(10)
    last_changed    TIMESTAMP,                    -- PIC X(26)
    broker_id       VARCHAR(10),                  -- PIC 9(10)
    brokers_ref     VARCHAR(10),                  -- PIC X(10)
    payment         DECIMAL(8,2),                 -- PIC 9(6) with implied decimal
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_policies_customer
        FOREIGN KEY (customer_num) REFERENCES customers(customer_num) ON DELETE CASCADE,
    CONSTRAINT chk_policy_type
        CHECK (policy_type IN ('E', 'H', 'M', 'C'))
);

CREATE INDEX idx_policies_policy_num ON policies(policy_num);
CREATE INDEX idx_policies_customer_num ON policies(customer_num);
CREATE INDEX idx_policies_type ON policies(policy_type);

-- ============================================================================
-- MOTOR POLICIES TABLE
-- Maps from DB2-MOTOR structure in lgpolicy.cpy
--
-- COBOL field mappings:
--   DB2-M-MAKE         PIC X(15)  -> make VARCHAR(15)
--   DB2-M-MODEL        PIC X(15)  -> model VARCHAR(15)
--   DB2-M-VALUE        PIC 9(6)   -> value DECIMAL(8,2)
--   DB2-M-REGNUMBER    PIC X(7)   -> reg_number VARCHAR(7)
--   DB2-M-COLOUR       PIC X(8)   -> colour VARCHAR(8)
--   DB2-M-CC           PIC 9(4)   -> cc INTEGER
--   DB2-M-MANUFACTURED PIC X(10)  -> manufactured DATE
--   DB2-M-PREMIUM      PIC 9(6)   -> premium DECIMAL(8,2)
--   DB2-M-ACCIDENTS    PIC 9(6)   -> accidents INTEGER
-- ============================================================================

CREATE TABLE motor_policies (
    id              SERIAL PRIMARY KEY,
    policy_num      VARCHAR(10) UNIQUE NOT NULL,  -- PIC 9(10)
    make            VARCHAR(15),                  -- PIC X(15)
    model           VARCHAR(15),                  -- PIC X(15)
    value           DECIMAL(8,2),                 -- PIC 9(6)
    reg_number      VARCHAR(7),                   -- PIC X(7)
    colour          VARCHAR(8),                   -- PIC X(8)
    cc              INTEGER,                      -- PIC 9(4)
    manufactured    DATE,                         -- PIC X(10) as date
    premium         DECIMAL(8,2),                 -- PIC 9(6)
    accidents       INTEGER DEFAULT 0,            -- PIC 9(6)
    CONSTRAINT fk_motor_policy
        FOREIGN KEY (policy_num) REFERENCES policies(policy_num) ON DELETE CASCADE
);

CREATE INDEX idx_motor_policies_policy_num ON motor_policies(policy_num);

-- ============================================================================
-- ENDOWMENT POLICIES TABLE
-- Maps from DB2-ENDOWMENT structure in lgpolicy.cpy
--
-- COBOL field mappings:
--   DB2-E-WITHPROFITS  PIC X      -> with_profits BOOLEAN
--   DB2-E-EQUITIES     PIC X      -> equities BOOLEAN
--   DB2-E-MANAGEDFUND  PIC X      -> managed_fund BOOLEAN
--   DB2-E-FUNDNAME     PIC X(10)  -> fund_name VARCHAR(10)
--   DB2-E-TERM         PIC 9(2)   -> term INTEGER
--   DB2-E-SUMASSURED   PIC 9(6)   -> sum_assured DECIMAL(10,2)
--   DB2-E-LIFEASSURED  PIC X(31)  -> life_assured VARCHAR(31)
-- ============================================================================

CREATE TABLE endowment_policies (
    id              SERIAL PRIMARY KEY,
    policy_num      VARCHAR(10) UNIQUE NOT NULL,  -- PIC 9(10)
    with_profits    BOOLEAN DEFAULT FALSE,        -- PIC X (Y/N)
    equities        BOOLEAN DEFAULT FALSE,        -- PIC X (Y/N)
    managed_fund    BOOLEAN DEFAULT FALSE,        -- PIC X (Y/N)
    fund_name       VARCHAR(10),                  -- PIC X(10)
    term            INTEGER,                      -- PIC 9(2)
    sum_assured     DECIMAL(10,2),                -- PIC 9(6)
    life_assured    VARCHAR(31),                  -- PIC X(31)
    CONSTRAINT fk_endowment_policy
        FOREIGN KEY (policy_num) REFERENCES policies(policy_num) ON DELETE CASCADE
);

CREATE INDEX idx_endowment_policies_policy_num ON endowment_policies(policy_num);

-- ============================================================================
-- HOUSE POLICIES TABLE
-- Maps from DB2-HOUSE structure in lgpolicy.cpy
--
-- COBOL field mappings:
--   DB2-H-PROPERTYTYPE PIC X(15)  -> property_type VARCHAR(15)
--   DB2-H-BEDROOMS     PIC 9(3)   -> bedrooms INTEGER
--   DB2-H-VALUE        PIC 9(8)   -> value DECIMAL(10,2)
--   DB2-H-HOUSENAME    PIC X(20)  -> house_name VARCHAR(20)
--   DB2-H-HOUSENUMBER  PIC X(4)   -> house_number VARCHAR(4)
--   DB2-H-POSTCODE     PIC X(8)   -> postcode VARCHAR(8)
-- ============================================================================

CREATE TABLE house_policies (
    id              SERIAL PRIMARY KEY,
    policy_num      VARCHAR(10) UNIQUE NOT NULL,  -- PIC 9(10)
    property_type   VARCHAR(15),                  -- PIC X(15)
    bedrooms        INTEGER,                      -- PIC 9(3)
    value           DECIMAL(10,2),                -- PIC 9(8)
    house_name      VARCHAR(20),                  -- PIC X(20)
    house_number    VARCHAR(4),                   -- PIC X(4)
    postcode        VARCHAR(8),                   -- PIC X(8)
    CONSTRAINT fk_house_policy
        FOREIGN KEY (policy_num) REFERENCES policies(policy_num) ON DELETE CASCADE
);

CREATE INDEX idx_house_policies_policy_num ON house_policies(policy_num);

-- ============================================================================
-- COMMERCIAL POLICIES TABLE
-- Maps from DB2-COMMERCIAL structure in lgpolicy.cpy
--
-- COBOL field mappings:
--   DB2-B-Address        PIC X(255) -> address TEXT
--   DB2-B-Postcode       PIC X(8)   -> postcode VARCHAR(8)
--   DB2-B-Latitude       PIC X(11)  -> latitude VARCHAR(11)
--   DB2-B-Longitude      PIC X(11)  -> longitude VARCHAR(11)
--   DB2-B-Customer       PIC X(255) -> customer TEXT
--   DB2-B-PropType       PIC X(255) -> property_type TEXT
--   DB2-B-FirePeril      PIC 9(4)   -> fire_peril INTEGER
--   DB2-B-FirePremium    PIC 9(8)   -> fire_premium DECIMAL(10,2)
--   DB2-B-CrimePeril     PIC 9(4)   -> crime_peril INTEGER
--   DB2-B-CrimePremium   PIC 9(8)   -> crime_premium DECIMAL(10,2)
--   DB2-B-FloodPeril     PIC 9(4)   -> flood_peril INTEGER
--   DB2-B-FloodPremium   PIC 9(8)   -> flood_premium DECIMAL(10,2)
--   DB2-B-WeatherPeril   PIC 9(4)   -> weather_peril INTEGER
--   DB2-B-WeatherPremium PIC 9(8)   -> weather_premium DECIMAL(10,2)
--   DB2-B-Status         PIC 9(4)   -> status INTEGER
--   DB2-B-RejectReason   PIC X(255) -> reject_reason TEXT
-- ============================================================================

CREATE TABLE commercial_policies (
    id              SERIAL PRIMARY KEY,
    policy_num      VARCHAR(10) UNIQUE NOT NULL,  -- PIC 9(10)
    address         TEXT,                         -- PIC X(255)
    postcode        VARCHAR(8),                   -- PIC X(8)
    latitude        VARCHAR(11),                  -- PIC X(11)
    longitude       VARCHAR(11),                  -- PIC X(11)
    customer        TEXT,                         -- PIC X(255)
    property_type   TEXT,                         -- PIC X(255)
    fire_peril      INTEGER,                      -- PIC 9(4)
    fire_premium    DECIMAL(10,2),                -- PIC 9(8)
    crime_peril     INTEGER,                      -- PIC 9(4)
    crime_premium   DECIMAL(10,2),                -- PIC 9(8)
    flood_peril     INTEGER,                      -- PIC 9(4)
    flood_premium   DECIMAL(10,2),                -- PIC 9(8)
    weather_peril   INTEGER,                      -- PIC 9(4)
    weather_premium DECIMAL(10,2),                -- PIC 9(8)
    status          INTEGER,                      -- PIC 9(4)
    reject_reason   TEXT,                         -- PIC X(255)
    CONSTRAINT fk_commercial_policy
        FOREIGN KEY (policy_num) REFERENCES policies(policy_num) ON DELETE CASCADE
);

CREATE INDEX idx_commercial_policies_policy_num ON commercial_policies(policy_num);

-- ============================================================================
-- CLAIMS TABLE
-- Maps from DB2-CLAIM structure in lgpolicy.cpy
--
-- COBOL field mappings:
--   DB2-C-Num          PIC 9(10)  -> claim_num VARCHAR(10)
--   DB2-C-Date         PIC X(10)  -> claim_date DATE
--   DB2-C-Paid         PIC 9(8)   -> paid DECIMAL(10,2)
--   DB2-C-Value        PIC 9(8)   -> value DECIMAL(10,2)
--   DB2-C-Cause        PIC X(255) -> cause TEXT
--   DB2-C-Observations PIC X(255) -> observations TEXT
-- ============================================================================

CREATE TABLE claims (
    id              SERIAL PRIMARY KEY,
    claim_num       VARCHAR(10) UNIQUE NOT NULL,  -- PIC 9(10)
    policy_num      VARCHAR(10) NOT NULL,         -- PIC 9(10)
    claim_date      DATE,                         -- PIC X(10)
    paid            DECIMAL(10,2),                -- PIC 9(8)
    value           DECIMAL(10,2),                -- PIC 9(8)
    cause           TEXT,                         -- PIC X(255)
    observations    TEXT,                         -- PIC X(255)
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_claims_policy
        FOREIGN KEY (policy_num) REFERENCES policies(policy_num) ON DELETE CASCADE
);

CREATE INDEX idx_claims_claim_num ON claims(claim_num);
CREATE INDEX idx_claims_policy_num ON claims(policy_num);

-- ============================================================================
-- COUNTERS TABLE
-- Replaces CICS Named Counter Server functionality
-- Used for tracking application-level counters and statistics
-- ============================================================================

CREATE TABLE counters (
    name            VARCHAR(50) PRIMARY KEY,
    value           BIGINT DEFAULT 0,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- HELPER FUNCTIONS
-- ============================================================================

-- Function to update timestamp on row modification
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply auto-update triggers
CREATE TRIGGER update_customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_policies_updated_at
    BEFORE UPDATE ON policies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_claims_updated_at
    BEFORE UPDATE ON claims
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_counters_updated_at
    BEFORE UPDATE ON counters
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- HELPER FUNCTIONS FOR COUNTER SEQUENCES
-- These functions provide atomic counter increment similar to CICS Named Counters
-- ============================================================================

-- Get next customer number (formatted with leading zeros)
CREATE OR REPLACE FUNCTION next_customer_num()
RETURNS VARCHAR(10) AS $$
BEGIN
    RETURN LPAD(nextval('customer_num_seq')::TEXT, 10, '0');
END;
$$ LANGUAGE plpgsql;

-- Get next policy number (formatted with leading zeros)
CREATE OR REPLACE FUNCTION next_policy_num()
RETURNS VARCHAR(10) AS $$
BEGIN
    RETURN LPAD(nextval('policy_num_seq')::TEXT, 10, '0');
END;
$$ LANGUAGE plpgsql;

-- Get next claim number (formatted with leading zeros)
CREATE OR REPLACE FUNCTION next_claim_num()
RETURNS VARCHAR(10) AS $$
BEGIN
    RETURN LPAD(nextval('claim_num_seq')::TEXT, 10, '0');
END;
$$ LANGUAGE plpgsql;

-- Atomic increment for named counters (like CICS Named Counter GET with UPDATE)
CREATE OR REPLACE FUNCTION increment_counter(counter_name VARCHAR(50), increment_by BIGINT DEFAULT 1)
RETURNS BIGINT AS $$
DECLARE
    new_value BIGINT;
BEGIN
    INSERT INTO counters (name, value)
    VALUES (counter_name, increment_by)
    ON CONFLICT (name) DO UPDATE
    SET value = counters.value + increment_by,
        updated_at = CURRENT_TIMESTAMP
    RETURNING value INTO new_value;

    RETURN new_value;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- INITIAL COUNTER VALUES
-- Initialize counters used for statistics (replaces LGSETUP functionality)
-- ============================================================================

INSERT INTO counters (name, value) VALUES
    ('customer_add_count', 0),
    ('customer_inq_count', 0),
    ('customer_upd_count', 0),
    ('policy_add_count', 0),
    ('policy_inq_count', 0),
    ('policy_upd_count', 0),
    ('policy_del_count', 0),
    ('claim_add_count', 0),
    ('total_transactions', 0);
