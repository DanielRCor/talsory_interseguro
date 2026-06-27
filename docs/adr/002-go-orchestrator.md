# ADR 002 - Go as the Orchestrator

## Decision

The client calls the Go API first, and the Go API calls the Node API.

## Rationale

The challenge statement assigns the original matrix input and the main matrix operation to the Go service.

## Consequence

The Go API becomes the main entrypoint and owns the final response contract.
