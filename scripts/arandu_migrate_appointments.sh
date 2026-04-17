#!/bin/bash

MIGRATION_FILE="internal/infrastructure/repository/sqlite/migrations/0013_add_appointments.up.sql"

if [ ! -f "$MIGRATION_FILE" ]; then
    echo "Migration not found: $MIGRATION_FILE"
    exit 1
fi

echo "Applying appointments migration to all tenants..."

TENANTS_DIR="storage/tenants"

for db_file in $(find "$TENANTS_DIR" -name "*.db"); do
    echo "Applying to: $db_file"
    sqlite3 "$db_file" < "$MIGRATION_FILE" 2>&1
done

echo "Done"
