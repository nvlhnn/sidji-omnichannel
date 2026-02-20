#!/bin/bash
# ==============================================
# Sidji Omnichannel - Initial SSL Setup
# ==============================================
# Run this ONCE on your VPS to get SSL certificates
# for both sidji.nvlhnn.dpdns.org and api.sidji.nvlhnn.dpdns.org
# Usage: ./scripts/init-ssl.sh [email]
# ==============================================

set -e

DOMAIN_FRONTEND="sidji.nvlhnn.dpdns.org"
DOMAIN_API="api.sidji.nvlhnn.dpdns.org"
EMAIL="${1:-admin@nvlhnn.dpdns.org}"

echo "🔒 Setting up SSL for:"
echo "   Frontend: $DOMAIN_FRONTEND"
echo "   API:      $DOMAIN_API"

# Create required directories
mkdir -p certbot/conf
mkdir -p certbot/www

# Step 1: Create a temporary nginx config (HTTP only, for Let's Encrypt)
cat > nginx/conf.d/temp-ssl.conf << TEMPCONF
server {
    listen 80;
    server_name $DOMAIN_FRONTEND $DOMAIN_API;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 200 'Sidji is setting up SSL...';
        add_header Content-Type text/plain;
    }
}
TEMPCONF

# Step 2: Start nginx with temporary config
echo "🚀 Starting nginx for SSL verification..."
mv nginx/conf.d/default.conf nginx/conf.d/default.conf.bak
docker compose -f docker-compose.prod.yml up -d nginx

# Wait for nginx to start
sleep 5

# Step 3: Get SSL certificate for BOTH domains (single cert)
echo "📜 Requesting SSL certificate for both domains..."
docker compose -f docker-compose.prod.yml run --rm certbot \
    certbot certonly --webroot \
    --webroot-path=/var/www/certbot \
    --email "$EMAIL" \
    --agree-tos \
    --no-eff-email \
    -d "$DOMAIN_FRONTEND" \
    -d "$DOMAIN_API"

# Step 4: Clean up and restore real config
echo "🔄 Restoring production nginx config..."
rm nginx/conf.d/temp-ssl.conf
mv nginx/conf.d/default.conf.bak nginx/conf.d/default.conf

# Stop nginx (deploy.sh will start everything)
docker compose -f docker-compose.prod.yml down

echo ""
echo "✅ SSL setup complete!"
echo "   Certificate covers both:"
echo "   - $DOMAIN_FRONTEND"
echo "   - $DOMAIN_API"
echo ""
echo "📌 Certificate will auto-renew via the certbot container."
echo "📌 Next: run ./scripts/deploy.sh --build --migrate"
