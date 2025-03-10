#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Load environment variables from .env file if it exists
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo "The .env file was not found. Please create it with the required environment variables"
    exit 1
fi

# Define required environment variables
REQUIRED_VARS=("POSTGRES_HOST" "POSTGRES_PORT" "POSTGRES_USER" "POSTGRES_PASSWORD" "POSTGRES_DB")

# Check if all required environment variables are set
for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        echo "The variable $var is not set"
        exit 1
    fi
done

# Install Go dependencies
echo "Installing Go dependencies..."
go mod tidy
echo "Dependencies installed"

# Check if PostgreSQL client is installed
if ! command -v psql &> /dev/null
then
    echo "PostgreSQL client is not installed. Please install it and try again"
    exit 1
fi

# Connect to PostgreSQL and execute migration script
echo "Connecting to PostgreSQL..."
echo "Creating the 'prices' table in the database $POSTGRES_DB..."

PGPASSWORD=$POSTGRES_PASSWORD psql -U $POSTGRES_USER -h $POSTGRES_HOST -p $POSTGRES_PORT -d $POSTGRES_DB \
-f internal/datasource/migrations/0001_create_prices_table.sql

echo "Database setup completed successfully"
