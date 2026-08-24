# ESSL HRMS

HRMS platform with:
- Go backend APIs (Gin, MySQL, Redis)
- React frontend dashboard
- Lead verification workflow with cron-compatible service logic

## Repository Structure

- `backend/` Go services and shared packages
- `frontend/` React app

## Prerequisites

Install these before running locally:
- Go 1.25.x (module toolchain is go1.25.0)
- Node.js 18+ and npm
- MySQL or MariaDB
- Redis
- Reacher service running at `http://localhost:8080`

## Local Setup

### 1. Clone and open the repo

```bash
git clone <your-repo-url>
cd essl-hrms
```

### 2. Backend setup

```bash
cd backend
go mod download
```

Create database and import dump:

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS scalent_hrms;"
mysql -u root -p scalent_hrms < DB/scalent_hrms_24_August.sql
```

### 3. Configure backend

Backend config is loaded from:
- `backend/services/core/app/config.yaml`

Update these values for your machine:
- `db.host`, `db.user`, `db.pass`, `db.name`
- `cache.host`, `cache.port`
- `reacher.baseURL`
- optional SMTP values under `service`

Default local server port is:
- `:8081`

### 4. Start Redis (if not already running)

Example with Docker:

```bash
docker run --name hrms-redis -p 6379:6379 -d redis:7
```

### 5. Start backend

From `backend/`:

```bash
make run-core-local
```

This runs Wire in `services/core/app` and starts the API server.

Health check example:

```bash
curl http://localhost:8081/scalent-hrms/home
```

### 6. Frontend setup and run

In a new terminal:

```bash
cd frontend
npm install
npm start
```

Frontend API base URL is configured in:
- `frontend/src/API/apiConstants.js`

Default points to:
- `http://localhost:8081/scalent-hrms`

## Running the Project Locally (Quick Start)

Use this sequence:
1. Start MySQL
2. Start Redis
3. Ensure Reacher is up on `localhost:8080`
4. Run backend: `cd backend && make run-core-local`
5. Run frontend: `cd frontend && npm start`

## Cron: Schedule and Run Locally

## Current Status

Cron logic exists in:
- `backend/pkg/cron/cron.go`
- `backend/pkg/cron/verify_pending_leads.go`

Important: it is not currently started by default in `main.go`.

## How the schedule works

`StartCronJobs()` runs a loop:
- checks if pending leads exist
- verifies one pending lead
- sleeps for a randomized delay based on the user verification interval

Delay formula:
- base interval = user setting seconds (default 15s when not set)
- next run delay = base interval + random(0..base interval)

So with 15s base interval, effective delay is between 15s and 30s.

## Enable cron locally

To run scheduled cron locally, wire it during server initialization:

1. In `backend/services/core/app/wire.go`, add provider:

```go
service.NewLeadServiceImpl,
wire.Bind(new(service.LeadService), new(*service.LeadServiceImpl)),
```

(already present)

2. In `backend/services/core/app/main.go`, initialize and start cron after server setup:

```go
import "github.com/scalent.io/scalent-hrms/pkg/cron"

// after initServer(...)
cron.Init(registry.Options.LeadService)
if err := cron.StartCronJobs(); err != nil {
    log.Print(err)
}
```

3. Restart backend:

```bash
cd backend
make run-core-local
```

You should see cron logs in the backend console.

## Optional: Test cron behavior quickly

- Insert or upload leads with `PENDING` verification status
- Ensure related email list and user settings exist
- Watch backend logs for:
  - pending check
  - verification run
  - next randomized delay

## Useful Commands

Backend compile check:

```bash
cd backend
go test ./services/core/...
```

Frontend production build:

```bash
cd frontend
npm run build
```

## Notes

- Keep secrets (SMTP password, DB credentials) out of committed config when possible.
- If login fails for older records due to legacy password storage, use the latest backend code where password migration/rehash compatibility is implemented in login flow.
