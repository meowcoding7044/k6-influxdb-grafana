#!/bin/bash
set -e

DB_NAME="k6"
DB_USER="admin"
DB_PASS="admin123"

echo "Creating database $DB_NAME..."
influx -username $DB_USER -password $DB_PASS -execute "CREATE DATABASE $DB_NAME"

echo "Granting privileges to $DB_USER..."
influx -username $DB_USER -password $DB_PASS -execute "GRANT ALL ON $DB_NAME TO $DB_USER"

echo "✅ InfluxDB initialization complete"
