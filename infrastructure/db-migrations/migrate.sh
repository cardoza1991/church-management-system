#!/bin/bash
set -e

# Check if we have credentials
if [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_HOST" ] || [ -z "$DB_NAME" ]; then
  echo "Database environment variables must be set"
  exit 1
fi

# Run all SQL files in order
for sql_file in $(ls *.sql | sort); do
  echo "Running migration: $sql_file"
  mysql -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$sql_file"
  echo "Completed migration: $sql_file"
done
