# Railway Deployment Troubleshooting Guide

**Last Updated:** 28 February 2026  
**Issue:** Backend deployed to Railway but not responding  
**Status:** 🔴 BLOCKING - Needs investigation

---

## 🚨 Current Problem

Backend code deployed to Railway at `https://cloakbe-production.up.railway.app`, but:

- Health endpoint returns NO RESPONSE (timeout)
- Web app gets XMLHttpRequest network error
- Mac app gets HTTPS connection error (sandbox)

---

## 🔍 Investigation Steps

### Step 1: Check Railway Logs (REQUIRED)

1. Go to https://railway.app
2. Login to your account
3. Navigate to your CLOAK project
4. Click on "api" service
5. Click "Logs" tab
6. Look for these patterns:

**Expected Success Logs:**

```
Starting server on :8080
Connected to database
```

**Common Error Patterns:**

#### Error: "DATABASE_URL not found"

```
panic: sql: unknown driver "postgres"
or
MissingRequiredField: DATABASE_URL
```

**Fix:** Add DATABASE_URL to Railway environment variables

- Go to Variables tab
- Add: `DATABASE_URL=postgresql://user:pass@host/dbname`

#### Error: "connection refused"

```
dial tcp: connection refused
error connecting to database
```

**Fix:** Database not started or wrong credentials

- Check if PostgreSQL service exists in Railway
- Verify DATABASE_URL matches actual database

#### Error: "port already in use"

```
listen tcp :8080: bind: address already in use
```

**Fix:** Port conflict

- Public domain already assigned
- Deploy to different app name

---

### Step 2: Check Environment Variables

1. In Railway dashboard, click "api" service
2. Click "Variables" tab
3. Verify these exist:

    ```
    PORT=8080
    ENVIRONMENT=production
    DATABASE_URL=postgresql://...
    JWT_SECRET=<unique-value>
    HMAC_SECRET=<unique-value>
    LOG_LEVEL=debug
    ```

4. **If DATABASE_URL is empty:**
    - Go to PostgreSQL service
    - Get connection string
    - Paste into "api" service Variables

---

### Step 3: Manual Test from Terminal

```bash
# Test if Railway endpoint responds at all
curl -v https://cloakbe-production.up.railway.app/health

# Expected response:
# HTTP/1.1 200 OK
# {"status":"ok"}

# If timeout, check Railway dashboard for deploy status
```

---

### Step 4: Check Deployment Status

1. Go to Railway project → "api" service
2. Look for blue "Deploy Status" badge
3. Expect to see: ✅ "Deployed" (green)

**If showing orange or red:**

- Click on the failed deployment
- View build logs
- Look for errors in Go compilation or Docker build

---

## ⚙️ Common Railway Issues & Fixes

### Issue 1: Missing Environment Variables

**Symptom:** App crashes on startup, config errors in logs  
**Solution:**

```bash
# Check .env.deployed template
# Make sure Railway vars match:
cat /Users/maxroth/Documents/Programming/Go/CLOAKBE/.env.deployed

# Add missing vars to Railway dashboard Variables tab
```

### Issue 2: Database Connection String Wrong

**Symptom:** "connection refused" errors  
**Solution:**

```bash
# Get PostgreSQL connection string from Railway
# In PostgreSQL service → Variables tab
# Copy DATABASE_URL (looks like: postgresql://default:...)

# Paste into "api" service → Variables tab
```

### Issue 3: Migrations Didn't Run

**Symptom:** Database tables don't exist, 500 errors  
**Solution:**

```bash
# Option A: Redeploy with migrations
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE
git push origin main  # Triggers redeploy

# Option B: Run migrations manually in Railway
# (Contact Railway support or use CLI)
```

### Issue 4: Port Not 8080

**Symptom:** Public URL works but gives 404/connection refused  
**Solution:**

```bash
# Check main.go listening on correct port
grep "app.Listen" /Users/maxroth/Documents/Programming/Go/CLOAKBE/cmd/api/main.go
# Should show: app.Listen(":8080")

# If different:
# 1. Update to :8080
# 2. Commit and push
# 3. Railway redeploys automatically
```

---

## 🚀 Step-by-Step Fix Process

### If Logs Show Errors:

```bash
# 1. Fix the code locally
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE

# 2. Test locally first
make start
# Wait for "Starting server on :8080"

# 3. Commit fix
git add .
git commit -m "Fix: [describe issue]"

# 4. Push to GitHub (Railway auto-redeploys)
git push origin main

# 5. Wait 2-3 minutes for redeploy
# Check Railway Deployments tab

# 6. Test new deployment
curl https://cloakbe-production.up.railway.app/health
```

### If No Logs or Empty Logs:

```bash
# 1. Check if app is actually running
# Go to Railway dashboard → api service
# Look for green "Running" status

# 2. If not running, click "Redeploy"
# This rebuilds Docker image and starts fresh

# 3. Wait 3-5 minutes

# 4. Check logs again
```

---

## 🧪 Testing After Fix

Once Railway endpoint responds:

### Test 1: Health Check

```bash
curl https://cloakbe-production.up.railway.app/health
# Expected: {"status":"ok"} with 200 status
```

### Test 2: Database Connection

```bash
# Try registering new business (creates DB record)
curl -X POST https://cloakbe-production.up.railway.app/api/v1/auth/business/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!",
    "business_name": "Test Business"
  }'

# Expected: {"token": "...", "user": {...}}
```

### Test 3: From Flutter App

```bash
# 1. Update Flutter API endpoint (already set to Railway URL)
# lib/core/constants/app_constants.dart:
# static const String apiBaseUrl = 'https://cloakbe-production.up.railway.app';

# 2. Run Flutter web (no sandbox issues)
cd /Users/maxroth/Documents/Programming/FlutterProjects/CLOAK
flutter run -d web

# 3. Try registration
# Check network tab in browser dev tools
# Should see requests going to Railway endpoint
```

---

## 📊 Debugging Checklist

- [ ] Checked Railway logs for errors
- [ ] Verified all environment variables set
- [ ] Tested health endpoint with curl
- [ ] Confirmed deployment status is green/running
- [ ] Fixed any config issues found
- [ ] Re-deployed if needed
- [ ] Waited 2-3 minutes for deploy
- [ ] Tested health endpoint again
- [ ] Tested registration endpoint
- [ ] Tested from Flutter app

---

## 🆘 If Still Not Working

### Option 1: Local Testing (Verification)

```bash
# Verify backend code works locally
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE
make start

# In another terminal
curl http://localhost:8080/health
# Should return {"status":"ok"}

# If this works, it's a Railway issue, not code issue
```

### Option 2: Change Deployment Provider

```bash
# If Railway continues to fail, try Render:
# 1. Create account at https://render.com
# 2. Follow docs/DEPLOYMENT_GUIDE.md for Render
# 3. Redeploy there instead
```

### Option 3: Local Tunnel (Temporary Solution)

```bash
# While troubleshooting Railway:
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE
make start       # Terminal 1
make tunnel      # Terminal 2 - gets ngrok URL

# Update Flutter to use ngrok URL temporarily
# lib/core/constants/app_constants.dart
# static const String apiBaseUrl = 'http://[ngrok-url]';

# Recompile Flutter
flutter run -d web
```

---

## 📝 For Next Session

**Before resuming:**

1. Check Railway logs at https://railway.app/project/[YOUR_PROJECT]
2. Note any error messages
3. Run this command to verify status:
    ```bash
    curl -v https://cloakbe-production.up.railway.app/health
    ```
4. If good, jump to "Testing After Fix" section
5. If bad, follow "Investigation Steps" section

**Key Resources:**

- Railway docs: https://docs.railway.app
- Health endpoint: https://cloakbe-production.up.railway.app/health
- Local fallback: `make start` then use localhost:8080
