#!/bin/bash

# Load data into PostgreSQL database
# This script seeds the lunch-menu database with restaurant data

set -e

# Default database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-lunch_menu}

echo "🌱 Seeding lunch-menu database..."
echo "   Host: $DB_HOST:$DB_PORT"
echo "   Database: $DB_NAME"
echo "   User: $DB_USER"
echo ""

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "❌ Error: psql is not installed or not in PATH"
    echo "   Please install PostgreSQL client tools"
    exit 1
fi

# Check if database is reachable
echo "🔍 Checking database connection..."
if ! psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ Error: Cannot connect to database"
    echo "   Please ensure:"
    echo "   - PostgreSQL is running"
    echo "   - Database '$DB_NAME' exists"
    echo "   - Connection parameters are correct"
    exit 1
fi

# Run the seed script
echo "📊 Loading restaurant data..."
if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f seed_data.sql; then
    echo ""
    echo "✅ Database seeded successfully!"
    echo ""
    echo "📈 Database stats:"
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
        SELECT 
            'Restaurants' as table_name, 
            COUNT(*) as count 
        FROM restaurants 
        WHERE is_active = true
        UNION ALL
        SELECT 
            'Menu Items' as table_name, 
            COUNT(*) as count 
        FROM menu_items 
        WHERE is_available = true;
    "
else
    echo "❌ Error: Failed to seed database"
    exit 1
fi

echo ""
echo "🚀 Your lunch-menu database is ready!"
echo "   Start the Go application with: go run main.go"
echo "   Or with Docker: docker-compose up"
