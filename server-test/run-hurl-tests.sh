#!/bin/bash

# AE SaaS Basic HURL Test Runner with Template Support
# Runs comprehensive API tests using HURL with unique identifiers per run
# Usage: ./run-hurl-tests.sh [test_number]
# Example: ./run-hurl-tests.sh 02   # Runs all tests starting with "02"

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
HOST="http://localhost:8080"
CONFIG_FILE="tests/hurl/hurl.config"
HURL_DIR="tests/hurl"
TEMPLATES_DIR="tests/hurl/templates"
PROCESSED_DIR="tests/hurl/processed"
RESULTS_DIR="test_results"

# Enable nullglob so globs that match nothing expand to empty
shopt -s nullglob

# Check for test filter argument
TEST_FILTER="$1"
VERBOSE_OUTPUT=""

if [ -n "$TEST_FILTER" ]; then
    echo -e "${BLUE}🔍 Running tests matching pattern: ${TEST_FILTER}*${NC}"
    VERBOSE_OUTPUT="true"
else
    echo -e "${BLUE}🚀 Running all tests${NC}"
fi

# Generate unique identifiers for this test run
TIMESTAMP=$(date +%s)
NANO_PART=$(date +%N | cut -c1-6)  # Get microseconds
PROCESS_ID=$$  # Current process ID
RANDOM_PART=$RANDOM
RANDOM_ID=$(echo "${TIMESTAMP}${NANO_PART}${PROCESS_ID}${RANDOM_PART}" | shasum | cut -c1-8)
UNIQUE_ID="${TIMESTAMP}_${RANDOM_ID}"
UNIQUE_USERNAME="testuser_${UNIQUE_ID}"
UNIQUE_EMAIL="test_${UNIQUE_ID}@example.com"
UNIQUE_CUSTOMER="customer_${UNIQUE_ID}"
UNIQUE_ORG="org_${UNIQUE_ID}"
UNIQUE_PASSWORD="Pass123_${RANDOM_ID}"

# Read host from config if it exists
if [ -f "$CONFIG_FILE" ]; then
    HOST=$(grep "^host" "$CONFIG_FILE" | cut -d' ' -f3 2>/dev/null || echo "$HOST")
fi
HOST=${HOST:-${TEST_HOST:-"http://localhost:8080"}}

# Create directories
mkdir -p "$RESULTS_DIR"
mkdir -p "$PROCESSED_DIR"
mkdir -p "$TEMPLATES_DIR"

echo -e "${BLUE}🚀 Starting AE SaaS Basic HURL Tests with Templating${NC}"
echo -e "${BLUE}Host: ${HOST}${NC}"
echo -e "${BLUE}Results Directory: ${RESULTS_DIR}${NC}"
echo -e "${BLUE}Unique Test ID: ${UNIQUE_ID}${NC}"
echo -e "${BLUE}Test Username: ${UNIQUE_USERNAME}${NC}"
echo -e "${BLUE}Test Email: ${UNIQUE_EMAIL}${NC}"
echo ""

# For tests we want to disable email verification so registration -> login works
export FEATURE_EMAIL_VERIFICATION=false

# Function to process template files
process_template() {
    local template_file="$1"
    local output_file="$2"
    
    sed -e "s|{{UNIQUE_ID}}|${UNIQUE_ID}|g" \
        -e "s|{{UNIQUE_USERNAME}}|${UNIQUE_USERNAME}|g" \
        -e "s|{{UNIQUE_EMAIL}}|${UNIQUE_EMAIL}|g" \
        -e "s|{{UNIQUE_CUSTOMER}}|${UNIQUE_CUSTOMER}|g" \
        -e "s|{{UNIQUE_ORG}}|${UNIQUE_ORG}|g" \
        -e "s|{{UNIQUE_PASSWORD}}|${UNIQUE_PASSWORD}|g" \
        -e "s|{{HOST}}|${HOST}|g" \
        -e "s|{{host}}|${HOST}|g" \
        "$template_file" > "$output_file"
}

# Function to check server availability
check_server() {
    echo -e "${YELLOW}🔍 Checking server availability...${NC}"
    if curl -s --max-time 5 "${HOST}/api/v1/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Server is running${NC}"
        return 0
    else
        echo -e "${RED}❌ Server is not responding at ${HOST}${NC}"
        echo -e "${YELLOW}💡 Make sure the server is running with: make run${NC}"
        return 1
    fi
}

# Function to run a single test
run_test() {
    local test_file="$1"
    local test_name=$(basename "$test_file" .hurl)
    
    echo -e "${BLUE}🧪 Running ${test_name}.hurl...${NC}"
    
    # Check if this is a template file or regular file
    local source_file="$PROCESSED_DIR/${test_name}.hurl"
    
    if [ -f "$TEMPLATES_DIR/${test_name}.hurl" ]; then
        # Process template
        process_template "$TEMPLATES_DIR/${test_name}.hurl" "$source_file"
    else
        # Use regular file and process it for variables
        process_template "$test_file" "$source_file"
    fi
    
    # Skip files that contain no HTTP requests (variables-only files)
    if ! grep -E "^(GET|POST|PUT|DELETE|PATCH|HEAD) " -q "$source_file"; then
        echo -e "${YELLOW}ℹ️  Skipping ${test_name}.hurl (no HTTP requests found)${NC}"
        return 0
    fi

    # Run the test with or without verbose output
    local hurl_output=""
    if [ -n "$VERBOSE_OUTPUT" ]; then
        # Show verbose output for filtered tests
        echo -e "${YELLOW}📄 Processed test file contents:${NC}"
        echo "----------------------------------------"
        cat "$source_file"
        echo "----------------------------------------"
        echo ""
        
        if hurl_output=$(hurl "$source_file" \
            --variable "host=${HOST}" \
            --variable "base_url=${HOST}" \
            --variable "auth_token=${AUTH_TOKEN}" \
            --test \
            --verbose 2>&1); then
            echo -e "${GREEN}✅ ${test_name}.hurl passed${NC}"
            echo -e "${BLUE}📋 Test output:${NC}"
            echo "$hurl_output"
            echo ""
            return 0
        else
            echo -e "${RED}❌ ${test_name}.hurl failed${NC}"
            echo -e "${RED}📝 Full error output:${NC}"
            echo "$hurl_output"
            echo ""
            return 1
        fi
    else
        # Regular quiet mode: always run and then inspect JSON for failures
        hurl "$source_file" \
            --variable "host=${HOST}" \
            --variable "base_url=${HOST}" \
            --variable "auth_token=${AUTH_TOKEN}" \
            --test \
            --json > "$RESULTS_DIR/${test_name}.json" 2>/dev/null || true

        if [ -f "$RESULTS_DIR/${test_name}.json" ]; then
            local failed_msg
            failed_msg=$(jq -r '.entries[]?.asserts[]? | select(.success == false) | .message' "$RESULTS_DIR/${test_name}.json" 2>/dev/null | head -1 || true)
            if [ -z "$failed_msg" ] || [ "$failed_msg" = "null" ]; then
                echo -e "${GREEN}✅ ${test_name}.hurl passed${NC}"
                return 0
            else
                echo -e "${RED}❌ ${test_name}.hurl failed${NC}"
                echo -e "${RED}📝 Error details:${NC}"
                echo "\"$failed_msg\""
                return 1
            fi
        else
            echo -e "${RED}❌ ${test_name}.hurl failed${NC}"
            echo -e "${RED}📝 Error details:${NC}"
            echo "\"No response\""
            return 1
        fi
    fi

}

# Check server availability first
if ! check_server; then
    exit 1
fi

# Perform a quick login to obtain a reusable auth token for template tests
echo -e "${YELLOW}🔑 Obtaining reusable auth token for template tests...${NC}"
AUTH_TOKEN=$(curl -s -X POST "${HOST}/api/v1/auth/login" -H 'Content-Type: application/json' -d '{"email":"testuser@unburdy.de","password":"newpass123"}' | jq -r '.data.token // empty')
if [ -z "$AUTH_TOKEN" ]; then
    echo -e "${RED}⚠️  Failed to obtain auth token; some template tests may fail${NC}"
else
    echo -e "${GREEN}✅ Obtained auth token${NC}"
fi

echo ""

# Track test results
total_tests=0
passed_tests=0
failed_tests=0
failed_test_names=()

echo -e "${GREEN}🔍 Collecting test files...${NC}"

# Get test files based on filter
test_files=()
if [ -n "$TEST_FILTER" ]; then
    # Filter tests by pattern
    echo -e "${YELLOW}📋 Filtering tests with pattern: ${TEST_FILTER}*${NC}"
    VERBOSE_OUTPUT=true  # Auto-enable verbose output for filtered tests
    
    # Look in templates directory first
    if [ -d "$TEMPLATES_DIR" ]; then
        for file in "$TEMPLATES_DIR/${TEST_FILTER}"*.hurl; do
            if [ -f "$file" ]; then
                test_files+=("$file")
            fi
        done
    fi
    
    # Also look in main hurl directory for direct files
    for file in "$HURL_DIR/${TEST_FILTER}"*.hurl; do
        if [ -f "$file" ]; then
            # Only add if not already added from templates
            basename_file=$(basename "$file")
            found_in_templates=false
            for template_file in "${test_files[@]}"; do
                if [ "$(basename "$template_file")" = "$basename_file" ]; then
                    found_in_templates=true
                    break
                fi
            done
            if [ "$found_in_templates" = false ]; then
                test_files+=("$file")
            fi
        fi
    done
    
    if [ ${#test_files[@]} -eq 0 ]; then
        echo -e "${RED}❌ No test files found matching pattern: ${TEST_FILTER}*${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}📋 Found ${#test_files[@]} test files matching pattern:${NC}"
    for file in "${test_files[@]}"; do
        echo "  $(basename "$file")"
    done
    echo ""
else
    # Get all test files: templates first (if any), then regular hurl files
    if [ -d "$TEMPLATES_DIR" ]; then
        for file in "$TEMPLATES_DIR"/*.hurl; do
            if [ -f "$file" ]; then
                test_files+=("$file")
            fi
        done
    fi

    # Always include regular hurl files from main directory
    if [ -d "$HURL_DIR" ]; then
        for file in "$HURL_DIR"/*.hurl; do
            if [ -f "$file" ]; then
                test_files+=("$file")
            fi
        done
    fi

    echo -e "${GREEN}📋 Found ${#test_files[@]} total test files${NC}"
fi

# Sort test files by name for consistent execution order
IFS=$'\n' test_files=($(sort <<<"${test_files[*]}"))
unset IFS

echo -e "${GREEN}🚀 Starting test execution...${NC}"
echo ""

# Run tests
for test_file in "${test_files[@]}"; do
    total_tests=$((total_tests + 1))
    
    if run_test "$test_file"; then
        passed_tests=$((passed_tests + 1))
    else
        failed_tests=$((failed_tests + 1))
        failed_test_names+=("$(basename "$test_file" .hurl)")
    fi
    echo ""
done

# Print summary
echo -e "${BLUE}📊 Test Summary${NC}"
echo -e "${BLUE}===============${NC}"
echo -e "${BLUE}Total Tests: ${total_tests}${NC}"
echo -e "${GREEN}Passed: ${passed_tests}${NC}"
echo -e "${RED}Failed: ${failed_tests}${NC}"

if [ $failed_tests -gt 0 ]; then
    echo -e "${RED}❌ ${failed_tests} test(s) failed:${NC}"
    for test_name in "${failed_test_names[@]}"; do
        echo -e "   ${RED}• ${test_name}${NC}"
    done
    echo -e "${YELLOW}💡 Check individual result files in ${RESULTS_DIR}/ for details${NC}"
    exit 1
else
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
fi