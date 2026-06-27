# API Examples

## Health checks

### Go API

```bash
curl http://localhost:8080/health
```

```json
{
  "status": "ok",
  "service": "go-qr-api"
}
```

### Node API

```bash
curl http://localhost:3000/health
```

```json
{
  "status": "ok",
  "service": "node-statistics-api"
}
```

## Generar JWT demo

```bash
curl -X POST http://localhost:8080/auth/demo-token
```

## Analisis QR

```bash
curl -X POST http://localhost:8080/api/v1/qr/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d "{\"matrix\":[[1,2],[3,4],[5,6]]}"
```

## Rotacion de matriz

```bash
curl -X POST http://localhost:8080/api/v1/matrix/rotate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d "{\"matrix\":[[1,2,3],[4,5,6]],\"direction\":\"counterclockwise\"}"
```

## Estadisticas directas en Node

```bash
curl -X POST http://localhost:3000/api/v1/statistics \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d "{\"matrices\":{\"q\":[[1,0],[0,1]],\"r\":[[2,3],[0,4]]}}"
```

## Error de validacion esperado

```bash
curl -X POST http://localhost:8080/api/v1/qr/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN" \
  -d "{\"matrix\":[[1,2,3],[4,5,6]]}"
```

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Matrix validation failed",
    "details": [
      "matrix must satisfy rows >= columns, got 2 rows and 3 columns"
    ]
  }
}
```
