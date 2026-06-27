# Talsory Interseguro Challenge

Aplicacion compuesta por dos APIs y un frontend para resolver y probar el flujo completo del reto.

- `api-go`: valida la matriz, calcula la factorizacion QR y coordina la respuesta final.
- `api-node`: recibe `Q` y `R` y calcula estadisticas.
- `frontend`: interfaz web para probar el flujo, generar un JWT demo y consumir ambos servicios.

## Flujo

1. El frontend envia una matriz a `api-go`.
2. `api-go` calcula la factorizacion QR.
3. `api-go` llama por HTTP a `api-node`.
4. `api-node` devuelve las estadisticas.
5. `api-go` responde con el resultado consolidado.

## Stack

- Go + Fiber
- Node.js + Express + TypeScript
- Vite + TypeScript
- Docker + Docker Compose

## Requisitos

- Go 1.26+
- Node.js 22+
- npm
- Docker Desktop

## Variables de entorno

Puedes tomar como base el archivo `.env.example`.

```env
GO_API_PORT=8080
NODE_API_PORT=3000
FRONTEND_PORT=5173
NODE_API_URL=http://api-node:3000
HTTP_CLIENT_TIMEOUT_MS=3000
ENABLE_AUTH=true
JWT_SECRET=change-me-only-if-auth-enabled
```
## Link del despliegue
https://talsory-interseguro-frontend.onrender.com/
## Ejecutar localmente

### Con Docker

```bash
docker compose up --build
```

Servicios:

- Frontend: `http://localhost:5173`
- API Go: `http://localhost:8080`
- API Node: `http://localhost:3000`

### Manual

Si es la primera vez, instala dependencias en `api-node` y `frontend`.

```bash
cd api-node
npm install
npm run dev
```

```bash
cd api-go
go run ./cmd/server
```

```bash
cd frontend
npm install
npm run dev
```

## Endpoints principales

### Go API

- `GET /`
- `GET /health`
- `POST /auth/demo-token`
- `POST /api/v1/qr/analyze`
- `POST /api/v1/matrix/rotate`

### Node API

- `GET /`
- `GET /health`
- `POST /api/v1/statistics`

## JWT

Con `ENABLE_AUTH=true`, las rutas bajo `/api/v1` quedan protegidas.

Para pruebas locales y demo, `api-go` expone:

```bash
POST /auth/demo-token
```

Ese endpoint genera un token temporal que el frontend puede cargar automaticamente.

## Pruebas

```bash
cd api-go && go test ./...
cd api-node && npm test
cd frontend && npm run build
```

## Despliegue

El repositorio incluye `render.yaml` para desplegar:

- `talsory-interseguro-api-node`
- `talsory-interseguro-api-go`
- `talsory-interseguro-frontend`

En Render:

1. Conecta el repositorio.
2. Elige `New > Blueprint`.
3. Usa la rama `main`.
4. Confirma `render.yaml`.
5. Ejecuta el deploy.

## Ejemplos de uso

Hay ejemplos de requests y respuestas en [docs/api-examples.md](C:/Users/theda/Documents/talsory_interseguro/repo/docs/api-examples.md).
