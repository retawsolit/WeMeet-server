# WeMeet Server - API Specification

## 1. Introduction

This document provides the definitive API specification for the WeMeet Server backend application. WeMeet Server exposes a RESTful API architecture using JSON (JavaScript Object Notation) as the primary data exchange format for all request and response payloads.

### Base URL
http://localhost:8080/api/v1


### Global Authentication Middleware
All protected endpoints require a stateless JSON Web Token (JWT) supplied in the HTTP request header:
Authorization: Bearer <your_jwt_token>


---

## 2. API Endpoints

### 2.1. Authentication Services (`/auth`)

#### User Login / Token Generation
- **Endpoint:** `POST /auth/login`
- **Description:** Validates user credentials and issues a secure JWT access token.
- **Request Payload:**
```json
{
  "username": "user@example.com",
  "password": "password123"
}
Success Response (200 OK):

JSON
{
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 600,
  "user": {
    "id": "user-123",
    "username": "user@example.com",
    "email": "user@example.com",
    "role": "moderator"
  }
}
Error Statuses:

400 Bad Request - Invalid credential structures

401 Unauthorized - Invalid username or password

500 Internal Server Error - Database connection fault

User Logout / Session Revocation
Endpoint: POST /auth/logout

Authentication: Required

Success Response (200 OK):

JSON
{
  "status": "success",
  "message": "Logged out successfully"
}
Token Security Verification
Endpoint: GET /auth/verify

Authentication: Required

Description: Verifies the cryptographic signature and expiration state of the current token.

Success Response (200 OK):

JSON
{
  "status": "success",
  "valid": true,
  "user_id": "user-123"
}
2.2. Room & Session Orchestration (/room)
Create Meeting Room
Endpoint: POST /room

Authentication: Required

Request Payload:

JSON
{
  "room_name": "Architecture Review Meeting",
  "max_participants": 50,
  "max_duration": 120
}
Success Response (200 OK):

JSON
{
  "status": "success",
  "room_id": "room-789",
  "room_name": "Architecture Review Meeting",
  "connection_url": "ws://localhost:8080/ws/room/room-789",
  "created_at": "2026-06-22T15:30:00Z"
}
Fetch Active Room State
Endpoint: GET /room/:roomId

Authentication: Required

Success Response (200 OK):

JSON
{
  "status": "success",
  "room": {
    "id": "room-789",
    "room_name": "Architecture Review Meeting",
    "status": "active",
    "participant_count": 5,
    "max_participants": 50
  }
}
Fetch Room Members List
Endpoint: GET /room/:roomId/members

Authentication: Required

Success Response (200 OK):

JSON
{
  "status": "success",
  "members": [
    {
      "user_id": "user-123",
      "username": "Nguyễn Sỹ Thủy",
      "role": "moderator",
      "joined_at": "2026-06-22T15:32:00Z"
    }
  ]
}
Request Room Join Token
Endpoint: POST /room/:roomId/join

Authentication: Required

Success Response (200 OK):

JSON
{
  "status": "success",
  "livekit_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.livekit-token-data..."
}
Terminate Room Session
Endpoint: DELETE /room/:roomId

Authentication: Required (Moderator Only)

Success Response (200 OK):

JSON
{
  "status": "success",
  "message": "Room session terminated successfully"
}
2.3. User Management (/user)
Register Identity Profile
Endpoint: POST /user

Request Payload:

JSON
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securepassword125"
}
Success Response (201 Created):

JSON
{
  "status": "success",
  "user_id": "user-456",
  "message": "User profile registered successfully"
}
Fetch Profile Metadata
Endpoint: GET /user/:userId

Authentication: Required

Success Response (200 OK):

JSON
{
  "status": "success",
  "user": {
    "id": "user-456",
    "username": "newuser",
    "email": "newuser@example.com",
    "role": "attendee"
  }
}
2.4. Asset & File Management (/file)
Upload Presentation Asset
Endpoint: POST /file/upload

Authentication: Required

Request Format: multipart/form-data

Payload Parameters: file (Binary Blob)

Success Response (200 OK):

JSON
{
  "status": "success",
  "file_id": "file-999",
  "file_name": "thesis_presentation.pdf",
  "storage_url": "http://localhost:8080/uploads/file-999.pdf"
}
Delete Session Asset
Endpoint: DELETE /file/:fileId

Authentication: Required

Success Response (200 OK):

JSON
{
  "status": "success",
  "message": "Asset deleted from persistent storage"
}
2.5. In-Room Polls Service (/room/:roomId/polls)
Create Interactive Poll
Endpoint: POST /room/:roomId/polls

Authentication: Required

Request Payload:

JSON
{
  "question": "Does the system meet performance constraints?",
  "options": ["Yes, fully", "Partially", "No"]
}
Success Response (201 Created):

JSON
{
  "status": "success",
  "poll_id": "poll-111",
  "question": "Does the system meet performance constraints?"
}
2.6. System Infrastructure Status
Application Health Check
Endpoint: GET /health

Authentication: None

Success Response (200 OK):

JSON
{
  "status": "healthy",
  "timestamp": "2026-06-22T15:45:00Z",
  "services": {
    "database": "connected",
    "redis": "connected",
    "nats": "connected"
  }
}
3. Global Error Handling Topography
When an API transaction encounters an error condition, the system guarantees a standardized response envelope:

Error Payload Format
JSON
{
  "status": "error",
  "error_code": "RESOURCE_NOT_FOUND",
  "message": "The requested entity room-999 does not exist in the active cluster state."
}
Standard Status Codes
400 Bad Request: Payload validation constraints failed.

401 Unauthorized: Invalid or expired JWT authentication token.

403 Forbidden: Insufficient Role-Based Access Control privileges.

404 Not Found: Target entity or route does not exist.

500 Internal Server Error: Unhandled failure within upstream dependency stacks.