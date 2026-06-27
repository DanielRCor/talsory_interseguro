# ADR 001 - Two APIs

## Decision

Keep `api-go` and `api-node` as separate services.

## Rationale

The challenge explicitly asks for one API in Go and one API in Node.js communicating over HTTP.

## Consequence

This adds a bit of integration complexity, but it demonstrates service boundaries and inter-service communication clearly.
