# ngrok Local Tunnel Setup

Quick guide to expose your local backend publicly for testing.

## What is ngrok?

ngrok creates a secure tunnel to your local backend, giving you a public HTTPS URL. Perfect for:

- Testing Flutter app without deploying
- Sharing API with team members
- Testing webhooks locally

## Setup

### Step 1: ng rok is already installed

```bash
ngrok --version
# Output: ngrok version 3.36.1
```

### Step 2: Start Local Backend

Keep your backend running:

```bash
make start
# or
make backend
```

### Step 3: Start ngrok Tunnel (In New Terminal)

```bash
ngrok http 8080
```

You'll see:

```
Session Status                outbound
Account                       Username
Version                       3.36.1
Region                        us
Latency                       34ms
Web Interface                 http://127.0.0.1:4040
Forwarding                    https://abc123def456.ngrok.io -> http://localhost:8080

Connections                   ttl     opn     rt1     rt5     p50     p95
                              0       0       0.00    0.00    0.00    0.00
```

### Step 4: Copy Your Tunnel URL

Your public URL is: **`https://abc123def456.ngrok.io`**

## Update Flutter App

### Option A: Quick Manual Update (For Testing)

Edit `lib/core/constants/app_constants.dart`:

```dart
static const String apiBaseUrl = 'https://abc123def456.ngrok.io';
```

Then:

```bash
flutter run -d macos
# Press R to hot reload
```

### Option B: Environment-Based URLs (Professional)

Create an environment constant file:

```dart
// lib/core/constants/api_config.dart
class ApiConfig {
  // Set this based on your environment
  static const String environment = 'local'; // 'local' or 'production'

  static String get baseUrl {
    switch (environment) {
      case 'production':
        return 'https://cloak-api-prod.railway.app';
      case 'local':
        return 'https://your-ngrok-url.ngrok.io';
      default:
        return 'http://localhost:8080';
    }
  }
}
```

Then use in AppConstants:

```dart
import '../api_config.dart';

abstract class AppConstants {
  static const String apiBaseUrl = ApiConfig.baseUrl;
  // ... rest of constants
}
```

## Full Workflow

### Terminal 1: Backend

```bash
cd /Users/maxroth/Documents/Programming/Go/CLOAKBE
make start
# Output: Backend running on http://localhost:8080
```

### Terminal 2: ngrok Tunnel

```bash
ngrok http 8080
# Output: https://abc123def456.ngrok.io
```

### Terminal 3: Flutter App

```bash
cd /Users/maxroth/Documents/Programming/FlutterProjects/CLOAK

# Update URL to ngrok URL
# (in lib/core/constants/app_constants.dart)

flutter run -d macos
# Press R to hot reload with new URL
```

### Test Registration

- Open Flutter app
- Go to register
- Enter: name=test, email=test@test.com, password=test
- Should connect!

## ngrok Tips

### See Request Details

Open ngrok web interface: http://127.0.0.1:4040

- See all requests
- Inspect headers/body
- Replay requests

### URL Changes on Restart

Each `ngrok http 8080` gives a new URL. To get a static URL:

1. Buy ngrok plan ($5/month)
2. Use reserved domain
3. Or: keep ngrok running continuously

### Test Without Flutter

```bash
# Test health endpoint through ngrok
curl https://abc123def456.ngrok.io/health

# Test registration through ngrok
curl -X POST https://abc123def456.ngrok.io/api/v1/auth/business/register \
  -H "Content-Type: application/json" \
  -d '{"name":"test","email":"test@test.com","password":"pass123"}'
```

## When ngrok URL Changes

Every time you restart `ngrok http 8080`, you get a new URL. To make this easier:

### Create Makefile Target

Add to Makefile:

```makefile
.PHONY: tunnel

tunnel:
	@echo "Starting ngrok tunnel..."
	@echo "Copy the HTTPS URL below and update lib/core/constants/app_constants.dart"
	ngrok http 8080
```

Then run:

```bash
make tunnel
```

## Recommended Development Workflow

### For Daily Development (Use ngrok)

```bash
# Terminal 1
make start

# Terminal 2
make tunnel

# Terminal 3
flutter run -d macos
```

### For Deployment (Use Railway)

1. Push to GitHub
2. Deploy on Railway
3. Update Flutter app with Railway URL
4. Done!

## Troubleshooting

### "ngrok: command not found"

```bash
brew install ngrok
ngrok --version
```

### "Connection refused"

- Verify backend is running: `curl http://localhost:8080/health`
- ngrok must be in same terminal session or separate terminal

### Flutter app still can't connect

1. Verify ngrok URL is correct
2. Check Flutter hot reload happened (Press R)
3. Test ngrok directly: `curl https://YOUR_URL/health`
4. Check ngrok web interface: http://127.0.0.1:4040

## Cost

- Free tier: 1 live tunnel, changes URL on restart
- Pro tier: $5/month, static URL + more features
