-- ============================================================================
-- GENAPP Seed Data
-- Test data for development and testing
--
-- This script populates the database with sample data matching
-- the types of records the COBOL application would handle.
-- ============================================================================

-- ============================================================================
-- CUSTOMERS
-- Sample customers with various profiles
-- ============================================================================

INSERT INTO customers (customer_num, first_name, last_name, date_of_birth, house_name, house_number, postcode, phone_home, phone_mobile, email_address) VALUES
    ('1000000001', 'John', 'Smith', '1985-03-15', 'Oak House', '42', 'AB12 3CD', '02012345678', '07123456789', 'john.smith@example.com'),
    ('1000000002', 'Jane', 'Doe', '1990-07-22', 'Rose Cottage', '15', 'EF45 6GH', '02023456789', '07234567890', 'jane.doe@example.com'),
    ('1000000003', 'Robert', 'Johnson', '1978-11-08', 'The Willows', '8', 'IJ78 9KL', '02034567890', '07345678901', 'r.johnson@example.com'),
    ('1000000004', 'Emily', 'Williams', '1995-01-30', '', '123', 'MN01 2OP', '02045678901', '07456789012', 'emily.w@example.com'),
    ('1000000005', 'Michael', 'Brown', '1982-06-18', 'Ivy Lodge', '7', 'QR34 5ST', '02056789012', '07567890123', 'mbrown@example.com'),
    ('1000000006', 'Sarah', 'Taylor', '1988-09-25', '', '256', 'UV67 8WX', '02067890123', '07678901234', 'sarah.taylor@example.com'),
    ('1000000007', 'David', 'Anderson', '1975-04-12', 'Elm House', '33', 'YZ90 1AB', '02078901234', '07789012345', 'david.a@example.com'),
    ('1000000008', 'Lisa', 'Thomas', '1992-12-05', '', '89', 'CD23 4EF', '02089012345', '07890123456', 'lisa.thomas@example.com'),
    ('1000000009', 'James', 'Jackson', '1980-08-20', 'Birch Villa', '12', 'GH56 7IJ', '02090123456', '07901234567', 'jjackson@example.com'),
    ('1000000010', 'Amanda', 'White', '1987-02-14', '', '45', 'KL89 0MN', '02001234567', '07012345678', 'a.white@example.com');

-- Update sequence to continue after seed data
SELECT setval('customer_num_seq', 1000000010, true);

-- ============================================================================
-- POLICIES - Motor Type (M)
-- ============================================================================

INSERT INTO policies (policy_num, customer_num, policy_type, issue_date, expiry_date, last_changed, broker_id, brokers_ref, payment) VALUES
    ('1000000001', '1000000001', 'M', '2024-01-01', '2025-01-01', CURRENT_TIMESTAMP, '0000000001', 'BRK001REF', 450.00),
    ('1000000002', '1000000002', 'M', '2024-02-15', '2025-02-15', CURRENT_TIMESTAMP, '0000000001', 'BRK002REF', 380.00),
    ('1000000003', '1000000003', 'M', '2024-03-20', '2025-03-20', CURRENT_TIMESTAMP, '0000000002', 'BRK003REF', 520.00);

INSERT INTO motor_policies (policy_num, make, model, value, reg_number, colour, cc, manufactured, premium, accidents) VALUES
    ('1000000001', 'Toyota', 'Camry', 25000.00, 'AB12CDE', 'Silver', 2000, '2023-06-01', 450.00, 0),
    ('1000000002', 'Honda', 'Civic', 18000.00, 'FG34HIJ', 'Blue', 1800, '2022-09-15', 380.00, 1),
    ('1000000003', 'Ford', 'Focus', 15000.00, 'KL56MNO', 'Red', 1600, '2021-03-20', 520.00, 2);

-- ============================================================================
-- POLICIES - Endowment Type (E)
-- ============================================================================

INSERT INTO policies (policy_num, customer_num, policy_type, issue_date, expiry_date, last_changed, broker_id, brokers_ref, payment) VALUES
    ('1000000004', '1000000004', 'E', '2024-01-10', '2044-01-10', CURRENT_TIMESTAMP, '0000000002', 'BRK004REF', 150.00),
    ('1000000005', '1000000005', 'E', '2024-04-01', '2049-04-01', CURRENT_TIMESTAMP, '0000000003', 'BRK005REF', 200.00);

INSERT INTO endowment_policies (policy_num, with_profits, equities, managed_fund, fund_name, term, sum_assured, life_assured) VALUES
    ('1000000004', TRUE, TRUE, FALSE, 'GROWTHFND', 20, 50000.00, 'Emily Williams'),
    ('1000000005', TRUE, FALSE, TRUE, 'BALANCED', 25, 75000.00, 'Michael Brown');

-- ============================================================================
-- POLICIES - House Type (H)
-- ============================================================================

INSERT INTO policies (policy_num, customer_num, policy_type, issue_date, expiry_date, last_changed, broker_id, brokers_ref, payment) VALUES
    ('1000000006', '1000000006', 'H', '2024-02-01', '2025-02-01', CURRENT_TIMESTAMP, '0000000001', 'BRK006REF', 280.00),
    ('1000000007', '1000000007', 'H', '2024-03-15', '2025-03-15', CURRENT_TIMESTAMP, '0000000002', 'BRK007REF', 350.00),
    ('1000000008', '1000000001', 'H', '2024-05-01', '2025-05-01', CURRENT_TIMESTAMP, '0000000001', 'BRK008REF', 420.00);

INSERT INTO house_policies (policy_num, property_type, bedrooms, value, house_name, house_number, postcode) VALUES
    ('1000000006', 'Semi-Detached', 3, 250000.00, '', '256', 'UV67 8WX'),
    ('1000000007', 'Detached', 4, 450000.00, 'Elm House', '33', 'YZ90 1AB'),
    ('1000000008', 'Terraced', 2, 180000.00, 'Oak House', '42', 'AB12 3CD');

-- ============================================================================
-- POLICIES - Commercial Type (C)
-- ============================================================================

INSERT INTO policies (policy_num, customer_num, policy_type, issue_date, expiry_date, last_changed, broker_id, brokers_ref, payment) VALUES
    ('1000000009', '1000000009', 'C', '2024-01-20', '2025-01-20', CURRENT_TIMESTAMP, '0000000003', 'BRK009REF', 1500.00),
    ('1000000010', '1000000010', 'C', '2024-04-10', '2025-04-10', CURRENT_TIMESTAMP, '0000000003', 'BRK010REF', 2200.00);

INSERT INTO commercial_policies (policy_num, address, postcode, latitude, longitude, customer, property_type, fire_peril, fire_premium, crime_peril, crime_premium, flood_peril, flood_premium, weather_peril, weather_premium, status, reject_reason) VALUES
    ('1000000009', '123 Business Park, Industrial Estate, London', 'E1 4AB', '51.5074', '-0.1278', 'ABC Trading Ltd', 'Warehouse', 3, 500.00, 2, 300.00, 1, 200.00, 2, 250.00, 1, ''),
    ('1000000010', '456 High Street, Manchester', 'M1 1AA', '53.4808', '-2.2426', 'XYZ Retail Corp', 'Retail Shop', 2, 400.00, 4, 600.00, 1, 150.00, 1, 100.00, 1, '');

-- Update sequence to continue after seed data
SELECT setval('policy_num_seq', 1000000010, true);

-- ============================================================================
-- CLAIMS
-- Sample claims against commercial policies
-- ============================================================================

INSERT INTO claims (claim_num, policy_num, claim_date, paid, value, cause, observations) VALUES
    ('1000000001', '1000000009', '2024-06-15', 5000.00, 8000.00, 'Fire damage to storage area', 'Partial payment made. Awaiting final assessment.'),
    ('1000000002', '1000000009', '2024-08-20', 0.00, 2500.00, 'Break-in and theft', 'Claim under investigation. Police report pending.'),
    ('1000000003', '1000000010', '2024-07-10', 1500.00, 1500.00, 'Storm damage to signage', 'Claim fully settled.');

-- Update sequence to continue after seed data
SELECT setval('claim_num_seq', 1000000003, true);

-- ============================================================================
-- UPDATE COUNTERS
-- Set initial statistics based on seed data
-- ============================================================================

UPDATE counters SET value = 10 WHERE name = 'customer_add_count';
UPDATE counters SET value = 10 WHERE name = 'policy_add_count';
UPDATE counters SET value = 3 WHERE name = 'claim_add_count';
UPDATE counters SET value = 23 WHERE name = 'total_transactions';

-- ============================================================================
-- VERIFICATION QUERIES (for manual testing)
-- ============================================================================

-- Uncomment to verify data loaded correctly:
-- SELECT 'Customers: ' || COUNT(*) FROM customers;
-- SELECT 'Policies: ' || COUNT(*) FROM policies;
-- SELECT 'Motor Policies: ' || COUNT(*) FROM motor_policies;
-- SELECT 'Endowment Policies: ' || COUNT(*) FROM endowment_policies;
-- SELECT 'House Policies: ' || COUNT(*) FROM house_policies;
-- SELECT 'Commercial Policies: ' || COUNT(*) FROM commercial_policies;
-- SELECT 'Claims: ' || COUNT(*) FROM claims;
