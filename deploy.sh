#!/bin/bash

# TODO:
# Implement options: deploy.sh reinstall & deploy.sh uninstall

set -e

echo "📦 Packaging and deploying Go Orchestrator..."

# 1. Compile the app
go build -o dockerOrchestrator ./cmd/dockerOrchestrator/dockerOrchestrator.go

# 2. Move the binary to the global system path
sudo mv dockerOrchestrator /usr/local/bin/dockerOrchestrator
sudo chmod +x /usr/local/bin/dockerOrchestrator

# 3. Create a default environment file if it doesn't exist
if [ ! -f /etc/dockerOrchestrator.env ]; then
  echo "Creating default environment file at /etc/dockerOrchestrator.env"
  sudo bash -c 'cat << EOF > /etc/dockerOrchestrator.env
APP_ENV=prod
APP_PATH=/srv/dockerOrchestrator
BLUEPRINT_PATH=${APP_PATH}/blueprints
DATA_PATH=${APP_PATH}/data
MEDIA_PATH=/srv/media
PUBLIC_PATH=/srv/public
EOF'
fi

# 4. Copy the systemd service file into place
sudo bash -c 'cat << EOF > /etc/systemd/system/dockerOrchestrator.service
[Unit]
Description=My Go Container & Disk Orchestrator
After=docker.service
Requires=docker.service

[Service]
Type=simple
ExecStart=/usr/local/bin/dockerOrchestrator
EnvironmentFile=/etc/dockerOrchestrator.env
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF'

# 5. Tell systemd to reload its configuration and restart the app
sudo systemctl daemon-reload
sudo systemctl enable --now dockerOrchestrator

echo "Service is running! Check logs with: journalctl -u dockerOrchestrator -f"
