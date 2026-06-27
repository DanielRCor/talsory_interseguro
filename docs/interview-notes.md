# Interview Notes

## What Was Built

Two decoupled APIs:

- Go/Fiber API for QR factorization and orchestration.
- Node/Express API for matrix statistics.
- Optional Vite frontend for local demo and interviewer walkthroughs.

## Why Go Does QR

The challenge assigns the primary matrix operation and original input handling to the Go service.

## Why Node Does Statistics

The challenge states that the second API should receive the matrices produced by Go and compute additional operations.

## Why No Database

There is no persistence requirement in the prompt, so adding a database would increase complexity without improving the solution.

## Why Docker

Docker was required by the challenge and makes both services reproducible locally and in future deployment environments.

## How the Rotation vs QR Ambiguity Was Handled

QR factorization was prioritized because it is the explicit and technically concrete requirement. The ambiguity is documented instead of broadening the scope.

## How Matrix Validation Works

- Matrix must exist and be non-empty.
- Rows must be non-empty and rectangular.
- Values must be finite numbers.
- The QR implementation supports `rows >= columns`.

## How Quality Was Ensured

- Unit tests for Node statistics logic.
- HTTP tests for the Node API.
- Unit tests for Go QR properties and validation.
- HTTP tests for the Go API handler.
- Docker setup for full local execution.

## Optional Features Added

- JWT auth can be enabled on both services with `ENABLE_AUTH=true`.
- A matrix rotation endpoint was added in Go to cover the ambiguous “rotation” mention without replacing QR as the primary operation.
- A lightweight frontend was added to make the end-to-end behavior easier to demonstrate without Postman or curl.
