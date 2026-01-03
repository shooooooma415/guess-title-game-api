#!/bin/bash

# Migration script for guess-title-game-api
# Usage: ./scripts/migration.sh [up|down|create|force] [args...]

set -e

# Load environment variables
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Database configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-guess_title_game}"
DB_SSL_MODE="${DB_SSL_MODE:-disable}"

# Migration path
MIGRATION_PATH="db/migrations"

# Database URL
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

# Check if migrate is installed
if ! command -v migrate &> /dev/null; then
    echo "Error: 'migrate' command not found"
    echo "Please install golang-migrate:"
    echo "  macOS: brew install golang-migrate"
    echo "  Linux: curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz && sudo mv migrate /usr/local/bin/"
    echo "  Or visit: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"
    exit 1
fi

# Command handler
case "$1" in
    up)
        echo "Running migrations up..."
        migrate -path ${MIGRATION_PATH} -database "${DATABASE_URL}" up
        echo "Migrations completed successfully!"
        ;;

    down)
        echo "Running migrations down..."
        if [ -z "$2" ]; then
            echo "Rolling back 1 migration..."
            migrate -path ${MIGRATION_PATH} -database "${DATABASE_URL}" down 1
        else
            echo "Rolling back $2 migrations..."
            migrate -path ${MIGRATION_PATH} -database "${DATABASE_URL}" down "$2"
        fi
        echo "Rollback completed successfully!"
        ;;

    force)
        if [ -z "$2" ]; then
            echo "Error: version number required"
            echo "Usage: ./scripts/migration.sh force <version>"
            exit 1
        fi
        echo "Forcing migration version to $2..."
        migrate -path ${MIGRATION_PATH} -database "${DATABASE_URL}" force "$2"
        echo "Version forced successfully!"
        ;;

    create)
        if [ -z "$2" ]; then
            echo "Error: migration name required"
            echo "Usage: ./scripts/migration.sh create <migration_name>"
            exit 1
        fi
        echo "Creating new migration: $2"
        migrate create -ext sql -dir ${MIGRATION_PATH} -seq "$2"
        echo "Migration files created successfully!"
        ;;

    version)
        echo "Checking current migration version..."
        migrate -path ${MIGRATION_PATH} -database "${DATABASE_URL}" version
        ;;

    drop)
        read -p "Are you sure you want to DROP ALL tables? This cannot be undone! (yes/no): " confirm
        if [ "$confirm" = "yes" ]; then
            echo "Dropping all tables..."
            migrate -path ${MIGRATION_PATH} -database "${DATABASE_URL}" drop -f
            echo "All tables dropped!"
        else
            echo "Operation cancelled."
        fi
        ;;

    *)
        echo "Usage: ./scripts/migration.sh [command] [args...]"
        echo ""
        echo "Commands:"
        echo "  up              Apply all pending migrations"
        echo "  down [n]        Rollback n migrations (default: 1)"
        echo "  force <version> Force set migration version"
        echo "  create <name>   Create a new migration file"
        echo "  version         Show current migration version"
        echo "  drop            Drop all tables (requires confirmation)"
        echo ""
        echo "Environment variables:"
        echo "  DB_HOST         Database host (default: localhost)"
        echo "  DB_PORT         Database port (default: 5432)"
        echo "  DB_USER         Database user (default: postgres)"
        echo "  DB_PASSWORD     Database password (default: postgres)"
        echo "  DB_NAME         Database name (default: guess_title_game)"
        echo "  DB_SSL_MODE     SSL mode (default: disable)"
        exit 1
        ;;
esac
