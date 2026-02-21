#!/bin/bash
# ==============================================
# Sidji Deploy Webhook Setup
# ==============================================
# Run this ONCE on VPS via Tencent Cloud web console
# It installs a tiny webhook listener that triggers
# deploy.sh when GitHub pushes
# ==============================================

set -e

DEPLOY_SECRET="${1:-$(openssl rand -hex 20)}"
PROJECT_DIR="/opt/sidji-omnichannel/sidji-omnichannel"
WEBHOOK_PORT=5001

echo "🔧 Setting up deploy webhook..."

# Install webhook tool
apt-get update && apt-get install -y webhook

# Create webhook config
mkdir -p /etc/webhook
cat > /etc/webhook/hooks.json << HOOKEOF
[
  {
    "id": "deploy",
    "execute-command": "${PROJECT_DIR}/scripts/deploy.sh",
    "command-working-directory": "${PROJECT_DIR}",
    "pass-arguments-to-command": [
      { "source": "string", "name": "--build" }
    ],
    "trigger-rule": {
      "match": {
        "type": "value",
        "value": "${DEPLOY_SECRET}",
        "parameter": {
          "source": "header",
          "name": "X-Deploy-Token"
        }
      }
    }
  }
]
HOOKEOF

# Create systemd service
cat > /etc/systemd/system/deploy-webhook.service << SVCEOF
[Unit]
Description=Sidji Deploy Webhook
After=network.target

[Service]
Type=simple
ExecStart=/usr/bin/webhook -hooks /etc/webhook/hooks.json -port ${WEBHOOK_PORT} -verbose
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
SVCEOF

# Start service
systemctl daemon-reload
systemctl enable deploy-webhook
systemctl start deploy-webhook

echo ""
echo "✅ Webhook listener running on port ${WEBHOOK_PORT}"
echo ""
echo "╔══════════════════════════════════════════════════╗"
echo "║  Add these GitHub Secrets:                       ║"
echo "╠══════════════════════════════════════════════════╣"
echo "║  DEPLOY_SECRET:      ${DEPLOY_SECRET}"
echo "║  DEPLOY_WEBHOOK_URL: http://43.157.203.189:${WEBHOOK_PORT}/hooks/deploy"
echo "╚══════════════════════════════════════════════════╝"
echo ""
echo "Test with:"
echo "  curl -X POST -H 'X-Deploy-Token: ${DEPLOY_SECRET}' http://43.157.203.189:${WEBHOOK_PORT}/hooks/deploy"
