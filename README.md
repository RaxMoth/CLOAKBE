# CLOAK Backend (CLOAKBE)

A high-performance Go/Fiber REST API backend for the CLOAK digital ticketing system. Provides complete ticket management, business authentication, and real-time slot reservation with PostgreSQL persistence.

## 🎯 Project Overview

**CLOAK** is a B2B SaaS platform for managing digital tickets at events and venues. CLOAKBE handles:

- **Business Management** - Registration, authentication, service creation
- **Ticket Lifecycle** - Check-in (QR generation), scanning/verification, release/cancellation
- **Slot Management** - Real-time seat/capacity management with row-level locking
- **QR Security** - HMAC-SHA256 signed QR codes for verification integrity

## 🛠️ Technology Stack

- **Go 1.22+** - Language
- **Fiber v2** - Web framework (fast, lightweight)
- **PostgreSQL 16** - Database
- **pgx/v5** - Database driver (prepared statements, async)
- **JWT (golang-jwt)** - Authentication tokens
- **bcrypt** - Password hashing
- **Docker** - Containerization

## 📋 Project Structure

```
CLOAKBE/
├── cmd/
│   └── api/
│       └── main.go                  # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go                # Environment & config management
│   ├── domain/
│   │   └── entities.go              # Domain models & repository interfaces
│   ├── apperror/
│   │   ├── errors.go                # Typed error system with HTTP status
│   │   └── helpers.go               # Error checking predicates
│   ├── database/
│   │   └── db.go                    # PostgreSQL connection pool
│   ├── qr/
│   │   └── payload.go               # QR signing & encoding/decoding
│   ├── usecase/
│   │   ├── auth_usecase.go          # Authentication logic
│   │   ├── ticket_usecase.go        # Ticket operations
│   │   └── service_usecase.go       # Service management
│   ├── repository/
│   │   ├── business_repo.go         # Business CRUD
│   │   ├── customer_repo.go         # Customer CRUD + upsert
│   │   ├── service_repo.go          # Service CRUD
│   │   ├── slot_repo.go             # Slot operations (read/update/release)
│   │   └── ticket_repo.go           # Ticket CRUD
│   ├── handler/
│   │   ├── auth_handler.go          # HTTP auth endpoints
│   │   ├── ticket_handler.go        # HTTP ticket endpoints
│   │   └── service_handler.go       # HTTP service endpoints
│   └── middleware/
│       └── auth.go                  # JWT validation & role enforcement
├── migrations/
│   ├── 000001_init_schema.up.sql    # Database schema creation
│   └── 000001_init_schema.down.sql  # Rollback script
├── docs/
│   ├── INDEX.md                     # Documentation index
│   ├── QUICK_START.md               # Quick start guide
│   ├── SESSION_SUMMARY.md           # What was built & current status
│   ├── ENVIRONMENTS.md              # Environment setup guide
│   ├── MAKEFILE_GUIDE.md            # All available make commands
│   ├── DEPLOYMENT_GUIDE.md          # Deployment instructions
│   └── TESTING_STATUS.md            # Testing coverage & status
├── .env.example                     # Example environment variables
├── Dockerfile                       # Docker image definition
├── docker-compose.yml               # Local development containers
├── Makefile                         # Build & run automation
├── go.mod                           # Go dependencies
├── go.sum                           # Dependency checksums
└── README.md                        # This file
```

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL 16+ (local or Docker)
- Make (usually pre-installed macOS/Linux)
- Docker & Docker Compose (optional, for containerized setup)

### Local Setup (Development)

**1. Clone and setup:**

```bash
cd CLOAKBE
cp .env.example .env
```

**2. Create PostgreSQL database:**

```bash
createdb cloak
```

**3. Apply database migrations:**

```bash
make migrate-up
# or manually:
# psql cloak < migrations/000001_init_schema.up.sql
```

**4. Run the API:**

```bash
make run
# or: go run cmd/api/main.go
```

API will be available at `http://localhost:8080`

### Docker Setup (Recommended)

```bash
# Build and start all containers (PostgreSQL + API)
make docker-up

# View logs
make docker-logs

# Stop containers
make docker-down
```

## 📡 API Endpoints

### Authentication

| Method | Endpoint                        | Body                                     | Auth? |
| ------ | ------------------------------- | ---------------------------------------- | ----- |
| POST   | `/api/v1/auth/business/register`| `{email, password, business_name}`       | No    |
| POST   | `/api/v1/auth/business/login`   | `{email, password}`                      | No    |
| POST   | `/api/v1/auth/customer/login`   | `{phone_number}`                         | No    |

**Response:** `{access_token, refresh_token, user_id, role}`

### Tickets (Business)

| Method | Endpoint                    | Body / Query                     | Auth? |
| ------ | --------------------------- | -------------------------------- | ----- |
| POST   | `/api/v1/tickets/checkin`   | `{service_id, customer_id}`      | Yes   |
| POST   | `/api/v1/tickets/scan`      | `{qr_payload, hmac_signature}`   | Yes   |
| POST   | `/api/v1/tickets/:id/release` | `-`                            | Yes   |

### Services (Business)

| Method | Endpoint                  | Body / Query              | Auth? |
| ------ | ------------------------- | ------------------------- | ----- |
| POST   | `/api/v1/services`        | `{name, capacity}`        | Yes   |
| GET    | `/api/v1/services`        | `?page=1&limit=10`        | Yes   |
| GET    | `/api/v1/services/:id`    | `-`                       | Yes   |
| GET    | `/api/v1/services/:id/stats` | `-`                     | Yes   |

### Health Check

| Method | Endpoint | Auth? |
| ------ | -------- | ----- |
| GET    | `/health`| No    |

## 📚 Documentation

For detailed information, see the `docs/` folder:

- [SESSION_SUMMARY.md](docs/SESSION_SUMMARY.md) - Overview of what's been completed
- [QUICK_START.md](docs/QUICK_START.md) - Detailed local + Flutter setup
- [ENVIRONMENTS.md](docs/ENVIRONMENTS.md) - Environment variables guide
- [MAKEFILE_GUIDE.md](docs/MAKEFILE_GUIDE.md) - All available make commands
- [DEPLOYMENT_GUIDE.md](docs/DEPLOYMENT_GUIDE.md) - Deployment to production
- [TESTING_STATUS.md](docs/TESTING_STATUS.md) - Testing & coverage info

## 🔐 Security

- **JWT Tokens**: 15-minute access, 7-day refresh tokens
- **Password Hashing**: bcrypt with 12 salt rounds
- **QR Verification**: HMAC-SHA256 signing with business-specific keys
- **Database**: Row-level locking on slot operations (prevents race conditions)
- **CORS**: Configured for frontend domain

## 🗄️ Database

### PostgreSQL Schema Highlights

- `businesses` - Business accounts with bcrypt passwords
- `customers` - Customer profiles
- `services` - Event/venue services with capacity
- `slots` - Individual capacity units (e.g., seats) with availability
- `tickets` - Ticket records with check-in status
- **Row-Level Locking**: Prevents race conditions on slot claims

See [migrations/000001_init_schema.up.sql](migrations/000001_init_schema.up.sql) for full schema.

## 📊 Clean Architecture

```
HTTP Request
    ↓
Handler (HTTP parsing, validation)
    ↓
Usecase (business logic, orchestration)
    ↓
Repository (data access, persistence)
    ↓
PostgreSQL Database
```

Each layer is independently testable and loosely coupled.

## 🔨 Available Make Commands

```bash
make help                  # Show all commands
make install               # Download Go dependencies
make run                   # Run API locally (go run)
make build                 # Build binary: ./bin/api
make test                  # Run all tests
make test-cover            # Run tests with coverage
make docker-build          # Build Docker image
make docker-up             # Start containers (Compose)
make docker-down           # Stop containers
make docker-logs           # View container logs
make migrate-up            # Apply migrations
make migrate-down          # Rollback migrations
make fmt                   # Format code (gofmt)
make lint                  # Run linter (golangci-lint)
make clean                 # Remove build artifacts
```

## 🧪 Testing

Run the test suite:

```bash
# All tests
make test

# With coverage reports
make test-cover

# View coverage in browser
make test-coverage-html
```

See [docs/TESTING_STATUS.md](docs/TESTING_STATUS.md) for test coverage details.

## 🚢 Deployment

### Docker Deployment

```bash
# Build image
docker build -t cloakbe .

# Run container
docker run -p 8080:8080 --env-file .env.prod cloakbe
```

### Cloud Deployment

Guides available for:
- **Railway** - See [docs/DEPLOYMENT_GUIDE.md](docs/DEPLOYMENT_GUIDE.md#railway)
- **Render** - See [docs/DEPLOYMENT_GUIDE.md](docs/DEPLOYMENT_GUIDE.md#render)
- **AWS EC2/ECS** - Standard Docker deployment
- **Google Cloud Run** - Containerized deployment

### Environment Setup for Production

```bash
cp .env.example .env.prod

# Edit with production values:
# - Strong JWT_SECRET & HMAC_SECRET (use: openssl rand -base64 32)
# - Prod database URL
# - Set ENVIRONMENT=production
# - Set LOG_LEVEL=warn
```

## 📝 Development Workflow

**Local development (3 terminals):**

1. **Terminal 1: PostgreSQL**
   ```bash
   make docker-up # Just DB
   ```

2. **Terminal 2: Go API**
   ```bash
   make run
   ```

3. **Terminal 3: Flutter Frontend**
   ```bash
   cd ../CLOAK
   flutter run -d web
   ```

Access app at `http://localhost:PORT` (shown by Flutter)

## 🐛 Troubleshooting

### Database Connection Errors

```bash
# Check if PostgreSQL is running
psql cloak -c "SELECT 1"

# Reset migrations if needed
make migrate-down
make migrate-up
```

### Port Already in Use

```bash
# Change PORT in .env
PORT=8081  # Use different port

# Or kill process on port 8080
lsof -i :8080 | grep -v PID | awk '{print $2}' | xargs kill
```

### See more

Check [docs/RAILWAY_TROUBLESHOOTING.md](docs/RAILWAY_TROUBLESHOOTING.md) for detailed troubleshooting.

## 📦 Dependencies

Run `go mod tidy` to update dependencies.

Key packages:
- `github.com/gofiber/fiber/v2` - Web framework
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/golang-jwt/jwt` - JWT tokens
- `golang.org/x/crypto` - Cryptography & bcrypt

## 📄 License

MIT License - See LICENSE file

## 🤝 Contributing

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Commit changes: `git commit -am 'Add feature'`
3. Push: `git push origin feature/your-feature`
4. Open PR

## ❓ Questions?

Refer to [docs/INDEX.md](docs/INDEX.md) for complete documentation index or check individual doc files.

---

**Built with ❤️ for the CLOAK ticketing system**
