# Workflow Canvas Rules

Use this file when changing `workflow/`, workflow metadata, canvas dispatch,
chain/group/chord behavior, Redis workflow state, or worker workflow
advancement.

## Responsibilities

- `workflow/` owns public canvas primitives: signatures, chains, groups,
  chords, workflow metadata, validation, and dispatch helpers.
- Canvas dispatch must use `producer.Producer`; do not duplicate producer
  enqueue logic in workflow code.
- Workers advance workflow state after task execution using metadata reserved by
  the canvas.
- Redis workflow storage lives behind backend contracts and is implemented in
  `redisbackend/`.

## Signature Dispatch

- Normalize signatures against the canvas default queue before dispatch.
- Preserve args, kwargs, metadata, timing, priority, and retry policy.
- Use copied slices and maps when passing user-provided values into producer
  calls.
- Reserved workflow metadata must override user metadata through
  `MergeMetadata` so control keys stay trustworthy.

## Chain Flow

1. Generate a workflow ID and first task ID.
2. Normalize and save the full chain record before dispatching the first
   signature.
3. Dispatch the first task with chain metadata containing kind, chain ID, and
   step index.
4. Worker advancement records completed step state atomically and dispatches the
   next signature only when the next step has not already been dispatched.
5. Failed chain steps do not dispatch later steps.

## Group Flow

1. Generate a group ID and task IDs for each child signature.
2. Save group state before dispatching child tasks.
3. Dispatch each child with group metadata containing kind, group ID, and index.
4. Workers record each child completion once, using task ID de-duplication.
5. Group completion succeeds only when all child tasks complete successfully.

## Chord Flow

- A chord is a header group plus a callback signature.
- Save the callback with group state before dispatching header tasks.
- Header tasks use chord/group metadata.
- Workers dispatch the callback once when all header tasks succeed.
- Callback tasks carry chord callback metadata.

## Tests

- Cover signature normalization, metadata merge precedence, chain first dispatch,
  chain next-step advancement, group completion, duplicate completion,
  chord callback dispatch, failure behavior, and scheduled workflow signatures.
- Add Redis integration coverage when workflow Lua scripts or Redis state shapes
  change.
