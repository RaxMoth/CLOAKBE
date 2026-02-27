# CLOAK - Quick Start Guide

Get up and running in 5 minutes.

## Prerequisites

```bash
# Check Flutter
flutter --version     # Should be 3.19.6+

# Check Go
go version            # Should be 1.22+

# Check Docker
docker --version      # Should be 4.0+

# Backend tools installed
make --version
migrate --version     # Should be 4.19.1+
ngrok --version       # Should be 3.36.1+
```

---

## Option 1: Local Development (Recommended for Testing)

### 1️⃣ Terminal 1 - Start Backend + Database

```bash
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE

# Start everything (database, migrations, API server)
make start

# You should see:
# ✓ Database running on localhost:5432
# ✓ Migrations complete
# ✓ Backend running on http://localhost:8080
```

### 2️⃣ Terminal 2 - Create Public Tunnel

```bash
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE

# Start ngrok tunnel
make tunnel

# You'll see output like:
# Forwarding    https://abc123def456.ngrok.io -> http://localhost:8080
#
# 📝 COPY THIS URL: https://abc123def456.ngrok.io
```

### 3️⃣ Update Flutter App URL

```bash
# Edit Flutter configuration
nano /Users/maxroth/Documents/Programming/FlutterProjects/CLOAK/lib/core/constants/app_constants.dart

# Find this line:
#   static const String apiBaseUrl = 'http://192.168.1.179:8080';
#
# Replace with your ngrok URL:
#   static const String apiBaseUrl = 'https://abc123def456.ngrok.io';
```

### 4️⃣ Terminal 3 - Run Flutter App

```bash
cd /Users/maxroth/Documents/Programming/FlutterProjects/CLOAK

# Start Flutter app
flutter run -d macos

# Press R to hot reload after URL change
```

### ✅ Test It

1. Open Flutter app
2. Click "Business Register" (or "Customer Login")
3. Enter: name=test, email=test@test.com, password=test
4. Should connect and register! ✓

---

## Option 2: Production Deployment (For Sharing/Stability)

### 1️⃣ Deploy Backend to Railway

**One-time setup:**

1. Go to [railway.app](https://railway.app)
2. Login with GitHub
3. Click "Create New Project"
4. Select "Deploy from GitHub"
5. Find and select `CLOAKBE` repo
6. Railway auto-detects Dockerfile and deploys ✓

**After deployment:**

- Railway gives you a URL: `https://cloak-api-xxx.railway.app`
- PostgreSQL is auto-provisioned ✓
- Environment variables auto-set ✓

### 2️⃣ Update Flutter App

```bash
nano /Users/maxroth/Documents/Programming/FlutterProjects/CLOAK/lib/core/constants/app_constants.dart

# Replace URL:
static const String apiBaseUrl = 'https://cloak-api-xxx.railway.app';
```

### 3️⃣ Hot Reload Flutter

```bash
# App is still running from Option 1, Terminal 3
# Just modify the file, then press R in Flutter terminal
```

### ✅ Test It

Same as Option 1 - should work seamlessly!

---

## Quick Command Reference

### Backend Management

```bash
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE

make start        # ▶️  Start backend + database
make stop         # ⏹️  Stop database
make backend      # ▶️  Just run API (db must be running)
make db-shell     # 📊 Open PostgreSQL terminal
```

### Tunneling (Local Development)

```bash
make tunnel       # 🌐 Start ngrok tunnel for Flutter testing
make tunnel-help  # ℹ️  Show tunnel workflow
```

### Environment Setup

```bash
make env-local    # 🏠 Use local environment (.env.local)
make env-dev      # 🔧 Use dev environment (.env.dev)
make env-prod     # 🚀 Use production environment (.env.prod)
```

### Testing

```bash
# Test health endpoint
curl http://localhost:8080/health

# Test through ngrok
curl https://YOUR_NGROK_URL/health

# Test registration
curl -X POST http://localhost:8080/api/v1/auth/business/register \
  -H "Content-Type: application/json" \
  -d '{"name":"test","email":"test@test.com","password":"test"}'
```

---

## Complete Local Workflow

### First Time Setup

```bash
# 1. Clone backend
git clone git@github.com:YOUR_USERNAME/CLOAKBE.git
cd CLOAKBE

# 2. Start everything
make start
# ✓ Database running
# ✓ Backend running

# 3. In another terminal, start tunnel
make tunnel
# ✓ ngrok running

# 4. Clone Flutter app (if not already)
cd /Users/maxroth/Documents/Programming/FlutterProjects/CLOAK
git clone git@github.com:YOUR_USERNAME/CLOAK.git

# 5. Update API URL and run
# (Edit lib/core/constants/app_constants.dart)
flutter run -d macos
```

### Daily Development

```bash
# Terminal 1: Backend
cd CLOAKBE && make start

# Terminal 2: Tunnel
make tunnel

# Terminal 3: Flutter
cd CLOAK && flutter run -d macos

# Edit code, Flutter hot reloads automatically
# Press R to reload if needed
```

---

## Troubleshooting

### Flutter Can't Connect to Backend

1. **Check backend is running:**

    ```bash
    curl http://localhost:8080/health
    # Should return: {"status":"ok"}
    ```

2. **Check ngrok tunnel:**

    ```bash
    curl https://YOUR_NGROK_URL/health
    # Should return: {"status":"ok"}
    ```

3. **Verify Flutter app has correct URL:**
    - Edit `lib/core/constants/app_constants.dart`
    - Confirm URL matches ngrok URL
    - Press R in Flutter terminal to reload

4. **Check ngrok web interface:**
    - Open http://127.0.0.1:4040 in browser
    - Should show your API requests
    - If empty, tunnel not receiving traffic

### "ngrok: command not found"

```bash
brew install ngrok
ngrok --version
# Should show: ngrok version 3.36.1
```

### Backend Won't Start

```bash
# Check if port 8080 is in use
lsof -i :8080

# If in use, kill the process
killall -9 api

# Try again
make start
```

### Database Connection Failed

```bash
# Verify Docker is running
docker ps

# Check if PostgreSQL container exists
docker ps -a | grep postgres

# Restart database
make db-down
make db-up
```

---

## Environment Files

### Use Local Environment (Default)

```bash
make env-local
# Uses .env.local with localhost:5432
```

### Use Dev Environment

```bash
make env-dev
# Uses .env.dev (update with your dev server)
nano .env.dev
```

### Use Production Environment

```bash
make env-prod
# Uses .env.prod (for Railway deployment)
```

---

## Next Steps

1. ✅ **Get local working** - Follow Option 1 above
2. ✅ **Test API** - Use curl commands to verify endpoints
3. ✅ **Test Flutter app** - Register/login from Flutter
4. ↪️ **Deploy to Railway** - Follow Option 2 for production
5. ↪️ **Connect to real domain** - Update CORS and DNS

---

## Useful Links

- 📚 [Full Environments Guide](./docs/ENVIRONMENTS.md)
- 🌐 [ngrok Setup Guide](./docs/NGROK_LOCAL_SETUP.md)
- 🚀 [Deployment Guide](./docs/DEPLOYMENT_GUIDE.md)
- 📖 [Makefile Reference](./docs/MAKEFILE_GUIDE.md)

---

## Support

All three components working?

- ✅ Backend running
- ✅ Flutter app running
- ✅ Database connected
- ✅ Tunnel active

**Congratulations!** 🎉 CLOAK is fully operational. Start building features!
