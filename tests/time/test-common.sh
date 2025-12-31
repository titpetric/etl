#!/bin/bash
# Common test functions for datetime portability tests

# Insert JSON data and report
insert_test_data() {
  local json_file=$1
  local description=$2
  local row_count=$3
  
  echo "Insert: $description"
  cat "$json_file" | etl insert datetime_test
  echo "✓ $description ($row_count rows)"
  echo
}

# Display output for a specific column type
show_column_output() {
  local column=$1
  local description=$2
  local row_range=$3
  
  echo "$description"
  if [ -z "$row_range" ]; then
    etl get --all datetime_test | jq ".[] | select(.$column != null) | {id, $column}"
  else
    etl get --all datetime_test | jq ".$row_range | .[] | {id, $column}"
  fi
  echo
}

# Test section header
section() {
  echo
  echo "=== $1 ==="
  echo
}

# Subsection header
subsection() {
  echo "$1"
}

# Show all columns for a row range
show_all_columns() {
  local row_range=$1
  local description=$2
  
  echo "$description"
  etl get --all datetime_test | jq ".$row_range | .[]"
  echo
}

# Format test results header
format_test_header() {
  echo "Testing input format: $1"
  echo "Expected output type: $2"
  echo
}

# Cleanup
cleanup() {
  rm -f /tmp/test_*.json 2>/dev/null
}
