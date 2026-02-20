---
description: How to deploy Sidji Omnichannel to VPS
---

# Deploy Sidji Omnichannel to VPS

## Server Info
- **Domain (frontend):** sidji.nvlhnn.dpdns.org
- **Domain (API):** api.sidji.nvlhnn.dpdns.org
- **VPS IP:** 43.157.203.189
- **Database:** Neon PostgreSQL (external)
- **Redis:** Not used

## Architecture

```
                                         ┌─── Neon PostgreSQL (external)
sidji.nvlhnn.dpdns.org → Nginx → Next.js :3000
api.sidji.nvlhnn.dpdns.org → Nginx → Go API :8080 ──┘
```

## Memory Budget (2GB RAM VPS)

| Service      | Limit  | Reserved |
|--------------|--------|----------|
| Go API       | 256MB  | 128MB    |
| Next.js Web  | 256MB  | 128MB    |
| Nginx        | 64MB   | 32MB     |
| OS + Docker  | ~800MB | -        |
| **Free RAM** | ~624MB | -        |
| **Swap**     | 2GB    | (disk)   |

## Step-by-Step

### 1. Setup VPS (run once)

// turbo-all

```bash
ssh root@43.157.203.189
```

Then run:

```bash
chmod +x scripts/setup-vps.sh
./scripts/setup-vps.sh
```

### 2. Clone the repository

```bash
cd /opt/sidji-omnichannel
git clone YOUR_GITHUB_REPO_URL -b hexagonal .
```

### 3. Create .env.production (interactive, no secrets in git)

```bash
chmod +x scripts/create-env.sh
./scripts/create-env.sh
```

This will prompt you for:
- DATABASE_URL (from Neon dashboard)
- META_APP_ID, META_APP_SECRET
- GEMINI_API_KEY
- AWS credentials (optional)

APP_SECRET is auto-generated.

### 4. Run database migrations

```bash
apt-get install -y postgresql-client

# Use the DATABASE_URL from your .env.production
source .env.production
for f in migrations/*.up.sql; do
    echo "Applying: $f"
    psql "$DATABASE_URL" -f "$f"
done
```

### 5. Setup SSL (run once)

```bash
chmod +x scripts/init-ssl.sh
./scripts/init-ssl.sh
```

### 6. Deploy

```bash
chmod +x scripts/deploy.sh
./scripts/deploy.sh --build
```

### 7. Verify

```bash
# Check services
docker compose -f docker-compose.prod.yml ps

# Check memory
docker stats --no-stream

# Test API
curl https://api.sidji.nvlhnn.dpdns.org/health

# View logs
docker compose -f docker-compose.prod.yml logs -f api
```

## Updating

```bash
cd /opt/sidji-omnichannel
./scripts/deploy.sh --build
# (deploy.sh does git pull automatically)
```

## Useful Commands

```bash
# View all logs
docker compose -f docker-compose.prod.yml logs -f

# Restart specific service
docker compose -f docker-compose.prod.yml restart api

# Connect to Neon database
source .env.production
psql "$DATABASE_URL"

# Check disk usage
df -h

# Check memory
free -h
```

## Troubleshooting

### Container won't start
```bash
docker compose -f docker-compose.prod.yml logs api
docker compose -f docker-compose.prod.yml logs web
```

### Database connection issues
```bash
source .env.production
psql "$DATABASE_URL" -c "SELECT 1"
```

### SSL Certificate Renewal
```bash
docker compose -f docker-compose.prod.yml run --rm certbot certbot renew --force-renewal
docker compose -f docker-compose.prod.yml restart nginx
```

### Database Backup
```bash
source .env.production
pg_dump "$DATABASE_URL" > backup_$(date +%Y%m%d).sql
```
