# WeMeet Server - Installation & Deployment Guide

This document provides a comprehensive technical manual for installing, compiling, and deploying the WeMeet Server backend application across various development environments.

## Table of Contents

1. [Quick Start with Docker Compose](#1-quick-start-with-docker-compose)
2. [Native Linux Installation (Ubuntu/Debian)](#2-native-linux-installation-ubuntudebian)
3. [Native macOS Installation](#3-native-mac-os-installation)
4. [Windows Workstation Installation](#4-windows-workstation-installation)
5. [Production Reverse Proxy Setup (Nginx)](#5-production-reverse-proxy-setup-nginx)
6. [Infrastructure Troubleshooting](#6-infrastructure-troubleshooting)

---

## 1. Quick Start with Docker Compose

The most efficient pipeline to initialize the local WeMeet testing ecosystem is leveraging Docker orchestration tools to boot the microservice architecture dependencies.

### Prerequisites
- Docker Engine 20.10+
- Docker Compose v2.0+

### Step-by-Step Execution
1. Clone the master code repository workspace:
   ```bash
   git clone [https://github.com/retawsolit/WeMeet-server.git](https://github.com/retawsolit/WeMeet-server.git)
   cd WeMeet-server

Boot the persistent database, caching nodes, messaging bus, and WebRTC SFU layers in detached background mode:

docker compose up -d db redis nats livekit

Validate container operational status runtime mapping:

docker compose ps

Expected active output matrix:

NAME                SERVICE             STATUS      PORTS
wemeet-redis        redis               Up          6379/tcp
wemeet-db           db                  Up          3306/tcp
wemeet-nats         nats                Up          4222/tcp, 8222/tcp
wemeet-livekit      livekit             Up          7880/tcp, 7881/tcp, 7882/udp
wemeet-api          wemeet-api          Up          8080/tcp

Perform a network socket ping verification to confirm application layer readiness:

curl http://localhost:8080/health

2. Native Linux Installation (Ubuntu/Debian)
Native Prerequisites Configuration
Ensure your operating system package list contains the necessary compiler stacks and service daemons:

sudo apt update
sudo apt install -y git golang-1.24 mariadb-server redis-server wget

# Append compiled Go binaries toolchain path variables
export PATH=$PATH:/usr/lib/go-1.24/bin

Step 1: Manage Core Dependencies

go mod download
go mod tidy

Step 2: Establish the MariaDB Schema Layer
Start the relational runtime daemon and create the authorization access rules:

sudo systemctl start mariadb

mysql -u root -p << EOF
CREATE DATABASE wemeet CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'wemeet'@'localhost' IDENTIFIED BY 'secure_password';
GRANT ALL PRIVILEGES ON wemeet.* TO 'wemeet'@'localhost';
FLUSH PRIVILEGES;
EOF

# Inject database layout structure blueprints
mysql -u wemeet -p wemeet < sql_dump/install.sql

Step 3: Run Standalone Signaling Middleware (NATS)

# Fetch official stable binary build distribution
wget [https://github.com/nats-io/nats-server/releases/download/v2.11.0/nats-server-v2.11.0-linux-amd64.tar.gz](https://github.com/nats-io/nats-server/releases/download/v2.11.0/nats-server-v2.11.0-linux-amd64.tar.gz)
tar -xzf nats-server-v2.11.0-linux-amd64.tar.gz
sudo mv nats-server-v2.11.0-linux-amd64/nats-server /usr/local/bin/

# Bind background service parameters
sudo mkdir -p /etc/nats
sudo cp nats_server.conf /etc/nats/

Step 4: Compiling the Native Binary Execution Target

# Execute local cross-platform make compilation routine
make linux-amd64
# Target binary is output to: bin/WeMeet-server-linux-amd64

3. Native macOS Installation
Prerequisites Setup via Homebrew

brew install go mariadb redis nats-server wget

Execution Strategy
Fire background resource runtime engines:

brew services start mariadb
brew services start redis
brew services start nats-server

Build the server instance using native system architecture flags:

cp config_sample.yaml config.yaml
go build -o WeMeet-server main.go
./WeMeet-server --config=config.yaml

4. Windows Workstation Installation
Preferred Architecture: Docker Desktop
The recommended workflow on Windows platforms is operating the stack directly within WSL2 abstraction kernels using Docker Desktop:

git clone [https://github.com/retawsolit/WeMeet-server.git](https://github.com/retawsolit/WeMeet-server.git)
cd WeMeet-server
docker compose up -d

5. Production Reverse Proxy Setup (Nginx)
To secure client signaling transport layers using SSL/TLS, terminate external traffic at an Nginx proxy gateway container node.

upstream wemeet_api_upstream {
    server localhost:8080;
}

server {
    listen 80;
    server_name wemeet.yourdomain.edu.vn;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name wemeet.yourdomain.edu.vn;

    ssl_certificate /etc/letsencrypt/live/wemeet.yourdomain.edu.vn/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/wemeet.yourdomain.edu.vn/privkey.pem;

    location / {
        proxy_pass http://wemeet_api_upstream;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Mandatory WebSocket Protocol Upgrade Links
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}

6. Infrastructure Troubleshooting
Relational Database Handshake Terminated (failed to connect to database)
Root Cause: MariaDB instance port wrapper boundaries (3306) are offline or blocked.

Resolution Check: Inspect server metrics logs using docker compose logs db or verify manual access credentials via mysql -h 127.0.0.1 -u wemeet -p.

Signaling Route Broker Missing (NATS connection refused)
Root Cause: The Go server driver is initializing faster than the upstream NATS broker network allocation layer.

Resolution Check: Configure manual sleep loops inside system systemd service blocks or use docker compose restart wemeet-api to reset active network bridges.

