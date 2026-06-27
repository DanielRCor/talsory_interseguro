# Talsory Interseguro Challenge

This repository implements the Interseguro coding challenge with two small backend services:

- `api-go`: Go + Fiber API that validates an input matrix, computes a reduced QR factorization with Modified Gram-Schmidt, calls the Node API over HTTP, and returns the combined result.
- `api-node`: Node.js + Express + TypeScript API that receives the `Q` and `R` matrices and computes aggregate statistics.

Optional extensions implemented:

- JWT protection for `/api/v1/*` routes in both services when `ENABLE_AUTH=true`
- Matrix rotation helper endpoint in the Go API: `POST /api/v1/matrix/rotate`

## Architecture

```txt
Client
  -> Go API (`POST /api/v1/qr/analyze`)
    -> QR factorization
    -> HTTP call to Node API (`POST /api/v1/statistics`)
  -> Combined JSON response
```

## Technologies

- Go 1.26
- Fiber v2
- Node.js 22
- Express 5
- TypeScript
- Jest + Supertest
- Docker + Docker Compose

## Prerequisites

- Go 1.26+
- Node.js 22+
- npm
- Docker Desktop / Docker Engine with Compose

## Environment Variables

Copy values from `.env.example` if you want to override defaults.

```env
GO_API_PORT=8080
NODE_API_PORT=3000
NODE_API_URL=http://api-node:3000
HTTP_CLIENT_TIMEOUT_MS=3000
ENABLE_AUTH=false
JWT_SECRET=change-me-only-if-auth-enabled
```

## Local Run

### Node API

```bash
cd api-node
npm install
npm test
npm run build
npm run dev
```

### Go API

```bash
cd api-go
go test ./...
go run ./cmd/server
```

By default the Go API expects the Node API at `http://localhost:3000`.

## Docker Run

```bash
docker compose up --build
```

Services:

- Go API: `http://localhost:8080`
- Node API: `http://localhost:3000`

## Tests

```bash
cd api-node && npm test
cd api-go && go test ./...
```

## Example Requests

### Health checks

```bash
curl http://localhost:8080/health
curl http://localhost:3000/health
```

### End-to-end QR analysis

```bash
curl -X POST http://localhost:8080/api/v1/qr/analyze \
  -H "Content-Type: application/json" \
  -d '{"matrix":[[1,2],[3,4],[5,6]]}'
```

### Optional rotation endpoint

```bash
curl -X POST http://localhost:8080/api/v1/matrix/rotate \
  -H "Content-Type: application/json" \
  -d '{"matrix":[[1,2,3],[4,5,6]],"direction":"counterclockwise"}'
```

### Optional JWT flow

Generate a local test token:

```bash
node -e "const jwt=require('jsonwebtoken'); console.log(jwt.sign({sub:'demo-user'}, 'change-me-only-if-auth-enabled', {algorithm:'HS256'}))"
```

Then call protected routes with:

```bash
curl -X POST http://localhost:8080/api/v1/qr/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"matrix":[[1,2],[3,4],[5,6]]}'
```

## Technical Decisions

- Two APIs are kept separate because the challenge explicitly asks for one Go API and one Node API communicating via HTTP.
- The Go API is the orchestrator because it receives the original matrix in the prompt.
- QR factorization was prioritized over matrix rotation because the statement is more specific about QR.
- Modified Gram-Schmidt was chosen to keep the algorithm explainable and dependency-light.
- No database was added because the challenge does not require persistence.

## Known Limitations

- The QR implementation only supports matrices with `rows >= columns`.
- Linearly dependent columns are rejected instead of using a rank-deficient decomposition strategy.
- JWT is implemented but disabled by default to preserve the simple local developer flow.

## Documentation

- [Preflight report](C:/Users/theda/Documents/talsory_interseguro/repo/docs/preflight-report.md)
- [API examples](C:/Users/theda/Documents/talsory_interseguro/repo/docs/api-examples.md)
- [Deployment plan](C:/Users/theda/Documents/talsory_interseguro/repo/docs/deployment-plan.md)
- [Interview notes](C:/Users/theda/Documents/talsory_interseguro/repo/docs/interview-notes.md)
