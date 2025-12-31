#!/bin/bash
set -e

source "$(dirname "$0")/test-common.sh"

# Use environment variable or default to Docker service
export ETL_DB_DSN="${ETL_DB_DSN:-postgres://postgres:postgres@localhost:15432/etl_test?sslmode=disable}"
echo "=== PostgreSQL DateTime Portability Test ==="
echo

# Setup
etl query datetime-postgres.sql

# Test 1: Single row with all column types
section "Test 1: Single Row with All Column Types"

cat > /tmp/test_basic.json << 'EOF'
[
  {
    "postgres_date": "2024-01-15",
    "postgres_timestamp": "2024-01-15 14:30:45",
    "postgres_timestamptz": "2024-01-15 14:30:45+00",
    "postgres_bigint": 1705335045
  }
]
EOF

insert_test_data "/tmp/test_basic.json" "Basic insert with all native column types" "1"

show_all_columns "[0:1]" "All columns for initial row:"

# Test 2: DATE column
section "Test 2: DATE Column"

subsection "DATE accepts: 2024-01-15"
show_column_output "postgres_date" "Output (date only, no time):" "[0:1]"

# Test 3: TIMESTAMP column (naive, session-dependent)
section "Test 3: TIMESTAMP Column (Naive - Session Timezone Dependent - RISKY)"

subsection "TIMESTAMP with input '2024-01-15 14:30:45':"
show_column_output "postgres_timestamp" "Output (depends on session timezone):" "[0:1]"

echo "⚠️  RISKY: Output depends on session timezone setting"
echo "    Same stored value appears different in different sessions"
echo

# Test 4: TIMESTAMPTZ column (timezone-aware, safe)
section "Test 4: TIMESTAMPTZ Column (Timezone-Aware - SAFE)"

subsection "TIMESTAMPTZ with input '2024-01-15 14:30:45+00':"
show_column_output "postgres_timestamptz" "Output (always with timezone):" "[0:1]"

echo "✓ SAFE: Stores UTC internally, always includes timezone in output"
echo

# Test 5: Multiple input formats for TIMESTAMP
section "Test 5: TIMESTAMP - Input Format Compatibility (RISKY)"

cat > /tmp/test_timestamp_formats.json << 'EOF'
[
  {
    "postgres_timestamp": "2024-01-15 14:30:45"
  },
  {
    "postgres_timestamp": "2024-01-15T14:30:45"
  },
  {
    "postgres_timestamp": "2024-01-15 14:30:45.123456"
  },
  {
    "postgres_timestamp": "2024-01-15T14:30:45Z"
  },
  {
    "postgres_timestamp": "2024-01-15T14:30:45+00:00"
  },
  {
    "postgres_timestamp": "2024-01-15 14:30:45.123456Z"
  }
]
EOF

insert_test_data "/tmp/test_timestamp_formats.json" "Multiple TIMESTAMP input formats" "6"

subsection "ISO8601 inputs (rows 2-4) - Naive, session-dependent:"
etl get --all datetime_test | jq '.[] | select(.postgres_timestamp != null) | {id, postgres_timestamp}' | head -9
echo

subsection "RFC3339 inputs (rows 5-7) - May be interpreted as session timezone:"
etl get --all datetime_test | jq '.[] | select(.postgres_timestamp != null) | {id, postgres_timestamp}' | tail -9
echo

# Test 6: Multiple input formats for TIMESTAMPTZ (safe)
section "Test 6: TIMESTAMPTZ - Input Format Compatibility (SAFE)"

cat > /tmp/test_timestamptz_formats.json << 'EOF'
[
  {
    "postgres_timestamptz": "2024-01-15T14:30:45Z"
  },
  {
    "postgres_timestamptz": "2024-01-15T14:30:45+00:00"
  },
  {
    "postgres_timestamptz": "2024-01-15T14:30:45.123456Z"
  }
]
EOF

insert_test_data "/tmp/test_timestamptz_formats.json" "Multiple TIMESTAMPTZ RFC3339 formats" "3"

subsection "RFC3339 inputs - Always parsed as UTC, always output with timezone:"
etl get --all datetime_test | jq '.[] | select(.postgres_timestamptz != null) | {id, postgres_timestamptz}'
echo

# Summary
section "Summary"

echo "✓ DATE: Stores date only (2024-01-15)"
echo "◐ TIMESTAMP: NAIVE - Session timezone dependent - AVOID for cross-DB"
echo "✓ TIMESTAMPTZ: TIMEZONE-AWARE - Safe, stores UTC, always returns with timezone"
echo "✓ BIGINT: Unix timestamp, safest for portability"
echo

echo "KEY DIFFERENCES:"
echo "  TIMESTAMP vs TIMESTAMPTZ:"
echo "    • TIMESTAMP is naive: 2024-01-15 14:30:45 (no TZ info)"
echo "    • TIMESTAMPTZ is aware: 2024-01-15 14:30:45+00 (explicit TZ)"
echo "    • Same TIMESTAMP value appears different in different sessions"
echo "    • TIMESTAMPTZ always shows the timezone"
echo

cleanup
echo "=== PostgreSQL Test Complete ==="
