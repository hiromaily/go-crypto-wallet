#!/bin/bash
# Initialize database schemas using Atlas migrations
# This script can be used as an alternative to the SQL-based initialization

set -e

echo "Initializing database schemas with Atlas migrations..."

# Wait for MySQL to be ready
echo "Waiting for MySQL to be ready..."
until mysqladmin ping -h"${MYSQL_HOST:-wallet-db}" -u"${MYSQL_USER:-root}" -p"${MYSQL_PASSWORD:-root}" --silent; do
  echo "MySQL is unavailable - sleeping"
  sleep 2
done

echo "MySQL is ready!"

# Create schemas if they don't exist
mysql -h"${MYSQL_HOST:-wallet-db}" -u"${MYSQL_USER:-root}" -p"${MYSQL_PASSWORD:-root}" <<EOF
CREATE DATABASE IF NOT EXISTS \`watch\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS \`keygen\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS \`sign\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EOF

# Check if Atlas is available
if ! command -v atlas &> /dev/null; then
  echo "Warning: Atlas CLI not found. Skipping migrations."
  echo "Install Atlas with: brew install arigaio/tap/atlas"
  echo "Or with Go: go install ariga.io/atlas/cmd/atlas@latest"
  exit 0
fi

# Apply migrations
echo "Applying Atlas migrations..."

# Watch schema
echo "Applying migrations for watch schema..."
atlas migrate apply \
  --dir file://tools/atlas/migrations/watch \
  --url "mysql://${MYSQL_USER:-root}:${MYSQL_PASSWORD:-root}@${MYSQL_HOST:-wallet-db}:3306/watch?charset=utf8mb4&parseTime=True&loc=Local" || {
  echo "Warning: Failed to apply watch schema migrations"
}

# Keygen schema
echo "Applying migrations for keygen schema..."
atlas migrate apply \
  --dir file://tools/atlas/migrations/keygen \
  --url "mysql://${MYSQL_USER:-root}:${MYSQL_PASSWORD:-root}@${MYSQL_HOST:-wallet-db}:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local" || {
  echo "Warning: Failed to apply keygen schema migrations"
}

# Sign schema
echo "Applying migrations for sign schema..."
atlas migrate apply \
  --dir file://tools/atlas/migrations/sign \
  --url "mysql://${MYSQL_USER:-root}:${MYSQL_PASSWORD:-root}@${MYSQL_HOST:-wallet-db}:3306/sign?charset=utf8mb4&parseTime=True&loc=Local" || {
  echo "Warning: Failed to apply sign schema migrations"
}

echo "Database initialization with Atlas completed!"

