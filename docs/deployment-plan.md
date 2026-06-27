# Deployment Plan

## Rule

No real deployment was executed. This file only prepares the manual path.

## Recommended Order

1. Deploy `api-node`.
2. Validate `GET /health`.
3. Copy the public Node URL.
4. Deploy `api-go` with `NODE_API_URL` pointing to the deployed Node API.
5. Validate `GET /health`.
6. Validate `POST /api/v1/qr/analyze`.

## Environment Variables

### Go API

```env
PORT=8080
GO_API_PORT=8080
NODE_API_URL=https://your-node-service-url
HTTP_CLIENT_TIMEOUT_MS=3000
ENABLE_AUTH=false
JWT_SECRET=change-me
```

### Node API

```env
PORT=3000
NODE_API_PORT=3000
ENABLE_AUTH=false
JWT_SECRET=change-me
```

## Platform Options Without Railway

### Render

- Good fit for Docker-based web services.
- Straightforward service-per-container setup.
- Risk: free-tier cold starts or quota constraints depending on current plan availability.

### Fly.io

- Good fit for containerized apps and private service-to-service networking.
- Risk: operational concepts are slightly more involved than Render for beginners.

### Google Cloud Run

- Good fit for stateless containers with HTTP health checks.
- Risk: requires more cloud setup and IAM familiarity than the simpler developer platforms.

### Vercel

- Not the preferred choice for this challenge because the requirement is centered on two Dockerized backend services.
- Could be considered only if the architecture is adapted away from the original Docker-first expectation.

## Manual Checklist Before Deployment

- Tests pass locally.
- `docker compose up --build` works locally.
- `.env.example` is complete.
- No real secrets are committed.
- README and examples are up to date.
- User explicitly approves any real deployment.
