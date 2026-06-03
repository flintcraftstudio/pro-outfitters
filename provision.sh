#!/usr/bin/env bash
# provision.sh — scaffold a new Flintcraft client project on the VPS
#
# Usage (run as root or with sudo):
#   sudo ./provision.sh <project> <domain>
#
# Example:
#   sudo ./provision.sh manifest manifest.flintcraft.studio
#
# What it does:
#   1. Allocates the next available port from /etc/flintcraft/port.state
#   2. Creates a per-project SSH deploy keypair with a restricted authorized_keys entry
#   3. Scaffolds /opt/{project}/ with docker-compose.yml and .env
#   4. Writes /etc/caddy/sites/{project}.caddy and reloads Caddy
#   5. Prints GitHub Actions secrets to stdout

set -euo pipefail

# ── Constants ──────────────────────────────────────────────────────────────────
DEPLOY_USER="deploy"
DEPLOY_HOME="/home/${DEPLOY_USER}"
STATE_DIR="/etc/flintcraft"
PORT_FILE="${STATE_DIR}/port.state"
CADDYFILE="/etc/caddy/Caddyfile"
CADDY_SITES_DIR="/etc/caddy/sites"
PORT_START=8100

# ── Args ───────────────────────────────────────────────────────────────────────
if [[ $# -ne 2 ]]; then
  echo "Usage: sudo $0 <project> <domain>" >&2
  echo "  Example: sudo $0 manifest manifest.flintcraft.studio" >&2
  exit 1
fi

PROJECT="$1"
DOMAIN="$2"

# Validate project name (lowercase letters, numbers, hyphens only)
if [[ ! "$PROJECT" =~ ^[a-z0-9-]+$ ]]; then
  echo "Error: project name must be lowercase letters, numbers, or hyphens only" >&2
  exit 1
fi

OPT_DIR="/opt/${PROJECT}"

# ── Checks ─────────────────────────────────────────────────────────────────────
if [[ "$EUID" -ne 0 ]]; then
  echo "Error: run as root or with sudo" >&2
  exit 1
fi

if ! id "$DEPLOY_USER" &>/dev/null; then
  echo "Error: '${DEPLOY_USER}' user does not exist — run the provisioning scripts first" >&2
  exit 1
fi

if [[ -d "$OPT_DIR" ]]; then
  echo "Error: ${OPT_DIR} already exists — project '${PROJECT}' may already be provisioned" >&2
  exit 1
fi

if ! command -v caddy &>/dev/null; then
  echo "Error: caddy not found on PATH" >&2
  exit 1
fi

if ! command -v docker &>/dev/null; then
  echo "Error: docker not found on PATH" >&2
  exit 1
fi

# ── Port allocation ────────────────────────────────────────────────────────────
mkdir -p "$STATE_DIR"

if [[ ! -f "$PORT_FILE" ]]; then
  echo "$PORT_START" > "$PORT_FILE"
fi

PORT=$(cat "$PORT_FILE")

if [[ "$PORT" -gt 9999 ]]; then
  echo "Error: port range exhausted (>9999)" >&2
  exit 1
fi

echo $((PORT + 1)) > "$PORT_FILE"

# ── Deploy keypair ─────────────────────────────────────────────────────────────
KEY_DIR="${OPT_DIR}/.ssh"
mkdir -p "$KEY_DIR"

ssh-keygen \
  -t ed25519 \
  -C "${PROJECT}-deploy@flintcraft" \
  -f "${KEY_DIR}/deploy_key" \
  -N ""

# Restricted authorized_keys entry — only allows the deploy command for this project
RESTRICTED_CMD="command=\"cd ${OPT_DIR} && docker compose pull && docker compose up -d\",no-port-forwarding,no-X11-forwarding,no-agent-forwarding"
PUB_KEY=$(cat "${KEY_DIR}/deploy_key.pub")

mkdir -p "${DEPLOY_HOME}/.ssh"
echo "${RESTRICTED_CMD} ${PUB_KEY}" >> "${DEPLOY_HOME}/.ssh/authorized_keys"
chmod 600 "${DEPLOY_HOME}/.ssh/authorized_keys"
chown "${DEPLOY_USER}:${DEPLOY_USER}" "${DEPLOY_HOME}/.ssh/authorized_keys"

# ── /opt/{project} scaffold ────────────────────────────────────────────────────
mkdir -p "$OPT_DIR"
chown "${DEPLOY_USER}:${DEPLOY_USER}" "$OPT_DIR"

cat > "${OPT_DIR}/docker-compose.yml" <<EOF
services:
  app:
    image: ghcr.io/flintcraftstudio/${PROJECT}:latest
    restart: unless-stopped
    env_file: .env
    ports:
      - "${PORT}:8080"
EOF

cat > "${OPT_DIR}/.env" <<EOF
# ${PROJECT} environment — add secrets here
# This file is read by docker compose at deploy time
PORT=${PORT}
EOF

chown "${DEPLOY_USER}:${DEPLOY_USER}" "${OPT_DIR}/docker-compose.yml" "${OPT_DIR}/.env"
chmod 600 "${OPT_DIR}/.env"

# ── Caddy site config ──────────────────────────────────────────────────────────
CADDY_SITE_FILE="${CADDY_SITES_DIR}/${PROJECT}.caddy"

mkdir -p "$CADDY_SITES_DIR"

if [[ -f "$CADDY_SITE_FILE" ]]; then
  echo "Error: ${CADDY_SITE_FILE} already exists — project '${PROJECT}' may already be provisioned" >&2
  exit 1
fi

cat > "$CADDY_SITE_FILE" <<EOF
# ${PROJECT}
${DOMAIN} {
    encode zstd gzip

    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "SAMEORIGIN"
        Referrer-Policy "strict-origin-when-cross-origin"
        -Server
    }

    reverse_proxy localhost:${PORT}
}
EOF

caddy reload --config "$CADDYFILE"

# ── Summary ────────────────────────────────────────────────────────────────────
PRIVATE_KEY=$(cat "${KEY_DIR}/deploy_key")

echo ""
echo "════════════════════════════════════════════════════════"
echo "  ✓ Project provisioned: ${PROJECT}"
echo "  ✓ Domain:              ${DOMAIN}"
echo "  ✓ Internal port:       ${PORT}"
echo "  ✓ Caddy site file:     ${CADDY_SITE_FILE}"
echo "  ✓ Caddy reloaded"
echo "════════════════════════════════════════════════════════"
echo ""
echo "Add these secrets to your GitHub repository:"
echo "  Settings → Secrets and variables → Actions → New repository secret"
echo ""
echo "  VPS_HOST  →  <your VPS IP or hostname>"
echo "  VPS_USER  →  ${DEPLOY_USER}"
echo "  VPS_SSH_KEY  →  (private key printed below)"
echo ""
echo "─────────────────── VPS_SSH_KEY ────────────────────────"
echo "$PRIVATE_KEY"
echo "────────────────────────────────────────────────────────"
echo ""
echo "After copying the key to GitHub, delete it from the VPS:"
echo "  rm ${KEY_DIR}/deploy_key"
echo ""
echo "The public key remains at: ${KEY_DIR}/deploy_key.pub"
echo "Edit app secrets at:       ${OPT_DIR}/.env"
echo ""