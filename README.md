# Talsory Interseguro Challenge

Implementacion de una prueba tecnica con dos APIs y un frontend simple para probar el flujo completo.

- `api-go`: recibe una matriz, valida la entrada, calcula la factorizacion QR y orquesta la respuesta final.
- `api-node`: recibe las matrices `Q` y `R` y calcula estadisticas.
- `frontend`: interfaz web para probar QR, rotacion, health checks y autenticacion JWT.

## Flujo del proyecto

1. El usuario envia una matriz desde el frontend.
2. La API en Go calcula la factorizacion QR.
3. La API en Go llama por HTTP a la API en Node.
4. La API en Node devuelve estadisticas sobre `Q` y `R`.
5. La API en Go responde con todo el resultado consolidado.

## Stack

- Go + Fiber
- Node.js + Express + TypeScript
- Vite + TypeScript
- Docker + Docker Compose

## Como levantarlo

### Opcion 1: Docker

Es la forma mas rapida de levantar todo junto.

1. Abre Docker Desktop y verifica que diga `Engine running`.
2. En la raiz del repo ejecuta:

```bash
docker compose up --build
```

Servicios disponibles:

- Frontend: `http://localhost:5173`
- API Go: `http://localhost:8080`
- API Node: `http://localhost:3000`

Si algun puerto ya esta ocupado, primero cierra los procesos locales que esten usando `5173`, `8080` o `3000`.

### Opcion 2: Manual

Si prefieres correr cada servicio por separado:

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

## Variables de entorno

Si quieres replicar la configuracion local recomendada, crea un archivo `.env` en la raiz usando como base `.env.example`.

Valores principales:

```env
GO_API_PORT=8080
NODE_API_PORT=3000
FRONTEND_PORT=5173
NODE_API_URL=http://api-node:3000
HTTP_CLIENT_TIMEOUT_MS=3000
ENABLE_AUTH=true
JWT_SECRET=change-me-only-if-auth-enabled
```

## Autenticacion JWT

Las rutas principales pueden protegerse con JWT cuando `ENABLE_AUTH=true`.

Para no depender de un proveedor externo, el proyecto expone un endpoint local de apoyo:

```bash
POST /auth/demo-token
```

Ese endpoint genera un token de prueba con expiracion corta para usar el frontend o probar las rutas protegidas.

## Endpoints principales

```bash
GET  /health
POST /auth/demo-token
POST /api/v1/qr/analyze
POST /api/v1/matrix/rotate
POST /api/v1/statistics
```

## Ejemplos rapidos

### Health check

```bash
curl http://localhost:8080/health
curl http://localhost:3000/health
```

### Analisis QR

```bash
curl -X POST http://localhost:8080/api/v1/qr/analyze \
  -H "Content-Type: application/json" \
  -d "{\"matrix\":[[1,2],[3,4],[5,6]]}"
```

### Token demo

```bash
curl -X POST http://localhost:8080/auth/demo-token
```

## Pruebas

```bash
cd api-go && go test ./...
cd api-node && npm test
cd frontend && npm run build
```

## Despliegue en Render

El repo ya incluye [render.yaml](C:/Users/theda/Documents/talsory_interseguro/repo/render.yaml) para crear los 3 servicios:

- `talsory-interseguro-api-node`
- `talsory-interseguro-api-go`
- `talsory-interseguro-frontend`

Pasos:

1. En Render, conecta este repositorio de GitHub.
2. Elige `New > Blueprint`.
3. Selecciona la rama `main`.
4. Confirma el archivo `render.yaml`.
5. Lanza el deploy.

Notas importantes:

- El frontend quedara como sitio estatico.
- Las dos APIs quedaran como web services free.
- `api-go` usara automaticamente la URL publica de `api-node`.
- El frontend usara automaticamente las URLs publicas de ambas APIs.
- En free tier, los servicios web pueden dormir si no reciben trafico por 15 minutos.

## Guia corta para reclutadores

Este proyecto busca mostrar tres cosas:

- comunicacion entre servicios en tecnologias distintas
- resolucion de un problema numerico real con Go
- una interfaz minima para demostrar el flujo sin depender solo de Postman

Decisiones principales:

- Go actua como servicio orquestador porque recibe la matriz original y ejecuta la parte numerica.
- Node se mantiene separado para cumplir el requerimiento de tener dos APIs comunicandose por HTTP.
- No se agrego base de datos porque la prueba no lo necesita.
- Se incluyo Docker para levantar todo con un solo comando.

## Documentacion adicional

- [API examples](C:/Users/theda/Documents/talsory_interseguro/repo/docs/api-examples.md)
- [Deployment plan](C:/Users/theda/Documents/talsory_interseguro/repo/docs/deployment-plan.md)
- [Interview notes](C:/Users/theda/Documents/talsory_interseguro/repo/docs/interview-notes.md)
- [Preflight report](C:/Users/theda/Documents/talsory_interseguro/repo/docs/preflight-report.md)
- [Quality report](C:/Users/theda/Documents/talsory_interseguro/repo/docs/quality-report.md)
