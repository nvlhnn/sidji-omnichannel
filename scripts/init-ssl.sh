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

# Step 1: Stop everything first
echo "⏹️  Stopping all containers..."
docker compose -f docker-compose.prod.yml down 2>/dev/null || true

# Step 2: Create a temporary nginx config (HTTP only)
echo "📝 Creating temporary HTTP-only nginx config..."
mv nginx/conf.d/default.conf nginx/conf.d/default.conf.bak

cat > nginx/conf.d/temp-ssl.conf << TEMPCONF
server {
    listen 80;
    server_name $DOMAIN_FRONTEND $DOMAIN_API;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 200 'Sidji SSL setup in progress...';
        add_header Content-Type text/plain;
    }
}
TEMPCONF

# Step 3: Start ONLY nginx (standalone, no dependencies)
echo "🚀 Starting nginx for SSL verification..."
docker run -d --name sidji-nginx-temp \
    -p 80:80 \
    -v "$(pwd)/nginx/nginx.conf:/etc/nginx/nginx.conf:ro" \
    -v "$(pwd)/nginx/conf.d:/etc/nginx/conf.d:ro" \
    -v "$(pwd)/certbot/www:/var/www/certbot:ro" \
    nginx:alpine

sleep 3

# Verify nginx is serving
echo "🔍 Testing HTTP access..."
curl -s http://localhost/ || echo "Warning: localhost test failed, but external access may still work"

# Step 4: Get SSL certificate using standalone certbot container
echo "📜 Requesting SSL certificate for both domains..."
docker run --rm \
    -v "$(pwd)/certbot/conf:/etc/letsencrypt" \
    -v "$(pwd)/certbot/www:/var/www/certbot" \
    certbot/certbot certonly --webroot \
    --webroot-path=/var/www/certbot \
    --email "$EMAIL" \
    --agree-tos \
    --no-eff-email \
    -d "$DOMAIN_FRONTEND" \
    -d "$DOMAIN_API"

# Step 5: Clean up
echo "🔄 Cleaning up..."
docker stop sidji-nginx-temp && docker rm sidji-nginx-temp
rm nginx/conf.d/temp-ssl.conf
mv nginx/conf.d/default.conf.bak nginx/conf.d/default.conf

echo ""
echo "✅ SSL setup complete!"
echo "   Certificate covers both:"
echo "   - $DOMAIN_FRONTEND"
echo "   - $DOMAIN_API"
echo ""
echo "📌 Now deploy with: ./scripts/deploy.sh --build"
