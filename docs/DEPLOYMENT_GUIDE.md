# CLOAK Backend - Deployment Guide

## Quick Deployment to Railway (Free Tier)

Railway.app provides free tier hosting perfect for testing. Follow these steps:

### Prerequisites

- GitHub account (for easy deployment)
- Railway account (free at railway.app)

### Step 1: Push to GitHub

```bash
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE

# Initialize git if needed
git init
git add .
git commit -m "Initial CLOAK backend"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/CLOAKBE.git
git push -u origin main
```

### Step 2: Deploy on Railway

1. Go to [railway.app](https://railway.app)
2. Click "Create New Project"
3. Select "Deploy from GitHub"
4. Connect your GitHub account
5. Select `CLOAKBE` repository
6. Railway will auto-detect the Dockerfile ✓

### Step 3: Add PostgreSQL Database

1. In Railway dashboard, click "Create New"
2. Select "Database" → "PostgreSQL"
3. Railway will automatically set `DATABASE_URL` environment variable

### Step 4: Update Environment Variables

In Railway dashboard, go to your deployment and set:

```env
ENVIRONMENT=production
JWT_SECRET=<generate-with-openssl-rand-base64-32>
HMAC_SECRET=<generate-with-openssl-rand-base64-32>
```

Generate secrets:

```bash
openssl rand -base64 32
```

### Step 5: Get Your Deployment URL

Railways shows your public URL like: `https://cloakbe-prod-xyz.railway.app`

Your API will be available at:

- Health: `https://cloakbe-prod-xyz.railway.app/health`
- Register: `https://cloakbe-prod-xyz.railway.app/api/v1/auth/business/register`

### Step 6: Update Flutter App

Update your Flutter app to use the deployed URL:

```dart
// lib/core/constants/app_constants.dart
static const String apiBaseUrl = 'https://cloakbe-prod-xyz.railway.app';
```

---

## Alternative: Deploy to Render (Also Free)

### Step 1: Connect GitHub

1. Go to [render.com](https://render.com)
2. Click "New" → "Web Service"
3. Connect GitHub repository
4. Select `CLOAKBE`

### Step 2: Configure Service

- **Name**: `cloak-api`
- **Region**: Choose closest to you
- **Branch**: `main`
- **Build Command**: `go build -o main ./cmd/api`
- **Start Command**: `./main`

### Step 3: Add PostgreSQL

- Click "Create Database"
- Use PostgreSQL v14+
- Render auto-sets `DATABASE_URL` in environment

### Step 4: Environment Variables

```env
ENVIRONMENT=production
JWT_SECRET=<your-secret>
HMAC_SECRET=<your-secret>
```

### Step 5: Deploy

- Render auto-deploys on git push
- Your URL: `https://cloak-api.onrender.com`

---

## Local Docker Deployment (For Testing)

If you want to deploy locally with Docker:

```bash
# Build image
docker build -t cloak-api:latest .

# Run with environment
docker run -e DATABASE_URL="postgres://user:pass@db:5432/cloak" \
           -e JWT_SECRET="your-secret" \
           -e HMAC_SECRET="your-secret" \
           -p 8080:8080 \
           cloak-api:latest
```

---

## Recommended Configuration

### For Development (Railway/Render Free Tier)

- **Database**: PostgreSQL (managed by Railway/Render)
- **Backend**: Deployed at `https://cloak-api-xxx.railway.app`
- **Cost**: FREE

### For Production

- **Database**: AWS RDS or Azure Database for PostgreSQL
- **Backend**: Kubernetes, AWS ECS, or Railway Pro Tier
- **Cost**: $20-50/month depending on traffic

---

## Troubleshooting

### "Connection failed" after deployment

1. Check environment variables are set (`ENVIRONMENT=production`, `DATABASE_URL`, etc.)
2. Verify database is running
3. Check logs: Railway dashboard → Logs tab

### Migrations not running

Currently migrations run manually via `make migrate`. For automatic:

- Add migration run to startup script
- Or use Render/Railway pre-deploy hooks

### "Invalid JWT Secret"

Generate new secret:

```bash
openssl rand -base64 32
```

---

## After Deployment

### Update Flutter App

Update the API endpoint in your Flutter app:

```dart
// lib/core/constants/app_constants.dart
static const String apiBaseUrl = 'https://your-deployed-url.railway.app';
```

Then hot reload:

```bash
flutter run -d macos
# Press R to hot reload
```

### Test Deployment

```bash
# Test health endpoint
curl https://your-deployed-url.railway.app/health

# Test business registration
curl -X POST https://your-deployed-url.railway.app/api/v1/auth/business/register \
  -H "Content-Type: application/json" \
  -d '{"name":"test","email":"test@example.com","password":"test123"}'
```

---

## Cost Breakdown

| Service              | Free Tier              | After Free          |
| -------------------- | ---------------------- | ------------------- |
| Railway Backend + DB | First $5/month         | $7/month (generous) |
| Render Backend       | Always free            | Pro pricing         |
| Render DB            | For 90 days            | $15/month           |
| AWS RDS              | t3.micro free (1 year) | $20-100/month       |

**Recommendation**: Start with Railway free tier for testing, upgrade only if needed.
