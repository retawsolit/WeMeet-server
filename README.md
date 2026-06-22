# WeMeet Server

<div align="center">

**A Custom Self-Hosted Web Conferencing Backend Engine**

*Graduation Thesis Project - School of Computer Science and Engineering*

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue)](https://golang.org)
[![Platform](https://img.shields.io/badge/Architecture-Microservices-orange)](#)

</div>

##  Table of Contents

- [Introduction](#-introduction)
- [Core Features](#-core-features)
- [System Requirements](#-system-requirements)
- [Project Structure](#-project-structure)
- [Installation & Deployment](#-installation--deployment)
- [Configuration Guide](#-configuration-guide)
- [API Endpoints](#-api-endpoints)
- [System Architecture](#-system-architecture)
- [Troubleshooting](#-troubleshooting)
- [Academic Identification](#-academic-identification)

##  Introduction

**WeMeet Server** is the core backend API engine for the WeMeet self-hosted web conferencing platform. Built strictly from the ground up using **Go (Golang)**, the system manages session orchestration, user token authentication, signaling middleware, and direct integration with WebRTC SFU topologies via LiveKit. 

This project is developed as part of a professional graduation thesis, focusing on high-performance concurrent request processing and optimal resource utilization in real-time communication.

### Tech Stack Details:
- **Core Engine**: Go 1.24+ (Fiber v2 Web Framework)
- **Database Layer**: MariaDB (GORM ORM Layer)
- **State & Token Cache**: Redis
- **Pub/Sub Messaging Broker**: NATS Server
- **WebRTC Core Middleware**: LiveKit SFU
- **Security Protocols**: JSON Web Tokens (JWT) + Cryptographic Authentication

##  Core Features

### Meeting Room & Session Orchestration
-  Dynamic session creation, updates, and terminations.
-  Real-time room participant state management.
-  Breakout Room sub-session isolation (Up to 16 concurrent sub-rooms).
-  Waiting Room queue architecture for moderated participant access.
-  Custom room constraints (Maximum duration limits and capacity thresholds).

### Authentication & Authorization Middleware
-  Stateless authentication using secure JWT tokens.
-  Role-Based Access Control (RBAC) handling different participant privileges (Moderator vs. Attendee).
-  System-wide security configurations via asymmetric/symmetric secret keys.

### Media & Asset Gateway
-  Secure upload pipelines for presentation files and documents.
-  Core media state synchronization via asynchronous handlers.

##  System Requirements

### Hardware Specifications
- **CPU**: Minimal 2 Cores (Recommended 4+ Cores for high WebRTC concurrency)
- **RAM**: Minimal 4GB (Recommended 8GB+)
- **Storage**: 50GB+ SSD 

### Software Prerequisites
- Go 1.24 or higher
- Docker & Docker Compose v2+
- MariaDB 10.3+ / MySQL 5.7+
- Redis 6.0+
- NATS Server 2.9+
- LiveKit Server 1.5+

##  Project Structure

WeMeet-server/
├── main.go                          # Application entry point
├── config.yaml                      # Centralized configuration file
├── go.mod & go.sum                 # Go dependencies
├── Dockerfile                       # Docker production image
├── Dockerfile.dev                   # Docker development image
├── docker-compose.yaml              # Docker Compose infrastructure setup
├── Makefile                         # Build commands
├── livekit.yaml                     # LiveKit standalone configuration
├── nats_server.conf                 # NATS server configuration
│
├── pkg/                             # Main application package
│   ├── config/                      # Configuration management
│   ├── controllers/                 # HTTP request handlers
│   │   ├── room.go                  # Room core management
│   │   ├── user.go                  # User identity management
│   │   ├── auth.go                  # JWT Authentication control
│   │   ├── file.go                  # Document asset handling
│   │   ├── lti_v1.go               # LTI integration middleware
│   │   └── polls.go                 # In-room voting features
│   ├── services/                    # Business logic layer
│   │   ├── db/                      # Database operations (GORM framework)
│   │   ├── livekit/                 # LiveKit connection client
│   │   ├── nats/                    # NATS topology messaging
│   │   └── redis/                   # Redis cache repository
│   ├── models/                      # Data models & DTOs
│   ├── dbmodels/                    # Relational database models
│   ├── routers/                     # Route definitions
│   │   └── app_routers.go          # Central multiplexer router setup
│   ├── factory/                     # Dependency injection container
│   └── helpers/                     # Internal utility functions
│
├── helpers/                         # Root level server infrastructure hooks
│   ├── startup.go                   # Bootstrapping and validation logic
│   └── close_connections.go         # Cleanup handler on graceful shutdown
│
├── version/                         # Build version metadata
└── sql_dump/                        # Relational schema migrations
└── install.sql                  # Database schema layout init

##  Installation & Deployment

### 1. Repository Setup
```bash
git clone https://github.com/retawsolit/WeMeet-server.git
cd WeMeet-server

2. Infrastructure Initialization (Docker Compose)
The easiest way to boot dependencies (MariaDB, Redis, NATS, LiveKit) is using the pre-configured compose stack:

docker compose up -d db redis nats livekit

3. Schema Migration
Initialize the persistent storage layer by executing the SQL setup scripts against your database instance:

mysql -u root -p wemeet < sql_dump/install.sql

4. Compiling and Executing the Native Server
Download external dependencies and launch the backend binary:

go mod download
go mod tidy
go run main.go --config=config.yaml

Configuration Guide
Key structural blocks in config.yaml:

Server Security Identity
client:
  port: 8080                    # Application binding port
  debug: true                   # Toggles verbosity and developer utilities
  api_key: "wemeet"             # Public gateway API identifier
  secret: "your_cryptographic_secret_key" # Validates JWT token payloads
  token_validity: 10m           # Living span of issued access tokens

WebRTC Core SFU Linkage
livekit_info:
  host: "http://localhost:7880"
  api_key: "your_livekit_api_key"
  secret: "your_livekit_secret"

Persistent Relational Database Infrastructure

database_info:
  driver_name: mysql
  host: localhost
  port: 3306
  username: root
  password: secure_password
  database: wemeet
  max_open_connections: 50
  max_idle_connections: 10

API Endpoints
Session & Room Management
HTTP Method	API Path URI	Functionality Target
POST	/api/v1/room	Initializes a unique meeting session room
GET	/api/v1/room/:roomId	Retrieves configuration states for an active session
GET	/api/v1/room/:roomId/members	Streams current room participant list
POST	/api/v1/room/:roomId/join	Authorizes token generation for entering a room
POST	/api/v1/room/:roomId/leave	Clears active token states upon exiting
DELETE	/api/v1/room/:roomId	Triggers hard termination of room instance
POST	/api/v1/room/:roomId/breakout	Distributes members into sub-session pools

Identity Profiles
HTTP Method	API Path URI	Functionality Target
POST	/api/v1/user	Creates non-admin user identity profile
GET	/api/v1/user/:userId	Retrieves profile states
PUT	/api/v1/user/:userId	Updates profile metadata
DELETE	/api/v1/user/:userId	Revokes and soft-deletes profile identity

Session Token Guards
HTTP Method	API Path URI	Functionality Target
POST	/api/v1/auth/login	Validates security payload and issues JWT keys
POST	/api/v1/auth/logout	Revokes current token scopes
POST	/api/v1/auth/refresh	Re-issues tokens using valid refresh chains
GET	/api/v1/auth/verify	Validates signature structure of inbound requests

Asset Management
HTTP Method	API Path URI	Functionality Target
POST	/api/v1/file/upload	Multipart file upload pipeline for session assets
GET	/api/v1/file/:fileId	Streams specific session files or whiteboards
DELETE	/api/v1/file/:fileId	Garbage collection of unreferenced files

System Status Check
HTTP Method	API Path URI	Functionality Target
GET	/health	Lightweight service ping status payload check

System Architecture

┌─────────────────────────────────────────────┐
│          Web Clients (Browsers)             │
│     (wemeet-ui, WeMeet-client)              │
└────────────────┬────────────────────────────┘
                 │ HTTP / WebSocket Signaling
┌────────────────▼─────────────────────────────┐
│         WeMeet-Server (Fiber API Engine)     │
├──────────────────────────────────────────────┤
│  Controllers  │  Models  │  Services Layer   │
└────┬──────────┴────┬─────────────┬───────────┘
     │               │             │
┌────▼──────┐   ┌─────▼──┐  ┌───────▼────┐
│ Database  │   │  NATS  │  │  LiveKit   │
│ (MariaDB) │   │ Broker │  │(WebRTC SFU)│
└───────────┘   └───┬────┘  └────────────┘
                    │
             ┌──────▼───────┐
             │  Redis Cache │
             └──────────────┘

Troubleshooting
Common Infrastructure Bottlenecks
Database Connection Failure (failed to connect to database)

Verify that MariaDB is bound correctly to your local host container ports (3306).

Ensure you have cleanly executed the foundational scheme inside sql_dump/install.sql.

Signaling Broker Offline (NATS connection refused)

Verify NATS server runtime status using docker compose ps nats.

Match connection addresses specified in config.yaml to ensure cross-container host routing matches.