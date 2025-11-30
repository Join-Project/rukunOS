#!/bin/bash

# Script untuk backup database
# Usage: ./backup-db.sh [production|development]

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"
cd "$PROJECT_DIR"

ENV="${1:-development}"

if [ "$ENV" = "production" ]; then
    COMPOSE_FILE="docker-compose.prod.yml"
    CONTAINER_NAME="rukunos-db-prod"
else
    COMPOSE_FILE="docker-compose.yml"
    CONTAINER_NAME="rukunos-db"
fi

# Load environment variables
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

DB_USER=${DB_USER:-rukunos_user}
DB_NAME=${DB_NAME:-rukunos_db}
BACKUP_DIR="$PROJECT_DIR/backups"

# Create backup directory if not exists
mkdir -p "$BACKUP_DIR"

# Generate backup filename
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/backup_${ENV}_${TIMESTAMP}.sql"

echo "Creating database backup..."
echo "Environment: $ENV"
echo "Database: $DB_NAME"
echo "Backup file: $BACKUP_FILE"

# Check if container is running
if ! docker ps --format "{{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
    echo "Error: Database container '$CONTAINER_NAME' is not running!"
    exit 1
fi

# Create backup
docker exec "$CONTAINER_NAME" pg_dump -U "$DB_USER" "$DB_NAME" > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "✓ Backup created successfully: $BACKUP_FILE"
    
    # Compress backup
    echo "Compressing backup..."
    gzip "$BACKUP_FILE"
    echo "✓ Backup compressed: ${BACKUP_FILE}.gz"
    
    # Show backup size
    BACKUP_SIZE=$(du -h "${BACKUP_FILE}.gz" | cut -f1)
    echo "Backup size: $BACKUP_SIZE"
    
    # Cleanup old backups (keep last 7 days)
    echo "Cleaning up old backups (keeping last 7 days)..."
    find "$BACKUP_DIR" -name "backup_${ENV}_*.sql.gz" -mtime +7 -delete
    echo "✓ Cleanup completed"
else
    echo "✗ Backup failed!"
    exit 1
fi

