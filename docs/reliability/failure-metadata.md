# Failure Metadata

Phase 5 records structured failure metadata on `TaskResult.Metadata` so callers
can inspect retry and DLQ decisions without parsing error strings.
