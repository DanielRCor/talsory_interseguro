# ADR 004 - Modified Gram-Schmidt

## Decision

Use the Modified Gram-Schmidt algorithm for reduced QR factorization.

## Rationale

It is straightforward to implement, explain in an interview, and sufficient for the challenge scope.

## Consequence

The implementation supports matrices with `rows >= columns` and rejects rank-deficient inputs when the column norm collapses below tolerance.
