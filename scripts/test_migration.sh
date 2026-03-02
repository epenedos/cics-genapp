#!/bin/bash
# ============================================================================
# Test Migration Script
# Runs PostgreSQL in Docker and applies migrations
# ============================================================================

set -e

CONTAINER_NAME="genapp_postgres_test"
DB_NAME="genapp_test"
DB_USER="genapp"
DB_PASSWORD="genapp_secret"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "=== GENAPP Migration Test ==="
echo ""

# Check if container already exists
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "Stopping and removing existing container..."
    docker stop "$CONTAINER_NAME" 2>/dev/null || true
    docker rm "$CONTAINER_NAME" 2>/dev/null || true
fi

# Start PostgreSQL container
echo "Starting PostgreSQL container..."
docker run -d \
    --name "$CONTAINER_NAME" \
    -e POSTGRES_USER="$DB_USER" \
    -e POSTGRES_PASSWORD="$DB_PASSWORD" \
    -e POSTGRES_DB="$DB_NAME" \
    -p 5433:5432 \
    postgres:15-alpine

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to start..."
sleep 3

for i in {1..30}; do
    if docker exec "$CONTAINER_NAME" pg_isready -U "$DB_USER" -d "$DB_NAME" > /dev/null 2>&1; then
        echo "PostgreSQL is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "ERROR: PostgreSQL did not start in time"
        docker logs "$CONTAINER_NAME"
        exit 1
    fi
    sleep 1
done

# Apply migrations
echo ""
echo "Applying migrations..."
docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" < "$PROJECT_DIR/migrations/001_initial_schema.sql"

# Verify tables were created
echo ""
echo "Verifying tables..."
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "\dt"

# Verify sequences
echo ""
echo "Verifying sequences..."
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "\ds"

# Test the helper functions
echo ""
echo "Testing helper functions..."
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "SELECT next_customer_num() AS next_customer;"
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "SELECT next_policy_num() AS next_policy;"
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "SELECT next_claim_num() AS next_claim;"

# Apply seed data
echo ""
echo "Applying seed data..."
docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" < "$PROJECT_DIR/scripts/seed.sql"

# Verify data was loaded
echo ""
echo "Verifying data counts..."
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "
SELECT 'customers' AS table_name, COUNT(*) AS count FROM customers
UNION ALL
SELECT 'policies', COUNT(*) FROM policies
UNION ALL
SELECT 'motor_policies', COUNT(*) FROM motor_policies
UNION ALL
SELECT 'endowment_policies', COUNT(*) FROM endowment_policies
UNION ALL
SELECT 'house_policies', COUNT(*) FROM house_policies
UNION ALL
SELECT 'commercial_policies', COUNT(*) FROM commercial_policies
UNION ALL
SELECT 'claims', COUNT(*) FROM claims
UNION ALL
SELECT 'counters', COUNT(*) FROM counters;
"

# Test a join query (customer with policies)
echo ""
echo "Testing join query (customers with policies)..."
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "
SELECT c.customer_num, c.first_name, c.last_name, COUNT(p.id) AS policy_count
FROM customers c
LEFT JOIN policies p ON c.customer_num = p.customer_num
GROUP BY c.customer_num, c.first_name, c.last_name
ORDER BY c.customer_num
LIMIT 5;
"

# Test counter increment
echo ""
echo "Testing counter increment function..."
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME" -c "
SELECT name, value FROM counters WHERE name = 'total_transactions';
SELECT increment_counter('total_transactions', 5) AS new_value;
SELECT name, value FROM counters WHERE name = 'total_transactions';
"

echo ""
echo "=== All tests passed! ==="
echo ""
echo "Container '$CONTAINER_NAME' is still running on port 5433."
echo "Connection string: postgresql://$DB_USER:$DB_PASSWORD@localhost:5433/$DB_NAME"
echo ""
echo "To stop the container: docker stop $CONTAINER_NAME && docker rm $CONTAINER_NAME"
