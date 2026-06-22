# WeMeet Server

<div align="center">

**Hệ thống hội thảo web mã nguồn mở, có thể mở rộng, hiệu suất cao**

[English](./README.md) | [Tiếng Việt](./README.vi.md)

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-green)](LICENSE)
[![Build Status](https://github.com/retawsolit/WeMeet-server/actions/workflows/test.yml/badge.svg)](https://github.com/retawsolit/WeMeet-server/actions)

</div>

##  Mục lục

- [Giới thiệu](#-giới-thiệu)
- [Tính năng chính](#-tính-năng-chính)
- [Yêu cầu hệ thống](#-yêu-cầu-hệ-thống)
- [Cấu trúc dự án](#-cấu-trúc-dự-án)
- [Cài đặt](#-cài-đặt)
- [Cấu hình](#-cấu-hình)
- [Chạy ứng dụng](#-chạy-ứng-dụng)
- [API Endpoints](#-api-endpoints)
- [Kiến trúc](#-kiến-trúc)
- [Phát triển](#-phát-triển)
- [Docker](#-docker)
- [Troubleshooting](#-troubleshooting)
- [Đóng góp](#-đóng-góp)

##  Giới thiệu

**WeMeet Server** là backend API cho nền tảng hội thảo web tự lưu trữ, được xây dựng bằng Go. Nó cung cấp các dịch vụ mạnh mẽ để quản lý phòng họp, xác thực người dùng, quản lý tài nguyên phương tiện, và tích hợp với các dịch vụ truyền phát như LiveKit.

### Thông tin chi tiết:
- **Nền tảng**: Go 1.24+
- **Framework Web**: Fiber v2
- **Cơ sở dữ liệu**: MySQL/MariaDB
- **Message Broker**: NATS
- **Cache**: Redis
- **WebRTC Server**: LiveKit
- **Chứng thực**: JWT + OAuth2

##  Tính năng chính

### Quản lý Phòng họp
-  Tạo, cập nhật, xóa phòng họp
-  Quản lý thành viên phòng họp
-  Hỗ trợ phòng học nhóm (Breakout Rooms)
-  Giới hạn thời gian và số lượng người tham gia
-  Chế độ chờ vào phòng (Waiting Room)

### Xác thực & Bảo mật
-  Xác thực dựa trên JWT Token
-  Hỗ trợ OAuth2 (Google, Microsoft, v.v.)
-  Tích hợp LTI v1.0 cho các nền tảng e-learning
-  Kiểm soát truy cập theo vai trò (RBAC)
-  Webhook hỗ trợ cho các sự kiện hệ thống

### Quản lý Phương tiện
-  Quản lý hình ảnh và tài liệu
-  Hỗ trợ tải lên tập tin
-  Chuyển đổi tài liệu (PDF, hình ảnh)
-  Streaming phương tiện tổng hợp

### Tính năng Nâng cao
-  Ghi âm phiên họp với LiveKit
-  Tích hợp Etherpad cho ghi chú cộng tác
-  Nhập RTMP cho các buổi phát trực tiếp
-  Text-to-Speech (Chuyển đổi từ văn bản thành giọng nói)
-  Bình chọn và thăm dò ý kiến
-  Phân tích & đo lường (Prometheus metrics)

##  Yêu cầu hệ thống

### Tối thiểu
- **CPU**: 2 cores
- **RAM**: 4GB
- **Disk**: 50GB
- **OS**: Linux (Ubuntu 18.04+), macOS hoặc Windows (WSL2)

### Khuyến nghị
- **CPU**: 4+ cores
- **RAM**: 8GB+
- **Disk**: 200GB+ (tùy theo nhu cầu ghi âm)
- **Kết nối**: Tối thiểu 1Mbps (khuyến nghị 5Mbps+)

### Phần mềm
- Go 1.24 hoặc cao hơn
- Docker & Docker Compose (cho triển khai containerized)
- MySQL 5.7+ hoặc MariaDB 10.3+
- Redis 6.0+
- NATS 2.9+
- LiveKit 1.5+

##  Cấu trúc dự án

```
WeMeet-server/
├── main.go                          # Entry point
├── config.yaml                      # Tệp cấu hình
├── go.mod & go.sum                 # Go dependencies
├── Dockerfile                       # Docker production image
├── Dockerfile.dev                   # Docker development image
├── docker-compose.yaml              # Docker Compose setup
├── Makefile                         # Build commands
├── livekit.yaml                     # LiveKit configuration
├── ingress.yaml                     # RTMP Ingress configuration
├── nats_server.conf                 # NATS server configuration
│
├── pkg/                             # Main application package
│   ├── config/                      # Configuration management
│   ├── controllers/                 # HTTP request handlers
│   │   ├── room.go                  # Room management
│   │   ├── user.go                  # User management
│   │   ├── auth.go                  # Authentication
│   │   ├── file.go                  # File handling
│   │   ├── recording.go             # Recording management
│   │   ├── etherpad.go              # Etherpad integration
│   │   ├── lti_v1.go               # LTI v1.0 integration
│   │   ├── polls.go                 # Polls feature
│   │   ├── analytics.go             # Analytics tracking
│   │   └── ... (other controllers)
│   ├── services/                    # Business logic layer
│   │   ├── db/                      # Database operations
│   │   ├── livekit/                 # LiveKit client
│   │   ├── nats/                    # NATS messaging
│   │   └── redis/                   # Redis cache
│   ├── models/                      # Data models
│   ├── dbmodels/                    # Database ORM models (GORM)
│   ├── routers/                     # Route definitions
│   │   └── app_routers.go          # Main router setup
│   ├── factory/                     # Dependency injection
│   └── helpers/                     # Utility functions
│
├── helpers/                         # Root level helpers
│   ├── startup.go                   # Server startup logic
│   └── close_connections.go         # Cleanup on shutdown
│
├── version/                         # Version management
├── sql_dump/                        # SQL migration files
│   └── install.sql                  # Database schema
│
└── docs/                            # Documentation
    └── ETHERPAD.md                  # Etherpad integration guide
```

##  Cài đặt

### 1. Clone Repository

```bash
git clone https://github.com/retawsolit/WeMeet-server.git
cd WeMeet-server
```

### 2. Cài đặt Dependencies

```bash
go mod download
go mod tidy
```

### 3. Thiết lập Cơ sở dữ liệu

```bash
# Tạo database
mysql -u root -p < sql_dump/install.sql

# Hoặc sử dụng Docker (được khuyến nghị)
docker compose up -d db
```

### 4. Cài đặt Dependencies Khác

```bash
# Redis
docker compose up -d redis

# NATS
docker compose up -d nats

# LiveKit
docker compose up -d livekit
```

##  Cấu hình

### Tệp config.yaml

Tệp `config.yaml` là file cấu hình chính. Dưới đây là các mục cấu hình quan trọng:

#### Client Configuration
```yaml
client:
  port: 8080                    # Server port
  debug: true                   # Debug mode
  path: "/app/client/dist"      # Client static files path
  api_key: "wemeet"             # API key (public ID)
  secret: "your_secret_key"     # Secret key (phải bảo mật)
  token_validity: 10m           # Token validity duration
  bbb_join_host: "http://your-host:8080"  # Host for BBB join URLs
```

#### Room Default Settings
```yaml
room_default_settings:
  max_duration: 0              # Max room duration (0 = unlimited)
  max_participants: 0          # Max participants (0 = unlimited)
  max_num_breakout_rooms: 6    # Max breakout rooms (max 16)
```

#### LiveKit Configuration
```yaml
livekit_info:
  host: "http://livekit-server:7880"
  api_key: "your_livekit_api_key"
  secret: "your_livekit_secret"
```

#### Database Configuration
```yaml
database_info:
  driver_name: mysql
  host: localhost
  port: 3306
  username: root
  password: password
  database: wemeet
  max_open_connections: 100
  max_idle_connections: 10
```

#### Redis Configuration
```yaml
redis_info:
  host: localhost:6379
  username: ""
  password: ""
  db: 0
  # Hoặc sử dụng Redis Sentinel:
  # sentinel_master_name: wemeet
  # sentinel_addresses:
  #   - redis-sentinel-1:26379
  #   - redis-sentinel-2:26379
```

#### NATS Configuration
```yaml
nats_info:
  host: "nats://localhost:4222"
  jwt_token: ""
  seed: ""
```

#### Logging Configuration
```yaml
log_settings:
  log_file: "./log/WeMeet.log"
  maxsize: 20                  # Max log file size (MB)
  maxbackups: 4                # Number of backup logs
  maxage: 2                    # Max age (days)
  log_level: "info"            # Log level (debug/info/warn/error)
```

#### Webhook Configuration (Tùy chọn)
```yaml
webhook_conf:
  enable: false
  url: "https://your-webhook-endpoint.com/wemeet"
  enable_for_per_meeting: false
```

#### Prometheus Configuration
```yaml
prometheus:
  enable: false
  metrics_path: "/metrics"
```

### Tạo Secret Key

Để tạo một secret key an toàn:

```bash
# Sử dụng OpenSSL
openssl rand -hex 32

# Hoặc
cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 36 | head -n 1
```

## Chạy ứng dụng

### Chạy từ Source

```bash
# Với config mặc định (config.yaml)
go run main.go

# Với config tùy chỉnh
go run main.go --config=/path/to/config.yaml
```

### Build Binary

```bash
# Tất cả platforms
make releases

# Chỉ Linux AMD64
make linux-amd64

# Chỉ Linux ARM64 (cho Raspberry Pi, v.v.)
make linux-arm64
```

### Chạy Binary

```bash
./bin/WeMeet-server-linux-amd64
./bin/WeMeet-server-linux-amd64 --config=/path/to/config.yaml
```

##  API Endpoints

### Room Management
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | `/api/v1/room` | Tạo phòng mới |
| GET | `/api/v1/room/:roomId` | Lấy thông tin phòng |
| GET | `/api/v1/room/:roomId/members` | Lấy danh sách thành viên |
| POST | `/api/v1/room/:roomId/join` | Tham gia phòng |
| POST | `/api/v1/room/:roomId/leave` | Rời phòng |
| DELETE | `/api/v1/room/:roomId` | Xóa phòng |
| POST | `/api/v1/room/:roomId/breakout` | Tạo phòng học nhóm |

### User Management
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | `/api/v1/user` | Tạo người dùng mới |
| GET | `/api/v1/user/:userId` | Lấy thông tin người dùng |
| PUT | `/api/v1/user/:userId` | Cập nhật người dùng |
| DELETE | `/api/v1/user/:userId` | Xóa người dùng |

### Authentication
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | `/api/v1/auth/login` | Đăng nhập |
| POST | `/api/v1/auth/logout` | Đăng xuất |
| POST | `/api/v1/auth/refresh` | Làm mới token |
| GET | `/api/v1/auth/verify` | Xác minh token |

### Recording
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | `/api/v1/recording/:roomId` | Lấy danh sách ghi âm |
| DELETE | `/api/v1/recording/:recordingId` | Xóa ghi âm |
| GET | `/api/v1/recording/:recordingId/download` | Tải xuống ghi âm |

### File Management
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | `/api/v1/file/upload` | Tải lên tập tin |
| GET | `/api/v1/file/:fileId` | Tải xuống tập tin |
| DELETE | `/api/v1/file/:fileId` | Xóa tập tin |

### Health Check
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | `/health` | Kiểm tra trạng thái server |

### Metrics (nếu bật Prometheus)
| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | `/metrics` | Prometheus metrics |

## Kiến trúc

### Tổng quan kiến trúc

```
┌─────────────────────────────────────────────┐
│          Web Clients (Browsers)             │
│     (wemeet-ui, WeMeet-client)              │
└────────────────┬────────────────────────────┘
                 │ HTTP/WebSocket
┌─────────────────▼────────────────────────────┐
│         WeMeet-Server (Fiber API)            │
├──────────────────────────────────────────────┤
│  Controllers  │  Models  │  Services         │
└────┬──────────┴────┬─────────────┬───────────┘
     │               │             │
┌────▼─────┐   ┌─────▼──┐  ┌──────▼────┐
│ Database  │  │  NATS  │  │  LiveKit  │
│(MySQL)    │  │ Broker │  │(WebRTC)   │
└───────────┘  └───┬────┘  └───────────┘
                   │
            ┌──────▼───────┐
            │  Redis Cache │
            └───────────────┘
```

### Layer Architecture

1. **Controller Layer** (`pkg/controllers/`)
   - Xử lý HTTP requests
   - Validation input
   - Trả về responses

2. **Service Layer** (`pkg/services/`)
   - Business logic
   - Database operations
   - Tích hợp với external services

3. **Model Layer** (`pkg/models/`, `pkg/dbmodels/`)
   - Data structures
   - ORM models (GORM)

4. **Factory Layer** (`pkg/factory/`)
   - Dependency injection
   - Component initialization

### Thành phần chính

#### 1. **Room Service**
- Quản lý vòng đời phòng họp
- Xử lý tham gia/rời phòng
- Quản lý phòng học nhóm

#### 2. **Authentication Service**
- JWT token generation
- Token validation
- OAuth2 integration
- LTI support

#### 3. **LiveKit Integration**
- Tạo access tokens cho WebRTC
- Quản lý media streams
- Recording management

#### 4. **NATS Messaging**
- Real-time event publishing
- Message queue cho xử lý async
- Broadcasting tin nhắn

#### 5. **Database Layer**
- User management
- Room configuration
- Recording metadata
- Session history

#### 6. **Cache Layer (Redis)**
- Token caching
- Session management
- Real-time data caching


### Service Ports

| Service | Port | URL |
|---------|------|-----|
| WeMeet API | 8080 | http://localhost:8080 |
| NATS | 4222 | nats://localhost:4222 |
| LiveKit | 7880 | http://localhost:7880 |
| Etherpad | 9001 | http://localhost:9001 |
| MariaDB | 3306 | localhost:3306 |
| Redis | 6379 | localhost:6379 |
