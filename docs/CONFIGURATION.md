# WeMeet Server - Configuration Guide

This guide provides a detailed specification of the centralized configuration topology for the WeMeet Server backend framework. The system processes environmental states utilizing a single structured file named `config.yaml`.

## Table of Contents

1. [Configuration Initialization](#1-configuration-initialization)
2. [Global Parameter Blueprint](#2-global-parameter-blueprint)
3. [Detailed Section Reference](#3-detailed-section-reference)
   - [Client & Gateway Settings](#client--gateway-settings)
   - [Room Default Parameters](#room-default-parameters)
   - [Database Connectivity](#database-connectivity)
   - [Distributed Cache Layer](#distributed-cache-layer)
   - [Message Broker Brokerage](#message-broker-brokerage)
   - [WebRTC SFU Infrastructure](#webrtc-sfu-infrastructure)
4. [Environment Variable Substitution](#4-environment-variable-substitution)
5. [Configuration Validation](#5-configuration-validation)

---

## 1. Configuration Initialization

By default, the compiled Go binary searches for a `config.yaml` file within the execution root folder. To boot the server with a custom structured configuration file, append the execution flag:

```bash
./WeMeet-server --config=/path/to/custom-config.yaml

2. Global Parameter Blueprint
The absolute layout schema of config.yaml for the core meeting infrastructure is defined below:
client:
  port: 8080
  debug: true
  api_key: "wemeet"
  secret: "your_jwt_signing_secret"
  token_validity: 10m

room_default_settings:
  max_participants: 50
  max_duration: 120m
  allow_polls: true

log_settings:
  level: "info"
  format: "json"

database_info:
  driver_name: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "secure_password"
  database: "wemeet"
  max_open_connections: 50
  max_idle_connections: 10

redis_info:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

nats_info:
  servers:
    - "nats://localhost:4222"

livekit_info:
  host: "http://localhost:7880"
  api_key: "devkey"
  secret: "secretkey"

  3. Detailed Section Reference
Client & Gateway Settings
Manages endpoint access bindings and security criteria for verifying core JavaScript Token payloads.

port (int): The network port the Fiber application listens to.

debug (bool): Toggles verbose output engine flags. Must be set to false in execution deployment.

api_key (string): Standard validation string for gateway filtering.

secret (string): The cryptographic seed phrase used to sign and unpack JWT tokens.

token_validity (duration): Lifespan threshold of temporary access signatures (e.g., 10m, 1h).

Room Default Parameters
Applies fallback limitations on runtime WebRTC channel allocations initiated by endpoint requests.

max_participants (int): Hard participant volume ceiling allowed inside a single sub-room.

max_duration (duration): Default expiration time window for active channels.

allow_polls (bool): Enables or disables the creation of interactive user survey structures.

Database Connectivity
Establishes state persistence mapping targeting the relational MariaDB service infrastructure layer.

driver_name (string): Database driver type, assigned strictly to "mysql".

host / port: Destination routing coordinates for the persistent server instance.

username / password: Security credentials for handling GORM ORM execution routines.

database (string): Target schema identifier assigned to hold WeMeet tables.

max_open_connections / max_idle_connections: Adjusts connection pooling thresholds to handle rapid real-time concurrency.

Distributed Cache Layer
Configures volatile key-value tracking to hold active authentication status arrays using Redis.

host / port: Location of the transient database memory engine instance.

password (string): Set to empty string if authentication is bypassed.

db (int): Target logical cache index allocation.

Message Broker Brokerage
Defines connection coordinates pointing to the asynchronous NATS Pub/Sub event streaming middleware core.

servers (array): A collection of qualified NATS server network addresses (nats://...) to stream events between the gateway microservices.

WebRTC SFU Infrastructure
Stores the validation secrets required to generate LiveKit interaction permissions for real-time video/audio streaming.

host (string): Web address mapping to the active LiveKit cluster instance.

api_key / secret: Key signatures generated from the livekit.yaml blueprint used to issue participant join tokens.

4. Environment Variable Substitution
To maintain operational security and prevent leakage of production secrets into version control systems, WeMeet Server allows replacing parameters inside config.yaml using system environment variables:

database_info:
  host: "${DB_HOST}"
  username: "${DB_USER}"
  password: "${DB_PASSWORD}"

livekit_info:
  api_key: "${LIVEKIT_API_KEY}"
  secret: "${LIVEKIT_SECRET}"
Before initializing the compiled execution file, export the variables to the active environment thread:

export DB_PASSWORD="my_production_password"
export LIVEKIT_SECRET="my_livekit_secure_sfu_secret"

5. Configuration Validation
To ensure structural consistency prior to launching the platform network stack, you can validate the syntax accuracy of your YAML tree design by utilizing standard parser formatting utilities such as yq:

# Verify YAML layout syntax
yq eval '.' config.yaml > /dev/null && echo "Structure Status: Valid"
