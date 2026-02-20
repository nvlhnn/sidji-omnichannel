#!/bin/bash
# ==============================================
# Sidji Omnichannel - Create .env.production
# ==============================================
# Run this ONCE on the VPS to create the env file.
# This file is NOT stored in git.
# It will prompt you for all secrets interactively.
#
# Usage: ./scripts/create-env.sh
# ==============================================

set -e

ENV_FILE=".env.production"

if [ -f "$ENV_FILE" ]; then
    echo "⚠️  $ENV_FILE already exists!"
    read -p "Overwrite? (y/N): " confirm
    if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
        echo "Cancelled."
        exit 0
    fi
fi

echo "╔══════════════════════════════════════════╗"
echo "║  🔧 Sidji Omnichannel - Environment Setup ║"
echo "╚══════════════════════════════════════════╝"
echo ""
echo "This will create .env.production with your secrets."
echo "Press Enter to use [default] values where shown."
echo ""

# Auto-generate APP_SECRET
APP_SECRET=$(openssl rand -hex 32)
echo "✅ APP_SECRET auto-generated"

# Database
echo ""
echo "── Database (Neon PostgreSQL) ──"
read -p "DATABASE_URL (full connection string): " DATABASE_URL
if [ -z "$DATABASE_URL" ]; then
    echo "❌ DATABASE_URL is required!"
    exit 1
fi

# Meta API
echo ""
echo "── Meta API (WhatsApp & Instagram) ──"
read -p "META_APP_ID: " META_APP_ID
read -p "META_APP_SECRET: " META_APP_SECRET
read -p "META_VERIFY_TOKEN [omni_secret_2026]: " META_VERIFY_TOKEN
META_VERIFY_TOKEN=${META_VERIFY_TOKEN:-omni_secret_2026}

# AI
echo ""
echo "── AI Configuration ──"
read -p "AI_PROVIDER [gemini]: " AI_PROVIDER
AI_PROVIDER=${AI_PROVIDER:-gemini}
read -p "GEMINI_API_KEY: " GEMINI_API_KEY
read -p "OPENAI_API_KEY (optional, press Enter to skip): " OPENAI_API_KEY

# AWS (optional)
echo ""
echo "── AWS S3 (optional, press Enter to skip) ──"
read -p "AWS_ACCESS_KEY_ID: " AWS_ACCESS_KEY_ID
read -p "AWS_SECRET_ACCESS_KEY: " AWS_SECRET_ACCESS_KEY
read -p "AWS_S3_BUCKET [sidji-omnichannel-media]: " AWS_S3_BUCKET
AWS_S3_BUCKET=${AWS_S3_BUCKET:-sidji-omnichannel-media}

# Write the file
cat > "$ENV_FILE" << EOF
# ==============================================
# Sidji Omnichannel - Production Environment
# ==============================================
# ⚠️  This file is NOT stored in git!
# Generated on: $(date)
# ==============================================

# Server
APP_ENV=production
APP_PORT=8080
APP_SECRET=${APP_SECRET}

# Database
DATABASE_URL=${DATABASE_URL}

# Meta API (WhatsApp & Instagram)
META_APP_ID=${META_APP_ID}
META_APP_SECRET=${META_APP_SECRET}
META_VERIFY_TOKEN=${META_VERIFY_TOKEN}

# AWS S3 (for media storage)
AWS_REGION=ap-southeast-1
AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
AWS_S3_BUCKET=${AWS_S3_BUCKET}

# Frontend URL (for CORS)
FRONTEND_URL=https://sidji.nvlhnn.dpdns.org

# AI Configuration
AI_PROVIDER=${AI_PROVIDER}
GEMINI_API_KEY=${GEMINI_API_KEY}
OPENAI_API_KEY=${OPENAI_API_KEY}
EOF

chmod 600 "$ENV_FILE"

echo ""
echo "✅ Created $ENV_FILE (permissions: 600 - owner only)"
echo ""
echo "   Review: nano $ENV_FILE"
echo ""
echo "   Next steps:"
echo "   1. ./scripts/init-ssl.sh"
echo "   2. ./scripts/deploy.sh --build --migrate"
