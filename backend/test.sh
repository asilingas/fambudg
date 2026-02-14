#!/bin/bash
# Test runner script for Fambudg backend

set -e

echo "🧪 Fambudg Test Runner"
echo "====================="
echo ""

# Check if PostgreSQL is running
if ! docker compose ps | grep -q "fambudg-db-1.*Up"; then
    echo "⚠️  PostgreSQL is not running. Starting..."
    docker compose up -d
    echo "⏳ Waiting for PostgreSQL to be ready..."
    sleep 5
fi

# Check if test database exists and run migrations
echo "📊 Setting up test database..."
DB_NAME=fambudg_test go run cmd/migrate/main.go up 2>/dev/null || echo "✓ Test database already migrated"
echo ""

# Parse command line argument
TEST_TYPE=${1:-all}

case $TEST_TYPE in
  unit)
    echo "🔬 Running Unit Tests..."
    echo "========================"
    go test ./internal/... -v
    ;;

  bdd)
    echo "🥒 Running BDD Tests..."
    echo "======================"
    go test ./tests/... --godog.format=pretty
    ;;

  coverage)
    echo "📈 Running Tests with Coverage..."
    echo "================================="
    go test ./... -coverprofile=coverage.out
    echo ""
    echo "📊 Coverage Summary:"
    go tool cover -func=coverage.out | tail -1
    echo ""
    echo "💡 Tip: Run 'go tool cover -html=coverage.out' to view detailed coverage in browser"
    ;;

  all)
    echo "🧪 Running All Tests..."
    echo "======================="
    go test ./... -v
    ;;

  *)
    echo "❌ Unknown test type: $TEST_TYPE"
    echo ""
    echo "Usage: ./test.sh [unit|bdd|coverage|all]"
    echo ""
    echo "Examples:"
    echo "  ./test.sh          # Run all tests"
    echo "  ./test.sh unit     # Run unit tests only"
    echo "  ./test.sh bdd      # Run BDD tests only"
    echo "  ./test.sh coverage # Run with coverage report"
    exit 1
    ;;
esac

echo ""
echo "✅ Tests completed!"
