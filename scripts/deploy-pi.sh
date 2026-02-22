#!/bin/bash
set -euo pipefail

# Deploy Fambudg on Raspberry Pi
# Usage: ssh pi@<pi-ip> 'cd /var/www/fambudg && ./scripts/deploy-pi.sh'

APP_DIR="/var/www/fambudg"
COMPOSE="docker compose -f docker-compose.prod.yml --env-file .env.prod"
BRANCH="${DEPLOY_BRANCH:-main}"

cd "$APP_DIR"

echo "==> Pulling latest code from ${BRANCH}..."
git fetch origin "$BRANCH"
git reset --hard "origin/${BRANCH}"

echo "==> Building Docker images..."
$COMPOSE build

echo "==> Starting services..."
$COMPOSE up -d

echo "==> Running migrations..."
$COMPOSE --profile migrate run --rm migrate

echo "==> Health check..."
for i in $(seq 1 6); do
  if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
    echo "Health check passed!"
    $COMPOSE ps
    echo "==> Deploy complete."
    exit 0
  fi
  echo "Attempt $i/6: waiting 5s..."
  sleep 5
done

echo "Health check failed after 6 attempts"
$COMPOSE logs --tail=20
exit 1
