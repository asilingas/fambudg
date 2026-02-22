#!/bin/bash
set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-/var/www/fambudg/backups}"
RETENTION_DAYS=7
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/fambudg_${TIMESTAMP}.sql.gz"

mkdir -p "$BACKUP_DIR"

echo "Starting backup: ${BACKUP_FILE}"

docker compose -f "${APP_DIR:-/var/www/fambudg}/docker-compose.prod.yml" exec -T db \
  pg_dump -U "${DB_USER:-fambudg}" "${DB_NAME:-fambudg}" | gzip > "$BACKUP_FILE"

echo "Backup completed: $(du -h "$BACKUP_FILE" | cut -f1)"

# Delete backups older than retention period
find "$BACKUP_DIR" -name "fambudg_*.sql.gz" -mtime +${RETENTION_DAYS} -delete

echo "Cleaned up backups older than ${RETENTION_DAYS} days"
