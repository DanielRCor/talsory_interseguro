# ADR 006 - Optional Rotation Endpoint

## Decision

Add a dedicated matrix rotation endpoint in the Go API instead of mixing rotation into the main QR contract.

## Rationale

The challenge text ambiguously mentions rotation, but QR is the primary required operation. A separate endpoint allows that ambiguity to be addressed without weakening the main QR-focused design.

## Consequence

The QR endpoint remains stable and interview-friendly, while the optional operation is still demonstrable if asked about the ambiguity.
