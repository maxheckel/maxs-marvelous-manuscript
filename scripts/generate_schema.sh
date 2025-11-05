#!/bin/bash

# Create a clean temp database for Jet code generation
DB_PATH="data/temp_schema.db"

# Remove old temp db if exists
rm -f "$DB_PATH"

# Create new database
sqlite3 "$DB_PATH" "VACUUM;"

# Apply migrations in order (only the Up parts)
for migration in migrations/*.sql; do
    echo "Applying $migration..."
    # Extract only the Up section
    awk '/\+migrate Up/,/\+migrate Down/' "$migration" | grep -v "migrate Down" | grep -v "migrate Up" | sqlite3 "$DB_PATH"
done

echo "Schema database created at $DB_PATH"
