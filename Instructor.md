# INSTRUCTIONS — Go Backend Agent

## Project: QRCheck — Digital Ticketing System (Backend)

You are a senior Go engineer building the backend for a B2B SaaS product called **QRCheck**. It replaces paper tickets (coat checks, bag storage, parking stubs, etc.) with signed dynamic QR codes. Businesses configure their services, employees issue digital tickets, and customers redeem them by showing a QR on their phone.

This is a real business product. Code quality, maintainability, and performance are non-negotiable. Write code as if it will be read by other senior engineers, extended by junior ones, and run in production serving thousands of concurrent users.

---

## Guiding Principles

- **Clarity over cleverness.** Every function should do one thing and be named for exactly that thing.
- **Explicit over implicit.** No magic, no global state. Dependencies are injected.
- **Fail loudly in development, gracefully in production.** Structured errors, never bare `panic`.
- **The database is the source of truth.** Slot state is always authoritative server-side.
- **Security is not an afterthought.** HMAC secrets never leave the server. Passwords are always bcrypt-hashed.

---

## Tech Stack

| Concern          | Choice                                                             | Reason                                        |
| ---------------- | ------------------------------------------------------------------ | --------------------------------------------- |
| Language         | Go 1.22+                                                           | Performance, simplicity, strong stdlib        |
| HTTP Framework   | [Fiber v2](https://github.com/gofiber/fiber)                       | Fast, Express-like, good middleware ecosystem |
| Database         | PostgreSQL 16                                                      | ACID, row-level locking for slot claims       |
| Query layer      | [sqlc](https://sqlc.dev/) + [pgx/v5](https://github.com/jackc/pgx) | Type-safe generated queries, no ORM magic     |
| Migrations       | [golang-migrate](https://github.com/golang-migrate/migrate)        | SQL-first, version-controlled schema          |
| Auth             | JWT (HS256) via [golang-jwt](https://github.com/golang-jwt/jwt)    | Stateless, simple two-role system             |
| QR Signing       | HMAC-SHA256 (stdlib `crypto/hmac`)                                 | Server-side only, payload integrity           |
| Config           | `.env` + [godotenv](https://github.com/joho/godotenv)              | Simple, 12-factor compliant                   |
| Logging          | [zerolog](https://github.com/rs/zerolog)                           | Structured JSON logs, zero alloc              |
| Containerization | Docker + docker-compose                                            | Reproducible local dev + easy deployment      |

---

## Project Structure

Follow this structure exactly. It is based on the standard Go project layout and clean architecture principles.

```
qrcheck-backend/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point. Wire dependencies, start server.
├── internal/
│   ├── config/
│   │   └── config.go                # Load and validate env vars into a typed Config struct
│   ├── database/
│   │   ├── db.go                    # pgx pool setup
│   │   ├── query/                   # sqlc-generated files (do not edit manually)
│   │   └── schema.sql               # Reference schema (source of truth for sqlc)
│   ├── domain/
│   │   ├── business.go              # Business entity + interface definitions
│   │   ├── service.go               # Service entity (e.g. "Cloakroom")
│   │   ├── slot.go                  # Slot entity
│   │   ├── ticket.go                # Ticket entity
│   │   └── customer.go              # Customer entity
│   ├── repository/
│   │   ├── business_repo.go         # DB operations for businesses
│   │   ├── service_repo.go          # DB operations for services + slot generation
│   │   ├── slot_repo.go             # DB operations for slots (with row locking)
│   │   ├── ticket_repo.go           # DB operations for tickets
│   │   └── customer_repo.go         # DB operations for customers
│   ├── usecase/
│   │   ├── auth_usecase.go          # Register, login, JWT issuance
│   │   ├── service_usecase.go       # Service CRUD + slot generation logic
│   │   ├── ticket_usecase.go        # Check-in, scan, release — core business logic
│   │   └── customer_usecase.go      # Customer registration and lookup
│   ├── handler/
│   │   ├── auth_handler.go          # HTTP handlers for auth routes
│   │   ├── service_handler.go       # HTTP handlers for service routes
│   │   ├── ticket_handler.go        # HTTP handlers for ticket routes
│   │   └── customer_handler.go      # HTTP handlers for customer routes
│   ├── middleware/
│   │   ├── auth.go                  # JWT validation middleware, role enforcement
│   │   ├── logger.go                # Request logging middleware
│   │   └── error_handler.go         # Centralized error → HTTP response mapping
│   ├── qr/
│   │   ├── payload.go               # QR payload struct definition
│   │   ├── signer.go                # HMAC sign + verify
│   │   └── encoder.go               # Encode payload to base64 string and back
│   └── apperror/
│       └── errors.go                # Typed application errors (NotFound, Conflict, Unauthorized, etc.)
├── migrations/
│   ├── 000001_init_schema.up.sql
│   ├── 000001_init_schema.down.sql
│   └── ...
├── docker-compose.yml
├── Dockerfile
├── .env.example
├── sqlc.yaml
├── Makefile                         # make migrate-up, make generate, make run, make test
└── README.md
```

---

## Database Schema

Use UUIDs everywhere. All timestamps are `TIMESTAMPTZ`. Design migrations to be reversible.

### `businesses`

```sql
CREATE TABLE businesses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    email       TEXT UNIQUE NOT NULL,
    password    TEXT NOT NULL,  -- bcrypt, min cost 12
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### `customers`

```sql
CREATE TABLE customers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       TEXT UNIQUE NOT NULL,
    password    TEXT NOT NULL,  -- bcrypt, min cost 12
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### `services`

A service is a configurable category a business offers. The cloakroom is one service. A festival might have "Bag Storage", "Bike Parking", and "Valuables Locker" as three separate services.

```sql
CREATE TABLE services (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    description TEXT,
    total_slots INT NOT NULL CHECK (total_slots > 0),
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### `slots`

Auto-generated when a service is created. Never manually inserted by handlers.

```sql
CREATE TABLE slots (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id  UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    slot_number INT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'free', -- CHECK: 'free' | 'occupied'
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (service_id, slot_number),
    CHECK (status IN ('free', 'occupied'))
);

CREATE INDEX idx_slots_service_status ON slots(service_id, status);
```

### `tickets`

```sql
CREATE TABLE tickets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id     UUID NOT NULL REFERENCES slots(id),
    service_id  UUID NOT NULL REFERENCES services(id),
    business_id UUID NOT NULL REFERENCES businesses(id),
    customer_id UUID REFERENCES customers(id), -- NULL for guests
    status      TEXT NOT NULL DEFAULT 'active', -- CHECK: 'active' | 'released'
    hmac        TEXT NOT NULL,
    issued_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    released_at TIMESTAMPTZ,
    CHECK (status IN ('active', 'released'))
);

CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_customer ON tickets(customer_id) WHERE customer_id IS NOT NULL;
```

---

## Domain Entities

Define clean structs in `internal/domain/`. These are the core types that flow through the whole application. Repository and usecase layers operate on these, not on raw DB types.

```go
// internal/domain/ticket.go
type Ticket struct {
    ID         uuid.UUID
    SlotID     uuid.UUID
    ServiceID  uuid.UUID
    BusinessID uuid.UUID
    CustomerID *uuid.UUID // nil for guests
    Status     TicketStatus
    HMAC       string
    IssuedAt   time.Time
    ReleasedAt *time.Time
}

type TicketStatus string
const (
    TicketStatusActive   TicketStatus = "active"
    TicketStatusReleased TicketStatus = "released"
)
```

Do the same for all entities. Use typed status strings, not raw strings in logic.

---

## QR Payload

The QR code in the Flutter app encodes a single URL-safe base64 string. This string is a signed JSON payload. The HMAC secret lives only on the server.

### Payload structure

```go
// internal/qr/payload.go
type Payload struct {
    Version    int       `json:"v"`    // Always 1 for now
    TicketID   string    `json:"tid"`
    ServiceID  string    `json:"sid"`
    BusinessID string    `json:"bid"`
    SlotNumber int       `json:"slot"`
    IssuedAt   int64     `json:"iat"`  // Unix timestamp
    HMAC       string    `json:"hmac"`
}
```

### Signing

HMAC is computed over a deterministic canonical string — field order is fixed:

```
v={v}&tid={tid}&sid={sid}&bid={bid}&slot={slot}&iat={iat}
```

```go
// internal/qr/signer.go
func Sign(p Payload, secret string) (string, error)
func Verify(p Payload, secret string) bool
```

### Encoding

```go
// internal/qr/encoder.go
func Encode(p Payload) (string, error)   // marshal JSON → base64 URL-safe string
func Decode(raw string) (Payload, error) // base64 decode → unmarshal JSON
```

The Flutter app encodes this string into a QR image. To scan, it reads the QR, gets the string back, sends it to `/tickets/scan` — the server decodes and verifies.

---

## Repository Layer

Repositories handle all database access. They accept a context and return domain types or `apperror` errors. They do not contain business logic.

Define interfaces in `internal/domain/` so usecases depend on the interface, not the concrete implementation. This makes testing trivial.

```go
// internal/domain/ticket.go
type TicketRepository interface {
    Create(ctx context.Context, t CreateTicketParams) (*Ticket, error)
    GetByID(ctx context.Context, id uuid.UUID) (*Ticket, error)
    Release(ctx context.Context, id uuid.UUID) (*Ticket, error)
    ListByCustomer(ctx context.Context, customerID uuid.UUID) ([]Ticket, error)
}
```

### Critical: Slot Claiming

Claiming a slot **must** use `SELECT ... FOR UPDATE SKIP LOCKED` to prevent two concurrent check-ins from claiming the same slot. This is the most important query in the system.

```sql
-- name: ClaimNextFreeSlot :one
UPDATE slots
SET status = 'occupied', updated_at = NOW()
WHERE id = (
    SELECT id FROM slots
    WHERE service_id = $1 AND status = 'free'
    ORDER BY slot_number
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;
```

Wrap the slot claim and ticket creation in a **single database transaction** in the usecase layer.

---

## Usecase Layer

Usecases contain all business logic. They orchestrate repositories, apply rules, and return domain types or errors. They know nothing about HTTP.

### `TicketUsecase`

```go
type TicketUsecase interface {
    CheckIn(ctx context.Context, req CheckInRequest) (*CheckInResult, error)
    Scan(ctx context.Context, qrPayload string) (*ScanResult, error)
    Release(ctx context.Context, ticketID uuid.UUID, businessID uuid.UUID) error
}
```

**CheckIn logic:**

1. Verify service exists and belongs to the requesting business
2. Begin transaction
3. Claim next free slot (SKIP LOCKED query)
4. If no free slot → return `apperror.Conflict("no free slots available")`
5. Build QR payload
6. Sign payload with HMAC
7. Create ticket record (store HMAC for audit)
8. Commit transaction
9. Encode payload to base64
10. Return `CheckInResult{TicketID, SlotNumber, QRPayload, IssuedAt}`

**Scan logic:**

1. Decode base64 QR payload
2. Verify HMAC — if invalid → `apperror.Unauthorized("invalid QR signature")`
3. Fetch ticket by ID from DB
4. If ticket status is `released` → `apperror.Conflict("ticket already released")`
5. Return ticket info (slot number, service name, issued at)

**Release logic:**

1. Fetch ticket, verify it belongs to the requesting business
2. Begin transaction
3. Set ticket `status = 'released'`, `released_at = NOW()`
4. Set slot `status = 'free'`, `updated_at = NOW()`
5. Commit

---

## Handler Layer

Handlers are thin. They parse + validate input, call the usecase, and serialize the response. No business logic here.

All handlers follow this pattern:

```go
func (h *TicketHandler) CheckIn(c *fiber.Ctx) error {
    var req CheckInRequest
    if err := c.BodyParser(&req); err != nil {
        return apperror.BadRequest("invalid request body")
    }
    if err := validate(req); err != nil {
        return err
    }
    businessID := extractBusinessID(c) // from JWT middleware
    result, err := h.usecase.CheckIn(c.Context(), req, businessID)
    if err != nil {
        return err // error handler middleware maps this to HTTP
    }
    return c.Status(fiber.StatusCreated).JSON(result)
}
```

---

## API Reference

All routes are prefixed `/api/v1`. Content-Type is `application/json`.

### Authentication

**POST `/api/v1/auth/business/register`**

```json
// Request
{ "name": "Fabric Berlin", "email": "ops@fabric.de", "password": "..." }
// Response 201
{ "business_id": "uuid", "token": "jwt" }
```

**POST `/api/v1/auth/business/login`**

```json
// Request
{ "email": "ops@fabric.de", "password": "..." }
// Response 200
{ "token": "jwt", "business": { "id": "uuid", "name": "Fabric Berlin" } }
```

**POST `/api/v1/auth/customer/register`**
**POST `/api/v1/auth/customer/login`**
Same pattern, role in JWT claim will be `customer`.

---

### Services `[Business JWT required]`

**POST `/api/v1/services`** — Create service + auto-generate slots

```json
// Request
{ "name": "Cloakroom", "description": "Ground floor", "total_slots": 150 }
// Response 201
{ "id": "uuid", "name": "Cloakroom", "total_slots": 150, "free_slots": 150 }
```

**GET `/api/v1/services`** — List all services for this business

```json
// Response 200
[
    {
        "id": "uuid",
        "name": "Cloakroom",
        "total_slots": 150,
        "free_slots": 112,
        "is_active": true
    }
]
```

**GET `/api/v1/services/:id`** — Service detail + slot summary

**PATCH `/api/v1/services/:id`** — Update name / description / is_active

**DELETE `/api/v1/services/:id`** — Soft delete (set is_active = false, preserve data)

---

### Tickets `[Business JWT required]`

**POST `/api/v1/tickets/checkin`**

```json
// Request
{ "service_id": "uuid", "customer_id": "uuid or omitted" }
// Response 201
{
  "ticket_id": "uuid",
  "slot_number": 42,
  "qr_payload": "base64url-encoded-signed-payload",
  "issued_at": "2024-03-10T22:00:00Z"
}
```

**POST `/api/v1/tickets/scan`**

```json
// Request
{ "qr_payload": "base64url-encoded-signed-payload" }
// Response 200
{
  "ticket_id": "uuid",
  "slot_number": 42,
  "status": "active",
  "service_name": "Cloakroom",
  "issued_at": "2024-03-10T22:00:00Z"
}
// Response 400 — invalid signature
// Response 409 — already released
```

**POST `/api/v1/tickets/release`**

```json
// Request
{ "ticket_id": "uuid" }
// Response 200
{ "ticket_id": "uuid", "slot_number": 42, "released_at": "2024-03-11T01:30:00Z" }
```

---

### Customer Wallet `[Customer JWT required]`

**GET `/api/v1/customer/tickets`** — List all active tickets for the logged-in customer

```json
[
    {
        "ticket_id": "uuid",
        "service_name": "Cloakroom",
        "business_name": "Fabric Berlin",
        "slot_number": 42,
        "status": "active",
        "qr_payload": "base64url-encoded-signed-payload",
        "issued_at": "2024-03-10T22:00:00Z"
    }
]
```

---

## Error Handling

Define typed errors in `internal/apperror/errors.go`. The Fiber error handler middleware maps these to HTTP responses.

```go
type AppError struct {
    Code    int
    Message string
    Err     error
}

func NotFound(msg string) *AppError       // → 404
func BadRequest(msg string) *AppError     // → 400
func Unauthorized(msg string) *AppError   // → 401
func Forbidden(msg string) *AppError      // → 403
func Conflict(msg string) *AppError       // → 409
func Internal(err error) *AppError        // → 500 (log the real error, return generic message)
```

Error response format:

```json
{ "error": "no free slots available", "code": 409 }
```

Never leak internal error details (stack traces, SQL errors) in production responses. Log them with zerolog.

---

## Middleware

### Auth Middleware

Extract and validate JWT from `Authorization: Bearer <token>`. Attach `business_id` or `customer_id` and `role` to the Fiber context. Reject missing or invalid tokens with 401.

### Role Middleware

```go
func RequireRole(role string) fiber.Handler
```

Wrap business-only routes with `RequireRole("business")` and customer-only with `RequireRole("customer")`.

### Logger Middleware

Log every request: method, path, status, latency, request ID. Use zerolog JSON format.

### Error Handler

Global Fiber error handler that maps `*AppError` to HTTP responses, and converts any other `error` to a 500 with a generic message (log the real one).

---

## Environment Variables

```env
# .env.example
APP_ENV=development          # development | production
APP_PORT=8080

DATABASE_URL=postgres://user:password@localhost:5432/qrcheck?sslmode=disable

JWT_SECRET=your-256-bit-secret-here
JWT_EXPIRY_HOURS=72

QR_HMAC_SECRET=your-separate-256-bit-secret-here

BCRYPT_COST=12
```

Always validate all required env vars on startup. Fail fast with a clear error message if any are missing.

---

## Makefile

Provide a `Makefile` with at minimum:

```makefile
run:          # go run cmd/server/main.go
build:        # go build -o bin/server cmd/server/main.go
test:         # go test ./...
migrate-up:   # run all pending migrations
migrate-down: # roll back last migration
generate:     # run sqlc generate
lint:         # golangci-lint run
docker-up:    # docker-compose up -d
docker-down:  # docker-compose down
```

---

## Testing Expectations

- Unit test every usecase method. Mock the repository interfaces.
- Integration test the critical path: check-in → scan → release, including the concurrent slot-claim scenario.
- Use `testcontainers-go` for integration tests that need a real postgres instance.
- Aim for >80% coverage on the `usecase` and `qr` packages.

---

## README

Write a `README.md` that covers:

1. Prerequisites (Go, Docker)
2. Local setup (`cp .env.example .env`, `make docker-up`, `make migrate-up`, `make run`)
3. API overview with curl examples for the key flows
4. Architecture overview (one paragraph per layer)
5. How to run tests
