#!/bin/bash
# ==============================================
# Sidji Omnichannel - Cloudflare Origin SSL Setup
# ==============================================
# This script helps you set up Cloudflare Origin
# Certificates for end-to-end encryption (Full Strict).
#
# BEFORE running this script, you must:
# 1. Go to Cloudflare Dashboard → SSL/TLS → Origin Server
# 2. Click "Create Certificate"
# 3. Choose: RSA (2048), hostnames: *.sidji.nvlhnn.dpdns.org, sidji.nvlhnn.dpdns.org
# 4. Certificate validity: 15 years (recommended)
# 5. Copy the Origin Certificate (PEM) and Private Key
#
# Usage: ./scripts/setup-ssl.sh
# ==============================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SSL_DIR="ssl"
CERT_FILE="$SSL_DIR/origin.pem"
KEY_FILE="$SSL_DIR/origin-key.pem"

echo -e "${BLUE}╔══════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  🔒 Cloudflare Origin SSL Setup              ║${NC}"
echo -e "${BLUE}║  Best practice: Full (Strict) encryption     ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════╝${NC}"
echo ""

# Create SSL directory
mkdir -p "$SSL_DIR"

# Check if certificates already exist
if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
    echo -e "${YELLOW}⚠️  SSL certificates already exist:${NC}"
    echo "   $CERT_FILE"
    echo "   $KEY_FILE"
    echo ""
    read -p "Overwrite? (y/N): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}Keeping existing certificates.${NC}"
        exit 0
    fi
fi

echo ""
echo -e "${YELLOW}📋 Step 1: Get your Origin Certificate from Cloudflare${NC}"
echo ""
echo "   1. Go to: https://dash.cloudflare.com"
echo "   2. Select domain: nvlhnn.dpdns.org"
echo "   3. Go to: SSL/TLS → Origin Server"
echo "   4. Click: 'Create Certificate'"
echo "   5. Settings:"
echo "      - Key type: RSA (2048)"
echo "      - Hostnames: *.sidji.nvlhnn.dpdns.org, sidji.nvlhnn.dpdns.org"
echo "      - Validity: 15 years"
echo "   6. Click 'Create'"
echo ""
echo -e "${YELLOW}📋 Step 2: Paste your Origin Certificate${NC}"
echo "   (Paste the PEM certificate, then press Enter + Ctrl+D)"
echo ""

cat > "$CERT_FILE"

echo ""
echo -e "${YELLOW}📋 Step 3: Paste your Private Key${NC}"
echo "   (Paste the private key, then press Enter + Ctrl+D)"
echo ""

cat > "$KEY_FILE"

# Set secure file permissions
chmod 600 "$KEY_FILE"
chmod 644 "$CERT_FILE"

echo ""
echo -e "${GREEN}✅ SSL certificates saved:${NC}"
echo "   Certificate: $CERT_FILE"
echo "   Private Key: $KEY_FILE"

# Validate certificate format
if grep -q "BEGIN CERTIFICATE" "$CERT_FILE" && grep -q "BEGIN" "$KEY_FILE"; then
    echo -e "${GREEN}✅ Certificate format looks valid${NC}"
else
    echo -e "${RED}⚠️  Warning: Certificate format might be invalid.${NC}"
    echo "   Make sure you pasted the full PEM content including"
    echo "   -----BEGIN CERTIFICATE----- and -----END CERTIFICATE-----"
fi

echo ""
echo -e "${YELLOW}📋 Step 4: Configure Cloudflare SSL mode${NC}"
echo ""
echo "   1. Go to: Cloudflare Dashboard → SSL/TLS → Overview"
echo "   2. Set encryption mode to: ${GREEN}Full (strict)${NC}"
echo ""
echo -e "${YELLOW}📋 Step 5: Open firewall ports${NC}"
echo ""
echo "   Run these commands:"
echo "     sudo ufw allow 80/tcp"
echo "     sudo ufw allow 443/tcp"
echo ""
echo -e "${YELLOW}📋 Step 6: Deploy${NC}"
echo ""
echo "   Run: ./scripts/deploy.sh --build"
echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  ✅ SSL setup complete!                       ║${NC}"
echo -e "${GREEN}║  Deploy to activate HTTPS.                   ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════╝${NC}"
