# API Examples

## Go API Health

```bash
curl http://localhost:8080/health
```

```json
{
  "status": "ok",
  "service": "go-qr-api"
}
```

## Node API Health

```bash
curl http://localhost:3000/health
```

```json
{
  "status": "ok",
  "service": "node-statistics-api"
}
```

## Node Statistics Endpoint

```bash
curl -X POST http://localhost:3000/api/v1/statistics \
  -H "Content-Type: application/json" \
  -d '{"matrices":{"q":[[1,0],[0,1]],"r":[[2,3],[0,4]]}}'
```

```json
{
  "max": 4,
  "min": 0,
  "average": 1.375,
  "sum": 11,
  "hasDiagonalMatrix": true,
  "diagonalMatrices": ["q"]
}
```

## Go QR Analyze Endpoint

```bash
curl -X POST http://localhost:8080/api/v1/qr/analyze \
  -H "Content-Type: application/json" \
  -d '{"matrix":[[1,2],[3,4],[5,6]]}'
```

Response shape:

```json
{
  "input": {
    "rows": 3,
    "columns": 2
  },
  "qr": {
    "q": [[0.0]],
    "r": [[0.0]]
  },
  "statistics": {
    "max": 0,
    "min": 0,
    "average": 0,
    "sum": 0,
    "hasDiagonalMatrix": false,
    "diagonalMatrices": []
  },
  "metadata": {
    "algorithm": "modified-gram-schmidt",
    "tolerance": 1e-9
  }
}
```

Exact matrix values vary with floating-point output.

## Common Error Example

```bash
curl -X POST http://localhost:8080/api/v1/qr/analyze \
  -H "Content-Type: application/json" \
  -d '{"matrix":[[1,2,3],[4,5,6]]}'
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
