# Task Model Rules

Use this file when changing `task/`, task-related root exports, payload codecs,
metadata, handlers, retry policy, timing, validation, or result/state types.

## Ownership

- `task/` is Redis-independent. It must not import `redisbackend`, Redis
  clients, worker runtime code, CLI code, or application configuration.
- Task envelopes are the SDK-level task invocation model. Backend messages are
  the serialized storage/delivery model.
- Keep validation in the task model where possible so producers, workers,
  schedulers, workflows, and admin paths share the same safety checks.

## Envelope And Message Flow

1. `NewTaskEnvelope` applies defaults for ID, priority, retry policy, created
   time, payload, and metadata.
2. `TaskEnvelope.Validate` verifies ID, task name, queue name, priority, retry
   policy, timing, and attempt count.
3. `TaskEnvelopeToMessage` encodes payloads with the selected `PayloadCodec` and
   clones mutable byte slices/maps.
4. `TaskMessageToEnvelope` validates message identity fields before decoding
   payloads and reconstructing a validated envelope.

## Validation Rules

- Keep task IDs, task names, queue names, namespace-compatible names, priorities,
  retry policies, timings, attempts, and states strictly validated.
- Preserve defensive copies for payload args, kwargs, metadata maps, and raw
  payload bytes.
- Avoid using `map[string]any` values directly after accepting user input unless
  they are copied or encoded through the codec boundary.
- Default retry policy and priority should remain centralized in `task/`.

## Handlers And Results

- Handler code receives `HandlerContext` and `TaskPayload`; it should not depend
  on Redis implementation details.
- Task results should carry public metadata keys for failure categories and
  retry/dead-letter status.
- Do not store secrets in payloads, results, or metadata. These values may be
  persisted and exposed through inspect/admin surfaces.

## Tests

- Add tests for every new validation rule, default, clone behavior, codec edge
  case, and handler/result contract.
- Prefer table tests for validation matrices.
- Include malformed message or malformed payload coverage when changing decode
  paths.

## Documentation

- Update `docs/docs/concepts/task-model.md` when envelope validation, defaults,
  timing, metadata, or result semantics change.
- Apply docs updates through `docs-sync.md` when user-facing task model behavior
  changes.
