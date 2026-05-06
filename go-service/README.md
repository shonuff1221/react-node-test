# Go PDF Report Service

Generates downloadable PDF reports for individual students.

## Setup

```bash
cd go-service
go mod download
```

## Running

```bash
# Set DATABASE_URL to match the backend's PostgreSQL connection
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/school_mgmt?sslmode=disable"
export PORT=8080

go run main.go
```

## API

- `GET /api/v1/students/{id}/report` — Generate and download a PDF report for a specific student
- `GET /healthz` — Health check endpoint

## Integration with Node Backend

The Node.js backend proxies `/api/v1/students/:id/report` to this Go service. Set the `GO_SERVICE_URL` environment variable in the backend `.env`:

```
GO_SERVICE_URL=http://localhost:8080
```

If the Go service is not running, the endpoint returns a 503 error.
