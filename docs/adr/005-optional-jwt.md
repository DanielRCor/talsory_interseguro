# ADR 005 - Optional JWT Protection

## Decision

Implement JWT authentication as an optional feature in both APIs, disabled by default.

## Rationale

Authentication was listed as optional in the challenge. Adding it behind `ENABLE_AUTH=true` extends the solution without complicating the default local flow.

## Consequence

Protected routes require a bearer token only when the flag is enabled, and documentation must show how to generate a local test token.
