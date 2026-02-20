#!/bin/bash
# ==============================================
# Sidji Omnichannel - VPS Initial Setup
# ==============================================
# Run this ONCE on a fresh Ubuntu/Debian VPS
# Usage: curl -sSL <url> | bash
# Or:    ./scripts/setup-vps.sh
# ==============================================

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔══════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  🖥️  Sidji VPS Initial Setup         ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════╝${NC}"
echo ""

# ==============================================
# Step 1: System updates & essentials
# ==============================================
echo -e "${YELLOW}📦 Step 1: Updating system...${NC}"
apt-get update && apt-get upgrade -y
apt-get install -y \
    curl \
    wget \
    git \
    ufw \
    htop \
    fail2ban \
    unattended-upgrades

# ==============================================
# Step 2: Setup swap (critical for 2GB RAM!)
# ==============================================
echo -e "${YELLOW}💾 Step 2: Setting up 2GB swap...${NC}"
if [ ! -f /swapfile ]; then
    fallocate -l 2G /swapfile
    chmod 600 /swapfile
    mkswap /swapfile
    swapon /swapfile
    echo '/swapfile none swap sw 0 0' >> /etc/fstab
    # Optimize swap for low-RAM server
    sysctl vm.swappiness=10
    echo 'vm.swappiness=10' >> /etc/sysctl.conf
    sysctl vm.vfs_cache_pressure=50
    echo 'vm.vfs_cache_pressure=50' >> /etc/sysctl.conf
    echo -e "${GREEN}  ✅ Swap enabled (2GB)${NC}"
else
    echo -e "${GREEN}  ✅ Swap already exists${NC}"
fi

# ==============================================
# Step 3: Install Docker
# ==============================================
echo -e "${YELLOW}🐳 Step 3: Installing Docker...${NC}"
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh

    # Add current user to docker group
    usermod -aG docker $USER

    # Enable Docker service
    systemctl enable docker
    systemctl start docker
    echo -e "${GREEN}  ✅ Docker installed${NC}"
else
    echo -e "${GREEN}  ✅ Docker already installed${NC}"
fi

# Install Docker Compose plugin (if not included)
if ! docker compose version &> /dev/null; then
    apt-get install -y docker-compose-plugin
    echo -e "${GREEN}  ✅ Docker Compose plugin installed${NC}"
fi

# ==============================================
# Step 4: Firewall setup
# ==============================================
echo -e "${YELLOW}🔥 Step 4: Configuring firewall...${NC}"
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80/tcp
ufw allow 443/tcp
echo "y" | ufw enable
echo -e "${GREEN}  ✅ Firewall configured (SSH, HTTP, HTTPS only)${NC}"

# ==============================================
# Step 5: Fail2ban setup
# ==============================================
echo -e "${YELLOW}🛡️  Step 5: Configuring fail2ban...${NC}"
cat > /etc/fail2ban/jail.local << 'EOF'
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5

[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
EOF

systemctl enable fail2ban
systemctl restart fail2ban
echo -e "${GREEN}  ✅ Fail2ban configured${NC}"

# ==============================================
# Step 6: Optimize kernel for Docker
# ==============================================
echo -e "${YELLOW}⚡ Step 6: Kernel optimization...${NC}"
cat >> /etc/sysctl.conf << 'EOF'

# Docker & network optimizations
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 1024
net.ipv4.ip_forward = 1
net.ipv4.tcp_rmem = 4096 87380 6291456
net.ipv4.tcp_wmem = 4096 16384 4194304
fs.file-max = 65536
EOF

sysctl -p
echo -e "${GREEN}  ✅ Kernel optimized${NC}"

# ==============================================
# Step 7: Setup Docker log rotation
# ==============================================
echo -e "${YELLOW}📝 Step 7: Docker log rotation...${NC}"
cat > /etc/docker/daemon.json << 'EOF'
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF

systemctl restart docker
echo -e "${GREEN}  ✅ Docker log rotation configured${NC}"

# ==============================================
# Step 8: Create project directory
# ==============================================
echo -e "${YELLOW}📁 Step 8: Creating project directory...${NC}"
mkdir -p /opt/sidji-omnichannel
echo -e "${GREEN}  ✅ Directory: /opt/sidji-omnichannel${NC}"

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  ✅ VPS Setup Complete!                      ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════╝${NC}"
echo ""
echo -e "Next steps:"
echo -e "  1. Clone your repo:     ${BLUE}cd /opt/sidji-omnichannel && git clone <repo_url> .${NC}"
echo -e "  2. Copy env file:       ${BLUE}cp .env.production.example .env.production${NC}"
echo -e "  3. Edit env:            ${BLUE}nano .env.production${NC}"
echo -e "  4. Setup SSL:           ${BLUE}./scripts/init-ssl.sh your-domain.com your@email.com${NC}"
echo -e "  5. Deploy:              ${BLUE}./scripts/deploy.sh --build --migrate${NC}"
echo ""
echo -e "${YELLOW}⚠️  NOTE: You may need to log out and back in for Docker group changes.${NC}"
