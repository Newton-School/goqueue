# Failure Metadata

Phase 5 records structured failure metadata on `TaskResult.Metadata` so callers
can inspect retry and DLQ decisions without parsing error strings.

## Keys

Metadata keys use the `goqueue.failure.*` prefix. Public constants include
`FailureMetadataCategoryKey`, `FailureMetadataAttemptKey`,
`FailureMetadataMaxAttemptsKey`, `FailureMetadataRetryableKey`,
`FailureMetadataNextRetryAtKey`, `FailureMetadataDeadLetteredKey`,
`FailureMetadataDeadLetteredAtKey`, and `FailureMetadataLastErrorKey`.

## Categories

Failure categories are stable public constants: `FailureExecution`,
`FailureMalformedMessage`, `FailureUnknownTask`, `FailureExpired`,
`FailureRetryExhausted`, and `FailureRetryScheduleFailed`.
