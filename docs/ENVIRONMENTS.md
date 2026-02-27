# CLOAK Backend - Environment Setup Guide

Professional multi-environment configuration for local, development, and production deployments.

## Overview

CLOAK Backend supports three environment configurations:

| Environment | Purpose | Database | Security | Logging |
|-------------|---------|----------|----------|---------|
| **Local** | Local development | Docker (localhost) | Development keys | Debug |
| **Dev** | Staging/testing | Dev server | Example keys | Info |
| **Prod** | Production | Production DB | Strong secrets | Warning |

## Quick Start

### 1. Local Development (Default)

```bash
# Setup local environment
make env-local

# Start everything
make start

# Backend runs on http://localhost:8080
```

**Configuration (`.env.local`):**
- Database: `localhost:5432` (Docker container)
- JWT Secret: `dev-jwt-super-secret-local-only-12345`
- HMAC Secret: `dev-hmac-super-secret-local-only-12345`
- Log Level: `debug` (verbose logging)

---

### 2. Development Environment

For staging/pre-production testing.

```bash
# Setup development environment
make env-dev

# Update .env.dev with your dev server details
nano .env.dev

# Start database and backend with dev config
make migrate
make backend
```

**Configuration (`.env.dev`):**
- Update `DATABASE_URL` to your dev database server
- Generate new JWT_SECRET: `openssl rand -base64 32`
- Generate new HMAC_SECRET: `openssl rand -base64 32`
- Database: Your dev Postgres server
- Log Level: `info` (normal logging)

**Example Dev Setup:**
```dotenv
DATABASE_URL=postgres://dev_user:dev_password@dev-db.example.com:5432/cloak_dev?sslmode=require
JWT_SECRET=your-generated-dev-secret-here-32-chars
HMAC_SECRET=your-generated-dev-hmac-here-32-chars
ENVIRONMENT=development
```

---

### 3. Production Environment

For live deployment (AWS, DigitalOcean, etc.).

```bash
# Setup production environment
make env-prod

# CRITICAL: Update .env.prod with production secrets
nano .env.prod

# Verify all secrets are set (no CHANGE_ME values)
grep CHANGE_ME .env.prod  # Should return nothing

# Deploy
docker run --env-file .env.prod myregistry/cloak-api:latest
```

**Configuration (`.env.prod`):**

```dotenv
# Server
PORT=8080
ENVIRONMENT=production

# Database (Use RDS or managed Postgres)
DATABASE_URL=postgres://prod_user:SECURE_PASSWORD@prod-db.example.com:5432/cloak_prod?sslmode=require

# Security - MUST be strong and unique
JWT_SECRET=<generate-with-openssl>
HMAC_SECRET=<generate-with-openssl>

# Logging
LOG_LEVEL=warn

# CORS
ALLOWED_ORIGINS=https://app.example.com,https://www.example.com
```

**Generate Strong Secrets:**
```bash
# Generate 32-byte base64 encoded secrets
openssl rand -base64 32

# Example output:
# T8kL9mNx2pQ5rJ3vWe7sUfH6zYaJ4dK8wB0xY1cZ9a2=
```

---

## Environment Files Reference

### `.env.local` (Local Development)
- Use for daily development on your machine
- Docker PostgreSQL container
- Development security keys (unsafe for production)
- Debug logging enabled
- **Gitignore**: ✓ Ignored automatically

### `.env.dev` (Development/Staging)
- Use for pre-production testing environment
- Points to dev database server
- Example secrets (update before deploying)
- Info level logging
- **Gitignore**: ✓ Ignored automatically

### `.env.prod` (Production)
- Use for live deployment only
- Points to production database (RDS/managed)
- **MUST** use strong, unique secrets
- **MUST NOT** use CHANGE_ME values
- Warning level logging (minimal)
- **Gitignore**: ✓ Ignored automatically
- **Security**: Store in secrets manager, not in git

### `.env.example` (Documentation)
- Reference template showing all configuration options
- Comments explaining each setting
- **Gitignore**: ✗ Committed to git for reference
- Safe to commit (no secrets included)

---

## Configuration Details

### Port Configuration
```dotenv
PORT=8080        # Backend API port (default: 8080)
```

### Database Configuration

**Local (Docker):**
```dotenv
DATABASE_URL=postgres://postgres:postgres@localhost:5432/cloak_db?sslmode=disable
```

**Development:**
```dotenv
DATABASE_URL=postgres://dev_user:dev_password@dev-db.example.com:5432/cloak_dev?sslmode=require
```

**Production (AWS RDS example):**
```dotenv
DATABASE_URL=postgres://prod_user:secure_pass@cloak-prod.abc123.us-east-1.rds.amazonaws.com:5432/cloak_prod?sslmode=require
```

### Security Keys

**JWT Secret** (User authentication)
- Used to sign and verify JWT tokens
- Minimum 32 characters recommended
- **Local**: Can use simple dev key
- **Dev**: Generate with `openssl rand -base64 32`
- **Prod**: Generate with `openssl rand -base64 32` (MUST be strong)

**HMAC Secret** (QR code signing)
- Used to sign QR codes for verification
- Minimum 32 characters recommended
- **Local**: Can use simple dev key
- **Dev**: Generate with `openssl rand -base64 32`
- **Prod**: Generate with `openssl rand -base64 32` (MUST be strong)

### Logging Levels

| Level | Use Case | Detail |
|-------|----------|--------|
| `debug` | Local development | All events, full details |
| `info` | Development/Staging | Important events only |
| `warn` | Production | Warnings and errors only |
| `error` | Critical issues | Errors only |

---

## Makefile Environment Commands

```bash
# Setup local environment for development
make env-local

# Setup development/staging environment
make env-dev

# Setup production environment
make env-prod
```

Each command:
1. Copies the environment-specific file to `.env`
2. Shows configuration summary
3. Warns about required updates (for dev/prod)

---

## Step-by-Step Deployment

### Local Development Workflow

```bash
# 1. Setup local environment
make env-local

# 2. Start everything
make start
# (equivalent to: make db-up && make migrate && make backend)

# 3. Run tests
make test

# 4. When done, stop database
make stop
```

### Development Server Deployment

```bash
# 1. Setup development environment
make env-dev

# 2. Update configuration
nano .env.dev
# Update: DATABASE_URL, JWT_SECRET, HMAC_SECRET

# 3. Start database (assumes remote DB)
make migrate

# 4. Start backend
make backend

# 5. Test
curl http://localhost:8080/health
```

### Production Deployment (Docker)

```bash
# 1. Setup production environment
make env-prod

# 2. Secure configuration
nano .env.prod
# Update all production secrets

# 3. Build Docker image
docker build -t myregistry/cloak-api:v1.0.0 .

# 4. Run container with .env.prod
docker run \
  --name cloak-api \
  --env-file .env.prod \
  -p 8080:8080 \
  myregistry/cloak-api:v1.0.0

# 5. Verify
curl http://production-server:8080/health
```

---

## Security Best Practices

### ✅ DO:
- ✓ Generate unique secrets for each environment
- ✓ Use SSL/TLS for database connections (sslmode=require)
- ✓ Keep .env files in .gitignore
- ✓ Store production secrets in a secrets manager
- ✓ Rotate secrets periodically
- ✓ Use strong database credentials
- ✓ Enable CORS restrictions in production
- ✓ Use HTTPS for all production APIs

### ❌ DON'T:
- ✗ Commit .env files with secrets to git
- ✗ Use CHANGE_ME values in production
- ✗ Share secrets via email or chat
- ✗ Use HTTP in production
- ✗ Use same secrets across environments
- ✗ Allow permissive CORS in production
- ✗ Use sslmode=disable in production

---

## Troubleshooting

### "DATABASE_URL environment variable is required"

**Problem**: Missing database URL configuration

**Solution**:
```bash
# Check current environment file
cat .env

# Setup correct environment
make env-local  # or env-dev, env-prod

# Verify configuration
grep DATABASE_URL .env
```

### "Can't connect to database"

**Problem**: Database server not running or URL is incorrect

**Solution**:
```bash
# For local development:
make db-up      # Start Docker container

# Verify connection:
make db-shell   # Shell into database

# For dev/prod:
# Verify DATABASE_URL points to correct server
nano .env.dev   # or .env.prod
```

### "Invalid JWT Secret" 

**Problem**: JWT secret is too short or not set

**Solution**:
```bash
# Generate new secret
openssl rand -base64 32

# Update .env file
JWT_SECRET=<your-generated-secret>

# Restart backend
make backend
```

---

## Environment-Specific Features

### Local Development (.env.local)
- ✓ Debug logging enabled
- ✓ Development security keys
- ✓ Docker database included
- ✓ No SSL required for database
- ✓ Permissive CORS

### Development (.env.dev)
- ✓ Info level logging
- ✓ Strong security keys recommended
- ✓ Remote dev database
- ✓ SSL required for database
- ✓ Limited CORS

### Production (.env.prod)
- ✓ Warning level logging only
- ✓ **Must use** strong security keys
- ✓ Production database required
- ✓ **Must use** SSL for database
- ✓ **Restricted CORS** to approved domains
- ✓ Minimal logging for performance

---

## Additional Resources

- [PostgreSQL Connection Strings](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [AWS RDS Connection Guide](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.DBInstanceClass.html)
- [Docker Secrets Management](https://docs.docker.com/engine/swarm/secrets/)

