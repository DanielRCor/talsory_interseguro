# Quality Report

## Tests Executed

- `cd api-node && npm test`
- `cd api-node && npm run build`
- `cd api-go && go test ./...`
- Local end-to-end run with Node on `:3000` and Go on `:8080`
- `docker compose up --build -d`

## Results

- Node tests: passed
- Node TypeScript build: passed
- Go tests: passed
- Local end-to-end HTTP integration: passed
- Docker Compose: could not be validated because the Docker Desktop Linux engine was not running (`open //./pipe/dockerDesktopLinuxEngine: The system cannot find the file specified`)

## Coverage Summary

- Node statistics calculations for normal, diagonal, decimal, negative, and invalid inputs
- Node HTTP health and validation behavior
- Go matrix validation rules
- Go QR orthogonality, reconstruction, and triangularity
- Go HTTP success, validation failure, and Node downstream failure handling
- Local Go -> Node integration with a real `POST /api/v1/qr/analyze` request

## Bugs Found and Corrected During Implementation

1. PowerShell execution policy blocked `npm.ps1`, so `npm.cmd` was used consistently.
2. The local Codex session did not inherit the updated Go `PATH`, so the absolute Go binary path was used.
3. TypeScript 6 raised a deprecation warning for module resolution, so `ignoreDeprecations` was added to keep builds green.

## Remaining Risks

- The Node dependency tree reports moderate audit findings from transitive dev dependencies.
- The QR implementation rejects rank-deficient matrices instead of computing a more advanced decomposition.
- JWT authentication was intentionally not enabled because it was optional and outside the core delivery scope.
