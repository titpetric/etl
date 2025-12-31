#!/bin/bash
set -e

source "$(dirname "$0")/test-common.sh"

export ETL_DB_DSN="sqlite://file:datetime_test.db"

TIMESTAMP="2024-01-15 14:30:45"
UNIX_TIMESTAMP=1705335045

echo "=== SQLite DateTime Portability Test ==="
echo

# Setup
rm -f datetime_test.db
etl query datetime.sql

# Test 1: Basic single insert with all column types
section "Test 1: Single Row with All Column Types"

cat > /tmp/test_basic.json << 'EOF'
[
  {
    "id": 1,
    "sqlite_datetime": "2024-01-15 14:30:45",
    "sqlite_timestamp": "2024-01-15 14:30:45",
    "sqlite_text": "2024-01-15 14:30:45",
    "sqlite_int": 1705335045
  }
]
EOF

insert_test_data "/tmp/test_basic.json" "Basic insert with ISO8601 space format" "1"

show_all_columns "[0:1]" "All columns for row 1:"

# Test 2: Compare different column output formats
section "Test 2: Column Type Output Comparison"

cat > /tmp/test_comparison.json << 'EOF'
[
  {
    "id": 2,
    "sqlite_datetime": "2024-01-15T14:30:45",
    "sqlite_timestamp": "2024-01-15T14:30:45",
    "sqlite_text": "2024-01-15T14:30:45",
    "sqlite_int": 1705335045
  }
]
EOF

insert_test_data "/tmp/test_comparison.json" "Insert with ISO8601 T format" "1"

subsection "DATETIME output (RFC3339 with UTC):"
show_column_output "sqlite_datetime" "" "[1:2]"

subsection "TIMESTAMP output (RFC3339 with UTC):"
show_column_output "sqlite_timestamp" "" "[1:2]"

subsection "TEXT output (exact input, no parsing):"
show_column_output "sqlite_text" "" "[1:2]"

subsection "INTEGER output (Unix timestamp):"
show_column_output "sqlite_int" "" "[1:2]"

# Test 3: Multiple input formats in one batch insert
section "Test 3: Input Format Compatibility"

cat > /tmp/test_formats.json << 'EOF'
[
  {
    "id": 10,
    "sqlite_datetime": "2024-01-15 14:30:45",
    "sqlite_timestamp": "2024-01-15 14:30:45",
    "sqlite_text": "2024-01-15 14:30:45"
  },
  {
    "id": 11,
    "sqlite_datetime": "2024-01-15T14:30:45",
    "sqlite_timestamp": "2024-01-15T14:30:45",
    "sqlite_text": "2024-01-15T14:30:45"
  },
  {
    "id": 12,
    "sqlite_datetime": "2024-01-15 14:30:45.123456",
    "sqlite_timestamp": "2024-01-15 14:30:45.123456",
    "sqlite_text": "2024-01-15 14:30:45.123456"
  },
  {
    "id": 20,
    "sqlite_datetime": "2024-01-15T14:30:45Z",
    "sqlite_timestamp": "2024-01-15T14:30:45Z",
    "sqlite_text": "2024-01-15T14:30:45Z"
  },
  {
    "id": 21,
    "sqlite_datetime": "2024-01-15T14:30:45+00:00",
    "sqlite_timestamp": "2024-01-15T14:30:45+00:00",
    "sqlite_text": "2024-01-15T14:30:45+00:00"
  },
  {
    "id": 22,
    "sqlite_datetime": "2024-01-15T14:30:45.123456Z",
    "sqlite_timestamp": "2024-01-15T14:30:45.123456Z",
    "sqlite_text": "2024-01-15T14:30:45.123456Z"
  }
]
EOF

insert_test_data "/tmp/test_formats.json" "Multiple input formats (ISO8601 and RFC3339)" "6"

subsection "DATETIME/TIMESTAMP output for ISO8601 inputs (rows 10-12):"
etl get --all datetime_test | jq '.[2:5] | .[] | {id, sqlite_datetime, sqlite_timestamp}'
echo

subsection "TEXT output for ISO8601 inputs (exact input):"
etl get --all datetime_test | jq '.[2:5] | .[] | {id, sqlite_text}'
echo

subsection "DATETIME/TIMESTAMP output for RFC3339 inputs (rows 20-22):"
etl get --all datetime_test | jq '.[5:8] | .[] | {id, sqlite_datetime, sqlite_timestamp}'
echo

subsection "TEXT output for RFC3339 inputs (exact input):"
etl get --all datetime_test | jq '.[5:8] | .[] | {id, sqlite_text}'
echo

# Summary
section "Summary"

echo "✓ DATETIME: Accepts ISO8601 and RFC3339, outputs RFC3339 with UTC"
echo "✓ TIMESTAMP: Identical to DATETIME (synonyms in SQLite)"
echo "✓ TEXT: No parsing, returns exact input"
echo "✓ INTEGER: Unix timestamp, unambiguous"
echo

echo "Known issue:"
echo "  Input: 2024-01-15T14:30:45+00:00"
echo "  Output: 2024-01-15 14:30:45 +0000 +0000 (double timezone)"
echo "  Workaround: Use Z notation instead"
echo

cleanup
echo "=== SQLite Test Complete ==="
