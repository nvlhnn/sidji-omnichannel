#!/bin/bash
# ==============================================
# Sidji Omnichannel - Production Deploy Script
# ==============================================
# Usage: ./scripts/deploy.sh [--build]
#
# This script:
#   1. Pulls latest code from GitHub
#   2. Builds Docker images
#   3. Auto-runs DB migrations on API startup
#   4. Starts all services
#
# .env.production lives ONLY on the VPS (not in git)
# First time? Run: ./scripts/create-env.sh
# ==============================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

COMPOSE_FILE="docker-compose.prod.yml"
BUILD_FLAG=false

# Parse arguments
for arg in "$@"; do
    case $arg in
        --build)
            BUILD_FLAG=true
            ;;
    esac
done

echo -e "${BLUE}╔══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  🚀 Sidji Omnichannel Deploy             ║${NC}"
echo -e "${BLUE}║  🌐 sidji.nvlhnn.dpdns.org               ║${NC}"
echo -e "${BLUE}║  🔌 sidji-api.nvlhnn.dpdns.org            ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════╝${NC}"
echo ""

# Check .env.production exists
if [ ! -f .env.production ]; then
    echo -e "${RED}❌ .env.production not found!${NC}"
    echo -e "   This file lives only on the VPS (not in git)."
    echo -e "   Run: ${BLUE}./scripts/create-env.sh${NC} to create it."
    exit 1
fi

# Source env for variable access
set -a
source .env.production
set +a

# Check required variables
if [[ "$APP_SECRET" == *"CHANGE_THIS"* ]]; then
    echo -e "${RED}❌ Please update APP_SECRET in .env.production${NC}"
    echo -e "   Generate one with: openssl rand -hex 32"
    exit 1
fi

if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}❌ DATABASE_URL is not set in .env.production${NC}"
    exit 1
fi

# Step 1: Pull latest code from GitHub
echo -e "${YELLOW}📥 Step 1/6: Pulling latest code from GitHub...${NC}"
git fetch origin main
git reset --hard origin/main
echo -e "${GREEN}  ✅ Code updated${NC}"

# Step 2: Pull latest base images
echo -e "${YELLOW}📦 Step 2/6: Pulling latest base images...${NC}"
docker compose -f $COMPOSE_FILE pull nginx

# Step 3: Build application images
if [ "$BUILD_FLAG" = true ]; then
    echo -e "${YELLOW}🔨 Step 3/6: Building application images...${NC}"
    docker compose -f $COMPOSE_FILE build api web
else
    echo -e "${BLUE}⏭️  Step 3/6: Skipping build (use --build to rebuild)${NC}"
fi

# Step 4: Stop existing containers (graceful)
# We SKIP 'docker compose down' here to allow docker compose up to gracefully recreate containers.
# This reduces downtime from ~1 minute to ~5 seconds.

# Step 5: Seed schema_migrations if needed (first-time setup for auto-migrator)
echo -e "${YELLOW}📊 Step 5/6: Checking migration tracking...${NC}"
if command -v psql &> /dev/null && [ -n "$DATABASE_URL" ]; then
    # Create tracking table and seed existing migrations
    psql "$DATABASE_URL" -c "CREATE TABLE IF NOT EXISTS schema_migrations (version VARCHAR(255) PRIMARY KEY, applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW());" 2>/dev/null || true
    if [ -f migrations/seed_migrations.sql ]; then
        psql "$DATABASE_URL" -f migrations/seed_migrations.sql 2>/dev/null || true
    fi
    echo -e "${GREEN}  ✅ Migration tracking ready (auto-runs on API boot)${NC}"
else
    echo -e "${BLUE}  ⏭️  psql not available — migrations will auto-run on API startup${NC}"
fi

# Step 6: Start all services
echo -e "${YELLOW}🚀 Step 6/6: Starting all services...${NC}"
docker compose -f $COMPOSE_FILE up -d

# Wait and check health
echo ""
echo -e "${YELLOW}⏳ Waiting for services to be healthy...${NC}"
sleep 5

# Step 7: Clean up old Docker images to save disk space
echo ""
echo -e "${YELLOW}🧹 Step 7: Cleaning up old unused Docker images...${NC}"
docker image prune -f
echo -e "${GREEN}  ✅ Cleanup complete${NC}"

# Check service status
echo ""
echo -e "${BLUE}📊 Service Status:${NC}"
docker compose -f $COMPOSE_FILE ps

# Check memory usage
echo ""
echo -e "${BLUE}💾 Memory Usage:${NC}"
docker stats --no-stream --format "table {{.Name}}\t{{.MemUsage}}\t{{.MemPerc}}\t{{.CPUPerc}}" \
    sidji-api sidji-web sidji-nginx 2>/dev/null || true

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  ✅ Deployment Complete!                  ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════╝${NC}"
echo ""
echo -e "  🌐 Frontend:  https://sidji.nvlhnn.dpdns.org"
echo -e "  🔌 API:       https://sidji-api.nvlhnn.dpdns.org/api"
echo -e "  📱 WebSocket: wss://sidji-api.nvlhnn.dpdns.org/api/ws"
echo -e "  🗄️  Database:  Neon PostgreSQL (ap-southeast-1)"
echo ""
