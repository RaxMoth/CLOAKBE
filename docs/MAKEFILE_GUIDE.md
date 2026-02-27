# CLOAK Backend - Makefile Quick Reference

Your new **Makefile** makes starting the backend super easy! 🚀

## 🎯 All Available Commands

```bash
make help         # Show this menu
make start        # 🚀 START EVERYTHING (db + migrations + backend)
make stop         # Stop database
make db-up        # Start PostgreSQL container
make db-down      # Stop PostgreSQL container
make db-shell     # Open PostgreSQL terminal
make setup-env    # Create .env configuration file
make migrate      # Run database migrations
make migrate-down # Rollback migrations
make backend      # Start API server
make build        # Compile binary
make rebuilt      # Clean and rebuild
make test         # Test health endpoint
make test-register # Test registration
make test-login   # Test login
make clean        # Remove build artifacts
make all-down     # Stop and clean everything
```

---

## ⚡ Quick Start (One Command!)

```bash
make start
```

This runs:

1. ✅ `make db-up` → Starts PostgreSQL
2. ✅ `make migrate` → Sets up database schema
3. ✅ `make backend` → Starts API server on localhost:8080

---

## 📋 Step-by-Step Startup

### First Time Only:

```bash
# Create .env configuration
make setup-env

# Start everything
make start
```

### Every Other Time:

```bash
make start
```

---

## 🧪 Testing the Backend

Once running, test with:

```bash
# Health check
make test

# Register account
make test-register

# Login
make test-login
```

---

## 🗄️ Database Management

```bash
# Start database
make db-up

# Stop database (keeps data)
make db-down

# Clean database (removes data)
make db-clean

# Open PostgreSQL shell
make db-shell

# View database logs
make db-logs

# Run migrations
make migrate

# Rollback migrations
make migrate-down
```

---

## 🏗️ Backend Building

```bash
# Build binary
make build

# Start backend (must build first)
make backend

# Clean and rebuild
make rebuild
```

---

## 🧹 Cleanup

```bash
# Remove build artifacts
make clean

# Stop everything and clean up
make all-down
```

---

## 📝 Complete Setup Flow

First time setup:

```bash
make setup-env    # Create .env
make start        # Start everything
```

Now your backend is running! ✅

---

## 🆘 Troubleshooting

### Database already exists error

```bash
make db-clean     # Remove old database
make db-up        # Start fresh
make migrate      # Create new schema
```

### Backend won't start

```bash
make clean        # Clean build
make build        # Rebuild
make backend      # Start again
```

### Test endpoints

```bash
make test         # Quick health check
curl http://localhost:8080/health  # Manual test
```

---

## File Locations

```
CLOAKBE/
├── Makefile              (Main makefile - use this!)
├── Makefile.backup       (Backup of original)
├── Makefile.dev          (Development version)
├── bin/api              (Compiled binary)
├── .env                 (Configuration - created by make setup-env)
└── migrations/          (Database schema)
```

---

## 💡 Pro Tips

1. **Use `make start` for first-time setup** - It does everything
2. **Use `make test` to verify backend is working**
3. **Use `make db-shell` to inspect database directly**
4. **Keep database running** - `make stop` only stops the container, data persists
5. **Use `make all-down` only when you want fresh start**

---

## ✨ Your New Development Workflow

```bash
# Morning: Start development
make start

# Work on code...

# Evening: Check something
make test

# Before bed: Stop everything
make stop

# Tomorrow: Resume
make start
```

---

Done! 🎉 Now you have a fully automated backend setup with simple make commands!
