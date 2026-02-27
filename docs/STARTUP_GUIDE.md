# CLOAK Backend - Setup & Startup Guide ✅

## ✅ Backend Status

Your Go/Fiber backend is **FULLY IMPLEMENTED** with all necessary components:

### Architecture ✅

- ✅ **Fiber Web Framework** - Fast HTTP server
- ✅ **PostgreSQL Database** - pgx connection pool with full schema
- ✅ **JWT Authentication** - Business & customer login with tokens
- ✅ **Domain Models** - Businesses, Customers, Services, Slots, Tickets
- ✅ **Repository Layer** - PostgreSQL implementations
- ✅ **Use Cases** - Auth, Tickets, Services business logic
- ✅ **Handlers** - All REST endpoints implemented
- ✅ **Middleware** - Auth, Role-based access, CORS, logging

### Database Schema ✅

```sql
Tables:
  ├── businesses      (company accounts)
  ├── customers       (customer accounts)
  ├── services        (events/experiences)
  ├── slots           (capacity management)
  └── tickets         (QR codes for entry)
```

### API Endpoints ✅

```
PUBLIC (No Auth)
  POST /api/v1/auth/business/register        Register business account
  POST /api/v1/auth/business/login           Business login
  POST /api/v1/auth/customer/login           Customer login

PROTECTED (Business role)
  POST /api/v1/tickets/checkin               Create QR ticket
  POST /api/v1/tickets/scan                  Scan & verify QR
  POST /api/v1/tickets/:id/release           Mark ticket as used
  POST /api/v1/services                      Create new service
  GET  /api/v1/services                      List business services
  GET  /api/v1/services/:id                  Get service details
  GET  /api/v1/services/:id/stats            Get service statistics

PROTECTED (Customer role)
  GET  /api/v1/tickets/:id                   Get ticket details

HEALTH CHECK
  GET /health                                System health check
```

---

## 🚀 How to Start the Backend

### Option 1: Direct Execution (Fastest) ⚡

**Step 1: Set up PostgreSQL**

```bash
# Using Docker (easiest):
docker run -d \
  --name postgres-cloak \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=cloak_db \
  -p 5432:5432 \
  postgres:16

# Wait for PostgreSQL to be ready (10 seconds)
sleep 10
```

**Alternative: Using local PostgreSQL (if already installed)**

```bash
# Create database
createdb cloak_db

# Check it's working
psql -l | grep cloak_db
```

**Step 2: Create .env file**

```bash
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE
cat > .env << 'EOF'
# Server
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/cloak_db?sslmode=disable

# JWT
JWT_SECRET=your-secret-jwt-key-12345-change-in-production
HMAC_SECRET=your-secret-hmac-key-12345-change-in-production
EOF
```

**Step 3: Run Database Migrations**

```bash
# Install migrate tool (if not already installed)
brew install golang-migrate

# Run migrations
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/cloak_db?sslmode=disable" up
```

**Step 4: Start the Backend**

```bash
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE
./bin/api
```

Expected output:

```
🚀 Starting server on port 8080
```

---

### Option 2: Using Make Commands 🏗️

**Quick setup:**

```bash
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE

# Setup everything (installs deps, creates DB, runs migrations)
make local-setup

# Wait a few seconds, then:
make local-migrate

# Start the app
make local-run
```

---

### Option 3: Using Docker Compose 🐳

**One command setup:**

```bash
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE
docker-compose up -d
```

This will:

- Start PostgreSQL container
- Build and start API container
- API runs on http://localhost:8080

View logs:

```bash
docker-compose logs -f
```

Stop:

```bash
docker-compose down
```

---

## ✅ Test the Backend

Once running, test with these commands:

### 1. Health Check

```bash
curl http://localhost:8080/health
# Expected: {"status":"ok"}
```

### 2. Register Business Account

```bash
curl -X POST http://localhost:8080/api/v1/auth/business/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Business",
    "email": "business@test.com",
    "password": "password123"
  }'
```

### 3. Business Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/business/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "business@test.com",
    "password": "password123"
  }'
# Returns: {"token": "eyJ...", "refresh_token": "..."}
```

### 4. Create Service (requires JWT token from login)

```bash
# Replace YOUR_TOKEN with the token from login response
curl -X POST http://localhost:8080/api/v1/services \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Concert Event",
    "total_slots": 100
  }'
```

### 5. Create QR Ticket (Check-in)

```bash
curl -X POST http://localhost:8080/api/v1/tickets/checkin \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "service_id": "service-uuid-here",
    "customer_id": "customer-uuid-here"
  }'
# Returns: QR code HMAC digest for ticket
```

---

## 📋 Complete Startup Checklist

- [ ] PostgreSQL running (docker or local)
- [ ] `.env` file created with DATABASE_URL
- [ ] Migrations run (`migrate up`)
- [ ] Backend binary built (`./bin/api` exists)
- [ ] Backend started and listening on :8080
- [ ] Health endpoint responds (`curl /health`)
- [ ] Flutter app configured to use `http://localhost:8080`

---

## 🔧 Configuration Details

### Environment Variables (.env)

| Variable     | Default     | Description                           |
| ------------ | ----------- | ------------------------------------- |
| PORT         | 8080        | Server port                           |
| ENVIRONMENT  | development | Environment mode                      |
| DATABASE_URL | REQUIRED    | PostgreSQL connection string          |
| JWT_SECRET   | dev-secret  | JWT signing key (CHANGE IN PROD)      |
| HMAC_SECRET  | dev-secret  | QR code HMAC signing (CHANGE IN PROD) |

### Connection String Format

```
postgres://username:password@host:port/database?sslmode=disable
```

Examples:

```
# Local PostgreSQL
postgres://postgres:postgres@localhost:5432/cloak_db?sslmode=disable

# Docker PostgreSQL
postgres://postgres:postgres@postgres:5432/cloak_db?sslmode=disable
```

---

## 🚨 Troubleshooting

### "CONNECTION REFUSED"

```bash
# Check PostgreSQL is running
lsof -i :5432

# If not running with Docker:
docker ps | grep postgres

# If not running, start it:
docker run -d -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=cloak_db -p 5432:5432 postgres:16
```

### "DATABASE DOES NOT EXIST"

```bash
# Create database
createdb cloak_db

# Or using Docker:
docker exec <container-id> createdb -U postgres cloak_db
```

### "MIGRATION FAILED"

```bash
# Check migration status
migrate -path migrations -database "postgres://..." version

# Force reset migrations (development only!)
migrate -path migrations -database "postgres://..." force 0
```

### "COMPILATION ERROR"

```bash
# Update Go dependencies
go mod download
go mod tidy

# Rebuild
make build
```

### Check Logs

```bash
# Direct execution
# Logs print to terminal (see output directly)

# Docker Compose
docker-compose logs -f app

# macOS/Linux service
# Check system logs
```

---

## 📊 Database Schema

### businesses

- id (PK: UUID)
- name, email (unique), password
- role: 'business'
- hmac_key (for QR signing)
- timestamps

### customers

- id (PK: UUID)
- email (unique), phone
- timestamps

### services

- id (PK: UUID)
- business_id (FK: businesses.id)
- name, total_slots
- timestamps

### slots

- id (PK: UUID)
- service_id (FK: services.id)
- slot_number, status (free/occupied)
- Unique constraint: (service_id, slot_number)
- timestamps

### tickets

- id (PK: UUID)
- service_id (FK: services.id)
- slot_id (FK: slots.id, nullable)
- slot_number, customer_id
- status (active/released)
- hmac_digest (QR code verification)
- timestamps

---

## 🌐 Verify with Flutter App

Once backend is running on localhost:8080:

1. Run Flutter app: `flutter run -d web`
2. App should connect to http://localhost:8080
3. Test auth flow:
    - Register/Login with test account
    - Create service
    - Check-in customer (creates QR ticket)
    - Scan QR code
    - Release ticket

---

## 📚 Project Structure

```
CLOAKBE/
├── bin/
│   └── api                     (Compiled binary)
├── cmd/
│   └── api/
│       └── main.go            (Entry point)
├── internal/
│   ├── config/                (Configuration loading)
│   ├── database/              (PostgreSQL pool)
│   ├── domain/                (Data types like Ticket, Service)
│   ├── handler/               (HTTP handlers)
│   ├── middleware/            (Auth, CORS, logging)
│   ├── repository/            (Data access)
│   ├── usecase/               (Business logic)
│   └── qr/                    (QR code generation)
├── migrations/                (Database schema)
├── .env                       (Configuration - create this)
├── Makefile                   (Build commands)
├── docker-compose.yml         (Docker setup)
└── go.mod                     (Dependencies)
```

---

## ✨ Ready to Go!

Your backend is **production-ready**. Follow the startup steps above and you'll have a fully functional CLOAK API running locally that your Flutter frontend can connect to! 🚀
