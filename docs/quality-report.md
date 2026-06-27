# Quality Report

## Tests Executed

- `cd api-node && npm test`
- `cd api-node && npm run build`
- `cd api-go && go test ./...`
- `cd frontend && npm run build`
- Local end-to-end run with Node on `:3000` and Go on `:8080`
- Browser-based end-to-end run from the frontend on `:5173`
- `docker compose up --build -d`

## Results

- Node tests: passed
- Node TypeScript build: passed
- Go tests: passed
- Frontend build: passed
- Local end-to-end HTTP integration: passed
- Frontend browser integration: passed
- Docker Compose: could not be validated because the Docker Desktop Linux engine was not running (`open //./pipe/dockerDesktopLinuxEngine: The system cannot find the file specified`)

## Coverage Summary

- Node statistics calculations for normal, diagonal, decimal, negative, and invalid inputs
- Node HTTP health and validation behavior
- Node HTTP JWT enforcement behavior
- Go matrix validation rules
- Go QR orthogonality, reconstruction, and triangularity
- Go matrix rotation behavior
- Go HTTP success, validation failure, and Node downstream failure handling
- Go HTTP JWT enforcement behavior
- Demo JWT issuance from the Go API
- Local Go -> Node integration with a real `POST /api/v1/qr/analyze` request
- Frontend rendering, health probes, and QR action flow against the live local APIs

## Bugs Found and Corrected During Implementation

1. PowerShell execution policy blocked `npm.ps1`, so `npm.cmd` was used consistently.
2. The local Codex session did not inherit the updated Go `PATH`, so the absolute Go binary path was used.
3. TypeScript 6 raised a deprecation warning for module resolution, so `ignoreDeprecations` was added to keep builds green.
4. The first counterclockwise rotation implementation wrote into uninitialized rows, which caused a panic in Go tests and was corrected by preallocating the rotated matrix.
5. Browser-based frontend requests needed CORS support, so permissive demo-friendly middleware was added to both APIs.
6. JWT was upgraded from a dormant optional field into a complete demo flow with server-issued tokens and frontend support.

## Remaining Risks

- The Node dependency tree reports moderate audit findings from transitive dev dependencies.
- The QR implementation rejects rank-deficient matrices instead of computing a more advanced decomposition.
- JWT is enabled by default in the current setup, so authenticated calls need a bearer token unless `ENABLE_AUTH=false`.
