#!/bin/bash
set -e

source "$(dirname "$0")/test-common.sh"

# Use environment variable or default to Docker service
export ETL_DB_DSN="${ETL_DB_DSN:-mysql://root:root@tcp(127.0.0.1:13306)/etl_test}"

echo "=== MySQL DateTime Portability Test ==="
echo

# Setup
etl query datetime-mysql.sql

# Test 1: Single row with all column types
section "Test 1: Single Row with All Column Types"

cat > /tmp/test_basic.json << 'EOF'
[
  {
    "mysql_date": "2024-01-15",
    "mysql_datetime": "2024-01-15 14:30:45",
    "mysql_datetime_fsp": "2024-01-15 14:30:45.123456",
    "mysql_timestamp": "2024-01-15 14:30:45",
    "mysql_bigint": 1705335045
  }
]
EOF

insert_test_data "/tmp/test_basic.json" "Basic insert with all native column types" "1"

show_all_columns "[0:1]" "All columns for initial row:"

# Test 2: DATE column behavior
section "Test 2: DATE Column"

subsection "DATE accepts: 2024-01-15"
show_column_output "mysql_date" "Output (date only, no time):" "[0:1]"

# Test 3: DATETIME column
section "Test 3: DATETIME and DATETIME(6) Columns"

subsection "DATETIME accepts ISO8601 and RFC3339:"
show_column_output "mysql_datetime" "Output (naive datetime):" "[0:1]"

subsection "DATETIME(6) accepts ISO8601 with microseconds:"
show_column_output "mysql_datetime_fsp" "Output (with microsecond precision):" "[0:1]"

# Test 4: TIMESTAMP column (risky)
section "Test 4: TIMESTAMP Column (Session Timezone Dependent - RISKY)"

subsection "Current session timezone context:"
echo "TIMESTAMP with input '2024-01-15 14:30:45':"
show_column_output "mysql_timestamp" "" "[0:1]"
echo "Note: TIMESTAMP may differ if session timezone changes. This is the 'RISKY' behavior."
echo

# Test 5: Multiple input formats
section "Test 5: Input Format Compatibility"

cat > /tmp/test_formats.json << 'EOF'
[
  {
    "mysql_datetime": "2024-01-15 14:30:45"
  },
  {
    "mysql_datetime": "2024-01-15T14:30:45"
  },
  {
    "mysql_datetime": "2024-01-15 14:30:45.999999"
  },
  {
    "mysql_datetime": "2024-01-15T14:30:45Z"
  },
  {
    "mysql_datetime": "2024-01-15T14:30:45+00:00"
  },
  {
    "mysql_datetime": "2024-01-15T14:30:45.123456Z"
  }
]
EOF

insert_test_data "/tmp/test_formats.json" "Multiple DATETIME input formats (ISO8601 and RFC3339)" "6"

subsection "ISO8601 inputs (rows 2-4) - Reliable:"
etl get --all datetime_test | jq '.[] | select(.mysql_datetime != null) | {id, mysql_datetime}' | head -9
echo

subsection "RFC3339 inputs (rows 5-7) - May have issues:"
etl get --all datetime_test | jq '.[] | select(.mysql_datetime != null) | {id, mysql_datetime}' | tail -9
echo

# Summary
section "Summary"

echo "✓ DATE: Stores date only (2024-01-15)"
echo "✓ DATETIME: Naive datetime, accepts ISO8601 (both space and T), RFC3339 support varies"
echo "◐ DATETIME(6): Like DATETIME but with microsecond precision"
echo "✗ TIMESTAMP: SESSION TIMEZONE DEPENDENT - AVOID for cross-DB use"
echo "✓ BIGINT: Unix timestamp, safest for portability"
echo

echo "RISKY patterns:"
echo "  • TIMESTAMP with any input: Output depends on session timezone"
echo "  • RFC3339 with timezone in DATETIME: May be ignored or cause errors"
echo "  • Different sessions will see different TIMESTAMP values"
echo

cleanup
echo "=== MySQL Test Complete ==="
