#!/bin/bash

# E2E Test Runner Script
# This script helps run the complete e2e test suite

set -e

echo "🚀 E2E Authentication Test Runner"
echo "================================="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    print_error "Please run this script from the base-app directory"
    exit 1
fi

# Check if playwright is installed
if ! command -v npx playwright &> /dev/null; then
    print_warning "Installing Playwright..."
    npm install
fi

# Check if server is running
echo ""
echo "🔍 Checking server status..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    print_status "Server is running on http://localhost:8080"
    
    echo ""
    echo "📋 Reminder: For E2E testing, ensure the server is started with:"
    echo "   RATE_LIMIT_ENABLED=false MOCK_EMAIL=true air"
    echo ""
else
    print_error "Server is not running on http://localhost:8080"
    echo ""
    echo "Please start the server first:"
    echo "  cd ../server-api"
    echo "  RATE_LIMIT_ENABLED=false MOCK_EMAIL=true air"
    echo ""
    exit 1
fi

# Check if base app is available (in case it's already running)
echo ""
echo "🔍 Checking if base app needs to be started..."
if curl -s http://localhost:3000 > /dev/null 2>&1; then
    print_status "Base app is running on http://localhost:3000"
else
    print_warning "Base app is not running. Starting it now..."
    echo "Starting dev server in background..."
    
    # Start the dev server in background
    npm run dev > /dev/null 2>&1 &
    DEV_SERVER_PID=$!
    
    # Wait for dev server to start
    echo "Waiting for dev server to start..."
    for i in {1..30}; do
        if curl -s http://localhost:3000 > /dev/null 2>&1; then
            print_status "Base app started successfully on http://localhost:3000"
            break
        fi
        sleep 2
    done
    
    # Check if it actually started
    if ! curl -s http://localhost:3000 > /dev/null 2>&1; then
        print_error "Failed to start base app"
        kill $DEV_SERVER_PID 2>/dev/null || true
        exit 1
    fi
fi

# Install playwright browsers if needed
echo ""
echo "🔍 Checking Playwright browsers..."
if ! npx playwright install --dry-run chromium > /dev/null 2>&1; then
    print_warning "Installing Playwright browsers..."
    npx playwright install chromium
fi

# Run the tests
echo ""
echo "🧪 Running E2E tests..."
echo ""

if npx playwright test; then
    print_status "All E2E tests passed!"
else
    print_error "Some tests failed. Check the output above for details."
    
    echo ""
    echo "💡 Debugging tips:"
    echo "  - Run with --headed to see browser: npx playwright test --headed"
    echo "  - Run specific test: npx playwright test e2e/auth-flow.spec.ts"
    echo "  - Debug mode: npx playwright test --debug"
fi

# Cleanup if we started the dev server
if [ ! -z "$DEV_SERVER_PID" ]; then
    print_status "Stopping dev server..."
    kill $DEV_SERVER_PID 2>/dev/null || true
fi

echo ""
echo "✨ E2E test run completed!"